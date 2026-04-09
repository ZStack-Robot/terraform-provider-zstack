// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &vmNicResource{}
	_ resource.ResourceWithConfigure   = &vmNicResource{}
	_ resource.ResourceWithImportState = &vmNicResource{}
)

type vmNicResource struct {
	client *client.ZSClient
}

type vmNicResourceModel struct {
	Uuid           types.String `tfsdk:"uuid"`
	L3NetworkUuid  types.String `tfsdk:"l3_network_uuid"`
	Ip             types.String `tfsdk:"ip"`
	Mac            types.String `tfsdk:"mac"`
	VmInstanceUuid types.String `tfsdk:"vm_instance_uuid"`
	Netmask        types.String `tfsdk:"netmask"`
	Gateway        types.String `tfsdk:"gateway"`
	State          types.String `tfsdk:"state"`
}

func VmNicResource() resource.Resource {
	return &vmNicResource{}
}

func (r *vmNicResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *vmNicResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vm_nic"
}

func (r *vmNicResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage VM NICs in ZStack. " +
			"A VM NIC is a virtual network interface card that connects a virtual machine to an L3 network.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the VM NIC.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"l3_network_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the L3 network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ip": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The IP address of the VM NIC. If not specified, ZStack will assign one automatically.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mac": schema.StringAttribute{
				Computed:    true,
				Description: "The MAC address of the VM NIC.",
			},
			"vm_instance_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the VM instance this NIC belongs to.",
			},
			"netmask": schema.StringAttribute{
				Computed:    true,
				Description: "The network mask of the VM NIC.",
			},
			"gateway": schema.StringAttribute{
				Computed:    true,
				Description: "The gateway address of the VM NIC.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the VM NIC.",
			},
		},
	}
}

func (r *vmNicResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan vmNicResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateVmNicParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVmNicParamDetail{
			L3NetworkUuid: plan.L3NetworkUuid.ValueString(),
			Ip:            stringPtrOrNil(plan.Ip.ValueString()),
		},
	}

	vmNic, err := r.client.CreateVmNic(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating VM NIC",
			"Could not create vm nic, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(vmNic.UUID)
	plan.L3NetworkUuid = types.StringValue(vmNic.L3NetworkUuid)
	plan.Ip = types.StringValue(vmNic.Ip)
	plan.Mac = types.StringValue(vmNic.Mac)
	plan.VmInstanceUuid = stringValueOrNull(vmNic.VmInstanceUuid)
	plan.Netmask = types.StringValue(vmNic.Netmask)
	plan.Gateway = types.StringValue(vmNic.Gateway)
	plan.State = types.StringValue(vmNic.State)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vmNicResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state vmNicResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	vmNic, err := findResourceByQuery(r.client.QueryVmNic, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		tflog.Warn(ctx, "Unable to query VM NICs. It may have been deleted.: "+err.Error())
		state = vmNicResourceModel{
			Uuid: types.StringValue(""),
		}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	state.Uuid = types.StringValue(vmNic.UUID)
	state.L3NetworkUuid = types.StringValue(vmNic.L3NetworkUuid)
	state.Ip = types.StringValue(vmNic.Ip)
	state.Mac = types.StringValue(vmNic.Mac)
	state.VmInstanceUuid = stringValueOrNull(vmNic.VmInstanceUuid)
	state.Netmask = types.StringValue(vmNic.Netmask)
	state.Gateway = types.StringValue(vmNic.Gateway)
	state.State = types.StringValue(vmNic.State)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vmNicResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	response.Diagnostics.AddError(
		"Update not supported",
		"VM NIC resource does not support updates. Please recreate the resource instead.",
	)
}

func (r *vmNicResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state vmNicResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "VM NIC UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteVmNic(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting VM NIC",
			"Could not delete vm nic, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *vmNicResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
