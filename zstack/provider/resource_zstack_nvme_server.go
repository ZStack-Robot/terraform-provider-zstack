// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
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
	_ resource.Resource                = &nvmeServerResource{}
	_ resource.ResourceWithConfigure   = &nvmeServerResource{}
	_ resource.ResourceWithImportState = &nvmeServerResource{}
)

type nvmeServerResource struct {
	client *client.ZSClient
}

type nvmeServerModel struct {
	Uuid      types.String `tfsdk:"uuid"`
	Name      types.String `tfsdk:"name"`
	Ip        types.String `tfsdk:"ip"`
	Port      types.Int64  `tfsdk:"port"`
	Transport types.String `tfsdk:"transport"`
	State     types.String `tfsdk:"state"`
}

func NvmeServerResource() resource.Resource {
	return &nvmeServerResource{}
}

func (r *nvmeServerResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *nvmeServerResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_nvme_server"
}

func (r *nvmeServerResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage NVMe servers in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the NVMe server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the NVMe server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ip": schema.StringAttribute{
				Required:    true,
				Description: "The IP address of the NVMe server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The port of the NVMe server.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"transport": schema.StringAttribute{
				Required:    true,
				Description: "The transport protocol of the NVMe server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the NVMe server.",
			},
		},
	}
}

func (r *nvmeServerResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan nvmeServerModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.AddNvmeServerParam{
		BaseParam: param.BaseParam{},
		Params: param.AddNvmeServerParamDetail{
			Name:      plan.Name.ValueString(),
			Ip:        plan.Ip.ValueString(),
			Transport: plan.Transport.ValueString(),
		},
	}

	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		p.Params.Port = intPtr(int(plan.Port.ValueInt64()))
	}

	server, err := r.client.AddNvmeServer(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating NVMe Server",
			"Could not create NVMe server, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(server.UUID)
	plan.Name = stringValueOrNull(server.Name)
	plan.Ip = types.StringValue(server.Ip)
	plan.Port = types.Int64Value(int64(server.Port))
	plan.Transport = types.StringValue(server.Transport)
	plan.State = stringValueOrNull(server.State)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *nvmeServerResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state nvmeServerModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	server, err := findResourceByQuery(r.client.QueryNvmeServer, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		tflog.Warn(ctx, "Unable to query NVMe server. It may have been deleted.: "+err.Error())
		state = nvmeServerModel{Uuid: types.StringValue("")}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	state.Uuid = types.StringValue(server.UUID)
	state.Name = stringValueOrNull(server.Name)
	state.Ip = types.StringValue(server.Ip)
	state.Port = types.Int64Value(int64(server.Port))
	state.Transport = types.StringValue(server.Transport)
	state.State = stringValueOrNull(server.State)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (r *nvmeServerResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	response.Diagnostics.AddError(
		"Update not supported",
		"NVMe Server resource does not support updates. Please recreate the resource instead.",
	)
}

func (r *nvmeServerResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state nvmeServerModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "NVMe server UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteNvmeServer(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("Error deleting NVMe Server", "Could not delete NVMe server UUID "+state.Uuid.ValueString()+": "+err.Error())
		return
	}
}

func (r *nvmeServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
