package resource_application_test

import (
	"testing"

	"github.com/Lenstra/terraform-provider-logto/internal/provider/provider_logto"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider_logto.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: provider_logto.ProviderConfig + `
				resource "logto_application" "test_app" {
					name 			  = "test"
					description = "test app"
					type 				= "MachineToMachine"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check the script job attributes
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description ", "test app"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type ", "MachineToMachine"),
				),
			},
			// // ImportState testing
			{
				ResourceName:      "logto_application.test_app",
				ImportState:       true,
				ImportStateVerify: true,
				// The filesha and  scriptfile attributes does not exist in SimpleMDM
				// API, therefore there is no value for it during import.
			},
			// Update and Read testing
			{
				Config: provider_logto.ProviderConfig + `
				resource "logto_application" "test_app" {
					name 			  = "test modified"
					description = "test app modified"
					type 				= "MachineToMachine"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check the updated script job attributes
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description ", "test app modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type ", "MachineToMachine"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
