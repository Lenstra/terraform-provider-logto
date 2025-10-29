package client

import (
	"context"
	"net/http"
	"path"
)

func (c *Client) RoleGet(ctx context.Context, id string) (*RoleModel, error) {
	if id == "" {
		return nil, errEmptyID
	}

	req := &request{
		method: http.MethodGet,
		path:   path.Join("api/roles", id),
	}

	res, err := expect(200, 404)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, nil
	}

	var role RoleModel
	if err := decode(res.Body, &role); err != nil {
		return nil, err
	}
	return &role, nil
}

func (c *Client) RoleScopesGet(ctx context.Context, roleId string) ([]ScopeModel, error) {
	if roleId == "" {
		return nil, errEmptyID
	}

	req := &request{
		method: http.MethodGet,
		path:   path.Join("api/roles", roleId, "scopes"),
	}

	res, err := expect(200, 404)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, nil
	}

	var roleScopes []ScopeModel
	if err := decode(res.Body, &roleScopes); err != nil {
		return nil, err
	}
	return roleScopes, nil
}

func (c *Client) RoleCreate(ctx context.Context, role *RoleModel) (*RoleModel, error) {
	req := &request{
		method: http.MethodPost,
		path:   "api/roles",
		body:   role,
	}

	res, err := expect(200)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var returnRole RoleModel
	if err := decode(res.Body, &returnRole); err != nil {
		return nil, err
	}
	return &returnRole, nil
}

func (c *Client) RoleDelete(ctx context.Context, id string) error {
	if id == "" {
		return errEmptyID
	}

	req := &request{
		method: http.MethodDelete,
		path:   path.Join("api/roles", id),
	}

	_, err := expect(204)(c.do(ctx, req))
	return err
}

func (c *Client) RoleUpdate(ctx context.Context, role *RoleModel) (*RoleModel, error) {
	if role.ID == "" {
		return nil, errEmptyID
	}

	req := &request{
		method: http.MethodPatch,
		path:   path.Join("api/roles", role.ID),
		body:   role,
	}

	res, err := expect(200)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var returnRole RoleModel
	if err := decode(res.Body, &returnRole); err != nil {
		return nil, err
	}
	return &returnRole, nil
}
