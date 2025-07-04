package resource_application

import (
	"context"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state ApplicationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	application := decodePlan(ctx, plan)
	application, err := r.client.ApplicationCreate(ctx, application)
	if err != nil {
		resp.Diagnostics.AddError("Error creating application", err.Error())
		return
	}

	diags = convertToTerraformModel(ctx, application, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ApplicationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	application, err := r.client.ApplicationGet(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading application", err.Error())
		return
	}

	if application == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = convertToTerraformModel(ctx, application, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApplicationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	application := decodePlan(ctx, plan)

	application, err := r.client.ApplicationUpdate(ctx, application)
	if err != nil {
		resp.Diagnostics.AddError("Error updating application", err.Error())
		return
	}

	diags = convertToTerraformModel(ctx, application, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ApplicationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.ApplicationDelete(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting application", err.Error())
	}
}

func decodePlan(ctx context.Context, plan ApplicationModel) *client.ApplicationModel {
	model := &client.ApplicationModel{
		ID:                 plan.Id.ValueString(),
		Name:               plan.Name.ValueString(),
		Type:               plan.Type.ValueString(),
		Description:        plan.Description.ValueString(),
		IsThirdParty:       plan.IsThirdParty.ValueBool(),
		OidcClientMetadata: &client.OidcClientMetadata{},
	}
	plan.RedirectUris.ElementsAs(ctx, &model.OidcClientMetadata.RedirectUris, true)
	plan.PostLogoutRedirectUris.ElementsAs(ctx, &model.OidcClientMetadata.PostLogoutRedirectUris, true)

	if !plan.CorsAllowedOrigins.IsNull() {
		model.CustomClientMetadata = &client.CustomClientMetadata{}
		plan.CorsAllowedOrigins.ElementsAs(ctx, &model.CustomClientMetadata.CorsAllowedOrigins, true)
	}

	return model
}

func convertToTerraformModel(ctx context.Context, app *client.ApplicationModel, model *ApplicationModel) (diags diag.Diagnostics) {
	*model = ApplicationModel{
		Id:           types.StringValue(app.ID),
		TenantId:     types.StringValue(app.TenantId),
		Name:         types.StringValue(app.Name),
		Description:  types.StringValue(app.Description),
		Type:         types.StringValue(app.Type),
		IsThirdParty: types.BoolValue(app.IsThirdParty),
		IsAdmin:      types.BoolValue(app.IsAdmin),
	}

	if app.OidcClientMetadata != nil {
		model.RedirectUris, diags = basetypes.NewListValueFrom(ctx, types.StringType, app.OidcClientMetadata.RedirectUris)
		if diags.HasError() {
			return
		}
		model.PostLogoutRedirectUris, diags = basetypes.NewListValueFrom(ctx, types.StringType, app.OidcClientMetadata.PostLogoutRedirectUris)
		if diags.HasError() {
			return
		}
	}

	if app.CustomClientMetadata != nil {
		model.CorsAllowedOrigins, diags = basetypes.NewListValueFrom(ctx, types.StringType, app.CustomClientMetadata.CorsAllowedOrigins)
		if diags.HasError() {
			return
		}
	}
	return
}
