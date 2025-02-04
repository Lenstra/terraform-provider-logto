package resource_application

import (
	"context"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &applicationResource{}
	_ resource.ResourceWithConfigure   = &applicationResource{}
	_ resource.ResourceWithImportState = &applicationResource{}
)

func ApplicationResource() resource.Resource {
	return &applicationResource{}
}

// appResource is the resource implementation.
type applicationResource struct {
	client *client.Client
}

// Schema implements resource.Resource.
func (r *applicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ApplicationResourceSchema(ctx)
}

// Configure adds the provider configured client to the resource.
func (r *applicationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.Client)
}

// Metadata returns the resource type name.
func (r *applicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *applicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new application
func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApplicationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	redirectUris, err := convertListToSlice(plan.RedirectUris)
	if err != nil {
		resp.Diagnostics.AddError("Redirect Uris Conversion Error", "Failed to convert redirect uris: "+err.Error())
		return
	}

	postLogoutRedirectUris, err := convertListToSlice(plan.PostLogoutRedirectUris)
	if err != nil {
		resp.Diagnostics.AddError("Post logout redirect uris Conversion Error", "Failed to convert post logout redirect uris: "+err.Error())
		return
	}

	var description string
	if !plan.Description.IsNull() {
		description = plan.Description.ValueString()
	}

	application, err := r.client.ApplicationCreate(
		plan.Name.ValueString(),
		description,
		plan.Type.ValueString(),
		redirectUris,
		postLogoutRedirectUris,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating application",
			"Could not create application, unexpected error: "+err.Error(),
		)
	}

	if application == nil {
		resp.Diagnostics.AddError(
			"Error creating application",
			"Received nil application from API but no error",
		)
		return
	}

	secretsMap := make(map[string]attr.Value)
	for k, v := range application.Secrets {
		secretsMap[k] = types.StringValue(v)
	}

	secretsValue, diags := types.MapValue(types.StringType, secretsMap)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = types.StringValue(application.LogtoDefaultStruct.Id)
	plan.TenantId = types.StringValue(application.LogtoDefaultStruct.TenantId)
	plan.Name = types.StringValue(application.LogtoDefaultStruct.Name)
	plan.Description = types.StringValue(application.LogtoDefaultStruct.Description)
	plan.Type = types.StringValue(application.Type)
	plan.Secrets = secretsValue

	if !plan.RedirectUris.IsNull() {
		plan.RedirectUris = stringSliceToList(application.OidcClientMetadata.RedirectUris)
	}
	if !plan.PostLogoutRedirectUris.IsNull() {
		plan.PostLogoutRedirectUris = stringSliceToList(application.OidcClientMetadata.PostLogoutRedirectUris)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete an application
func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ApplicationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.ApplicationDelete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Logto application",
			"Could not delete application ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}
}

func (r *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ApplicationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	application, err := r.client.ApplicationGet(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Logto application",
			"Could not read application ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	secretsMap := make(map[string]attr.Value)
	for k, v := range application.Secrets {
		secretsMap[k] = types.StringValue(v)
	}

	secretsValue, diags := types.MapValue(types.StringType, secretsMap)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Name = types.StringValue(application.LogtoDefaultStruct.Name)
	state.Description = types.StringValue(application.LogtoDefaultStruct.Description)
	state.TenantId = types.StringValue(application.LogtoDefaultStruct.TenantId)
	state.Type = types.StringValue(application.Type)
	state.Secrets = secretsValue

	if len(application.OidcClientMetadata.RedirectUris) == 0 {
		state.RedirectUris = types.ListNull(types.StringType)
	} else {
		state.RedirectUris = stringSliceToList(application.OidcClientMetadata.RedirectUris)
	}

	if len(application.OidcClientMetadata.PostLogoutRedirectUris) == 0 {
		state.PostLogoutRedirectUris = types.ListNull(types.StringType)
	} else {
		state.PostLogoutRedirectUris = stringSliceToList(application.OidcClientMetadata.PostLogoutRedirectUris)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApplicationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	redirectUris, err := convertListToSlice(plan.RedirectUris)
	if err != nil {
		resp.Diagnostics.AddError("Redirect Uris Conversion Error", "Failed to convert redirect uris: "+err.Error())
		return
	}

	postLogoutRedirectUris, err := convertListToSlice(plan.PostLogoutRedirectUris)
	if err != nil {
		resp.Diagnostics.AddError("Post logout redirect uris Conversion Error", "Failed to convert post logout redirect uris: "+err.Error())
		return
	}

	application, err := r.client.ApplicationUpdate(
		state.Id.ValueString(),
		plan.Name.ValueString(),
		plan.Description.ValueString(),
		redirectUris,
		postLogoutRedirectUris,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating application",
			"Could not update application, unexpected error: "+err.Error(),
		)
		return
	}

	secretsMap := make(map[string]attr.Value)
	for k, v := range application.Secrets {
		secretsMap[k] = types.StringValue(v)
	}

	secretsValue, diags := types.MapValue(types.StringType, secretsMap)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = types.StringValue(application.Id)
	plan.TenantId = types.StringValue(application.TenantId)
	plan.Name = types.StringValue(application.Name)
	plan.Description = types.StringValue(application.Description)
	plan.Type = types.StringValue(application.Type)
	plan.Secrets = secretsValue

	if len(application.OidcClientMetadata.RedirectUris) == 0 {
		plan.RedirectUris = types.ListNull(types.StringType)
	} else {
		plan.RedirectUris = stringSliceToList(application.OidcClientMetadata.RedirectUris)
	}

	if len(application.OidcClientMetadata.PostLogoutRedirectUris) == 0 {
		plan.PostLogoutRedirectUris = types.ListNull(types.StringType)
	} else {
		plan.PostLogoutRedirectUris = stringSliceToList((application.OidcClientMetadata.PostLogoutRedirectUris))
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func convertListToSlice(list types.List) ([]string, error) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var result []string
	list.ElementsAs(context.Background(), &result, false)
	return result, nil
}

func stringSliceToList(slice []string) types.List {
	values := make([]attr.Value, len(slice))
	for i, s := range slice {
		values[i] = types.StringValue(s)
	}
	return types.ListValueMust(types.StringType, values)
}
