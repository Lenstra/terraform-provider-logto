// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider_logto

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
							resource "logto_application" "test" {
								name 											= "test"
								description 							= "test app"
								type 											= "SPA"
								redirect_uris 						= ["http://test.test.fr", "http://test.test.com"]
								post_logout_redirect_uris = ["http://redirect.test.fr", "http://redirect.test.com"]
							}

							data "logto_application" "test" {
								id = logto_application.test.id
							}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.logto_application.test", "name", "test"),
					resource.TestCheckResourceAttr("data.logto_application.test", "description", "test app"),
					resource.TestCheckResourceAttr("data.logto_application.test", "type", "SPA"),
					resource.TestCheckResourceAttr("data.logto_application.test", "redirect_uris.#", "2"),
					resource.TestCheckResourceAttr("data.logto_application.test", "redirect_uris.0", "http://test.test.fr"),
					resource.TestCheckResourceAttr("data.logto_application.test", "redirect_uris.1", "http://test.test.com"),
					resource.TestCheckResourceAttr("data.logto_application.test", "post_logout_redirect_uris.#", "2"),
					resource.TestCheckResourceAttr("data.logto_application.test", "post_logout_redirect_uris.0", "http://redirect.test.fr"),
					resource.TestCheckResourceAttr("data.logto_application.test", "post_logout_redirect_uris.1", "http://redirect.test.com"),
				),
			},
		},
	})
}
