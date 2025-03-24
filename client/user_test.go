package client

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()
	config.Logger = zerolog.New(os.Stdout)
	client, err := NewClient(config)
	require.NoError(t, err)

	user, err := client.UserGet(ctx, "not-found")
	require.NoError(t, err)
	require.Nil(t, user)

	expected := &UserModel{
		Username: "test",
		Name:     "test",
		Profile: &Profile{
			FamilyName: "test",
			GivenName:  "test",
			MiddleName: "test",
			Nickname:   "test",
		},
	}
	user, err = client.UserCreate(
		ctx,
		&UserModel{
			Username: "test",
			Name:     "test",
			Profile: &Profile{
				FamilyName: "test",
				GivenName:  "test",
				MiddleName: "test",
				Nickname:   "test",
			},
		},
	)
	require.NoError(t, err)
	require.NotEmpty(t, user.ID)

	userId := user.ID
	user.ID = ""
	require.Equal(t, expected, user)

	user, err = client.UserGet(ctx, userId)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.NotEmpty(t, userId)
	require.Equal(t, "test", user.Username)
	require.Equal(t, "test", user.Name)

	user.Name = "test update"
	user, err = client.UserUpdate(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.NotEmpty(t, user.ID)
	require.Equal(t, "test", user.Username)
	require.Equal(t, "test update", user.Name)

	err = client.UserDelete(ctx, user.ID)
	require.NoError(t, err)
}
