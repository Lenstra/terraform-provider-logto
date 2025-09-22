package client

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestAssignRoleToUser(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()
	config.Logger = zerolog.New(os.Stdout)
	client, err := NewClient(config)
	require.NoError(t, err)

	// Initialize objects for testing
	user, err := client.UserCreate(ctx, &UserModel{
		PrimaryEmail: "clientTestAssignRole@test.fr",
		Username:     "clientTestAssignRole",
		Name:         "clientTestAssignRole",
	})
	require.NoError(t, err)
	require.NotEmpty(t, user.ID)

	role, err := client.RoleCreate(ctx, &RoleModel{
		Name:        "clientTestAssignRole",
		Description: "A role to test the assignation with a user",
	})
	require.NoError(t, err)
	require.NotEmpty(t, role.ID)

	role1, err := client.RoleCreate(ctx, &RoleModel{
		Name:        "clientTestAssignRole1",
		Description: "A role to test the assignation with a user",
	})
	require.NoError(t, err)
	require.NotEmpty(t, role1.ID)

	roleId := role.ID
	roleId1 := role1.ID

	// Check that association works
	err = client.AssignRolesForUser(ctx, &RoleIdsModel{
		RoleIds: []string{
			role.ID,
		},
	}, user.ID)
	require.NoError(t, err)

	roles, err := client.GetRolesForUser(ctx, user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, roles)
	require.Equal(t, roleId, roles[0].ID)

	// Check that update works
	err = client.UpdateRolesForUser(ctx, &RoleIdsModel{
		RoleIds: []string{
			role1.ID,
		},
	}, user.ID)
	require.NoError(t, err)

	roles, err = client.GetRolesForUser(ctx, user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, roles)

	require.Equal(t, roleId1, roles[0].ID)

	// Check that deletion works
	err = client.DeleteRolesForUser(ctx, roleId1, user.ID)
	require.NoError(t, err)

	roles, err = client.GetRolesForUser(ctx, user.ID)
	require.NoError(t, err)
	require.Empty(t, roles)

	// Remove roles and user
	err = client.RoleDelete(ctx, roleId)
	require.NoError(t, err)

	err = client.RoleDelete(ctx, roleId1)
	require.NoError(t, err)

	err = client.UserDelete(ctx, user.ID)
	require.NoError(t, err)
}
