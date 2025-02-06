package provider_logto

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Logto client is properly configured.
	// It is also possible to use the LOGTO_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	ProviderConfig = `
provider "logto" {
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"logto": providerserver.NewProtocol6WithError(New("test")()),
	}
)
