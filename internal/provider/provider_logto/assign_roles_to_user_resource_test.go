package provider_logto

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAssignRolesToUserResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
					resource "logto_user" "test_user" {
							name              = "tf_user_link_roles"
							primary_email     = "tf_user_link_roles@test.fr"
					}

					resource "logto_role" "test_role" {
							name 				= "tf_test_link_roles"
							description = "tf_test_link_roles_description"
							type				= "User"
					}

					resource "logto_assign_roles_to_user" "link_roles" {
						role_ids = [
							logto_role.test_role.id
						]
						user_id = logto_user.test_user.id
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_assign_roles_to_user.link_roles", "id"),
					resource.TestCheckResourceAttrSet("logto_assign_roles_to_user.link_roles", "user_id"),

					// Check length
					resource.TestCheckResourceAttr("logto_assign_roles_to_user.link_roles", "role_ids.#", "1"),
				),
			},
			// Cannot test the import of the resource dynamically.
			// Update and Read testing
			{
				Config: ProviderConfig + `
					resource "logto_user" "test_user" {
							name              = "tf_user_link_roles"
					}

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

					resource "logto_assign_roles_to_user" "link_roles" {
						role_ids = [
							logto_role.test_role.id,
							logto_role.test_role2.id
						]
						user_id = logto_user.test_user.id
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_assign_roles_to_user.link_roles", "id"),
					resource.TestCheckResourceAttrSet("logto_assign_roles_to_user.link_roles", "user_id"),

					// Check length
					resource.TestCheckResourceAttr("logto_assign_roles_to_user.link_roles", "role_ids.#", "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccImportOfAssignRolesToUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// ImportState testing
			{
				// Cannot test the import of the resource dynamically so use predefined resources.
				Config: ProviderConfig + `
				resource "logto_api_resource_scope" "test_assign_roles_to_user_for_import" {
					roles_ids        = [
						"4wytu0yi21r4a6dwf9dtk",
						"wz548z4mwyzanyotp01fo"
					]
					user_id = mkgcd5t7j06q
				}
    `,
				ResourceName:      "logto_api_resource_scope.test_api_resource_scope_for_import",
				ImportState:       true,
				ImportStateVerify: true,                                                       // Verification disabled because API returns computed fields (TenantId, CreatedAt) not present in HCL
				ImportStateId:     "mkgcd5t7j06q/4wytu0yi21r4a6dwf9dtk-wz548z4mwyzanyotp01fo", // user_id/role_id_1-role_id_2
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
