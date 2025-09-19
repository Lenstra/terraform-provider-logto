package resource_assign_roles_to_user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func AssignRolesToUserResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"role_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Required:            true,
				Description:         "An array of API resource role IDs to assign.",
				MarkdownDescription: "An array of API resource role IDs to assign.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.StringAttribute{
				Required:            true,
				Description:         "The unique identifier of the user.",
				MarkdownDescription: "The unique identifier of the user.",
			},
		},
	}
}

type AssignRolesToUserModel struct {
	Id      types.String `tfsdk:"id"`
	RoleIds types.Set    `tfsdk:"role_ids"`
	UserId  types.String `tfsdk:"user_id"`
}
