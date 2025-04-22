package resource_role

import (
	"context"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	var plan RoleModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	role := &client.RoleModel{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		role.Type = plan.Type.ValueString()
	}

	if !plan.IsDefault.IsNull() && !plan.IsDefault.IsUnknown() {
		role.IsDefault = plan.IsDefault.ValueBool()
	}

	if !plan.ScopeIds.IsNull() && !plan.ScopeIds.IsUnknown() {
		var scopeIDs []string
		diags := plan.ScopeIds.ElementsAs(ctx, &scopeIDs, false)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			role.ScopeIds = scopeIDs
		}
	}

	role, err := r.client.RoleCreate(ctx, role)
	if err != nil {
		resp.Diagnostics.AddError("Error creating role", err.Error())
		return
	}

	r.updateRoleState(role, &plan)

	diags = resp.State.Set(ctx, plan)
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

	roleScopes, err := r.client.RoleScopesGet(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading role scopes", err.Error())
		return
	}

	if len(roleScopes) == 0 {
		state.ScopeIds = types.ListNull(types.StringType)
	} else {
		scopeIds := make([]string, 0, len(roleScopes))
		for _, roleScope := range roleScopes {
			scopeIds = append(scopeIds, roleScope.ID)
		}
		state.ScopeIds = stringSliceToList(scopeIds)
	}

	r.updateRoleState(role, &state)

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

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	role := &client.RoleModel{
		ID:          state.Id.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		role.Type = plan.Type.ValueString()
	}

	if !plan.IsDefault.IsNull() && !plan.IsDefault.IsUnknown() {
		role.IsDefault = plan.IsDefault.ValueBool()
	}

	role, err := r.client.RoleUpdate(ctx, role)
	if err != nil {
		resp.Diagnostics.AddError("Error updating role", err.Error())
		return
	}

	r.updateRoleState(role, &plan)

	diags = resp.State.Set(ctx, plan)
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

func stringSliceToList(slice []string) types.List {
	values := make([]attr.Value, len(slice))
	for i, s := range slice {
		values[i] = types.StringValue(s)
	}
	return types.ListValueMust(types.StringType, values)
}

func (r *roleResource) updateRoleState(role *client.RoleModel, model *RoleModel) {
	model.Id = types.StringValue(role.ID)
	model.Name = types.StringValue(role.Name)
	model.Description = types.StringValue(role.Description)

	if role.Type != "" {
		model.Type = types.StringValue(role.Type)
	} else {
		model.Type = types.StringNull()
	}

	if role.IsDefault {
		model.IsDefault = types.BoolValue(role.IsDefault)
	} else {
		model.IsDefault = types.BoolNull()
	}
}
