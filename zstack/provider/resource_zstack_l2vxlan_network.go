// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &l2vxlanNetworkResource{}
	_ resource.ResourceWithConfigure   = &l2vxlanNetworkResource{}
	_ resource.ResourceWithImportState = &l2vxlanNetworkResource{}
)

type l2vxlanNetworkResource struct {
	client *client.ZSClient
}

type l2vxlanNetworkResourceModel struct {
	Uuid              types.String `tfsdk:"uuid"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Vni               types.Int64  `tfsdk:"vni"`
	PoolUuid          types.String `tfsdk:"pool_uuid"`
	ZoneUuid          types.String `tfsdk:"zone_uuid"`
	PhysicalInterface types.String `tfsdk:"physical_interface"`
	Type              types.String `tfsdk:"type"`
}

func L2VxlanNetworkResource() resource.Resource {
	return &l2vxlanNetworkResource{}
}

func (r *l2vxlanNetworkResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*client.ZSClient)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", request.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *l2vxlanNetworkResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_l2vxlan_network"
}

func (r *l2vxlanNetworkResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Manages an L2 VXLAN network in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the L2 VXLAN network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the L2 VXLAN network.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the L2 VXLAN network.",
			},
			"vni": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The VXLAN Network Identifier (VNI).",
			},
			"pool_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the VXLAN pool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"zone_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The UUID of the zone.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"physical_interface": schema.StringAttribute{
				Computed:    true,
				Description: "The physical interface of the L2 VXLAN network.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the L2 VXLAN network.",
			},
		},
	}
}

func (r *l2vxlanNetworkResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan l2vxlanNetworkResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	params := param.CreateL2VxlanNetworkParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateL2VxlanNetworkParamDetail{
			Vni:         intPtrFromInt64OrNil(plan.Vni),
			PoolUuid:    plan.PoolUuid.ValueString(),
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			ZoneUuid:    stringPtrOrNil(plan.ZoneUuid.ValueString()),
		},
	}

	network, err := r.client.CreateL2VxlanNetwork(params)
	if err != nil {
		response.Diagnostics.AddError("Error creating L2 VXLAN network", err.Error())
		return
	}

	plan.Uuid = types.StringValue(network.UUID)
	plan.Name = types.StringValue(network.Name)
	plan.Description = stringValueOrNull(network.Description)
	plan.Vni = types.Int64Value(int64(network.Vni))
	plan.PoolUuid = types.StringValue(network.PoolUuid)
	plan.ZoneUuid = stringValueOrNull(network.ZoneUuid)
	plan.PhysicalInterface = stringValueOrNull(network.PhysicalInterface)
	plan.Type = stringValueOrNull(network.Type)

	diags = response.State.Set(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *l2vxlanNetworkResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state l2vxlanNetworkResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	network, err := findResourceByQuery(r.client.QueryL2VxlanNetwork, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError("Error reading L2 VXLAN network", err.Error())
		return
	}

	state.Name = types.StringValue(network.Name)
	state.Description = stringValueOrNull(network.Description)
	state.Vni = types.Int64Value(int64(network.Vni))
	state.PoolUuid = types.StringValue(network.PoolUuid)
	state.ZoneUuid = stringValueOrNull(network.ZoneUuid)
	state.PhysicalInterface = stringValueOrNull(network.PhysicalInterface)
	state.Type = stringValueOrNull(network.Type)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *l2vxlanNetworkResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	response.Diagnostics.AddError(
		"Update Not Supported",
		"L2 VXLAN networks cannot be updated. To modify the network, please delete and recreate it.",
	)
}

func (r *l2vxlanNetworkResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state l2vxlanNetworkResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteVxlanL2Network(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("Error deleting L2 VXLAN network", err.Error())
		return
	}
}

func (r *l2vxlanNetworkResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), request, response)
}
