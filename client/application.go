package client

import (
	"context"
	"net/http"
	"path"
)

func (c *Client) ApplicationGet(ctx context.Context, id string) (*ApplicationModel, error) {
	if id == "" {
		return nil, errEmptyID
	}

	req := &request{
		method: http.MethodGet,
		path:   path.Join("api/applications", id),
	}
	res, err := expect(200, 404)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, nil
	}

	var application ApplicationModel
	if err := decode(res.Body, &application); err != nil {
		return nil, err
	}
	return &application, nil
}

func (c *Client) ApplicationCreate(ctx context.Context, app *ApplicationModel) (*ApplicationModel, error) {
	req := &request{
		method: http.MethodPost,
		path:   "api/applications",
		body:   app,
	}

	res, err := expect(200)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var application ApplicationModel
	if err := decode(res.Body, &application); err != nil {
		return nil, err
	}
	return &application, nil
}

func (c *Client) ApplicationDelete(ctx context.Context, id string) error {
	if id == "" {
		return errEmptyID
	}

	req := &request{
		method: http.MethodDelete,
		path:   path.Join("api/applications", id),
	}
	_, err := expect(204)(c.do(ctx, req))
	return err
}

func (c *Client) ApplicationUpdate(ctx context.Context, app *ApplicationModel) (*ApplicationModel, error) {
	if app.ID == "" {
		return nil, errEmptyID
	}

	req := &request{
		method: http.MethodPatch,
		path:   path.Join("api/applications", app.ID),
		body:   app,
	}

	res, err := expect(200)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var application ApplicationModel
	if err := decode(res.Body, &application); err != nil {
		return nil, err
	}
	return &application, nil
}

func (c *Client) GetApplicationSecrets(ctx context.Context, id string) ([]Secret, error) {
	if id == "" {
		return nil, errEmptyID
	}

	req := &request{
		method: http.MethodGet,
		path:   path.Join("api/applications", id, "secrets"),
	}

	res, err := expect(200)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var secrets []Secret
	if err := decode(res.Body, &secrets); err != nil {
		return nil, err
	}
	return secrets, nil
}
