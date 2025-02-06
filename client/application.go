package client

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) ApplicationGet(ctx context.Context, id string) (*ApplicationModel, error) {
	req := &request{
		method: http.MethodGet,
		path:   "api/applications/" + id,
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
	req := &request{
		method: http.MethodDelete,
		path:   "api/applications/" + id,
	}
	_, err := expect(204)(c.do(ctx, req))
	return err
}

func (c *Client) ApplicationUpdate(ctx context.Context, app *ApplicationModel) (*ApplicationModel, error) {
	req := &request{
		method: http.MethodPatch,
		path:   "api/applications/" + app.ID,
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

func (c *Client) GetApplicationSecrets(ctx context.Context, applicationId string) ([]Secret, error) {
	req := &request{
		method: http.MethodGet,
		path:   fmt.Sprintf("api/applications/%s/secrets", applicationId),
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
