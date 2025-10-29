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
		expected.ResourceId,
		&ScopeModel{
			Name:        "test_scope",
			Description: "test_scope_description",
		},
	)
	require.NoError(t, err)
	require.NotEmpty(t, scope.ID)
	require.NotEmpty(t, scope.TenantId)

	require.Equal(t, expected.Name, scope.Name)
	require.Equal(t, expected.Description, scope.Description)

	queryParams := map[string]string{
		"page":      "1",  // Keep defaults
		"page_size": "20", // Keep defaults
	}

	apiResourceScopes, err := client.ApiResourceScopesList(ctx, apiResource.ID, queryParams)
	require.NoError(t, err)
	require.NotNil(t, apiResourceScopes)

	var foundScope *ScopeModel
	for _, s := range apiResourceScopes {
		if s.ID == scope.ID {
			foundScope = &s
			break
		}
	}

	require.NotNil(t, foundScope)
	require.Equal(t, "test_scope", foundScope.Name)
	require.Equal(t, "test_scope_description", foundScope.Description)

	foundScope.Name = "test_scope_update"
	foundScope.Description = "test_scope_description_update"
	updatedScope, err := client.ApiResourceScopeUpdate(ctx, foundScope)
	require.NoError(t, err)
	require.NotNil(t, updatedScope)
	require.NotEmpty(t, updatedScope.ID)
	require.Equal(t, "test_scope_update", updatedScope.Name)
	require.Equal(t, "test_scope_description_update", updatedScope.Description)

	err = client.ApiResourceDelete(ctx, apiResource.ID)
	require.NoError(t, err)
}
