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
	tenantId    string // 7rx5ya
	accessToken string // pat_bcAPnZs7RvafFzhbxRx8OqNCQkBdDUP0
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

// 200 response
func (c *Client) RequestResponse200(req *http.Request) ([]byte, error) {
	req.SetBasicAuth(c.accessToken, "")
	time.Sleep(1 * time.Second)
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

// 200 response, body (mobielconfig) return + SHA256 from header
func (c *Client) RequestResponse200Profile(req *http.Request) (string, string, error) {
	req.SetBasicAuth(c.accessToken, "")
	time.Sleep(1 * time.Second)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody := new(bytes.Buffer)
		_, err = resBody.ReadFrom(res.Body)
		if err != nil {
			return "", "", fmt.Errorf("got a non 204 status code: %v", res.StatusCode)
		}
		return "", "", fmt.Errorf("got a non 204 status code: %v - %s", res.StatusCode, resBody.String())
	}

	sha := res.Header.Get("etag")[3:35]
	resBody := new(bytes.Buffer)
	_, err = resBody.ReadFrom(res.Body)
	if err != nil {
		return "", "", fmt.Errorf("error open file: %v", err)
	}

	bodyString := resBody.String()

	return bodyString, sha, nil
}

// 201 response code
func (c *Client) RequestResponse201(req *http.Request) ([]byte, error) {
	req.SetBasicAuth(c.accessToken, "")
	time.Sleep(1 * time.Second)
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

// 202 response code
func (c *Client) RequestResponse202(req *http.Request) ([]byte, error) {
	req.SetBasicAuth(c.accessToken, "")
	time.Sleep(1 * time.Second)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusAccepted {
		resBody := new(bytes.Buffer)
		_, err = resBody.ReadFrom(res.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 202 status code: %v", res.StatusCode)
		}
		return nil, fmt.Errorf("got a non 202 status code: %v - %s - %s", res.StatusCode, req.URL, resBody.String())
	}

	return body, nil
}

// 204 response
func (c *Client) RequestResponse204(req *http.Request) ([]byte, error) {
	req.SetBasicAuth(c.accessToken, "")
	time.Sleep(1 * time.Second)
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
		return nil, fmt.Errorf("got a non 204 status code: %v - %s", res.StatusCode, resBody.String())
	}

	return body, nil
}

// 204 or 409 response
func (c *Client) RequestResponse204or409(req *http.Request) ([]byte, error) {
	req.SetBasicAuth(c.accessToken, "")
	time.Sleep(1 * time.Second)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusConflict {
		//profiles was not present or already assigned
		return nil, nil
	}

	if res.StatusCode != http.StatusNoContent {
		resBody := new(bytes.Buffer)
		_, err = resBody.ReadFrom(res.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 204 or 409 status code: %v", res.StatusCode)
		}
		return nil, fmt.Errorf("got a non 204 or 409 status code: %v - %s", res.StatusCode, resBody.String())
	}

	return body, nil
}
