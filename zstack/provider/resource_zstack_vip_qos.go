// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
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
	_ resource.Resource                = &vipQosResource{}
	_ resource.ResourceWithConfigure   = &vipQosResource{}
	_ resource.ResourceWithImportState = &vipQosResource{}
)

type vipQosResource struct {
	client *client.ZSClient
}

type vipQosModel struct {
	Uuid              types.String `tfsdk:"uuid"`
	VipUuid           types.String `tfsdk:"vip_uuid"`
	Port              types.Int64  `tfsdk:"port"`
	OutboundBandwidth types.Int64  `tfsdk:"outbound_bandwidth"`
	InboundBandwidth  types.Int64  `tfsdk:"inbound_bandwidth"`
	Type              types.String `tfsdk:"type"`
}

func VipQosResource() resource.Resource {
	return &vipQosResource{}
}

func (r *vipQosResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vip_qos"
}

func (r *vipQosResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages VIP QoS settings in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the VIP QoS.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vip_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the VIP to configure QoS for.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The port number for the QoS rule.",
			},
			"outbound_bandwidth": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The outbound bandwidth limit in bits per second.",
			},
			"inbound_bandwidth": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The inbound bandwidth limit in bits per second.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the VIP QoS.",
			},
		},
	}
}

func (r *vipQosResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *vipQosResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vipQosModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	setParam := param.SetVipQosParam{
		BaseParam: param.BaseParam{},
		Params:    param.SetVipQosParamDetail{},
	}

	if !plan.Port.IsNull() {
		port := int(plan.Port.ValueInt64())
		setParam.Params.Port = &port
	}
	if !plan.OutboundBandwidth.IsNull() {
		bandwidth := plan.OutboundBandwidth.ValueInt64()
		setParam.Params.OutboundBandwidth = &bandwidth
	}
	if !plan.InboundBandwidth.IsNull() {
		bandwidth := plan.InboundBandwidth.ValueInt64()
		setParam.Params.InboundBandwidth = &bandwidth
	}

	result, err := r.client.SetVipQos(plan.VipUuid.ValueString(), setParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating VIP QoS",
			"Could not create VIP QoS, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.VipUuid = types.StringValue(result.VipUuid)
	plan.Port = types.Int64Value(int64(result.Port))
	plan.OutboundBandwidth = types.Int64Value(result.OutboundBandwidth)
	plan.InboundBandwidth = types.Int64Value(result.InboundBandwidth)
	plan.Type = stringValueOrNull(result.Type)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "VIP QoS created", map[string]interface{}{
		"uuid":     result.UUID,
		"vip_uuid": result.VipUuid,
	})
}

func (r *vipQosResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vipQosModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetVipQos(state.VipUuid.ValueString())
	if err != nil {
		tflog.Warn(ctx, "VIP QoS not found, removing from state: "+err.Error())
		resp.State.RemoveResource(ctx)
		return
	}

	state.Uuid = types.StringValue(result.UUID)
	state.VipUuid = types.StringValue(result.VipUuid)
	state.Port = types.Int64Value(int64(result.Port))
	state.OutboundBandwidth = types.Int64Value(result.OutboundBandwidth)
	state.InboundBandwidth = types.Int64Value(result.InboundBandwidth)
	state.Type = stringValueOrNull(result.Type)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vipQosResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan vipQosModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	setParam := param.SetVipQosParam{
		BaseParam: param.BaseParam{},
		Params:    param.SetVipQosParamDetail{},
	}

	if !plan.Port.IsNull() {
		port := int(plan.Port.ValueInt64())
		setParam.Params.Port = &port
	}
	if !plan.OutboundBandwidth.IsNull() {
		bandwidth := plan.OutboundBandwidth.ValueInt64()
		setParam.Params.OutboundBandwidth = &bandwidth
	}
	if !plan.InboundBandwidth.IsNull() {
		bandwidth := plan.InboundBandwidth.ValueInt64()
		setParam.Params.InboundBandwidth = &bandwidth
	}

	result, err := r.client.SetVipQos(plan.VipUuid.ValueString(), setParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating VIP QoS",
			"Could not update VIP QoS, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.VipUuid = types.StringValue(result.VipUuid)
	plan.Port = types.Int64Value(int64(result.Port))
	plan.OutboundBandwidth = types.Int64Value(result.OutboundBandwidth)
	plan.InboundBandwidth = types.Int64Value(result.InboundBandwidth)
	plan.Type = stringValueOrNull(result.Type)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "VIP QoS updated", map[string]interface{}{
		"uuid":     result.UUID,
		"vip_uuid": result.VipUuid,
	})
}

func (r *vipQosResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vipQosModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteVipQos(state.VipUuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting VIP QoS",
			"Could not delete VIP QoS, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "VIP QoS deleted", map[string]interface{}{
		"vip_uuid": state.VipUuid.ValueString(),
	})
}

func (r *vipQosResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("vip_uuid"), req, resp)
}
