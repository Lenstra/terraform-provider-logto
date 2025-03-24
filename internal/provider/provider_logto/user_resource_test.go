package provider_logto

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
					resource "logto_user" "test_user" {
						name     = "test_user"
						username = "test_username"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_user.test_user", "name", "test_user"),
					resource.TestCheckResourceAttr("logto_user.test_user", "username", "test_username"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "logto_user.test_user",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: ProviderConfig + `
					resource "logto_user" "test_user" {
							name        = "test_user_modified"

							profile = {
								family_name = "test_family_name_modified"
								given_name = "test_given_name_modified"
								middle_name = "test_middle_name_modified"
								nickname = "test_nickname_modified"
							}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_user.test_user", "name", "test_user_modified"),

					resource.TestCheckResourceAttr("logto_user.test_user", "profile.%", "4"),
					resource.TestCheckResourceAttr("logto_user.test_user", "profile.family_name", "test_family_name_modified"),
					resource.TestCheckResourceAttr("logto_user.test_user", "profile.given_name", "test_given_name_modified"),
					resource.TestCheckResourceAttr("logto_user.test_user", "profile.middle_name", "test_middle_name_modified"),
					resource.TestCheckResourceAttr("logto_user.test_user", "profile.nickname", "test_nickname_modified"),
				),
			},
			{
				Config: ProviderConfig + `
					resource "logto_user" "test_user" {
							name        = "test_user_modified"
							username    = "test_username_modified"

							profile = {
								family_name = "test_family_name_modified"
							}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_user.test_user", "name", "test_user_modified"),
					resource.TestCheckResourceAttr("logto_user.test_user", "username", "test_username_modified"),

					resource.TestCheckResourceAttr("logto_user.test_user", "profile.%", "4"),
					resource.TestCheckResourceAttr("logto_user.test_user", "profile.family_name", "test_family_name_modified"),
					resource.TestCheckResourceAttr("logto_user.test_user", "profile.given_name", ""),
					resource.TestCheckResourceAttr("logto_user.test_user", "profile.middle_name", ""),
					resource.TestCheckResourceAttr("logto_user.test_user", "profile.nickname", ""),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
