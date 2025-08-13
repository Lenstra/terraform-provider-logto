package resource_api_resource_scope

import (
	"context"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func (r *apiResourceScopeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApiResourceScopeModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResourceScope := &client.ScopeModel{
		ResourceId: plan.ResourceId.ValueString(),
		Name:       plan.Name.ValueString(),
	}

	if !plan.Description.IsUnknown() && !plan.Description.IsNull() {
		apiResourceScope.Description = plan.Description.ValueString()
	}

	apiResourceScope, err := r.client.ApiResourceScopeCreate(ctx, apiResourceScope.ResourceId, apiResourceScope)
	if err != nil {
		resp.Diagnostics.AddError("Error creating api_resource_scope", err.Error())
		return
	}

	r.updateApiResourceScopeState(apiResourceScope, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *apiResourceScopeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ApiResourceScopeModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queryParams := map[string]string{
		"page":      "1",  // Keep defaults
		"page_size": "20", // Keep defaults
	}

	if !state.Name.IsUnknown() && !state.Name.IsNull() {
		queryParams["search"] = state.Name.ValueString()
	}

	if state.ResourceId.IsNull() || state.ResourceId.IsUnknown() {
		resp.Diagnostics.AddError("Invalid Resource ID", "ResourceId is null or unknown.")
		return
	}

	apiResourceScopes, err := r.client.ApiResourceScopesGetWithParams(ctx, state.ResourceId.ValueString(), queryParams)
	if err != nil {
		resp.Diagnostics.AddError("Error reading api_resource_scope", err.Error())
		return
	}

	if len(*apiResourceScopes) > 1 {
		resp.Diagnostics.AddError("Error reading api_resource_scope",
			"The API returned more than one result. Expected only one result for the query.")
		return
	}

	apiResourceScope := (*apiResourceScopes)[0]
	r.updateApiResourceScopeState(&apiResourceScope, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *apiResourceScopeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApiResourceScopeModel
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

	apiResourceScope := &client.ScopeModel{
		ResourceId: state.ResourceId.ValueString(),
		ID:         state.Id.ValueString(),
		Name:       plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		apiResourceScope.Description = plan.Description.ValueString()
	}

	scope, err := r.client.ApiResourceScopeUpdate(ctx, apiResourceScope)
	if err != nil {
		resp.Diagnostics.AddError("Error updating api_resource_scope", err.Error())
		return
	}

	r.updateApiResourceScopeState(scope, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *apiResourceScopeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ApiResourceScopeModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.ApiResourceScopeDelete(ctx, state.ResourceId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting api_resource_scope", err.Error())
	}
}

func (r *apiResourceScopeResource) updateApiResourceScopeState(scope *client.ScopeModel, model *ApiResourceScopeModel) {

	model.TenantId = types.StringValue(scope.TenantId)
	model.Id = types.StringValue(scope.ID)
	model.ResourceId = types.StringValue(scope.ResourceId)
	model.Name = types.StringValue(scope.Name)
	model.Description = types.StringValue(scope.Description)
	model.CreatedAt = basetypes.NumberValue(types.Float64Value(*scope.CreatedAt))

}
