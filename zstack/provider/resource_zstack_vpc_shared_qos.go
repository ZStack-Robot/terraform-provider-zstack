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
	_ resource.Resource                = &vpcSharedQosResource{}
	_ resource.ResourceWithConfigure   = &vpcSharedQosResource{}
	_ resource.ResourceWithImportState = &vpcSharedQosResource{}
)

type vpcSharedQosResource struct {
	client *client.ZSClient
}

type vpcSharedQosModel struct {
	Uuid          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	VpcUuid       types.String `tfsdk:"vpc_uuid"`
	L3NetworkUuid types.String `tfsdk:"l3_network_uuid"`
	Bandwidth     types.Int64  `tfsdk:"bandwidth"`
}

func VpcSharedQosResource() resource.Resource {
	return &vpcSharedQosResource{}
}

func (r *vpcSharedQosResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_shared_qos"
}

func (r *vpcSharedQosResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages VPC Shared QoS resources in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the VPC Shared QoS.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VPC Shared QoS.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the VPC Shared QoS.",
			},
			"vpc_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the VPC.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"l3_network_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the L3 network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bandwidth": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The bandwidth limit in bits per second.",
			},
		},
	}
}

func (r *vpcSharedQosResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *vpcSharedQosResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpcSharedQosModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createParam := param.CreateVpcSharedQosParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVpcSharedQosParamDetail{
			Name:          plan.Name.ValueString(),
			VpcUuid:       plan.VpcUuid.ValueString(),
			L3NetworkUuid: plan.L3NetworkUuid.ValueString(),
		},
	}

	if !plan.Description.IsNull() {
		createParam.Params.Description = stringPtr(plan.Description.ValueString())
	}
	if !plan.Bandwidth.IsNull() {
		bandwidth := plan.Bandwidth.ValueInt64()
		createParam.Params.Bandwidth = &bandwidth
	}

	result, err := r.client.CreateVpcSharedQos(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating VPC Shared QoS",
			"Could not create VPC Shared QoS, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = stringValueOrNull(result.Description)
	plan.VpcUuid = types.StringValue(result.VpcUuid)
	plan.L3NetworkUuid = types.StringValue(result.L3NetworkUuid)
	plan.Bandwidth = types.Int64Value(result.Bandwidth)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "VPC Shared QoS created", map[string]interface{}{
		"uuid": result.UUID,
		"name": result.Name,
	})
}

func (r *vpcSharedQosResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vpcSharedQosModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	results, err := r.client.QueryVpcSharedQos(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query VPC Shared QoS: "+err.Error())
		resp.State.RemoveResource(ctx)
		return
	}

	found := false
	for _, result := range results {
		if result.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(result.UUID)
			state.Name = types.StringValue(result.Name)
			state.Description = stringValueOrNull(result.Description)
			state.VpcUuid = types.StringValue(result.VpcUuid)
			state.L3NetworkUuid = types.StringValue(result.L3NetworkUuid)
			state.Bandwidth = types.Int64Value(result.Bandwidth)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "VPC Shared QoS not found, removing from state")
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vpcSharedQosResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan vpcSharedQosModel
	var state vpcSharedQosModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateParam := param.UpdateVpcSharedQosParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateVpcSharedQosParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	result, err := r.client.UpdateVpcSharedQos(plan.Uuid.ValueString(), updateParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating VPC Shared QoS",
			"Could not update VPC Shared QoS, unexpected error: "+err.Error(),
		)
		return
	}

	if !plan.Bandwidth.Equal(state.Bandwidth) && !plan.Bandwidth.IsNull() {
		result, err = r.client.ChangeVpcSharedQosBandwidth(plan.Uuid.ValueString(), param.ChangeVpcSharedQosBandwidthParam{
			BaseParam: param.BaseParam{},
			Params: param.ChangeVpcSharedQosBandwidthParamDetail{
				Bandwidth: plan.Bandwidth.ValueInt64(),
			},
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating VPC Shared QoS bandwidth",
				"Could not update VPC Shared QoS bandwidth, unexpected error: "+err.Error(),
			)
			return
		}
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = stringValueOrNull(result.Description)
	plan.VpcUuid = types.StringValue(result.VpcUuid)
	plan.L3NetworkUuid = types.StringValue(result.L3NetworkUuid)
	plan.Bandwidth = types.Int64Value(result.Bandwidth)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "VPC Shared QoS updated", map[string]interface{}{
		"uuid": result.UUID,
		"name": result.Name,
	})
}

func (r *vpcSharedQosResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vpcSharedQosModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteVpcSharedQos(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting VPC Shared QoS",
			"Could not delete VPC Shared QoS, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "VPC Shared QoS deleted", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
}

func (r *vpcSharedQosResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
