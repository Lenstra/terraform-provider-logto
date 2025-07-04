package client

import (
	"context"
	"net/http"
)

func (c *Client) SignInExperienceGet(ctx context.Context) (*SignInExperienceModel, error) {
	req := &request{
		method: http.MethodGet,
		path:   "api/sign-in-exp",
	}
	res, err := expect(200, 404)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, nil
	}

	var signInExperience SignInExperienceModel
	if err := decode(res.Body, &signInExperience); err != nil {
		return nil, err
	}
	return &signInExperience, nil
}

func (c *Client) SignInExperienceUpdate(ctx context.Context, signInExperience *SignInExperienceModel) (*SignInExperienceModel, error) {
	req := &request{
		method: http.MethodPatch,
		path:   "api/sign-in-exp",
		body:   signInExperience,
	}

	res, err := expect(200)(c.do(ctx, req))
	if err != nil {
		return nil, err
	}

	var returnSignInExp SignInExperienceModel
	if err := decode(res.Body, &returnSignInExp); err != nil {
		return nil, err
	}

	return &returnSignInExp, nil
}
