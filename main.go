// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

	"github.com/Lenstra/terraform-provider-logto/internal/provider/provider_logto"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Update the Provider schema.
//go:generate go run ./scripts/terraform-generator.go

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can also pass the specific commit if you want
	// commit  string = "" .
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "github.com/Lenstra/terraform-provider-logto",
		Debug:   debugMode,
	}

	err := providerserver.Serve(context.Background(), provider_logto.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
