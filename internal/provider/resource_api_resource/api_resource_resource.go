package resource_api_resource

import (
	"context"
	"math/big"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func (r *apiResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApiResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResource, diags := decodePlan(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResource, err := r.client.ApiResourceCreate(ctx, apiResource)
	if err != nil {
		resp.Diagnostics.AddError("Error creating api_resource", err.Error())
		return
	}

	diags = convertToTerraformModel(ctx, apiResource, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *apiResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ApiResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResource, err := r.client.ApiResourceGet(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading api_resource", err.Error())
		return
	}

	if apiResource == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = convertToTerraformModel(ctx, apiResource, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *apiResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApiResourceModel
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

	apiResource, diags := decodePlan(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResource, err := r.client.ApiResourceUpdate(ctx, apiResource)
	if err != nil {
		resp.Diagnostics.AddError("Error updating api_resource", err.Error())
		return
	}

	diags = convertToTerraformModel(ctx, apiResource, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *apiResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ApiResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.ApiResourceDelete(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting api_resource", err.Error())
	}

	
}

func decodePlan(ctx context.Context, plan ApiResourceModel) (*client.ApiResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := &client.ApiResourceModel{
		ID:        plan.Id.ValueString(),
		Name:      plan.Name.ValueString(),
		Indicator: plan.Indicator.ValueString(),
	}

	isDefault := plan.IsDefault.ValueBoolPointer()
	if isDefault == nil {
		falseVal := false
		isDefault = &falseVal
	}
	model.IsDefault = isDefault

	accessTokenTtlValue, _ := plan.AccessTokenTtl.ValueBigFloat().Float64()
	model.AccessTokenTtl = &accessTokenTtlValue

	var scopesList = make([]client.ScopeModel, 0)
	if !plan.Scopes.IsNull() && !plan.Scopes.IsUnknown() {
		var scopes []ScopesValue
		diags.Append(plan.Scopes.ElementsAs(ctx, &scopes, false)...)
		if diags.HasError() {
			return nil, diags
		}

		for _, s := range scopes {
			scope := client.ScopeModel{
				TenantId:    s.TenantId.ValueString(),
				ID:          s.Id.ValueString(),
				ResourceId:  s.ResourceId.ValueString(),
				Name:        s.Name.ValueString(),
				Description: s.Description.ValueString(),
			}

			if !s.CreatedAt.IsNull() && !s.CreatedAt.IsUnknown() {
				val, _ := s.CreatedAt.ValueBigFloat().Float64()
				scope.CreatedAt = &val
			}

			scopesList = append(scopesList, scope)
		}
	}
	model.Scopes = &scopesList

	return model, diags
}

func convertToTerraformModel(ctx context.Context, apiResource *client.ApiResourceModel, model *ApiResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	*model = ApiResourceModel{
		Id:        types.StringValue(apiResource.ID),
		Name:      types.StringValue(apiResource.Name),
		Indicator: types.StringValue(apiResource.Indicator),
		IsDefault: types.BoolPointerValue(apiResource.IsDefault),
	}

	if apiResource.AccessTokenTtl != nil {
		model.AccessTokenTtl = types.NumberValue(big.NewFloat(*apiResource.AccessTokenTtl))
	} else {
		model.AccessTokenTtl = types.NumberNull()
	}

	// Conversion des scopes
	var scopeValues []attr.Value

	if apiResource.Scopes != nil {
		for _, s := range *apiResource.Scopes {
			var createdAt basetypes.NumberValue
			if s.CreatedAt != nil {
				createdAt = basetypes.NewNumberValue(big.NewFloat(*s.CreatedAt))
			} else {
				createdAt = basetypes.NewNumberNull()
			}

			scopeValue := ScopesValue{
				CreatedAt:   createdAt,
				Description: basetypes.NewStringValue(s.Description),
				Id:          basetypes.NewStringValue(s.ID),
				Name:        basetypes.NewStringValue(s.Name),
				ResourceId:  basetypes.NewStringValue(s.ResourceId),
				TenantId:    basetypes.NewStringValue(s.TenantId),
				state:       attr.ValueStateKnown,
			}

			objVal, d := scopeValue.ToObjectValue(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}

			scopeValues = append(scopeValues, objVal)
		}
	}

	listType := ScopesValue{}.Type(ctx)
	listVal, d := types.ListValue(listType, scopeValues)
	diags.Append(d...)
	model.Scopes = listVal

	return diags
}
