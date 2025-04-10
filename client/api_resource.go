package client

import (
	"context"
	"net/http"
)

func (c *Client) ApiResourceGet(ctx context.Context, id string) (*ApiResourceModel, error) {
	if id == "" {
		return nil, errEmptyID
	}

	req := &request{
		method: http.MethodGet,
		path:   "api/resources/" + id,
	}
	res, err := expect(200, 404)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, nil
	}

	var apiResource ApiResourceModel
	if err := decode(res.Body, &apiResource); err != nil {
		return nil, err
	}
	return &apiResource, nil
}

func (c *Client) ApiResourceCreate(ctx context.Context, ApiResource *ApiResourceModel) (*ApiResourceModel, error) {
	req := &request{
		method: http.MethodPost,
		path:   "api/resources",
		body:   ApiResource,
	}

	res, err := expect(201)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var apiResource ApiResourceModel
	if err := decode(res.Body, &apiResource); err != nil {
		return nil, err
	}
	return &apiResource, nil
}

func (c *Client) ApiResourceDelete(ctx context.Context, id string) error {
	if id == "" {
		return errEmptyID
	}

	req := &request{
		method: http.MethodDelete,
		path:   "api/resources/" + id,
	}
	_, err := expect(204)(c.do(ctx, req))
	return err
}

func (c *Client) ApiResourceUpdate(ctx context.Context, ApiResource *ApiResourceModel) (*ApiResourceModel, error) {
	req := &request{
		method: http.MethodPatch,
		path:   "api/resources/" + ApiResource.ID,
		body:   ApiResource,
	}

	res, err := expect(200)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var returnApiResource ApiResourceModel
	if err := decode(res.Body, &returnApiResource); err != nil {
		return nil, err
	}
	return &returnApiResource, nil
}
