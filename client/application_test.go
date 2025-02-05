package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplication(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()
	client := NewClient(config)

	app, err := client.ApplicationGet(ctx, "not-found")
	require.NoError(t, err)
	require.Nil(t, app)

	app = &ApplicationModel{
		Name: "test",
		Type: "MachineToMachine",
	}
	app, err = client.ApplicationCreate(ctx, app)
	require.NoError(t, err)
	require.NotNil(t, app)
	require.NotEmpty(t, app.ID)
	require.Equal(t, "test", app.Name)
	require.Equal(t, "MachineToMachine", app.Type)
	require.Equal(t, "", app.Description)

	app, err = client.ApplicationGet(ctx, app.ID)
	require.NoError(t, err)
	require.NotNil(t, app)
	require.NotEmpty(t, app.ID)
	require.Equal(t, "test", app.Name)
	require.Equal(t, "MachineToMachine", app.Type)

	secrets, err := client.GetApplicationSecrets(ctx, app.ID)
	require.NoError(t, err)
	require.Len(t, secrets, 1)
	require.Equal(t, "Default secret", secrets[0].Name)

	app.Description = "test update"
	app, err = client.ApplicationUpdate(ctx, app)
	require.NoError(t, err)
	require.NotNil(t, app)
	require.NotEmpty(t, app.ID)
	require.Equal(t, "test", app.Name)
	require.Equal(t, "MachineToMachine", app.Type)
	require.Equal(t, "test update", app.Description)

	err = client.ApplicationDelete(ctx, app.ID)
	require.NoError(t, err)
}
