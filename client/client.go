package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"
)

const (
	tokenType = "Bearer"
)

type Config struct {
	Hostname          string
	Resource          string
	ApplicationID     string
	ApplicationSecret string
	HttpClient        *http.Client
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
		return nil, fmt.Errorf("Missing Logto hostname")
	}

	return &Client{
		conf: config,
	}, nil
}

type request struct {
	method string
	path   string
	body   any
}

func (r *request) toHttpRequest(ctx context.Context, conf *Config) (*http.Request, error) {
	url := fmt.Sprintf("https://%s/%s", conf.Hostname, r.path)
	req, err := http.NewRequestWithContext(ctx, r.method, url, nil)
	if err != nil {
		return nil, err
	}

	if r.body != nil {
		body, err := json.Marshal(r.body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(body))

		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) getAccessToken(ctx context.Context) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("resource", fmt.Sprintf("https://%s/api", c.conf.Hostname))
	data.Set("scope", "all")
	data.Encode()

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("https://%s/oidc/token", c.conf.Hostname),
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.conf.ApplicationID, c.conf.ApplicationSecret)

	resp, err := expect(200)(c.conf.HttpClient.Do(req))
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

	return decodedResponse.AccessToken, nil
}

func (c *Client) do(ctx context.Context, r *request) (*http.Response, error) {
	req, err := r.toHttpRequest(ctx, c.conf)
	if err != nil {
		return nil, err
	}

	accessToken, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	return c.conf.HttpClient.Do(req)
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
			return res, fmt.Errorf("got status code %s, expected status code in %v", res.Status, codes)
		}
		return res, nil
	}
}
