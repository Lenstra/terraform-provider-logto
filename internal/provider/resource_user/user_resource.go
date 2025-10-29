package resource_user

import (
	"context"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state UserModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	user, roleIds, diags := decodePlan(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.UserCreate(ctx, user)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", err.Error())
		return
	}

	// Put the user into the state before assigning roles in case of error during roles assignment
	convertToTerraformModel(ctx, user, nil, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

	if roleIds != nil {
		err = r.client.AssignRolesForUser(ctx, roleIds, user.ID)
		if err != nil {
			resp.Diagnostics.AddError("Error during assignation of role(s) for user", err.Error())
			return
		}
	}

	convertToTerraformModel(ctx, user, roleIds, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.UserGet(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading user", err.Error())
		return
	}

	if user == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	var rolesIds *client.RoleIdsModel
	if !state.RoleIds.IsNull() && !state.RoleIds.IsUnknown() {
		roles, err := r.client.GetRolesForUser(ctx, user.ID)
		if err != nil {
			resp.Diagnostics.AddError("Error reading role(s) of user", err.Error())
			return
		}

		rolesIds = &client.RoleIdsModel{}
		for _, r := range roles {
			rolesIds.RoleIds = append(rolesIds.RoleIds, r.ID)
		}
	}

	convertToTerraformModel(ctx, user, rolesIds, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state UserModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	user, roleIds, diags := decodePlan(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.UserUpdate(ctx, user)
	if err != nil {
		resp.Diagnostics.AddError("Error updating user", err.Error())
		return
	}

	if roleIds != nil {
		err := r.client.UpdateRolesForUser(ctx, roleIds, user.ID)
		if err != nil {
			resp.Diagnostics.AddError("Error updating role(s) of user", err.Error())
			return
		}
	}

	convertToTerraformModel(ctx, user, roleIds, &state)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.RoleIds.IsNull() && !state.RoleIds.IsUnknown() {
		roleIdslist, diag := convertSetToSlice(state.RoleIds)
		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, roleId := range roleIdslist {
			err := r.client.DeleteRolesForUser(ctx, roleId, state.Id.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("Error when removing role from user", err.Error())
			}
		}
	}

	err := r.client.UserDelete(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
	}
}

func decodePlan(_ context.Context, plan UserModel) (*client.UserModel, *client.RoleIdsModel, diag.Diagnostics) {
	var clientRolIds *client.RoleIdsModel

	if !plan.RoleIds.IsNull() && !plan.RoleIds.IsUnknown() {
		list, diags := convertSetToSlice(plan.RoleIds)
		if diags.HasError() {
			return nil, nil, diags
		}
		clientRolIds = &client.RoleIdsModel{
			RoleIds: list,
		}
	}

	user := &client.UserModel{
		ID:           plan.Id.ValueString(),
		PrimaryEmail: plan.PrimaryEmail.ValueString(),
		Username:     plan.Username.ValueString(),
		Name:         plan.Name.ValueString(),
		Profile: &client.Profile{
			FamilyName: plan.Profile.FamilyName.ValueString(),
			GivenName:  plan.Profile.GivenName.ValueString(),
			MiddleName: plan.Profile.MiddleName.ValueString(),
			Nickname:   plan.Profile.Nickname.ValueString(),
		},
	}

	return user, clientRolIds, nil
}

func convertToTerraformModel(_ context.Context, user *client.UserModel, roleIds *client.RoleIdsModel, model *UserModel) {
	*model = UserModel{
		Id:           types.StringValue(user.ID),
		PrimaryEmail: types.StringValue(user.PrimaryEmail),
		Username:     types.StringValue(user.Username),
		Name:         types.StringValue(user.Name),
	}

	if user.Profile != nil {
		model.Profile = ProfileValue{
			FamilyName: types.StringValue(user.Profile.FamilyName),
			GivenName:  types.StringValue(user.Profile.GivenName),
			MiddleName: types.StringValue(user.Profile.MiddleName),
			Nickname:   types.StringValue(user.Profile.Nickname),
			state:      attr.ValueStateKnown,
		}
	}

	if roleIds == nil {
		model.RoleIds = types.SetNull(types.StringType)
	} else {
		var roleVals []attr.Value
		for _, r := range roleIds.RoleIds {
			roleVals = append(roleVals, types.StringValue(r))
		}
		model.RoleIds = types.SetValueMust(types.StringType, roleVals)
	}
}

func convertSetToSlice(set types.Set) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	elems := set.Elements()
	result := make([]string, 0, len(elems))

	for _, e := range elems {
		s, ok := e.(types.String)
		if !ok {
			diags.AddError(
				"Error converting element in set to string",
				"Expected a types.String but got a different type",
			)
			continue
		}
		result = append(result, s.ValueString())
	}

	return result, diags
}
