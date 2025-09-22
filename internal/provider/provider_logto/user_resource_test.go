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
				resource "logto_role" "test_role" {
					name 				= "tf_test_link_roles"
					description = "tf_test_link_roles_description"
					type				= "User"
				}

				resource "logto_user" "test_user" {
					name              = "test_user_modified"
					primary_email     = "test_user@test.fr"

					profile = {
						family_name = "test_family_name_modified"
						given_name  = "test_given_name_modified"
						middle_name = "test_middle_name_modified"
						nickname    = "test_nickname_modified"
					}

					role_ids = [
						logto_role.test_role.id
					]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_user.test_user", "name", "test_user_modified"),
					resource.TestCheckResourceAttr("logto_user.test_user", "primary_email", "test_user@test.fr"),

					resource.TestCheckResourceAttr("logto_user.test_user", "profile.%", "4"),
					resource.TestCheckResourceAttr("logto_user.test_user", "role_ids.#", "1"),

					resource.TestCheckResourceAttr("logto_user.test_user", "profile.family_name", "test_family_name_modified"),
					resource.TestCheckResourceAttr("logto_user.test_user", "profile.given_name", "test_given_name_modified"),
					resource.TestCheckResourceAttr("logto_user.test_user", "profile.middle_name", "test_middle_name_modified"),
					resource.TestCheckResourceAttr("logto_user.test_user", "profile.nickname", "test_nickname_modified"),
				),
			},
			// Re-create the user with a profile
			{
				Taint: []string{"logto_user.test_user"},
				Config: ProviderConfig + `
				resource "logto_role" "test_role" {
					name 				= "tf_test_link_roles"
					description = "tf_test_link_roles_description"
					type				= "User"
				}

				resource "logto_role" "test_role2" {
					name 				= "tf_test_link_roles2"
					description = "tf_test_link_roles_description2"
					type				= "User"
				}

				resource "logto_user" "test_user" {
					name              = "test_user_modified"
					primary_email     = "test_user@test.fr"

					profile = {
						family_name = "test_family_name_modified"
						given_name  = "test_given_name_modified"
						middle_name = "test_middle_name_modified"
						nickname    = "test_nickname_modified"
					}

					role_ids = [
						logto_role.test_role.id,
						logto_role.test_role2.id
					]
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_user.test_user", "role_ids.#", "2"),
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
