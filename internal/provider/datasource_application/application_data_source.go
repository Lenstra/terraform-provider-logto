package datasource_application

import (
	"context"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/Lenstra/terraform-provider-logto/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &applicationDataSource{}

type applicationDataSource struct {
	client *client.Client
}

func ApplicationDataSource() datasource.DataSource {
	return &applicationDataSource{}
}

func (d *applicationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (d *applicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ApplicationDataSourceSchema(ctx)
}

func (d *applicationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*client.Client); ok {
		d.client = client
	}
}

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
	state.Secrets = secretsValue
	state.Id = types.StringValue(application.ID)
	state.TenantId = types.StringValue(application.TenantId)
	state.Name = types.StringValue(application.Name)
	state.Description = types.StringValue(application.Description)
	state.Type = types.StringValue(application.Type)
	if application.OidcClientMetadata != nil {
		if len(application.OidcClientMetadata.RedirectUris) == 0 {
			state.RedirectUris = types.ListNull(types.StringType)
		} else {
			state.RedirectUris = utils.StringSliceToList(application.OidcClientMetadata.RedirectUris)
		}

		if len(application.OidcClientMetadata.PostLogoutRedirectUris) == 0 {
			state.PostLogoutRedirectUris = types.ListNull(types.StringType)
		} else {
			state.PostLogoutRedirectUris = utils.StringSliceToList(application.OidcClientMetadata.PostLogoutRedirectUris)
		}
	}
	if application.CustomClientMetadata != nil {
		if len(application.CustomClientMetadata.CorsAllowedOrigins) == 0 {
			state.CorsAllowedOrigins = types.ListNull(types.StringType)
		} else {
			state.CorsAllowedOrigins = utils.StringSliceToList(application.CustomClientMetadata.CorsAllowedOrigins)
		}
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
