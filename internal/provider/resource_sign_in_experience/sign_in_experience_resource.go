package resource_sign_in_experience

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *signInExperienceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state SignInExperienceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	signInExperienceModel, diags := NewSignInExperienceBuilder(ctx).FromTfPlan(&plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	signInExperienceModel, err := r.client.SignInExperienceUpdate(ctx, signInExperienceModel)
	if err != nil {
		resp.Diagnostics.AddError("Error updating sign-in experience", err.Error())
		return
	}

	diags = convertToTerraformModel(ctx, signInExperienceModel, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *signInExperienceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SignInExperienceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	signInExperience, err := r.client.SignInExperienceGet(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading sign-in experience", err.Error())
		return
	}

	diags = convertToTerraformModel(ctx, signInExperience, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *signInExperienceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state SignInExperienceModel
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

	signInExperienceModel, diags := NewSignInExperienceBuilder(ctx).FromTfPlan(&plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	signInExperienceModel, err := r.client.SignInExperienceUpdate(ctx, signInExperienceModel)
	if err != nil {
		resp.Diagnostics.AddError("Error updating sign-in experience", err.Error())
		return
	}

	diags = convertToTerraformModel(ctx, signInExperienceModel, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *signInExperienceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"Deleting this resource is not supported. The resource will be stay in place on the API but at the next terraform implementation, the resource will adapt the input to the API content.",
	)

	resp.State.RemoveResource(ctx)
}

func convertListToSlice(ctx context.Context, list types.List) ([]string, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return []string{}, nil
	}
	var result []string
	diags := list.ElementsAs(ctx, &result, false)

	return result, diags
}

func stringSliceToList(slice []string) types.List {
	values := make([]attr.Value, len(slice))
	for i, s := range slice {
		values[i] = types.StringValue(s)
	}
	return types.ListValueMust(types.StringType, values)
}
