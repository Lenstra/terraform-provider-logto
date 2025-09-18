package client

import (
	"context"
	"net/http"
	"path"
)

func (c *Client) ApiResourceGet(ctx context.Context, id string) (*ApiResourceModel, error) {
	if id == "" {
		return nil, errEmptyID
	}

	req := &request{
		method: http.MethodGet,
		path:   path.Join("api/resources", id),
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

func (c *Client) ApiResourceList(ctx context.Context, query_params map[string]string) (*[]ApiResourceModel, error) {
	req := &request{
		method:           http.MethodGet,
		path:             path.Join("api/resources"),
		query_parameters: query_params,
	}
	res, err := expect(200)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var apiResources []ApiResourceModel
	if err := decode(res.Body, &apiResources); err != nil {
		return nil, err
	}
	return &apiResources, nil
}

func (c *Client) ApiResourceCreate(ctx context.Context, apiResource *ApiResourceModel) (*ApiResourceModel, error) {
	req := &request{
		method: http.MethodPost,
		path:   "api/resources",
		body:   apiResource,
	}

	res, err := expect(201)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var returnApiResource ApiResourceModel
	if err := decode(res.Body, &returnApiResource); err != nil {
		return nil, err
	}
	return &returnApiResource, nil
}

func (c *Client) ApiResourceDelete(ctx context.Context, id string) error {
	if id == "" {
		return errEmptyID
	}

	req := &request{
		method: http.MethodDelete,
		path:   path.Join("api/resources", id),
	}
	_, err := expect(204)(c.do(ctx, req))
	return err
}

func (c *Client) ApiResourceUpdate(ctx context.Context, apiResource *ApiResourceModel) (*ApiResourceModel, error) {
	req := &request{
		method: http.MethodPatch,
		path:   path.Join("api/resources/", apiResource.ID),
		body:   apiResource,
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
