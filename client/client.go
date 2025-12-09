package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	utils "github.com/Lenstra/go-utils/http"
	"github.com/rs/zerolog"
)

const (
	tokenType = "Bearer"
)

var (
	errEmptyID = errors.New("id should not be empty")
)

type Config struct {
	Hostname          string
	Resource          string
	ApplicationID     string
	ApplicationSecret string

	Logger     zerolog.Logger
	HttpClient *http.Client
}

func DefaultConfig() *Config {

	return &Config{
		Hostname:          os.Getenv("LOGTO_HOSTNAME"),
		Resource:          os.Getenv("LOGTO_RESOURCE"),
		ApplicationID:     os.Getenv("LOGTO_APPLICATION_ID"),
		ApplicationSecret: os.Getenv("LOGTO_APPLICATION_SECRET"),
		HttpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type Client struct {
	conf *Config

	accessTokenLock    sync.Mutex
	accessToken        string
	accessTokenExpires time.Time
}

func NewClient(config *Config) (*Client, error) {
	defConfig := DefaultConfig()
	if config.Hostname == "" {
		config.Hostname = defConfig.Hostname
	}
	switch {
	case config.Resource != "":
		break
	case defConfig.Resource != "":
		config.Resource = defConfig.Resource
	default:
		config.Resource = fmt.Sprintf("https://%s/api", config.Hostname)
	}
	if config.ApplicationID == "" {
		config.ApplicationID = defConfig.ApplicationID
	}
	if config.ApplicationSecret == "" {
		config.ApplicationSecret = defConfig.ApplicationSecret
	}
	if config.HttpClient == nil {
		config.HttpClient = defConfig.HttpClient
	}

	if config.Hostname == "" {
		return nil, fmt.Errorf("missing Logto hostname")
	}

	return &Client{
		conf: config,
	}, nil
}

type request struct {
	method                             string
	path                               string
	body                               any
	rawBody                            io.Reader
	headers                            map[string]string
	queryParameters                    map[string]string
	application_id, application_secret string
}

func (r *request) toHttpRequest(ctx context.Context, conf *Config) (*http.Request, error) {
	reqUrl := fmt.Sprintf("https://%s/%s", conf.Hostname, r.path)

	req, err := http.NewRequestWithContext(ctx, r.method, reqUrl, nil)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	for key, value := range r.queryParameters {
		params.Add(key, value)
	}
	req.URL.RawQuery = params.Encode()

	if r.body != nil {
		body, err := json.Marshal(r.body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(body))

		req.Header.Set("Content-Type", "application/json")
	}

	if r.rawBody != nil {
		req.Body = io.NopCloser(r.rawBody)
	}

	for key, value := range r.headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func (c *Client) getAccessToken(ctx context.Context) (string, error) {
	c.accessTokenLock.Lock()
	defer c.accessTokenLock.Unlock()
	if c.accessToken != "" && c.accessTokenExpires.Before(time.Now()) {
		return c.accessToken, nil
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("resource", c.conf.Resource)
	data.Set("scope", "all")
	data.Encode()

	req := &request{
		method:             "POST",
		path:               "oidc/token",
		application_id:     c.conf.ApplicationID,
		application_secret: c.conf.ApplicationSecret,
		rawBody:            strings.NewReader(data.Encode()),
		headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	}

	resp, err := expect(200)(c.do(ctx, req))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	type Response struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	var decodedResponse Response
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&decodedResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if decodedResponse.TokenType != tokenType {
		return "", fmt.Errorf("unexpected token type %q, expected %q", decodedResponse.TokenType, tokenType)
	}

	c.accessToken = decodedResponse.AccessToken
	c.accessTokenExpires = time.Now().Add(time.Duration(float64(decodedResponse.ExpiresIn) * 0.7))
	return c.accessToken, nil
}

func (c *Client) do(ctx context.Context, r *request) (*http.Response, error) {
	req, err := r.toHttpRequest(ctx, c.conf)
	if err != nil {
		return nil, err
	}

	if r.application_id != "" {
		req.SetBasicAuth(r.application_id, r.application_secret)
	} else {
		accessToken, err := c.getAccessToken(ctx)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	err = utils.LogRequest(c.conf.Logger.Trace(), req, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.conf.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	err = utils.LogResponse(c.conf.Logger.Trace(), resp, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func decode(r io.ReadCloser, out any) error {
	defer r.Close()
	dec := json.NewDecoder(r)
	return dec.Decode(out)
}

func expect(codes ...int) func(*http.Response, error) (*http.Response, error) {
	return func(res *http.Response, err error) (*http.Response, error) {
		if err != nil {
			return nil, err
		}
		if !slices.Contains(codes, res.StatusCode) {
			message := fmt.Sprintf("got status code %s, expected status code in %v", res.Status, codes)
			if res.StatusCode >= 400 {
				content, _ := io.ReadAll(res.Body)
				message += fmt.Sprintf(": %s", string(content))
			}
			return res, errors.New(message)
		}
		return res, nil
	}
}
