package client

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestApiResource(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()
	config.Logger = zerolog.New(os.Stdout)
	client, err := NewClient(config)
	require.NoError(t, err)

	apiResource, err := client.ApiResourceGet(ctx, "not-found")
	require.NoError(t, err)
	require.Nil(t, apiResource)

	expected := &ApiResourceModel{
		Name:      "test_api_resource",
		Indicator: "https://indicator.test",
	}

	apiResource, err = client.ApiResourceCreate(
		ctx,
		expected,
	)
	require.NoError(t, err)
	require.NotEmpty(t, apiResource.ID)
	require.NotEmpty(t, apiResource.TenantId)

	apiResourceId := apiResource.ID
	apiResource.ID = ""
	require.Equal(t, expected.Name, apiResource.Name)
	require.Equal(t, expected.Indicator, apiResource.Indicator)

	queryParams := map[string]string{
		"page":      "1",
		"page_size": "20",
	}
	apiResources, err := client.ApiResourceGetAll(ctx, queryParams)
	require.NoError(t, err)
	require.NotNil(t, apiResources)
	require.NotEmpty(t, apiResources)

	found := false
	for _, res := range *apiResources {
		if res.ID == apiResourceId {
			found = true
			require.Equal(t, "test_api_resource", res.Name)
			break
		}
	}
	require.True(t, found, "The created resource was not found in the list")

	apiResource, err = client.ApiResourceGet(ctx, apiResourceId)
	require.NoError(t, err)
	require.NotNil(t, apiResource)
	require.NotEmpty(t, apiResource)
	require.Equal(t, "test_api_resource", apiResource.Name)
	require.Equal(t, "https://indicator.test", apiResource.Indicator)

	apiResource.Name = "test_api_resource_update"
	apiResource, err = client.ApiResourceUpdate(ctx, apiResource)
	require.NoError(t, err)
	require.NotNil(t, apiResource)
	require.NotEmpty(t, apiResource.ID)
	require.Equal(t, "test_api_resource_update", apiResource.Name)

	err = client.ApiResourceDelete(ctx, apiResource.ID)
	require.NoError(t, err)
}
