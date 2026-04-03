// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

type userResource struct {
	client *client.ZSClient
}

type userResourceModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Password    types.String `tfsdk:"password"`
	Description types.String `tfsdk:"description"`
	AccountUuid types.String `tfsdk:"account_uuid"`
}

func UserResource() resource.Resource {
	return &userResource{}
}

func (r *userResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a User in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the user",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the user",
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The password for the user. This is a write-only sensitive field.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the user",
			},
			"account_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the account the user belongs to",
			},
		},
	}
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createParam := param.CreateUserParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateUserParamDetail{
			Name:        plan.Name.ValueString(),
			Password:    plan.Password.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	result, err := r.client.CreateUser(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			"Could not create user, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = stringValueOrNull(result.Description)
	plan.AccountUuid = types.StringValue(result.AccountUuid)
	// Password is write-only, preserve from plan

	tflog.Trace(ctx, "created a user")

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	queryParam.AddQ("uuid=" + state.Uuid.ValueString())

	users, err := r.client.QueryUser(&queryParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user",
			"Could not read user UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	if len(users) == 0 {
		tflog.Warn(ctx, "user not found, removing from state", map[string]interface{}{
			"uuid": state.Uuid.ValueString(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	user := users[0]

	state.Name = types.StringValue(user.Name)
	state.Description = stringValueOrNull(user.Description)
	state.AccountUuid = types.StringValue(user.AccountUuid)
	// Password is not returned by query, preserve from state

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state userResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateParam := param.UpdateUserParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateUserParamDetail{
			Uuid:        state.Uuid.ValueString(),
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	// Update password only if changed
	if !plan.Password.Equal(state.Password) {
		updateParam.Params.Password = stringPtr(plan.Password.ValueString())
	}

	result, err := r.client.UpdateUser(updateParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating user",
			"Could not update user, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = stringValueOrNull(result.Description)
	plan.AccountUuid = types.StringValue(result.AccountUuid)
	// Password is write-only, preserve from plan

	tflog.Trace(ctx, "updated a user")

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state userResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting user",
			"Could not delete user, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted a user")
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
