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
	_ resource.Resource                = &baremetalPxeServerResource{}
	_ resource.ResourceWithConfigure   = &baremetalPxeServerResource{}
	_ resource.ResourceWithImportState = &baremetalPxeServerResource{}
)

type baremetalPxeServerResource struct {
	client *client.ZSClient
}

type baremetalPxeServerModel struct {
	Uuid             types.String `tfsdk:"uuid"`
	Name             types.String `tfsdk:"name"`
	ZoneUuid         types.String `tfsdk:"zone_uuid"`
	Description      types.String `tfsdk:"description"`
	Hostname         types.String `tfsdk:"hostname"`
	SshUsername      types.String `tfsdk:"ssh_username"`
	SshPassword      types.String `tfsdk:"ssh_password"`
	SshPort          types.Int64  `tfsdk:"ssh_port"`
	StoragePath      types.String `tfsdk:"storage_path"`
	DhcpInterface    types.String `tfsdk:"dhcp_interface"`
	DhcpRangeBegin   types.String `tfsdk:"dhcp_range_begin"`
	DhcpRangeEnd     types.String `tfsdk:"dhcp_range_end"`
	DhcpRangeNetmask types.String `tfsdk:"dhcp_range_netmask"`
	State            types.String `tfsdk:"state"`
	Status           types.String `tfsdk:"status"`
}

func BaremetalPxeServerResource() resource.Resource {
	return &baremetalPxeServerResource{}
}

func (r *baremetalPxeServerResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *baremetalPxeServerResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_baremetal_pxe_server"
}

func (r *baremetalPxeServerResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage baremetal PXE servers in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the PXE server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the PXE server.",
			},
			"zone_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The zone UUID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description for the PXE server.",
			},
			"hostname": schema.StringAttribute{
				Required:    true,
				Description: "The PXE server hostname.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ssh_username": schema.StringAttribute{
				Required:    true,
				Description: "The SSH username.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ssh_password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The SSH password.",
			},
			"ssh_port": schema.Int64Attribute{
				Optional:    true,
				Description: "The SSH port.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"storage_path": schema.StringAttribute{
				Required:    true,
				Description: "The storage path.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dhcp_interface": schema.StringAttribute{
				Required:    true,
				Description: "The DHCP interface.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dhcp_range_begin": schema.StringAttribute{
				Optional:    true,
				Description: "DHCP range begin.",
			},
			"dhcp_range_end": schema.StringAttribute{
				Optional:    true,
				Description: "DHCP range end.",
			},
			"dhcp_range_netmask": schema.StringAttribute{
				Optional:    true,
				Description: "DHCP range netmask.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status.",
			},
		},
	}
}

