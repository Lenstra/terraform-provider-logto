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
									name 				= "test"
									description = "test app"
									type 				= "MachineToMachine"
							}

							output "test_app_secrets" {
								value = logto_application.test_app.secrets
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description", "test app"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type", "MachineToMachine"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttrSet("logto_application.test_app", "secrets.%"),
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
									name 				= "test modified"
									description = "test app modified"
									type 				= "MachineToMachine"
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description", "test app modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type", "MachineToMachine"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttrSet("logto_application.test_app", "secrets.%"),
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
									name 				= "test"
									description = "test app"
									type				= "MachineToMachine"
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description", "test app"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type", "MachineToMachine"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttrSet("logto_application.test_app", "secrets.%"),
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
									name 				= "test modified"
									description = "test app modified"
									type 				= "Native"
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description", "test app modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type", "Native"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttrSet("logto_application.test_app", "secrets.%"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccApplicationResourceWithRedirectUrisAndNotPostLogoutRedirectUris(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
							resource "logto_application" "test_app" {
									name 				  = "test"
									description 	= "test app"
									type				  = "SPA"
									redirect_uris = ["http://test.test.fr", "http://test.test.com"]
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description", "test app"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type", "SPA"),

					resource.TestCheckResourceAttr("logto_application.test_app", "redirect_uris.#", "2"),
					resource.TestCheckResourceAttr("logto_application.test_app", "redirect_uris.0", "http://test.test.fr"),
					resource.TestCheckResourceAttr("logto_application.test_app", "redirect_uris.1", "http://test.test.com"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttrSet("logto_application.test_app", "redirect_uris.#"),
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
									name 						= "test modified"
									description 		= "test app modified"
									type 						= "SPA"
									redirect_uris 	= ["http://test_modified.test.fr"]
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description", "test app modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type", "SPA"),

					resource.TestCheckResourceAttr("logto_application.test_app", "redirect_uris.#", "1"),
					resource.TestCheckResourceAttr("logto_application.test_app", "redirect_uris.0", "http://test_modified.test.fr"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttrSet("logto_application.test_app", "redirect_uris.#"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccApplicationResourceWithRedirectUrisAndPostLogoutRedirectUris(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
							resource "logto_application" "test_app" {
									name 				  						= "test"
									description 							= "test app"
									type				  						= "SPA"
									redirect_uris 						= ["http://test.test.fr", "http://test.test.com"]
									post_logout_redirect_uris = ["http://redirect.test.fr", "http://redirect.test.com"]
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description", "test app"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type", "SPA"),

					resource.TestCheckResourceAttr("logto_application.test_app", "redirect_uris.#", "2"),
					resource.TestCheckResourceAttr("logto_application.test_app", "redirect_uris.0", "http://test.test.fr"),
					resource.TestCheckResourceAttr("logto_application.test_app", "redirect_uris.1", "http://test.test.com"),

					resource.TestCheckResourceAttr("logto_application.test_app", "post_logout_redirect_uris.#", "2"),
					resource.TestCheckResourceAttr("logto_application.test_app", "post_logout_redirect_uris.0", "http://redirect.test.fr"),
					resource.TestCheckResourceAttr("logto_application.test_app", "post_logout_redirect_uris.1", "http://redirect.test.com"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttrSet("logto_application.test_app", "redirect_uris.#"),
					resource.TestCheckResourceAttrSet("logto_application.test_app", "post_logout_redirect_uris.#"),
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
									name 						= "test modified"
									description 		= "test app modified"
									type 						= "SPA"
									redirect_uris 						= ["http://test_modified.test.fr", "http://test_modified.test.com"]
									post_logout_redirect_uris = ["http://redirect_modified.test.fr", "http://redirect_modified.test.com"]
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description", "test app modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type", "SPA"),

					resource.TestCheckResourceAttr("logto_application.test_app", "redirect_uris.#", "2"),
					resource.TestCheckResourceAttr("logto_application.test_app", "redirect_uris.0", "http://test_modified.test.fr"),
					resource.TestCheckResourceAttr("logto_application.test_app", "redirect_uris.1", "http://test_modified.test.com"),

					resource.TestCheckResourceAttr("logto_application.test_app", "post_logout_redirect_uris.#", "2"),
					resource.TestCheckResourceAttr("logto_application.test_app", "post_logout_redirect_uris.0", "http://redirect_modified.test.fr"),
					resource.TestCheckResourceAttr("logto_application.test_app", "post_logout_redirect_uris.1", "http://redirect_modified.test.com"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttrSet("logto_application.test_app", "redirect_uris.#"),
					resource.TestCheckResourceAttrSet("logto_application.test_app", "post_logout_redirect_uris.#"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccApplicationResourceWithNotRedirectUrisAndPostLogoutRedirectUris(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
							resource "logto_application" "test_app" {
									name 				  = "test"
									description 	= "test app"
									type				  = "SPA"
									post_logout_redirect_uris = ["http://redirect_modified.test.fr", "http://redirect_modified.test.com"]
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description", "test app"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type", "SPA"),

					resource.TestCheckResourceAttr("logto_application.test_app", "post_logout_redirect_uris.#", "2"),
					resource.TestCheckResourceAttr("logto_application.test_app", "post_logout_redirect_uris.0", "http://redirect_modified.test.fr"),
					resource.TestCheckResourceAttr("logto_application.test_app", "post_logout_redirect_uris.1", "http://redirect_modified.test.com"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttrSet("logto_application.test_app", "post_logout_redirect_uris.#"),
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
									name 						= "test modified"
									description 		= "test app modified"
									type 						= "SPA"
									redirect_uris 						= ["http://test_modified.test.fr", "http://test_modified.test.com"]
									post_logout_redirect_uris = ["http://redirect_modified.test.fr", "http://redirect_modified.test.com"]
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description", "test app modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type", "SPA"),

					resource.TestCheckResourceAttr("logto_application.test_app", "redirect_uris.#", "2"),
					resource.TestCheckResourceAttr("logto_application.test_app", "redirect_uris.0", "http://test_modified.test.fr"),
					resource.TestCheckResourceAttr("logto_application.test_app", "redirect_uris.1", "http://test_modified.test.com"),

					resource.TestCheckResourceAttr("logto_application.test_app", "post_logout_redirect_uris.#", "2"),
					resource.TestCheckResourceAttr("logto_application.test_app", "post_logout_redirect_uris.0", "http://redirect_modified.test.fr"),
					resource.TestCheckResourceAttr("logto_application.test_app", "post_logout_redirect_uris.1", "http://redirect_modified.test.com"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttrSet("logto_application.test_app", "redirect_uris.#"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccApplicationResourceWithCorsAllowedOrigins(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ProviderConfig + `
							resource "logto_application" "test_app" {
								name                      = "test"
								description               = "test app"
								type                      = "SPA"
								cors_allowed_origins      = ["http://cors_allowed_origin_test.fr", "http://cors_allowed_origin_test.com"]
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description", "test app"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type", "SPA"),

					resource.TestCheckResourceAttr("logto_application.test_app", "cors_allowed_origins.#", "2"),
					resource.TestCheckResourceAttr("logto_application.test_app", "cors_allowed_origins.0", "http://cors_allowed_origin_test.fr"),
					resource.TestCheckResourceAttr("logto_application.test_app", "cors_allowed_origins.1", "http://cors_allowed_origin_test.com"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttrSet("logto_application.test_app", "cors_allowed_origins.#"),
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
								name                      = "test modified"
								description               = "test app modified"
								type                      = "SPA"
								cors_allowed_origins      = ["http://cors_allowed_origin_test.com"]
							}
							`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("logto_application.test_app", "name", "test modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "description", "test app modified"),
					resource.TestCheckResourceAttr("logto_application.test_app", "type", "SPA"),

					resource.TestCheckResourceAttr("logto_application.test_app", "cors_allowed_origins.#", "1"),
					resource.TestCheckResourceAttr("logto_application.test_app", "cors_allowed_origins.0", "http://cors_allowed_origin_test.com"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("logto_application.test_app", "id"),
					resource.TestCheckResourceAttrSet("logto_application.test_app", "cors_allowed_origins.#"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
