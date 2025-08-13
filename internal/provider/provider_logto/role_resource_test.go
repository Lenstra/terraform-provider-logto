package provider_logto

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleRessource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
							resource "logto_role" "test_role" {
									name 				= "tf_test_role"
									description = "tf_test_role_description"
									type				= "User"
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_role.test_role", "name", "tf_test_role"),
					resource.TestCheckResourceAttr("logto_role.test_role", "description", "tf_test_role_description"),
					resource.TestCheckResourceAttr("logto_role.test_role", "type", "User"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_role.test_role", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "logto_role.test_role",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: ProviderConfig + `
							resource "logto_role" "test_role" {
									name        = "tf_test_role_modified"
									description = "tf_test_role_description_modified"
									type				= "User"
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_role.test_role", "name", "tf_test_role_modified"),
					resource.TestCheckResourceAttr("logto_role.test_role", "description", "tf_test_role_description_modified"),
					resource.TestCheckResourceAttr("logto_role.test_role", "type", "User"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_role.test_role", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRoleRessourceRecreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
							resource "logto_api_resource" "api_resource" {
									name 				     = "tf_api_resource"
									indicator        = "https://api-resource.test"
									access_token_ttl = 3600
							}

							resource "logto_api_resource_scope" "api_resource_scope" {
									name 				= "tf_scope"
									resource_id = logto_api_resource.api_resource.id
									description = "test_scope_description"

									depends_on  = [logto_api_resource.api_resource]
							}

							resource "logto_role" "test_role" {
									name 				= "tf_test_role"
									description = "tf_test_role_description"
									type				= "User"
									scope_ids   = [
										logto_api_resource_scope.api_resource_scope.id
									]

									depends_on  = [
										logto_api_resource_scope.api_resource_scope
									]
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_role.test_role", "name", "tf_test_role"),
					resource.TestCheckResourceAttr("logto_role.test_role", "description", "tf_test_role_description"),
					resource.TestCheckResourceAttr("logto_role.test_role", "type", "User"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_role.test_role", "id"),
					resource.TestCheckResourceAttrSet("logto_role.test_role", "scope_ids.#"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "logto_role.test_role",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: ProviderConfig + `
							resource "logto_api_resource" "api_resource" {
									name 				     = "tf_api_resource"
									indicator        = "https://api-resource.test"
									access_token_ttl = 3600
							}

							resource "logto_api_resource_scope" "api_resource_scope" {
									name 				= "tf_scope"
									resource_id = logto_api_resource.api_resource.id
									description = "test_scope_description"

									depends_on  = [logto_api_resource.api_resource]
							}

							resource "logto_api_resource_scope" "api_resource_scope_two" {
									name 				= "tf_scope_two"
									resource_id = logto_api_resource.api_resource.id
									description = "test_scope_description_two"

									depends_on  = [logto_api_resource.api_resource]
							}

							resource "logto_role" "test_role" {
									name        = "tf_test_role_modified"
									description = "tf_test_role_description_modified"
									type				= "User"
									scope_ids   = [
										logto_api_resource_scope.api_resource_scope_two.id
									]

									depends_on  = [
										logto_api_resource_scope.api_resource_scope_two
									]
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_role.test_role", "name", "tf_test_role_modified"),
					resource.TestCheckResourceAttr("logto_role.test_role", "description", "tf_test_role_description_modified"),
					resource.TestCheckResourceAttr("logto_role.test_role", "type", "User"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_role.test_role", "id"),
					resource.TestCheckResourceAttrSet("logto_role.test_role", "scope_ids.#"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRoleRessourceAddScopeIdsAfterCreation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
							resource "logto_role" "test_role" {
									name 				= "tf_test_role"
									description = "tf_test_role_description"
									type				= "User"
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_role.test_role", "name", "tf_test_role"),
					resource.TestCheckResourceAttr("logto_role.test_role", "description", "tf_test_role_description"),
					resource.TestCheckResourceAttr("logto_role.test_role", "type", "User"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_role.test_role", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "logto_role.test_role",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: ProviderConfig + `
							resource "logto_api_resource" "api_resource" {
									name 				     = "tf_api_resource"
									indicator        = "https://api-resource.test"
									access_token_ttl = 3600
							}

							resource "logto_api_resource_scope" "api_resource_scope" {
									name 				= "tf_scope"
									resource_id = logto_api_resource.api_resource.id
									description = "test_scope_description"

									depends_on  = [logto_api_resource.api_resource]
							}

							resource "logto_role" "test_role" {
									name        = "tf_test_role_modified"
									description = "tf_test_role_description_modified"
									type				= "User"
									scope_ids   = [
										logto_api_resource_scope.api_resource_scope.id
									]

									depends_on  = [
										logto_api_resource_scope.api_resource_scope
									]
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_role.test_role", "name", "tf_test_role_modified"),
					resource.TestCheckResourceAttr("logto_role.test_role", "description", "tf_test_role_description_modified"),
					resource.TestCheckResourceAttr("logto_role.test_role", "type", "User"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_role.test_role", "id"),
					resource.TestCheckResourceAttrSet("logto_role.test_role", "scope_ids.#"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
