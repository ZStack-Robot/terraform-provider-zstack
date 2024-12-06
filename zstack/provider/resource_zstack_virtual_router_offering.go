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
		tflog.Warn(ctx, "virtual router image uuid is empty, so nothing to delete, skip it")
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
			"Error getting virtual router image uuid", "Could not read image uuid"+virtual_router.Name+": "+err.Error(),
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
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the virtual router image. Automatically generated by ZStack.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the virtual router image. This is a mandatory field.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description of the virtual router image, providing additional context or details.",
			},
			"cpu_num": schema.Int64Attribute{
				Required:    true,
				Description: "The URL where the virtual router image is located. This can be a file path or an HTTP link.",
			},
			"memory_size": schema.Int64Attribute{
				Required:    true,
				Description: "The type of media for the image. Examples include 'ISO' or 'Template'.",
			},
			"management_network_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The virtual router OS type that the image is optimized for. ",
			},
			"zone_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The platform type of the image. Defaults to Linux.",
			},
			"image_uuid": schema.StringAttribute{
				Required:    true,
				Description: "A list of UUIDs for the backup storages where the image is stored.",
			},
			"is_default": schema.BoolAttribute{
				Optional:    true,
				Description: "The architecture of the image, such as 'x86_64' or 'arm64'.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Indicates if the VirtIO drivers are required for the image.",
			},
		},
	}
}

func (r *virtualRouterOfferingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}
