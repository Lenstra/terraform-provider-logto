package provider_logto

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApiResourceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
					resource "logto_api_resource" "test_api_resource" {
						name             = "tf_test_api_resource"
						indicator        = "https://test-api-resource.test"
						access_token_ttl = 3600
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_api_resource.test_api_resource", "name", "tf_test_api_resource"),
					resource.TestCheckResourceAttr("logto_api_resource.test_api_resource", "indicator", "https://test-api-resource.test"),
					resource.TestCheckResourceAttr("logto_api_resource.test_api_resource", "access_token_ttl", "3600"),
					resource.TestCheckResourceAttrSet("logto_api_resource.test_api_resource", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "logto_api_resource.test_api_resource",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Idempotence check (reapply same config, no changes expected)
			{
				Config: ProviderConfig + `
					resource "logto_api_resource" "test_api_resource" {
						name             = "tf_test_api_resource"
						indicator        = "https://test-api-resource.test"
						access_token_ttl = 3600
					}
				`,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Update and Read testing
			{
				Config: ProviderConfig + `
					resource "logto_api_resource" "test_api_resource" {
						name             = "tf_test_api_resource modified"
						indicator        = "https://test-api-resource-modified.test"
						access_token_ttl = 4000
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_api_resource.test_api_resource", "name", "tf_test_api_resource modified"),
					resource.TestCheckResourceAttr("logto_api_resource.test_api_resource", "indicator", "https://test-api-resource-modified.test"),
					resource.TestCheckResourceAttr("logto_api_resource.test_api_resource", "access_token_ttl", "4000"),
					resource.TestCheckResourceAttrSet("logto_api_resource.test_api_resource", "id"),
				),
			},
			// Create with optional fields omitted (test defaults)
			{
				Config: ProviderConfig + `
					resource "logto_api_resource" "test_api_resource_defaults" {
						name      = "tf_test_api_resource_defaults"
						indicator = "https://default-api-resource.test"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("logto_api_resource.test_api_resource_defaults", "name", "tf_test_api_resource_defaults"),
					resource.TestCheckResourceAttr("logto_api_resource.test_api_resource_defaults", "indicator", "https://default-api-resource.test"),
					resource.TestCheckResourceAttr("logto_api_resource.test_api_resource_defaults", "access_token_ttl", "3600"),
					resource.TestCheckResourceAttr("logto_api_resource.test_api_resource_defaults", "scopes.#", "0"),

					resource.TestCheckResourceAttrSet("logto_api_resource.test_api_resource_defaults", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
