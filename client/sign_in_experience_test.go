package client

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestSignInExperience(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()
	config.Logger = zerolog.New(os.Stdout)
	client, err := NewClient(config)
	require.NoError(t, err)

	expected := &SignInExperienceModel{
		Color: &Color{
			PrimaryColor: "#ff4017",
			IsDarkModeEnabled: false,
			DarkPrimaryColor: "#686868",
		},
		Branding: &Branding{
			LogoUrl: "https://tom-siouan.com/static/img/lenstra-logo.svg",
			DarkLogoUrl: "https://tom-siouan.com/static/img/lenstra-logo.svg",
			Favicon: "https://tom-siouan.com/static/img/lenstra-logo.svg",
			DarkFavicon: "https://tom-siouan.com/static/img/lenstra-logo.svg",
		},
		LanguageInfo: &LanguageInfo{

		}
	}

	app, err := client.SignInExperienceUpdate(ctx)
	require.NoError(t, err)
	require.NotNil(t, app)

}
