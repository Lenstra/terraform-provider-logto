package datasource_application

import (
	"context"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/Lenstra/terraform-provider-logto/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func (d *applicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ApplicationModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	application, err := d.client.ApplicationGet(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading application", err.Error())
		return
	}

	secretsValue, diags := utils.GetSecrets(ctx, d.client, application.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = convertToTerraformModel(ctx, application, secretsValue, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func convertToTerraformModel(ctx context.Context, app *client.ApplicationModel, secretsValue basetypes.MapValue, model *ApplicationModel) (diags diag.Diagnostics) {
	*model = ApplicationModel{
		Id:          types.StringValue(app.ID),
		TenantId:    types.StringValue(app.TenantId),
		Name:        types.StringValue(app.Name),
		Description: types.StringValue(app.Description),
		Type:        types.StringValue(app.Type),
		Secrets:     secretsValue,
	}

	if app.OidcClientMetadata != nil {
		model.RedirectUris, diags = utils.ConvertList(ctx, types.StringType, app.OidcClientMetadata.RedirectUris)
		if diags.HasError() {
			return
		}
		model.PostLogoutRedirectUris, diags = utils.ConvertList(ctx, types.StringType, app.OidcClientMetadata.PostLogoutRedirectUris)
		if diags.HasError() {
			return
		}
	}

	var corsAllowedOrigins []string
	if app.CustomClientMetadata != nil {
		corsAllowedOrigins = app.CustomClientMetadata.CorsAllowedOrigins
	}
	model.CorsAllowedOrigins, diags = utils.ConvertList(ctx, types.StringType, corsAllowedOrigins)
	if diags.HasError() {
		return
	}

	return
}
