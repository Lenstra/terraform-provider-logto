package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"time"
)

type Config struct {
	TenantID    string
	AccessToken string
	HttpClient  *http.Client
}

func DefaultConfig() *Config {
	return &Config{
		TenantID:    os.Getenv("LOGTO_TENANT_ID"),
		AccessToken: os.Getenv("LOGTO_ACCESS_TOKEN"),
		HttpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type Client struct {
	conf *Config
}

func NewClient(config *Config) *Client {
	defConfig := DefaultConfig()
	if config.TenantID == "" {
		config.TenantID = defConfig.TenantID
	}
	if config.AccessToken == "" {
		config.AccessToken = defConfig.AccessToken
	}
	if config.HttpClient == nil {
		config.HttpClient = defConfig.HttpClient
	}

	return &Client{
		conf: config,
	}
}

type request struct {
	method string
	path   string
	body   any
}

func (r *request) toHttpRequest(ctx context.Context, conf *Config) (*http.Request, error) {
	url := fmt.Sprintf("https://%s.logto.app/%s", conf.TenantID, r.path)
	req, err := http.NewRequestWithContext(ctx, r.method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", conf.AccessToken))

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

func (c *Client) do(ctx context.Context, r *request) (*http.Response, error) {
	req, err := r.toHttpRequest(ctx, c.conf)
	if err != nil {
		return nil, err
	}
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
