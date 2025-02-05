package resource_application

import (
	"context"
	"fmt"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type applicationResource struct {
	client *client.Client
}

func ApplicationResource() resource.Resource {
	return &applicationResource{}
}

func (r *applicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *applicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ApplicationResourceSchema(ctx)
}

func (r *applicationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.Client)
}

func (r *applicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApplicationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	application := &client.ApplicationModel{
		Name: plan.Name.ValueString(),
		Type: plan.Type.ValueString(),
	}

	if !plan.Description.IsNull() {
		application.Description = plan.Description.ValueString()
	}

	oidcClientMetadata, customClientMetadata, err := r.buildClientMetadata(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error building metadata", err.Error())
		return
	}

	application.OidcClientMetadata = oidcClientMetadata
	application.CustomClientMetadata = customClientMetadata

	application, err = r.client.ApplicationCreate(ctx, application)
	if err != nil {
		resp.Diagnostics.AddError("Error creating application", err.Error())
		return
	}

	secretsValue, diags := r.getSecrets(ctx, application.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Secrets = secretsValue

	r.updateApplicationState(application, &plan)

	diags = resp.State.Set(ctx, plan)
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

	secretsValue, diags := r.getSecrets(ctx, application.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Secrets = secretsValue

	r.updateApplicationState(application, &state)

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

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	application := &client.ApplicationModel{
		ID:   state.Id.ValueString(),
		Name: plan.Name.ValueString(),
		Type: state.Type.ValueString(),
	}

	if !plan.Description.IsNull() {
		application.Description = plan.Description.ValueString()
	}

	oidcClientMetadata, customClientMetadata, err := r.buildClientMetadata(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error building metadata", err.Error())
		return
	}

	application.OidcClientMetadata = oidcClientMetadata
	application.CustomClientMetadata = customClientMetadata

	application, err = r.client.ApplicationUpdate(ctx, application)
	if err != nil {
		resp.Diagnostics.AddError("Error updating application", err.Error())
		return
	}

	secretsValue, diags := r.getSecrets(ctx, application.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Secrets = secretsValue

	r.updateApplicationState(application, &plan)

	diags = resp.State.Set(ctx, plan)
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

func (r *applicationResource) updateApplicationState(application *client.ApplicationModel, plan *ApplicationModel) {
	plan.Id = types.StringValue(application.ID)
	plan.TenantId = types.StringValue(application.TenantId)
	plan.Name = types.StringValue(application.Name)
	plan.Description = types.StringValue(application.Description)
	plan.Type = types.StringValue(application.Type)

	if application.OidcClientMetadata != nil {
		updateListField(application.OidcClientMetadata.RedirectUris, &plan.RedirectUris)
		updateListField(application.OidcClientMetadata.PostLogoutRedirectUris, &plan.PostLogoutRedirectUris)
	}

	if application.CustomClientMetadata != nil {
		updateListField(application.CustomClientMetadata.CorsAllowedOrigins, &plan.CorsAllowedOrigins)
	}
}

func updateListField(slice []string, plan *types.List) {
	if len(slice) == 0 {
		*plan = types.ListNull(types.StringType)
	} else {
		*plan = stringSliceToList(slice)
	}
}

func (r *applicationResource) buildClientMetadata(plan ApplicationModel) (*client.OidcClientMetadata, *client.CustomClientMetadata, error) {
	oidcClientMetadata := &client.OidcClientMetadata{
		RedirectUris:           []string{},
		PostLogoutRedirectUris: []string{},
	}

	if !plan.RedirectUris.IsNull() {
		redirectUris, err := convertListToSlice(plan.RedirectUris)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert redirect uris: %w", err)
		}
		oidcClientMetadata.RedirectUris = redirectUris
	}

	if !plan.PostLogoutRedirectUris.IsNull() {
		postLogoutRedirectUris, err := convertListToSlice(plan.PostLogoutRedirectUris)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert post logout redirect uris: %w", err)
		}
		oidcClientMetadata.PostLogoutRedirectUris = postLogoutRedirectUris
	}

	var customClientMetadata *client.CustomClientMetadata
	if !plan.CorsAllowedOrigins.IsNull() {
		corsAllowedOrigins, err := convertListToSlice(plan.CorsAllowedOrigins)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert cors allowed origins: %w", err)
		}
		customClientMetadata = &client.CustomClientMetadata{
			CorsAllowedOrigins: corsAllowedOrigins,
		}
	}

	return oidcClientMetadata, customClientMetadata, nil
}

func (r *applicationResource) getSecrets(ctx context.Context, applicationID string) (types.Map, diag.Diagnostics) {
	secrets, err := r.client.GetApplicationSecrets(ctx, applicationID)
	if err != nil {
		return types.MapNull(types.StringType), diag.Diagnostics{
			diag.NewErrorDiagnostic("Error getting secrets", err.Error()),
		}
	}

	secretsMap := make(map[string]attr.Value, len(secrets))
	for _, v := range secrets {
		secretsMap[v.Name] = types.StringValue(v.Value)
	}

	return types.MapValue(types.StringType, secretsMap)
}
