package client

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestApplication(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()
	config.Logger = zerolog.New(os.Stdout)
	client, err := NewClient(config)
	require.NoError(t, err)

	app, err := client.ApplicationGet(ctx, "not-found")
	require.NoError(t, err)
	require.Nil(t, app)

	expected := &ApplicationModel{
		Name:                 "test",
		Description:          "test",
		Type:                 "Traditional",
		OidcClientMetadata:   &OidcClientMetadata{RedirectUris: []string{}, PostLogoutRedirectUris: []string{}},
		CustomClientMetadata: &CustomClientMetadata{},
		CustomData:           map[string]interface{}{},
		ProtectedAppMetadata: nil,
		IsAdmin:              false,
		IsThirdParty:         true,
	}
	app, err = client.ApplicationCreate(
		ctx,
		&ApplicationModel{
			Name:         "test",
			Type:         "Traditional",
			Description:  "test",
			IsThirdParty: true,
		},
	)
	require.NoError(t, err)
	require.NotEmpty(t, app.ID)
	require.NotEmpty(t, app.TenantId)

	appId := app.ID
	app.ID = ""
	app.TenantId = ""
	require.Equal(t, expected, app)

	app, err = client.ApplicationGet(ctx, appId)
	require.NoError(t, err)
	require.NotNil(t, app)
	require.NotEmpty(t, appId)
	require.Equal(t, "test", app.Name)
	require.Equal(t, "Traditional", app.Type)

	secrets, err := client.GetApplicationSecrets(ctx, appId)
	require.NoError(t, err)
	require.Len(t, secrets, 1)
	require.Equal(t, "Default secret", secrets[0].Name)

	app.Description = "test update"
	app, err = client.ApplicationUpdate(ctx, app)
	require.NoError(t, err)
	require.NotNil(t, app)
	require.NotEmpty(t, app.ID)
	require.Equal(t, "test", app.Name)
	require.Equal(t, "Traditional", app.Type)
	require.Equal(t, "test update", app.Description)

	err = client.ApplicationDelete(ctx, app.ID)
	require.NoError(t, err)
}
