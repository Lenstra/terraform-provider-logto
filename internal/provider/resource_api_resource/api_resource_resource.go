package resource_api_resource

import (
	"context"
	"math/big"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &apiResourceResource{}
	_ resource.ResourceWithConfigure   = &apiResourceResource{}
	_ resource.ResourceWithImportState = &apiResourceResource{}
)

type apiResourceResource struct {
	client *client.Client
}

func ApiResourceResource() resource.Resource {
	return &apiResourceResource{}
}

func (r *apiResourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_resource"
}

func (r *apiResourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ApiResourceResourceSchema(ctx)
}

func (r *apiResourceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		return
	}
	r.client = client
}

func (r *apiResourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *apiResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApiResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResource := &client.ApiResourceModel{
		Name:      plan.Name.ValueString(),
		Indicator: plan.Indicator.ValueString(),
	}

	if !plan.AccessTokenTtl.IsUnknown() && !plan.AccessTokenTtl.IsNull() {
		accessTokenTtlValue, _ := plan.AccessTokenTtl.ValueBigFloat().Float64()
		apiResource.AccessTokenTtl = &accessTokenTtlValue
	}

	apiResource, err := r.client.ApiResourceCreate(ctx, apiResource)
	if err != nil {
		resp.Diagnostics.AddError("Error creating api_resource", err.Error())
		return
	}

	diag := r.updateApiResourceState(apiResource, &plan)
	resp.Diagnostics.Append(diag...)

	diags = resp.State.Set(ctx, plan)
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

	diag := r.updateApiResourceState(apiResource, &state)
	resp.Diagnostics.Append(diag...)

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

	apiResource := &client.ApiResourceModel{
		ID:   plan.Id.ValueString(),
		Name: plan.Name.ValueString(),
	}

	if !plan.AccessTokenTtl.IsUnknown() && !plan.AccessTokenTtl.IsNull() {
		ttl, _ := plan.AccessTokenTtl.ValueBigFloat().Float64()
		apiResource.AccessTokenTtl = &ttl
	}

	apiResource, err := r.client.ApiResourceUpdate(ctx, apiResource)
	if err != nil {
		resp.Diagnostics.AddError("Error updating api_resource", err.Error())
		return
	}

	diags = r.updateApiResourceState(apiResource, &state)
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

func (r *apiResourceResource) updateApiResourceState(apiResource *client.ApiResourceModel, model *ApiResourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	var scopesList []attr.Value

	var ScopeAttrTypes = map[string]attr.Type{
		"created_at":  types.NumberType,
		"description": types.StringType,
		"name":        types.StringType,
		"id":          types.StringType,
		"resource_id": types.StringType,
		"tenant_id":   types.StringType,
	}

	var ScopesType = types.ObjectType{
		AttrTypes: ScopeAttrTypes,
	}

	if apiResource.Scopes != nil {
		for _, scope := range *apiResource.Scopes {
			scopeObj, objDiags := types.ObjectValue(
				ScopeAttrTypes,
				map[string]attr.Value{
					"created_at":  types.NumberValue(big.NewFloat(*scope.CreatedAt)),
					"description": types.StringValue(scope.Description),
					"name":        types.StringValue(scope.Name),
					"id":          types.StringValue(scope.ID),
					"resource_id": types.StringValue(scope.ResourceId),
					"tenant_id":   types.StringValue(scope.TenantId),
				},
			)
			diags.Append(objDiags...)

			scopesList = append(scopesList, scopeObj)
		}
	}

	scopesListValue, listDiags := types.ListValue(ScopesType, scopesList)

	diags.Append(listDiags...)

	model.AccessTokenTtl = basetypes.NumberValue(types.Float64Value(*apiResource.AccessTokenTtl))
	model.Id = types.StringValue(apiResource.ID)
	model.Indicator = types.StringValue(apiResource.Indicator)
	model.IsDefault = types.BoolValue(*apiResource.IsDefault)
	model.Name = types.StringValue(apiResource.Name)
	model.Scopes = scopesListValue

	return diags
}
