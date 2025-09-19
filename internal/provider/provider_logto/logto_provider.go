package provider_logto

import (
	"context"
	"os"

	"github.com/Lenstra/terraform-provider-logto/internal/provider/resource_api_resource"
	"github.com/Lenstra/terraform-provider-logto/internal/provider/resource_api_resource_scope"
	"github.com/Lenstra/terraform-provider-logto/internal/provider/resource_application"
	"github.com/Lenstra/terraform-provider-logto/internal/provider/resource_assign_roles_to_user"
	"github.com/Lenstra/terraform-provider-logto/internal/provider/resource_role"
	"github.com/Lenstra/terraform-provider-logto/internal/provider/resource_user"
	"github.com/rs/zerolog"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
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

// Metadata returns the provider type name.
func (p *logtoProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "logto"
}

// Schema defines the provider-level schema for configuration data.
func (p *logtoProvider) Schema(ctx context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = LogtoProviderSchema(ctx)
}

func (p *logtoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring logto client")

	//Retrieve provider data from configuration
	var config LogtoModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.Hostname.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("hostname"),
			"Unknown Logto hostname",
			"The provider cannot create the logto API client as there is an unknown configuration value for the Logto hostname. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LOGTO_HOSTNAME environment variable.",
		)
	}
	if config.Resource.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("resource"),
			"Unknown Logto resource",
			"The provider cannot create the logto API client as there is an unknown configuration value for the Logto resource. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LOGTO_RESOURCE environment variable.",
		)
	}
	if config.ApplicationId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("application_id"),
			"Unknown Logto application ID",
			"The provider cannot create the logto API client as there is an unknown configuration value for the Logto application ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LOGTO_APPLICATION_ID environment variable.",
		)
	}
	if config.ApplicationSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("application_secret"),
			"Unknown Logto application secret",
			"The provider cannot create the logto API client as there is an unknown configuration value for the Logto application secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LOGTO_APPLICATION_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	hostname := os.Getenv("LOGTO_HOSTNAME")
	resource := os.Getenv("LOGTO_RESOURCE")
	applicationID := os.Getenv("LOGTO_APPLICATION_ID")
	applicationSecret := os.Getenv("LOGTO_APPLICATION_SECRET")

	if !config.Hostname.IsNull() {
		hostname = config.Hostname.ValueString()
	}
	if !config.Resource.IsNull() {
		resource = config.Resource.ValueString()
	}
	if !config.ApplicationId.IsNull() {
		applicationID = config.ApplicationId.ValueString()
	}
	if !config.ApplicationSecret.IsNull() {
		applicationSecret = config.ApplicationSecret.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if hostname == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("hostname"),
			"Missing Logto hostname",
			"The provider cannot create the Logto API client as there is a missing or empty value for the Logto hostname. "+
				"Set the access_token value in the configuration or use the LOGTO_HOSTNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if applicationID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("application_id"),
			"Missing Logto application ID",
			"The provider cannot create the Logto API client as there is a missing or empty value for the Logto application ID. "+
				"Set the access_token value in the configuration or use the LOGTO_APPLICATION_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if applicationSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("application_secret"),
			"Missing Logto application secret",
			"The provider cannot create the Logto API client as there is a missing or empty value for the Logto application secret. "+
				"Set the access_token value in the configuration or use the LOGTO_APPLICATION_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "logto_hostname", hostname)
	ctx = tflog.SetField(ctx, "logto_application_id", applicationID)
	ctx = tflog.SetField(ctx, "logto_application_secret", applicationSecret)

	tflog.Debug(ctx, "Creating Logto client")

	conf := &client.Config{
		Hostname:          hostname,
		Resource:          resource,
		ApplicationID:     applicationID,
		ApplicationSecret: applicationSecret,
	}

	if os.Getenv("TF_PROVIDER_LOGTO_LOG") != "" {
		conf.Logger = zerolog.New(os.Stdout)
	}

	apiClient, err := client.NewClient(conf)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build Logto client", err.Error())
		return
	}

	// Make the Logto client available during DataSource and Resource
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
		resource_application.ApplicationResource,
		resource_user.UserResource,
		resource_api_resource.ApiResourceResource,
		resource_api_resource_scope.ApiResourceScopeResource,
		resource_role.RoleResource,
		resource_assign_roles_to_user.AssignRolesToUserResource,
	}
}
