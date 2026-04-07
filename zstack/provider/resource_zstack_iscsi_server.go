// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &iscsiServerResource{}
	_ resource.ResourceWithConfigure   = &iscsiServerResource{}
	_ resource.ResourceWithImportState = &iscsiServerResource{}
)

type iscsiServerResource struct {
	client *client.ZSClient
}

type iscsiServerModel struct {
	Uuid             types.String `tfsdk:"uuid"`
	Name             types.String `tfsdk:"name"`
	Ip               types.String `tfsdk:"ip"`
	Port             types.Int64  `tfsdk:"port"`
	ChapUserName     types.String `tfsdk:"chap_user_name"`
	ChapUserPassword types.String `tfsdk:"chap_user_password"`
	State            types.String `tfsdk:"state"`
}

func IscsiServerResource() resource.Resource {
	return &iscsiServerResource{}
}

func (r *iscsiServerResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*client.ZSClient)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", request.ProviderData),
		)

		return
	}
	r.client = client
}

func (r *iscsiServerResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_iscsi_server"
}

func (r *iscsiServerResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage iSCSI servers in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the iSCSI server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the iSCSI server.",
			},
			"ip": schema.StringAttribute{
				Required:    true,
				Description: "The IP address of the iSCSI server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The port of the iSCSI server.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"chap_user_name": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The CHAP username for the iSCSI server.",
			},
			"chap_user_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The CHAP user password for the iSCSI server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the iSCSI server.",
			},
		},
	}
}

func (r *iscsiServerResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan iscsiServerModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.AddIscsiServerParam{
		BaseParam: param.BaseParam{},
		Params: param.AddIscsiServerParamDetail{
			Name: plan.Name.ValueString(),
			Ip:   plan.Ip.ValueString(),
		},
	}

	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		p.Params.Port = intPtr(int(plan.Port.ValueInt64()))
	}

	if !plan.ChapUserName.IsNull() && !plan.ChapUserName.IsUnknown() {
		p.Params.ChapUserName = stringPtr(plan.ChapUserName.ValueString())
	}

	if !plan.ChapUserPassword.IsNull() && !plan.ChapUserPassword.IsUnknown() {
		p.Params.ChapUserPassword = stringPtr(plan.ChapUserPassword.ValueString())
	}

	server, err := r.client.AddIscsiServer(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to create iSCSI server",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(server.UUID)
	plan.Name = types.StringValue(server.Name)
	plan.Ip = types.StringValue(server.Ip)
	plan.Port = types.Int64Value(int64(server.Port))
	plan.ChapUserName = stringValueOrNull(server.ChapUserName)
	plan.State = stringValueOrNull(server.State)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *iscsiServerResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state iscsiServerModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	serverList, err := r.client.QueryIscsiServer(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query iSCSI server. It may have been deleted.: "+err.Error())
		state = iscsiServerModel{Uuid: types.StringValue("")}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false
	for _, server := range serverList {
		if server.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(server.UUID)
			state.Name = types.StringValue(server.Name)
			state.Ip = types.StringValue(server.Ip)
			state.Port = types.Int64Value(int64(server.Port))
			state.ChapUserName = stringValueOrNull(server.ChapUserName)
			state.State = stringValueOrNull(server.State)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "iSCSI server not found. It might have been deleted outside of Terraform.")
		state = iscsiServerModel{Uuid: types.StringValue("")}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (r *iscsiServerResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan iscsiServerModel
	var state iscsiServerModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateIscsiServerParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateIscsiServerParamDetail{
			Name: plan.Name.ValueString(),
		},
	}

	if !plan.ChapUserName.IsNull() && !plan.ChapUserName.IsUnknown() {
		p.Params.ChapUserName = stringPtr(plan.ChapUserName.ValueString())
	}

	if !plan.ChapUserPassword.IsNull() && !plan.ChapUserPassword.IsUnknown() {
		p.Params.ChapUserPassword = stringPtr(plan.ChapUserPassword.ValueString())
	}

	server, err := r.client.UpdateIscsiServer(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to update iSCSI server",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(server.UUID)
	plan.Name = types.StringValue(server.Name)
	plan.Ip = types.StringValue(server.Ip)
	plan.Port = types.Int64Value(int64(server.Port))
	plan.ChapUserName = stringValueOrNull(server.ChapUserName)
	plan.State = stringValueOrNull(server.State)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *iscsiServerResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state iscsiServerModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "iSCSI server UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteIscsiServer(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("fail to delete iSCSI server", err.Error())
		return
	}
}

func (r *iscsiServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
