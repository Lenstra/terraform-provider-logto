package resource_assign_roles_to_user

import (
	"context"
	"strings"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *assignRolesToUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state AssignRolesToUserModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleIds, userId := decodePlan(plan)

	err := r.client.AssignRolesForUser(ctx, roleIds, userId)
	if err != nil {
		resp.Diagnostics.AddError("Error assagning role(s) to user", err.Error())
		return
	}

	state = plan // API return nothing so just keep plan values

	id := convertToStateId(userId, roleIds.RoleIds)
	state.Id = types.StringValue(id)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *assignRolesToUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AssignRolesToUserModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Id.IsNull() || state.Id.IsUnknown() {
		resp.Diagnostics.AddError("Invalid Id", "Resource Id is null or unknown.")
		return
	}

	parts := strings.Split(state.Id.ValueString(), "/")
	userId := parts[0]

	clientRoles, err := r.client.GetRolesForUser(ctx, userId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading user roles", err.Error())
		return
	}

	var roleVals []attr.Value
	for _, r := range clientRoles {
		roleVals = append(roleVals, types.StringValue(r.ID))
	}

	tfModel := &AssignRolesToUserModel{
		Id:      state.Id,
		RoleIds: types.SetValueMust(types.StringType, roleVals),
		UserId:  state.UserId,
	}

	diags = resp.State.Set(ctx, tfModel)
	resp.Diagnostics.Append(diags...)
}

func (r *assignRolesToUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state AssignRolesToUserModel
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

	roleIds, userId := decodePlan(plan)

	err := r.client.UpdateRolesForUser(ctx, roleIds, userId)
	if err != nil {
		resp.Diagnostics.AddError("Error updating assigned role(s) for user", err.Error())
		return
	}

	state = plan // API return nothing so just keep plan values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *assignRolesToUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AssignRolesToUserModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleIdslist := convertSetToSlice(state.RoleIds)
	for _, roleId := range roleIdslist {
		err := r.client.DeleteRolesForUser(ctx, roleId, state.UserId.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error when removing role from user", err.Error())
		}
	}
}

func decodePlan(plan AssignRolesToUserModel) (*client.RoleIdsModel, string) {
	list := convertSetToSlice(plan.RoleIds)

	roleIds := &client.RoleIdsModel{
		RoleIds: list,
	}
	userId := plan.UserId.ValueString()

	return roleIds, userId
}

func convertSetToSlice(set types.Set) []string {
	elems := set.Elements()
	result := make([]string, len(elems))
	for i, e := range elems {
		result[i] = e.(types.String).ValueString()
	}

	return result
}

func convertToStateId(userId string, roleIds []string) string {
	return userId + "/" + strings.Join(roleIds, "-")
}