func (r *baremetalPxeServerResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan baremetalPxeServerModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.CreateBaremetalPxeServerParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateBaremetalPxeServerParamDetail{
			ZoneUuid:         plan.ZoneUuid.ValueString(),
			Name:             plan.Name.ValueString(),
			Description:      stringPtrOrNil(plan.Description.ValueString()),
			Hostname:         plan.Hostname.ValueString(),
			SshUsername:      plan.SshUsername.ValueString(),
			SshPassword:      plan.SshPassword.ValueString(),
			StoragePath:      plan.StoragePath.ValueString(),
			DhcpInterface:    plan.DhcpInterface.ValueString(),
			DhcpRangeBegin:   stringPtrOrNil(plan.DhcpRangeBegin.ValueString()),
			DhcpRangeEnd:     stringPtrOrNil(plan.DhcpRangeEnd.ValueString()),
			DhcpRangeNetmask: stringPtrOrNil(plan.DhcpRangeNetmask.ValueString()),
		},
	}

	if !plan.SshPort.IsNull() && !plan.SshPort.IsUnknown() {
		p.Params.SshPort = intPtr(int(plan.SshPort.ValueInt64()))
	}

	pxeServer, err := r.client.CreateBaremetalPxeServer(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to create baremetal PXE server",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(pxeServer.UUID)
	plan.Name = types.StringValue(pxeServer.Name)
	plan.ZoneUuid = types.StringValue(pxeServer.ZoneUuid)
	plan.Description = stringValueOrNull(pxeServer.Description)
	plan.Hostname = types.StringValue(pxeServer.Hostname)
	plan.SshUsername = types.StringValue(pxeServer.SshUsername)
	plan.SshPassword = types.StringValue(pxeServer.SshPassword)
	plan.SshPort = types.Int64Value(int64(pxeServer.SshPort))
	plan.StoragePath = types.StringValue(pxeServer.StoragePath)
	plan.DhcpInterface = types.StringValue(pxeServer.DhcpInterface)
	plan.DhcpRangeBegin = stringValueOrNull(pxeServer.DhcpRangeBegin)
	plan.DhcpRangeEnd = stringValueOrNull(pxeServer.DhcpRangeEnd)
	plan.DhcpRangeNetmask = stringValueOrNull(pxeServer.DhcpRangeNetmask)
	plan.State = stringValueOrNull(pxeServer.State)
	plan.Status = stringValueOrNull(pxeServer.Status)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *baremetalPxeServerResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state baremetalPxeServerModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	pxeServers, err := r.client.QueryBaremetalPxeServer(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query baremetal PXE servers. It may have been deleted.: "+err.Error())
		state = baremetalPxeServerModel{Uuid: types.StringValue("")}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false
	for _, pxeServer := range pxeServers {
		if pxeServer.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(pxeServer.UUID)
			state.Name = types.StringValue(pxeServer.Name)
			state.ZoneUuid = types.StringValue(pxeServer.ZoneUuid)
			state.Description = stringValueOrNull(pxeServer.Description)
			state.Hostname = types.StringValue(pxeServer.Hostname)
			state.SshUsername = types.StringValue(pxeServer.SshUsername)
			state.SshPassword = types.StringValue(pxeServer.SshPassword)
			state.SshPort = types.Int64Value(int64(pxeServer.SshPort))
			state.StoragePath = types.StringValue(pxeServer.StoragePath)
			state.DhcpInterface = types.StringValue(pxeServer.DhcpInterface)
			state.DhcpRangeBegin = stringValueOrNull(pxeServer.DhcpRangeBegin)
			state.DhcpRangeEnd = stringValueOrNull(pxeServer.DhcpRangeEnd)
			state.DhcpRangeNetmask = stringValueOrNull(pxeServer.DhcpRangeNetmask)
			state.State = stringValueOrNull(pxeServer.State)
			state.Status = stringValueOrNull(pxeServer.Status)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Baremetal PXE server not found. It might have been deleted outside of Terraform.")
		state = baremetalPxeServerModel{Uuid: types.StringValue("")}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (r *baremetalPxeServerResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan baremetalPxeServerModel
	var state baremetalPxeServerModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateBaremetalPxeServerParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateBaremetalPxeServerParamDetail{
			Name:             plan.Name.ValueString(),
			Description:      stringPtrOrNil(plan.Description.ValueString()),
			DhcpRangeBegin:   stringPtrOrNil(plan.DhcpRangeBegin.ValueString()),
			DhcpRangeEnd:     stringPtrOrNil(plan.DhcpRangeEnd.ValueString()),
			DhcpRangeNetmask: stringPtrOrNil(plan.DhcpRangeNetmask.ValueString()),
		},
	}

	pxeServer, err := r.client.UpdateBaremetalPxeServer(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to update baremetal PXE server",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(pxeServer.UUID)
	plan.Name = types.StringValue(pxeServer.Name)
	plan.ZoneUuid = types.StringValue(pxeServer.ZoneUuid)
	plan.Description = stringValueOrNull(pxeServer.Description)
	plan.Hostname = types.StringValue(pxeServer.Hostname)
	plan.SshUsername = types.StringValue(pxeServer.SshUsername)
	plan.SshPassword = types.StringValue(pxeServer.SshPassword)
	plan.SshPort = types.Int64Value(int64(pxeServer.SshPort))
	plan.StoragePath = types.StringValue(pxeServer.StoragePath)
	plan.DhcpInterface = types.StringValue(pxeServer.DhcpInterface)
	plan.DhcpRangeBegin = stringValueOrNull(pxeServer.DhcpRangeBegin)
	plan.DhcpRangeEnd = stringValueOrNull(pxeServer.DhcpRangeEnd)
	plan.DhcpRangeNetmask = stringValueOrNull(pxeServer.DhcpRangeNetmask)
	plan.State = stringValueOrNull(pxeServer.State)
	plan.Status = stringValueOrNull(pxeServer.Status)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *baremetalPxeServerResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state baremetalPxeServerModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Baremetal PXE server UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteBaremetalPxeServer(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("fail to delete baremetal PXE server", err.Error())
		return
	}
}

func (r *baremetalPxeServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
