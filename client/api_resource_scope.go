// Permissions contains in an api_resource and used with their Id's by roles directly

package client

import (
	"context"
	"net/http"
	"path"
)

func (c *Client) ApiResourceScopesList(ctx context.Context, resourceId string, query_params map[string]string) ([]ScopeModel, error) {
	req := &request{
		method:          http.MethodGet,
		path:            path.Join("api/resources", resourceId, "scopes"),
		queryParameters: query_params,
	}

	res, err := expect(200)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var returnScope []ScopeModel
	if err := decode(res.Body, &returnScope); err != nil {
		return nil, err
	}

	return returnScope, nil
}

func (c *Client) ApiResourceScopeCreate(ctx context.Context, resourceId string, scope *ScopeModel) (*ScopeModel, error) {
	req := &request{
		method: http.MethodPost,
		path:   path.Join("api/resources", resourceId, "scopes"),
		body:   scope,
	}

	res, err := expect(201)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var returnScope ScopeModel
	if err := decode(res.Body, &returnScope); err != nil {
		return nil, err
	}

	return &returnScope, nil
}

func (c *Client) ApiResourceScopeDelete(ctx context.Context, resourceId string, scopeId string) error {
	if resourceId == "" || scopeId == "" {
		return errEmptyID
	}

	req := &request{
		method: http.MethodDelete,
		path:   path.Join("api/resources", resourceId, "scopes", scopeId),
	}

	_, err := expect(204)(c.do(ctx, req))
	return err
}

func (c *Client) ApiResourceScopeUpdate(ctx context.Context, scope *ScopeModel) (*ScopeModel, error) {
	if scope.ResourceId == "" || scope.ID == "" {
		return nil, errEmptyID
	}

	req := &request{
		method: http.MethodPatch,
		path:   path.Join("api/resources", scope.ResourceId, "scopes", scope.ID),
		body:   scope,
	}

	res, err := expect(200)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var returnScope ScopeModel
	if err := decode(res.Body, &returnScope); err != nil {
		return nil, err
	}
	return &returnScope, nil
}
