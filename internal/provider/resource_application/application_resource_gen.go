// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package resource_application

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ApplicationResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The unique identifier of the application.",
				MarkdownDescription: "The unique identifier of the application.",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
				},
			},
			"tenant_id": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Native",
						"SPA",
						"Traditional",
						"MachineToMachine",
						"Protected",
						"SAML",
					),
				},
			},
		},
	}
}

type ApplicationModel struct {
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	TenantId    types.String `tfsdk:"tenant_id"`
	Type        types.String `tfsdk:"type"`
}
