// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
)

var (
	_ resource.Resource              = &virtualRouterOfferingResource{}
	_ resource.ResourceWithConfigure = &virtualRouterOfferingResource{}
)

type virtualRouterOfferingResource struct {
	client *client.ZSClient
}

type virtualRouterOfferingResourceModel struct {
	Uuid                  types.String `tfsdk:"uuid"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	CpuNum                types.Int64  `tfsdk:"cpu_num"`
	MemorySize            types.Int64  `tfsdk:"memory_size"`
	ManagementNetworkUuid types.String `tfsdk:"management_network_uuid"`
	ZoneUuid              types.String `tfsdk:"zone_uuid"`
	ImageUuid             types.String `tfsdk:"image_uuid"`
	IsDefault             types.Bool   `tfsdk:"is_default"`
	Type                  types.String `tfsdk:"type"`
}

// Configure implements resource.ResourceWithConfigure.
func (r *virtualRouterOfferingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", req.ProviderData),
		)
		return
	}

	r.client = client
}

func VirtualRouterOfferingResource() resource.Resource {
	return &virtualRouterOfferingResource{}
}

// Create implements resource.Resource.
func (r *virtualRouterOfferingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan virtualRouterOfferingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Configuring ZStack client")
	offerParam := param.CreateVirtualRouterOfferingParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVirtualRouterOfferingDetailParam{
			Name:                  plan.Name.ValueString(),
			Description:           plan.Description.ValueString(),
			CpuNum:                int(plan.CpuNum.ValueInt64()),
			MemorySize:            plan.MemorySize.ValueInt64(),
			ManagementNetworkUuid: plan.ManagementNetworkUuid.ValueString(),
			ZoneUuid:              plan.ZoneUuid.ValueString(),
			ImageUuid:             plan.ImageUuid.ValueString(),
			IsDefault:             bool(plan.IsDefault.ValueBool()),
			Type:                  "VirtualRouter",
		},
	}

	virtual_router, err := r.client.CreateVirtualRouterOffering(offerParam)
	tflog.Debug(ctx, "Received virtual router offering", map[string]interface{}{
		"virtual_router": virtual_router,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Add virtual router offering to ZStack", "Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(virtual_router.UUID)
	plan.Name = types.StringValue(virtual_router.Name)
	plan.Description = types.StringValue(virtual_router.Description)
	plan.CpuNum = types.Int64Value(int64(virtual_router.CpuNum))
	plan.MemorySize = types.Int64Value(int64(virtual_router.MemorySize))
	plan.Type = types.StringValue(virtual_router.Type)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete implements resource.Resource.
func (r *virtualRouterOfferingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state virtualRouterOfferingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "virtual router offering uuid is empty, so nothing to delete, skip it")
		return
	}

	err := r.client.DeleteInstanceOffering(state.Uuid.ValueString(), param.DeleteModeEnforcing)

	if err != nil {
		resp.Diagnostics.AddError("fail to delete virtual router image", ""+err.Error())
		return
	}
}

// Metadata implements resource.Resource.
func (r *virtualRouterOfferingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_router_offer"
}

// Read implements resource.Resource.
func (r *virtualRouterOfferingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state virtualRouterOfferingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	virtual_router, err := r.client.GetVirtualRouterOffering(state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting virtual router offering uuid", "Could not read virtual router uuid: "+err.Error(),
		)
		return
	}

	state.Type = types.StringValue(virtual_router.Type)
	if virtual_router.Type == "" {
		state.Type = types.StringValue("VirtualRouter")
	} else {
		state.Type = types.StringValue(virtual_router.Type)
	}

	state.Uuid = types.StringValue(virtual_router.UUID)
	state.Description = types.StringValue(virtual_router.Description)
	state.Name = types.StringValue(virtual_router.Name)
	state.CpuNum = types.Int64Value(int64(virtual_router.CpuNum))
	state.MemorySize = types.Int64Value(virtual_router.MemorySize)
	state.Type = types.StringValue(virtual_router.Type)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements resource.Resource.
func (r *virtualRouterOfferingResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource allows you to manage virtual router offerings in ZStack. " +
			"A virtual router offering defines the configuration and resource settings for virtual router instances, such as CPU, memory, and management network. " +
			"You can define the offering's properties, such as its name, description, CPU and memory allocation, and the associated management network.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier (UUID) of the virtual router offering. Automatically generated by ZStack.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the virtual router offering. This is a mandatory field and must be unique.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description of the virtual router offering, providing additional context or details about the configuration.",
			},
			"cpu_num": schema.Int64Attribute{
				Required:    true,
				Description: "The number of CPUs allocated to the virtual router offering. This is a mandatory field.",
			},
			"memory_size": schema.Int64Attribute{
				Required:    true,
				Description: "The amount of memory (in bytes) allocated to the virtual router offering. This is a mandatory field.",
			},
			"management_network_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the management network associated with the virtual router offering. This is a mandatory field.",
			},
			"zone_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the zone where the virtual router offering is deployed. This is a mandatory field.",
			},
			"image_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the image used by the virtual router offering. This is a mandatory field.",
			},
			"is_default": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates whether this virtual router offering is the default offering. Defaults to `false` if not specified.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of the virtual router offering. Defaults to 'VirtualRouter' if not specified.",
			},
		},
	}
}

func (r *virtualRouterOfferingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}
