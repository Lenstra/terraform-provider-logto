package provider_logto

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApiResourceScopeResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
							resource "logto_api_resource" "test_api_resource" {
									name 				     = "tf_test_api_resource_scope"
									indicator        = "https://test-api-resource.test"
									access_token_ttl = 3600
							}

							resource "logto_api_resource_scope" "test_api_resource_scope" {
									name 				= "tf_test_scope"
									resource_id = logto_api_resource.test_api_resource.id
									description = "test_scope_description"

									depends_on  = [logto_api_resource.test_api_resource]
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_api_resource_scope.test_api_resource_scope", "name", "tf_test_scope"),
					resource.TestCheckResourceAttr("logto_api_resource_scope.test_api_resource_scope", "description", "test_scope_description"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_api_resource_scope.test_api_resource_scope", "id"),
					resource.TestCheckResourceAttrSet("logto_api_resource_scope.test_api_resource_scope", "resource_id"),
					resource.TestCheckResourceAttrSet("logto_api_resource_scope.test_api_resource_scope", "tenant_id"),
					resource.TestCheckResourceAttrSet("logto_api_resource_scope.test_api_resource_scope", "created_at"),
				),
			},
			// Cannot test the import of the resource dynamically.
			// Update and Read testing
			{
				Config: ProviderConfig + `
					resource "logto_api_resource" "test_api_resource" {
							name 				     = "tf_test_api_resource_scope"
							indicator        = "https://test-api-resource.test"
							access_token_ttl = 3600
					}

					resource "logto_api_resource_scope" "test_api_resource_scope" {
							name 				= "tf_test_scope_modified"
							resource_id = logto_api_resource.test_api_resource.id
							description = "test_scope_description_modified"

							depends_on  = [logto_api_resource.test_api_resource]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_api_resource_scope.test_api_resource_scope", "name", "tf_test_scope_modified"),
					resource.TestCheckResourceAttr("logto_api_resource_scope.test_api_resource_scope", "description", "test_scope_description_modified"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_api_resource_scope.test_api_resource_scope", "id"),
					resource.TestCheckResourceAttrSet("logto_api_resource_scope.test_api_resource_scope", "resource_id"),
					resource.TestCheckResourceAttrSet("logto_api_resource_scope.test_api_resource_scope", "tenant_id"),
					resource.TestCheckResourceAttrSet("logto_api_resource_scope.test_api_resource_scope", "created_at"),
				),
			},
			// Modifying resource_id
			{
				Config: ProviderConfig + `
					resource "logto_api_resource" "test_api_resource_two" {
							name 				     = "tf_test_api_resource_scope_two"
							indicator        = "https://test-api-resource.testtwo"
							access_token_ttl = 3600
					}

					resource "logto_api_resource_scope" "test_api_resource_scope" {
							name 				= "tf_test_scope_modified"
							resource_id = logto_api_resource.test_api_resource_two.id
							description = "test_scope_description_modified"

							depends_on  = [logto_api_resource.test_api_resource_two]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_api_resource_scope.test_api_resource_scope", "name", "tf_test_scope_modified"),
					resource.TestCheckResourceAttr("logto_api_resource_scope.test_api_resource_scope", "description", "test_scope_description_modified"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_api_resource_scope.test_api_resource_scope", "id"),
					resource.TestCheckResourceAttrSet("logto_api_resource_scope.test_api_resource_scope", "resource_id"),
					resource.TestCheckResourceAttrSet("logto_api_resource_scope.test_api_resource_scope", "tenant_id"),
					resource.TestCheckResourceAttrSet("logto_api_resource_scope.test_api_resource_scope", "created_at"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
func TestAccImportOfApiResourceScopeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// ImportState testing
			{
				// Cannot test the import of the resource dynamically so use predefined resources.
				Config: ProviderConfig + `
				resource "logto_api_resource_scope" "test_api_resource_scope_for_import" {
					name        = "test_api_resource_scope_for_import"
					resource_id = logto_api_resource.test_api_resource_for_import.id
				}
    `,
				ResourceName:      "logto_api_resource_scope.test_api_resource_scope_for_import",
				ImportState:       true,
				ImportStateVerify: false, // Verification disabled because API returns computed fields (TenantId, CreatedAt) not present in HCL
				ImportStateId:     "4vt7qldcnbutcmqfbcvuv/vyvo1yxu1u3yutqmzz3rb",
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
