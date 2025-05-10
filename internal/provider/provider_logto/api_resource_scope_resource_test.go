package provider_logto

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccApiResourceScopeResource(t *testing.T) {
	var importID string

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
			{
				// Just to extract state values
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
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["logto_api_resource_scope.test_api_resource_scope"]
					if !ok {
						return fmt.Errorf("resource not found in state")
					}
					resourceID := rs.Primary.Attributes["resource_id"]
					name := rs.Primary.Attributes["name"]
					if resourceID == "" || name == "" {
						return fmt.Errorf("resource_id or name is empty in state")
					}
					importID = fmt.Sprintf("%s:%s", resourceID, name)
					return nil
				},
			},
			// ImportState testing
			{
				ResourceName:      "logto_api_resource_scope.test_api_resource_scope",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(_ *terraform.State) (string, error) {
					return importID, nil
				},
			},
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
