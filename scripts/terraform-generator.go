package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"os/exec"

	"github.com/hashicorp/terraform-plugin-codegen-spec/code"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/hashicorp/terraform-plugin-codegen-spec/spec"
)

func Main() error {
	resp, err := http.Get("https://openapi.logto.io/source.yaml")
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("expected 200 status code, got %s", resp.Status)
	}

	f, err := os.CreateTemp("", "openapi")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	f.Close()

	cmd := exec.Command(
		"tfplugingen-openapi",
		"generate",
		"--config=./config/generator_config.yml",
		"--output=./provider_code_spec.json",
		f.Name(),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	spec, err := readSpec("provider_code_spec.json")
	if err != nil {
		return err
	}
	extra, err := readSpec("config/provider_code_extra.json")
	if err != nil {
		return err
	}

	// use the provider given in our extra conf
	spec.Provider = extra.Provider

	resources := map[string]resource.Resource{}
	for _, r := range extra.Resources {
		resources[r.Name] = r
	}
	datasources := map[string]datasource.DataSource{}
	for _, r := range extra.DataSources {
		datasources[r.Name] = r
	}

	for _, resource := range spec.Resources {
		for _, a := range resource.Schema.Attributes {
			if a.String != nil && a.Name == "id" {
				a.String.PlanModifiers = append(a.String.PlanModifiers, schema.StringPlanModifier{
					Custom: &schema.CustomPlanModifier{
						Imports: []code.Import{
							{Path: "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"},
						},
						SchemaDefinition: "stringplanmodifier.UseStateForUnknown()",
					},
				})
			}
			if a.List != nil && a.List.ComputedOptionalRequired == schema.Optional {
				a.List.ComputedOptionalRequired = schema.ComputedOptional
				a.List.PlanModifiers = append(a.List.PlanModifiers, schema.ListPlanModifier{
					Custom: &schema.CustomPlanModifier{
						Imports: []code.Import{
							{Path: "github.com/Lenstra/terraform-provider-logto/internal/provider/planmodifiers/listplanmodifier"},
						},
						SchemaDefinition: "listplanmodifier.NullIsEmpty()",
					},
				})
			}
		}
		extraResource, found := resources[resource.Name]
		if found {
			resource.Schema.Attributes = append(
				resource.Schema.Attributes,
				extraResource.Schema.Attributes...,
			)
		}
	}
	for _, datasource := range spec.DataSources {
		extraDatasource, found := datasources[datasource.Name]
		if found {
			datasource.Schema.Attributes = append(
				datasource.Schema.Attributes,
				extraDatasource.Schema.Attributes...,
			)
		}
	}

	content, err := json.MarshalIndent(spec, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile("provider_code_spec.json", content, 0o644)
	if err != nil {
		return err
	}

	cmd = exec.Command(
		"tfplugingen-framework",
		"generate",
		"all",
		"--input=./provider_code_spec.json",
		"--output=internal/provider",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func readSpec(path string) (*spec.Specification, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var extra spec.Specification
	if err := json.Unmarshal(content, &extra); err != nil {
		return nil, err
	}
	return &extra, nil
}

func main() {
	err := Main()
	if err != nil {
		log.Fatalf("err: %v", err)
		os.Exit(1)
	}
}
