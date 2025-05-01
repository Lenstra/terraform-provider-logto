package client

import (
	"context"
	"net/http"
	"path"
)

func (c *Client) UserGet(ctx context.Context, id string) (*UserModel, error) {
	if id == "" {
		return nil, errEmptyID
	}

	req := &request{
		method: http.MethodGet,
		path:   path.Join("api/users/", id),
	}
	res, err := expect(200, 404)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, nil
	}

	var user UserModel
	if err := decode(res.Body, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *Client) UserCreate(ctx context.Context, user *UserModel) (*UserModel, error) {
	req := &request{
		method: http.MethodPost,
		path:   "api/users",
		body:   user,
	}

	res, err := expect(200)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var User UserModel
	if err := decode(res.Body, &User); err != nil {
		return nil, err
	}
	return &User, nil
}

func (c *Client) UserDelete(ctx context.Context, id string) error {
	if id == "" {
		return errEmptyID
	}

	req := &request{
		method: http.MethodDelete,
		path:   path.Join("api/users/", id),
	}
	_, err := expect(204)(c.do(ctx, req))
	return err
}

func (c *Client) UserUpdate(ctx context.Context, user *UserModel) (*UserModel, error) {
	req := &request{
		method: http.MethodPatch,
		path:   path.Join("api/users/", user.ID),
		body:   user,
	}

	res, err := expect(200)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var returnUser UserModel
	if err := decode(res.Body, &returnUser); err != nil {
		return nil, err
	}
	return &returnUser, nil
}

func (c *Client) UserAssignRole(ctx context.Context, userId string, roleIds *RoleIds) error {
	req := &request{
		method: http.MethodPost,
		path:   path.Join("/api/users", userId, "/roles"),
		body:   roleIds,
	}

	_, err := expect(201)(c.do(ctx, req))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UserGetRoles(ctx context.Context, userId string) (*[]RoleModel, error) {
	req := &request{
		method: http.MethodPost,
		path:   path.Join("/api/users", userId, "/roles"),
	}

	res, err := expect(200)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var roles []RoleModel
	if err := decode(res.Body, &roles); err != nil {
		return nil, err
	}
	return &roles, nil
}
