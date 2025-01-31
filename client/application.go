package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) ApplicationGet(id string) (*ApplicationModel, error) {
	url := fmt.Sprintf("https://%s.logto.app/api/applications/%s", c.tenantId, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.requestResponse200(req)
	if err != nil {
		return nil, err
	}

	application := ApplicationModel{}
	err = json.Unmarshal(body, &application)
	if err != nil {
		return nil, err
	}

	return &application, nil
}

func (c *Client) ApplicationCreate(name, description, appType string) (*ApplicationModel, error) {
	url := fmt.Sprintf("https://%s.logto.app/api/applications", c.tenantId)

	jsonBody, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
		"type":        appType,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	body, err := c.requestResponse200(req)
	if err != nil {
		return nil, err
	}

	application := ApplicationModel{}
	err = json.Unmarshal(body, &application)
	if err != nil {
		return nil, err
	}

	return &application, nil
}

func (c *Client) ApplicationDelete(id string) error {
	url := fmt.Sprintf("https://%s.logto.app/api/applications/%s", c.tenantId, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	_, err = c.requestResponse204(req)
	return err
}

func (c *Client) ApplicationUpdate(id string, name string, description string) (*ApplicationModel, error) {
	fmt.Printf("ID: %+v\n", id)

	url := fmt.Sprintf("https://%s.logto.app/api/applications/%s", c.tenantId, id)

	jsonBody, err := json.Marshal(map[string]string{
		"name":        name,
		"description": description,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	body, err := c.requestResponse200(req)
	if err != nil {
		return nil, err
	}

	application := &ApplicationModel{}
	err = json.Unmarshal(body, application)
	if err != nil {
		return nil, err
	}

	return application, nil
}
