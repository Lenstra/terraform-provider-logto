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
			PrimaryColor:      "#000000",
			IsDarkModeEnabled: false,
			DarkPrimaryColor:  "#686868",
		},
		Branding: &Branding{
			LogoUrl:     "https://tom-siouan.com/static/img/lenstra-logo.svg",
			DarkLogoUrl: "https://tom-siouan.com/static/img/lenstra-logo.svg",
		},
		LanguageInfo: &LanguageInfo{
			AutoDetect:       true,
			FallbackLanguage: "en",
		},
		SignIn: &SignIn{
			Methods: []Methods{{
				Identifier:        "email",
				Password:          true,
				VerificationCode:  false,
				IsPasswordPrimary: false,
			}, {
				Identifier:        "phone",
				Password:          true,
				VerificationCode:  false,
				IsPasswordPrimary: true,
			}},
		},

		// NEED TO ENABLE CONNECTORS TO USE SignUp
		// SignUp: &SignUp{
		// 	Identifiers: []string{"email"},
		// 	Password:    false,
		// 	Verify:      false,
		// },

		Mfa: &Mfa{
			Factors:                       []string{"Totp", "WebAuthn"},
			Policy:                        "NoPrompt",
			OrganizationRequiredMfaPolicy: "NoPrompt",
		},

		SingleSignOnEnabled: true,
	}

	signInExp, err := client.SignInExperienceUpdate(ctx, expected)
	require.NoError(t, err)
	require.NotNil(t, signInExp)

	getSignInExp, err := client.SignInExperienceGet(ctx)
	require.NoError(t, err)
	require.NotNil(t, getSignInExp)

	require.Equal(t, expected.Color, getSignInExp.Color)
	require.Equal(t, expected.Branding, getSignInExp.Branding)
	require.Equal(t, expected.LanguageInfo, getSignInExp.LanguageInfo)
	require.Equal(t, expected.SignIn, getSignInExp.SignIn)
	require.Equal(t, expected.SingleSignOnEnabled, getSignInExp.SingleSignOnEnabled)
}
