package resource_role

import (
	"context"
	"fmt"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &roleResource{}
	_ resource.ResourceWithConfigure   = &roleResource{}
	_ resource.ResourceWithImportState = &roleResource{}
)

type roleResource struct {
	client *client.Client
}

func RoleResource() resource.Resource {
	return &roleResource{}
}

func (r *roleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (r *roleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = RoleResourceSchema(ctx)
}

func (r *roleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		return
	}
	r.client = client
}

func (r *roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state RoleModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	role := decodePlan(ctx, plan)

	fmt.Printf("role %v : ", role)

	role, err := r.client.RoleCreate(ctx, role)
	if err != nil {
		resp.Diagnostics.AddError("Error creating role", err.Error())
		return
	}

	diags = convertToTerraformModel(ctx, role, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RoleModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	role, err := r.client.RoleGet(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading role", err.Error())
		return
	}
	if role == nil {
		resp.State.RemoveResource(ctx)
	}

	roleScopes, err := r.client.RoleScopesGet(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading role scopes", err.Error())
		return
	}

	if len(roleScopes) == 0 {
		role.ScopeIds = nil
	} else {
		scopeIds := make([]string, 0, len(roleScopes))
		for _, roleScope := range roleScopes {
			scopeIds = append(scopeIds, roleScope.ID)
		}
		role.ScopeIds = scopeIds
	}

	diags = convertToTerraformModel(ctx, role, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state RoleModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	role := decodePlan(ctx, plan)

	role, err := r.client.RoleUpdate(ctx, role)
	if err != nil {
		resp.Diagnostics.AddError("Error updating role", err.Error())
		return
	}

	diags = convertToTerraformModel(ctx, role, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RoleModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.RoleDelete(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting role", err.Error())
	}
}

func convertToTerraformModel(ctx context.Context, role *client.RoleModel, model *RoleModel) (diags diag.Diagnostics) {
	*model = RoleModel{
		Id:          types.StringValue(role.ID),
		Name:        types.StringValue(role.Name),
		Description: types.StringValue(role.Description),
		Type:        types.StringValue(role.Type),
		IsDefault:   types.BoolValue(role.IsDefault),
	}

	if role.ScopeIds != nil {
		model.ScopeIds, diags = basetypes.NewListValueFrom(ctx, types.StringType, role.ScopeIds)
		if diags.HasError() {
			return diags
		}
	}
	return
}

func decodePlan(ctx context.Context, plan RoleModel) *client.RoleModel {
	model := &client.RoleModel{
		ID:          plan.Id.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		model.Type = plan.Type.ValueString()
	}

	if !plan.IsDefault.IsNull() && !plan.IsDefault.IsUnknown() {
		model.IsDefault = plan.IsDefault.ValueBool()
	}

	if !plan.ScopeIds.IsNull() && !plan.ScopeIds.IsUnknown() {
		var scopeIDs []string
		diags := plan.ScopeIds.ElementsAs(ctx, &scopeIDs, false)
		if !diags.HasError() {
			model.ScopeIds = scopeIDs
		}
	}

	return model
}
