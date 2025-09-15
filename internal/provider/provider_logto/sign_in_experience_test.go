package provider_logto

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSignInExperienceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
					resource "logto_sign_in_experience" "test_experience" {
						color = {
							primary_color 		 	 = "#f54242"
							dark_primary_color 	 = "#5442f5"
							is_dark_mode_enabled = false
						}

						branding = {
							logo_url = "https://logo_url.test"
							dark_logo_url = "https://dark_logo_url.test"
							favicon = "https://favicon.test"
							dark_favicon = "https://dark_favicon.test"
						}

						language_info = {
							auto_detect = true
							fallback_language = "en"
						}

						terms_of_use_url = "https://example.com/terms"
						privacy_policy_url = "https://example.com/privacy"
						agree_to_terms_policy = "Automatic"

						sign_in = {
              methods = [
              {
                identifier = "email"
                password = true
                verification_code = false
                is_password_primary = true
              },
              {
                identifier = "phone"
                password = true
                verification_code = false
                is_password_primary = true
              },
              ]
            }

						social_sign_in = {
							automatic_account_linking = true
						}

						sign_in_mode = "SignIn"
						custom_css = ".logto_signature{visibility:hidden;}"

						password_policy = {
							length = {
								min = 8
								max = 128
							}

							character_types = {
								min = 2
							}

							rejects = {
								pwned                     = true
								repetition_and_sequence   = true
								user_info                 = true
								words                     = ["password", "123456", "admin"]
							}
						}

						mfa = {
							factors = ["Totp"]
							policy  = "UserControlled"
						}

						single_sign_on_enabled = true

						support_email = "support@email.foo"
						support_website_url = "https://support_website_url.foo"
						unknown_session_redirect_url = "https://unknown_session_redirect_url.foo"

						captcha_policy = {
							enabled = true
						}

						sentinel_policy = {
							max_attempts = 10
							lockout_duration = 3600
						}

						email_blocklist_policy = {
							block_disposable_addresses = true
    					block_subaddressing = true
						}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify basic attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "terms_of_use_url", "https://example.com/terms"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "privacy_policy_url", "https://example.com/privacy"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sign_in_mode", "SignIn"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "support_email", "support@email.foo"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "support_website_url", "https://support_website_url.foo"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "unknown_session_redirect_url", "https://unknown_session_redirect_url.foo"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "social_sign_in_connector_targets.#", "0"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "single_sign_on_enabled", "true"),

					// Verify color attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "color.primary_color", "#f54242"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "color.dark_primary_color", "#5442f5"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "color.is_dark_mode_enabled", "false"),

					// Verify branding attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "branding.logo_url", "https://logo_url.test"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "branding.dark_logo_url", "https://dark_logo_url.test"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "branding.favicon", "https://favicon.test"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "branding.dark_favicon", "https://dark_favicon.test"),

					// Verify language_info attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "language_info.auto_detect", "true"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "language_info.fallback_language", "en"),

					//Verify sign_in attributes sign_in.methods attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sign_in.methods.#", "2"),

					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sign_in.methods.0.identifier", "email"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sign_in.methods.0.password", "true"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sign_in.methods.0.verification_code", "false"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sign_in.methods.0.is_password_primary", "true"),

					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sign_in.methods.1.identifier", "phone"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sign_in.methods.1.password", "true"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sign_in.methods.1.verification_code", "false"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sign_in.methods.1.is_password_primary", "true"),

					// Verify sign_up attributes
					// NEED TO ENABLE CONNECTORS TO TEST THIS RESOURCE

					//Verify social_sign_in attribute
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "social_sign_in.automatic_account_linking", "true"),

					//Verify custom_css attribute
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "custom_css", ".logto_signature{visibility:hidden;}"),

					// Verify password_policy attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "password_policy.length.min", "8"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "password_policy.length.max", "128"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "password_policy.character_types.min", "2"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "password_policy.rejects.pwned", "true"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "password_policy.rejects.repetition_and_sequence", "true"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "password_policy.rejects.user_info", "true"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "password_policy.rejects.words.#", "3"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "password_policy.rejects.words.0", "password"),

					// Verify MFA attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "mfa.factors.#", "1"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "mfa.factors.0", "Totp"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "mfa.policy", "UserControlled"),

					// Verify captcha Policy
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "captcha_policy.enabled", "true"),

					// Verify sentinel_policy attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sentinel_policy.max_attempts", "10"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sentinel_policy.lockout_duration", "3600"),

					// Verify email_blocklist_policy attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "email_blocklist_policy.block_disposable_addresses", "true"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "email_blocklist_policy.block_subaddressing", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "logto_sign_in_experience.test_experience",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: ProviderConfig + `
					resource "logto_sign_in_experience" "test_experience" {
						color = {
							primary_color 		 	 = "#4287f6"
							dark_primary_color 	 = "#5442f6"
							is_dark_mode_enabled = true
						}

						branding = {
							logo_url = "https://logo_url_update.test"
							dark_logo_url = "https://dark_logo_url_update.test"
							favicon = "https://favicon_update.test"
							dark_favicon = "https://dark_favicon_update.test"
						}

						language_info = {
							auto_detect = false
							fallback_language = "fr"
						}

						terms_of_use_url = "https://example.com/terms_update"
						privacy_policy_url = "https://example.com/privacy_update"
						agree_to_terms_policy = "Automatic"

						sign_in = {
              methods = [
              {
                identifier = "email"
                password = true
                verification_code = false
                is_password_primary = false
              },
              ]
            }

						// NEED TO ENABLE CONNECTORS TO USE SignUp
						// SignUp: &SignUp{
						// 	Identifiers: []string{"email"},
						// 	Password:    false,
						// 	Verify:      false,
						// },

						social_sign_in = {
							automatic_account_linking = false
						}

						sign_in_mode = "SignInAndRegister"
						custom_css = ".logto_signature{visibility:visible;}"

						password_policy = {
							length = {
								min = 10
								max = 120
							}

							character_types = {
								min = 3
							}

							rejects = {
								pwned                     = false
								repetition_and_sequence   = false
								user_info                 = false
								words                     = ["password", "123456"]
							}
						}

						mfa = {
							factors = ["Totp", "WebAuthn"]
							policy  = "NoPrompt"
							organization_required_mfa_policy = "NoPrompt"
						}

						single_sign_on_enabled = true // false not work in the case of this test, the API always return 'true' value so make it 'false' return a terraform error

						support_email = "support_update@email.foo"
						support_website_url = "https://support_website_url_update.foo"
						unknown_session_redirect_url = "https://unknown_session_redirect_url_update.foo"

						captcha_policy = {
							enabled = false
						}

						sentinel_policy = {
							max_attempts = 12
							lockout_duration = 4600
						}

						email_blocklist_policy = {
							block_disposable_addresses = false
    					block_subaddressing = false
						}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify basic attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "terms_of_use_url", "https://example.com/terms_update"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "privacy_policy_url", "https://example.com/privacy_update"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sign_in_mode", "SignInAndRegister"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "support_email", "support_update@email.foo"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "support_website_url", "https://support_website_url_update.foo"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "unknown_session_redirect_url", "https://unknown_session_redirect_url_update.foo"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "single_sign_on_enabled", "true"),

					// Verify color attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "color.primary_color", "#4287f6"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "color.dark_primary_color", "#5442f6"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "color.is_dark_mode_enabled", "true"),

					// Verify branding attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "branding.logo_url", "https://logo_url_update.test"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "branding.dark_logo_url", "https://dark_logo_url_update.test"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "branding.favicon", "https://favicon_update.test"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "branding.dark_favicon", "https://dark_favicon_update.test"),

					// Verify language_info attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "language_info.auto_detect", "false"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "language_info.fallback_language", "fr"),

					//Verify sign_in attributes sign_in.methods attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sign_in.methods.#", "1"),

					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sign_in.methods.0.identifier", "email"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sign_in.methods.0.password", "true"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sign_in.methods.0.verification_code", "false"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sign_in.methods.0.is_password_primary", "false"),

					// Verify sign_up attributes
					// NEED TO ENABLE CONNECTORS TO TEST THIS RESOURCE

					//Verify social_sign_in attribute
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "social_sign_in.automatic_account_linking", "false"),

					//Verify custom_css attribute
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "custom_css", ".logto_signature{visibility:visible;}"),

					// Verify password_policy attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "password_policy.length.min", "10"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "password_policy.length.max", "120"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "password_policy.character_types.min", "3"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "password_policy.rejects.pwned", "false"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "password_policy.rejects.repetition_and_sequence", "false"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "password_policy.rejects.user_info", "false"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "password_policy.rejects.words.#", "2"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "password_policy.rejects.words.0", "password"),

					// Verify MFA attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "mfa.factors.#", "2"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "mfa.factors.0", "Totp"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "mfa.factors.1", "WebAuthn"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "mfa.policy", "NoPrompt"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "mfa.organization_required_mfa_policy", "NoPrompt"),

					// Verify captcha Policy
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "captcha_policy.enabled", "false"),

					// Verify sentinel_policy attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sentinel_policy.max_attempts", "12"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "sentinel_policy.lockout_duration", "4600"),

					// Verify email_blocklist_policy attributes
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "email_blocklist_policy.block_disposable_addresses", "false"),
					resource.TestCheckResourceAttr("logto_sign_in_experience.test_experience", "email_blocklist_policy.block_subaddressing", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
