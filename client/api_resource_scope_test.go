package client

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestApiResourceScope(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()
	config.Logger = zerolog.New(os.Stdout)
	client, err := NewClient(config)
	require.NoError(t, err)

	apiResource, err := client.ApiResourceCreate(
		ctx,
		&ApiResourceModel{
			Name:      "test_api_resource",
			Indicator: "https://indicator.test",
		},
	)
	require.NoError(t, err)

	expected := &ScopeModel{
		Name:        "test_scope",
		Description: "test_scope_description",
		ResourceId:  apiResource.ID,
	}

	scope, err := client.ApiResourceScopeCreate(
		ctx,
		apiResource.ID,
		&ScopeModel{
			Name:        "test_scope",
			Description: "test_scope_description",
		},
	)
	require.NoError(t, err)
	require.NotEmpty(t, scope.ID)
	require.NotEmpty(t, scope.TenantId)

	scope.ID = ""
	require.Equal(t, expected.Name, scope.Name)
	require.Equal(t, expected.Description, scope.Description)

	queryParams := map[string]string{
		"page":      "1",
		"page_size": "20",
		"search":    "test_scope",
	}

	scopes, err := client.ApiResourceScopesGet(ctx, expected.ResourceId)
	require.NoError(t, err)
	require.NotNil(t, scopes)
	require.NotEmpty(t, scopes)
	require.NotNil(t, (*scopes)[0].Name)
	require.NotNil(t, (*scopes)[0].Description)
	require.NotEmpty(t, (*scopes)[0].Name)
	require.NotEmpty(t, (*scopes)[0].Description)
	require.Equal(t, "test_scope", (*scopes)[0].Name)
	require.Equal(t, "test_scope_description", (*scopes)[0].Description)

	scopes, err = client.ApiResourceScopesGetWithParams(ctx, expected.ResourceId, queryParams)
	require.NoError(t, err)
	require.NotNil(t, scopes)
	require.NotEmpty(t, scopes)
	require.NotNil(t, (*scopes)[0].Name)
	require.NotNil(t, (*scopes)[0].Description)
	require.NotEmpty(t, (*scopes)[0].Description)
	require.Equal(t, "test_scope", (*scopes)[0].Name)
	require.Equal(t, "test_scope_description", (*scopes)[0].Description)

	scope.Name = "test_scope_update"
	scope.Description = "test_scope_description_update"
	scope.ID = (*scopes)[0].ID
	scope, err = client.ApiResourceScopeUpdate(ctx, scope)
	require.NoError(t, err)
	require.NotNil(t, scope)
	require.NotEmpty(t, scope.ID)
	require.Equal(t, "test_scope_update", scope.Name)
	require.Equal(t, "test_scope_description_update", scope.Description)

	err = client.ApiResourceDelete(ctx, apiResource.ID)
	require.NoError(t, err)
}
