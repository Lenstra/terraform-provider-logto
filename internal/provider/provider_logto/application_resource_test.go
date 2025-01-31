package provider_logto

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationResourceWithoutTypeUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
							resource "logto_application" "test_app" {
									name = "test"
									description = "test app"
									type = "MachineToMachine"
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description", "test app"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type", "MachineToMachine"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "logto_application.test_app",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: ProviderConfig + `
							resource "logto_application" "test_app" {
									name = "test modified"
									description = "test app modified"
									type = "MachineToMachine"
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description", "test app modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type", "MachineToMachine"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccApplicationResourceWithTypeUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
							resource "logto_application" "test_app" {
									name = "test"
									description = "test app"
									type = "MachineToMachine"
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description", "test app"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type", "MachineToMachine"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "logto_application.test_app",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: ProviderConfig + `
							resource "logto_application" "test_app" {
									name = "test modified"
									description = "test app modified"
									type = "Native"
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description", "test app modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type", "Native"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
