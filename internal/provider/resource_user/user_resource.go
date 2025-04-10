package resource_user

import (
	"context"

	"github.com/Lenstra/terraform-provider-logto/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

type userResource struct {
	client *client.Client
}

func UserResource() resource.Resource {
	return &userResource{}
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = UserResourceSchema(ctx)
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		return
	}
	r.client = client
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	user, err := r.client.UserGet(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving user during import", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), user.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("primary_email"), user.PrimaryEmail)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("username"), user.Username)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), user.Name)...)

	if user.Profile != nil {
		profilePath := path.Root("profile")

		if user.Profile.FamilyName != "" {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, profilePath.AtName("family_name"), user.Profile.FamilyName)...)
		}

		if user.Profile.GivenName != "" {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, profilePath.AtName("given_name"), user.Profile.GivenName)...)
		}

		if user.Profile.MiddleName != "" {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, profilePath.AtName("middle_name"), user.Profile.MiddleName)...)
		}

		if user.Profile.Nickname != "" {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, profilePath.AtName("nickname"), user.Profile.Nickname)...)
		}
	}
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile := &client.Profile{
		FamilyName: "",
		GivenName:  "",
		MiddleName: "",
		Nickname:   "",
	}

	user := &client.UserModel{
		Name:         plan.Name.ValueString(),
		PrimaryEmail: plan.PrimaryEmail.ValueString(),
		Username:     plan.Username.ValueString(),
		Profile:      profile,
	}

	user, err := r.client.UserCreate(ctx, user)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", err.Error())
		return
	}

	state := UserModel{
		Id:           types.StringValue(user.ID),
		Name:         types.StringValue(user.Name),
		PrimaryEmail: types.StringValue(user.PrimaryEmail),
		Username:     types.StringValue(user.Username),
		Profile: ProfileValue{ // ProfileValue generated by terraform generator
			FamilyName: types.StringValue(user.Profile.FamilyName),
			GivenName:  types.StringValue(user.Profile.GivenName),
			MiddleName: types.StringValue(user.Profile.MiddleName),
			Nickname:   types.StringValue(user.Profile.Nickname),
		},
	}

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

	r.updateUserState(user, &state)

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

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := &client.UserModel{
		ID:           state.Id.ValueString(),
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

	user, err := r.client.UserUpdate(ctx, user)
	if err != nil {
		resp.Diagnostics.AddError("Error updating user", err.Error())
		return
	}

	r.updateUserState(user, &plan)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UserDelete(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
	}
}

func (r *userResource) updateUserState(user *client.UserModel, model *UserModel) {
	model.Id = types.StringValue(user.ID)
	model.PrimaryEmail = types.StringValue(user.PrimaryEmail)
	model.Username = types.StringValue(user.Username)
	model.Name = types.StringValue(user.Name)

	if user.Profile != nil {
		if !model.Profile.FamilyName.IsNull() {
			model.Profile.FamilyName = types.StringValue(user.Profile.FamilyName)
		}
		if !model.Profile.GivenName.IsNull() {
			model.Profile.GivenName = types.StringValue(user.Profile.GivenName)
		}
		if !model.Profile.MiddleName.IsNull() {
			model.Profile.MiddleName = types.StringValue(user.Profile.MiddleName)
		}
		if !model.Profile.Nickname.IsNull() {
			model.Profile.Nickname = types.StringValue(user.Profile.Nickname)
		}
	}
}
