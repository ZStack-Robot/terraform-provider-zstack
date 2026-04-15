// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"

	"terraform-provider-zstack/zstack/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &virtualRouterOfferingResource{}
	_ resource.ResourceWithConfigure   = &virtualRouterOfferingResource{}
	_ resource.ResourceWithImportState = &virtualRouterOfferingResource{}
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
	PublicNetworkUuid     types.String `tfsdk:"public_network_uuid"`
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
		Params: param.CreateVirtualRouterOfferingParamDetail{
			Name:                  plan.Name.ValueString(),
			Description:           stringPtr(plan.Description.ValueString()),
			CpuNum:                int(plan.CpuNum.ValueInt64()),
			MemorySize:            utils.MBToBytes(plan.MemorySize.ValueInt64()),
			ManagementNetworkUuid: plan.ManagementNetworkUuid.ValueString(),
			ZoneUuid:              plan.ZoneUuid.ValueString(),
			ImageUuid:             plan.ImageUuid.ValueString(),
			Type:                  stringPtr("VirtualRouter"),
		},
	}

	if !plan.PublicNetworkUuid.IsNull() && plan.PublicNetworkUuid.ValueString() != "" {
		offerParam.Params.PublicNetworkUuid = stringPtr(plan.PublicNetworkUuid.ValueString())
	}
	if !plan.IsDefault.IsNull() {
		offerParam.Params.IsDefault = boolPtr(plan.IsDefault.ValueBool())
	}

	virtual_router, err := r.client.CreateVirtualRouterOffering(offerParam)
	tflog.Debug(ctx, "Received virtual router offering", map[string]interface{}{
		"virtual_router": virtual_router,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Virtual Router Offering",
			"Could not create virtual router offering, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(virtual_router.UUID)
	plan.Name = types.StringValue(virtual_router.Name)
	plan.Description = types.StringValue(virtual_router.Description)
	plan.CpuNum = types.Int64Value(int64(virtual_router.CpuNum))
	plan.MemorySize = types.Int64Value(utils.BytesToMB(virtual_router.MemorySize))
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


	err := r.client.DeleteInstanceOffering(state.Uuid.ValueString(), param.DeleteModeEnforcing)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Virtual Router Offering",
			"Could not delete virtual router offering, unexpected error: "+err.Error(),
		)
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
	virtual_router, err := findResourceByGet(r.client.GetVirtualRouterOffering, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Virtual Router Offering",
			"Could not read virtual router offering UUID "+state.Uuid.ValueString()+": "+err.Error(),
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
	state.MemorySize = types.Int64Value(utils.BytesToMB(virtual_router.MemorySize))
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the virtual router offering. This is a mandatory field and must be unique.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description of the virtual router offering, providing additional context or details about the configuration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cpu_num": schema.Int64Attribute{
				Required:    true,
				Description: "The number of CPUs allocated to the virtual router offering. This is a mandatory field.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"memory_size": schema.Int64Attribute{
				Required:    true,
				Description: "The amount of memory  allocated to the virtual router offering. This is a mandatory field, in megabytes (MB)",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"management_network_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the management network associated with the virtual router offering. This is a mandatory field.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"public_network_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the public network associated with the virtual router offering. If not specified, it will share the same network UUID as the management network or vice versa, depending on the configuration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"zone_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the zone where the virtual router offering is deployed. This is a mandatory field.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"image_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the image used by the virtual router offering. This is a mandatory field.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_default": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates whether this virtual router offering is the default offering. Defaults to `false` if not specified.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of the virtual router offering. Defaults to 'VirtualRouter' if not specified.",
				Validators: []validator.String{
					stringvalidator.OneOf("VirtualRouter"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *virtualRouterOfferingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"Virtual Router Offering resource does not support updates. Please recreate the resource instead.",
	)
}

func (r *virtualRouterOfferingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
