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
	_ resource.Resource                = &accessKeyResource{}
	_ resource.ResourceWithConfigure   = &accessKeyResource{}
	_ resource.ResourceWithImportState = &accessKeyResource{}
)

type accessKeyResource struct {
	client *client.ZSClient
}

type accessKeyModel struct {
	Uuid            types.String `tfsdk:"uuid"`
	AccountUuid     types.String `tfsdk:"account_uuid"`
	UserUuid        types.String `tfsdk:"user_uuid"`
	Description     types.String `tfsdk:"description"`
	AccessKeyID     types.String `tfsdk:"access_key_id"`
	AccessKeySecret types.String `tfsdk:"access_key_secret"`
	State           types.String `tfsdk:"state"`
}

func AccessKeyResource() resource.Resource {
	return &accessKeyResource{}
}

func (r *accessKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_key"
}

func (r *accessKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZStack Access Key resource.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Description: "The UUID of the access key.",
				Computed:    true,
			},
			"account_uuid": schema.StringAttribute{
				Description: "The UUID of the account.",
				Required:    true,
			},
			"user_uuid": schema.StringAttribute{
				Description: "The UUID of the user.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the access key.",
				Optional:    true,
				Computed:    true,
			},
			"access_key_id": schema.StringAttribute{
				Description: "The access key ID.",
				Computed:    true,
				Sensitive:   true,
			},
			"access_key_secret": schema.StringAttribute{
				Description: "The access key secret. Only available after creation.",
				Computed:    true,
				Sensitive:   true,
			},
			"state": schema.StringAttribute{
				Description: "The state of the access key.",
				Computed:    true,
			},
		},
	}
}

func (r *accessKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *accessKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan accessKeyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createParam := param.CreateAccessKeyParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateAccessKeyParamDetail{
			AccountUuid: plan.AccountUuid.ValueString(),
			UserUuid:    plan.UserUuid.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	result, err := r.client.CreateAccessKey(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating access key",
			"Could not create access key, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.AccountUuid = types.StringValue(result.AccountUuid)
	plan.UserUuid = types.StringValue(result.UserUuid)
	plan.Description = stringValueOrNull(result.Description)
	plan.AccessKeyID = types.StringValue(result.AccessKeyID)
	plan.AccessKeySecret = types.StringValue(result.AccessKeySecret)
	plan.State = types.StringValue(result.State)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Access key created", map[string]interface{}{
		"uuid": result.UUID,
	})
}

func (r *accessKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state accessKeyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	queryParam.AddQ("uuid=" + state.Uuid.ValueString())

	accessKeys, err := r.client.QueryAccessKey(&queryParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading access key",
			"Could not read access key UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	if len(accessKeys) == 0 {
		tflog.Warn(ctx, "Access key not found, removing from state", map[string]interface{}{
			"uuid": state.Uuid.ValueString(),
		})
		state.Uuid = types.StringValue("")
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		return
	}

	accessKey := accessKeys[0]
	state.Uuid = types.StringValue(accessKey.UUID)
	state.AccountUuid = types.StringValue(accessKey.AccountUuid)
	state.UserUuid = types.StringValue(accessKey.UserUuid)
	state.Description = stringValueOrNull(accessKey.Description)
	state.AccessKeyID = types.StringValue(accessKey.AccessKeyID)
	// Preserve AccessKeySecret from state as it's only returned on Create
	state.State = types.StringValue(accessKey.State)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *accessKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan accessKeyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Access keys do not support updates, so just copy plan to state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Access key update (no-op)", map[string]interface{}{
		"uuid": plan.Uuid.ValueString(),
	})
}

func (r *accessKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state accessKeyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAccessKey(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting access key",
			"Could not delete access key, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "Access key deleted", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
}

func (r *accessKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
