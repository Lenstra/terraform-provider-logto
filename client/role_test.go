package client

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestRole(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()
	config.Logger = zerolog.New(os.Stdout)
	client, err := NewClient(config)
	require.NoError(t, err)

	role, err := client.RoleGet(ctx, "not-found")
	require.NoError(t, err)
	require.Nil(t, role)

	expected := &RoleModel{
		Name:        "test_role",
		Description: "A role to test the Terraform provider.",
	}
	role, err = client.RoleCreate(
		ctx,
		&RoleModel{
			Name:        "test_role",
			Description: "A role to test the Terraform provider.",
		},
	)
	require.NoError(t, err)
	require.NotEmpty(t, role.ID)
	require.NotEmpty(t, role.TenantId)
	require.NotEmpty(t, role.Name)
	require.NotEmpty(t, role.Description)
	require.NotEmpty(t, role.Type)
	require.NotNil(t, role.IsDefault)

	roleId := role.ID
	role.ID = ""
	require.Equal(t, expected.Name, role.Name)
	require.Equal(t, expected.Description, role.Description)

	role, err = client.RoleGet(ctx, roleId)
	require.NoError(t, err)
	require.NotNil(t, role)
	require.NotEmpty(t, role)
	require.Equal(t, expected.Name, role.Name)
	require.Equal(t, expected.Description, role.Description)

	role.Name = "test_role_update"
	role, err = client.RoleUpdate(ctx, role)
	require.NoError(t, err)
	require.NotNil(t, role)
	require.NotEmpty(t, role.ID)
	require.NotEmpty(t, role.TenantId)
	require.NotEmpty(t, role.Name)
	require.NotEmpty(t, role.Description)
	require.NotEmpty(t, role.Type)
	require.NotNil(t, role.IsDefault)
	require.Equal(t, "test_role_update", role.Name)

	err = client.RoleDelete(ctx, role.ID)
	require.NoError(t, err)
}
