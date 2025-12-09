package tfgen

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/code"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/hashicorp/terraform-plugin-codegen-spec/spec"
)

func Update(spec, extra *spec.Specification) {
	// use the provider given in our extra conf
	spec.Provider = extra.Provider

	resources := map[string]resource.Resource{}
	for _, r := range extra.Resources {
		resources[r.Name] = r
	}
	datasources := map[string]datasource.DataSource{}
	for _, d := range extra.DataSources {
		datasources[d.Name] = d
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

}
