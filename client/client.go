package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client holds all the information required to connect to a server
type Client struct {
	tenantId    string
	accessToken string
	httpClient  *http.Client
}

func NewClient(tenantId string, accessToken string) *Client {
	return &Client{
		tenantId:    tenantId,
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// 200 response code
func (c *Client) requestResponse200(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		resBody := new(bytes.Buffer)
		_, err = resBody.ReadFrom(res.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", res.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s - %s", res.StatusCode, req.URL, resBody.String())
	}

	return body, nil
}

// 201 response code
//
//nolint:unused
func (c *Client) requestResponse201(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		resBody := new(bytes.Buffer)
		_, err = resBody.ReadFrom(res.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 201 status code: %v", res.StatusCode)
		}
		return nil, fmt.Errorf("got a non 201 status code: %v - %s - %s", res.StatusCode, req.URL, resBody.String())
	}

	return body, nil
}

// 204 response code
func (c *Client) requestResponse204(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusNoContent {
		resBody := new(bytes.Buffer)
		_, err = resBody.ReadFrom(res.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 204 status code: %v", res.StatusCode)
		}
		return nil, fmt.Errorf("got a non 204 status code: %v - %s - %s", res.StatusCode, req.URL, resBody.String())
	}

	return body, nil
}
