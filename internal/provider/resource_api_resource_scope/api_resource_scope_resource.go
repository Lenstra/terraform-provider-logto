package resource_api_resource_scope

import (
	"context"
	"math/big"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func (r *apiResourceScopeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApiResourceScopeModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResourceScope := decodePlan(plan)

	apiResourceScope, err := r.client.ApiResourceScopeCreate(ctx, apiResourceScope.ResourceId, apiResourceScope)
	if err != nil {
		resp.Diagnostics.AddError("Error creating api_resource_scope", err.Error())
		return
	}

	convertToTerraformModel(apiResourceScope, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *apiResourceScopeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ApiResourceScopeModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Id.IsNull() || state.Id.IsUnknown() {
		resp.Diagnostics.AddError("Invalid Id", "Resource Id is null or unknown.")
		return
	}

	queryParams := map[string]string{
		"includeScopes": "yes",
		"page":          "1",  // Keep defaults
		"page_size":     "20", // Keep defaults
	}

	// Get scope
	apiResources, err := r.client.ApiResourceGetAll(ctx, queryParams)
	if err != nil {
		resp.Diagnostics.AddError("Error reading api_resource_scope", err.Error())
		return
	}

	// Get all ApiResource, then filter them because the Logto API does not provide a way
	// to directly GET the API_Resource with its Scopes. Generally, Scopes and API_Resources are few,
	// so just get all and filter to find the one we need.
	var resourceScope *client.ScopeModel
	for _, resource := range *apiResources {
		for _, scope := range *resource.Scopes {
			if scope.ID == state.Id.ValueString() {
				resourceScope = &scope
			}
		}
	}

	if resourceScope == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	convertToTerraformModel(resourceScope, &state)

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

	apiResourceScope := decodePlan(plan)

	apiResourceScope, err := r.client.ApiResourceScopeUpdate(ctx, apiResourceScope)
	if err != nil {
		resp.Diagnostics.AddError("Error updating api_resource_scope", err.Error())
		return
	}

	convertToTerraformModel(apiResourceScope, &state)

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

func decodePlan(plan ApiResourceScopeModel) *client.ScopeModel {
	model := &client.ScopeModel{
		ID:          plan.Id.ValueString(),
		TenantId:    plan.TenantId.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		ResourceId:  plan.ResourceId.ValueString(),
	}

	if !plan.CreatedAt.IsNull() && !plan.CreatedAt.IsUnknown() {
		createdAtBigFloat := plan.CreatedAt.ValueBigFloat()
		f, _ := createdAtBigFloat.Float64()
		model.CreatedAt = &f
	}

	return model
}

func convertToTerraformModel(apiResourceScope *client.ScopeModel, model *ApiResourceScopeModel) {
	*model = ApiResourceScopeModel{
		Id:          types.StringValue(apiResourceScope.ID),
		TenantId:    types.StringValue(apiResourceScope.TenantId),
		Name:        types.StringValue(apiResourceScope.Name),
		Description: types.StringValue(apiResourceScope.Description),
		ResourceId:  types.StringValue(apiResourceScope.ResourceId),
	}

	if apiResourceScope.CreatedAt != nil {
		createdAtBigFloat := new(big.Float).SetFloat64(*apiResourceScope.CreatedAt)
		model.CreatedAt = basetypes.NewNumberValue(createdAtBigFloat)
	}
}
