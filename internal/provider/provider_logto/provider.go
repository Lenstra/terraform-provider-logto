package provider_logto

import (
	"context"
	"os"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &logtoProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &logtoProvider{
			version: version,
		}
	}
}

// logtoProvider is the provider implementation.
type logtoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// logtoProviderModel maps provider schema data to a Go type.
type logtoProviderModel struct {
	TenantId    types.String `tfsdk:"tenantId"`
	AccessToken types.String `tfsdk:"accessKey"`
}

// Metadata returns the provider type name.
func (p *logtoProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "logto"
}

// Schema defines the provider-level schema for configuration data.
func (p *logtoProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Logto terraform provider developed by Lenstra.",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Required:    true,
				Description: "API tenant_id for you instance, can be set as environment variable LOGTO_TENANT_ID",
			},
			"access_token": schema.StringAttribute{
				Required:    true,
				Description: "API key for you instance, can be set as environment variable LOGTO_ACCESS_TOKEN",
			},
		},
	}
}

func (p *logtoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring logto client")

	//Retrieve provider data from configuration
	var config logtoProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.TenantId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("tenant_id"),
			"Unknown Logto tenant_id",
			"The provider cannot create the logto API client as there is an unknown configuration value for the Logto tenant_id. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LOGTO_TENANT_ID environment variable.",
		)
	}

	if config.AccessToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_token"),
			"Unknown Logto access token",
			"The provider cannot create the logto API client as there is an unknown configuration value for the Logto access token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LOGTO_ACCESS_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	tenant_id := os.Getenv("LOGTO_TENANT_ID")
	access_token := os.Getenv("LOGTO_ACCESS_TOKEN")

	if !config.TenantId.IsNull() {
		tenant_id = config.TenantId.ValueString()
	}

	if !config.AccessToken.IsNull() {
		access_token = config.AccessToken.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if tenant_id == "" {
		tenant_id = "default"
	}

	if access_token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_token"),
			"Missing Logto access token",
			"The provider cannot create the Logto API client as there is a missing or empty value for the Logto access token. "+
				"Set the access_token value in the configuration or use the LOGTO_ACCESS_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "logto_host", tenant_id)
	ctx = tflog.SetField(ctx, "logto_access_token", access_token)

	tflog.Debug(ctx, "Creating Logto client")

	apiClient := client.NewClient(tenant_id, access_token)

	// Make the SimpleMDM client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient

	tflog.Info(ctx, "Configured Logto client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *logtoProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *logtoProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		//ApplicationResource,
	}
}
