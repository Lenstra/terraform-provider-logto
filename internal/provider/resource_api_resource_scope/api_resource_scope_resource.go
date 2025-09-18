package resource_api_resource_scope

import (
	"context"
	"math/big"
	"strings"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.ResourceWithImportState = &apiResourceScopeResource{}
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
	if state.ResourceId.IsNull() || state.ResourceId.IsUnknown() {
		resp.Diagnostics.AddError("Invalid ResourceId", "ResourceId is null or unknown.")
		return
	}

	queryParams := map[string]string{
		"page":      "1",
		"page_size": "20",
	}

	resourceScopes, err := r.client.ApiResourceScopesList(ctx, state.ResourceId.ValueString(), queryParams)
	if err != nil {
		resp.Diagnostics.AddError("Error reading api_resource_scopes", err.Error())
		return
	}

	var foundScope *client.ScopeModel
	for _, scope := range resourceScopes {
		if scope.ID == state.Id.ValueString() {
			foundScope = &scope
			break
		}
	}

	if foundScope == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	convertToTerraformModel(foundScope, &state)

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
		Name:       plan.Name.ValueString(),
		ResourceId: plan.ResourceId.ValueString(),
		ID:         plan.Id.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		model.Description = plan.Description.ValueString()
	}

	return model
}

func convertToTerraformModel(apiResourceScope *client.ScopeModel, model *ApiResourceScopeModel) {
	*model = ApiResourceScopeModel{
		Id:         types.StringValue(apiResourceScope.ID),
		Name:       types.StringValue(apiResourceScope.Name),
		ResourceId: types.StringValue(apiResourceScope.ResourceId),
		TenantId:   types.StringValue(apiResourceScope.TenantId),
	}

	if apiResourceScope.Description != "" {
		model.Description = types.StringValue(apiResourceScope.Description)
	} else {
		model.Description = types.StringNull()
	}

	if apiResourceScope.CreatedAt != nil {
		createdAtBigFloat := new(big.Float).SetFloat64(*apiResourceScope.CreatedAt)
		model.CreatedAt = basetypes.NewNumberValue(createdAtBigFloat)
	} else {
		model.CreatedAt = basetypes.NewNumberNull()
	}
}

func (r *apiResourceScopeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected import identifier",
			"Expected format: <resource_id>/<scope_id>",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("resource_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
