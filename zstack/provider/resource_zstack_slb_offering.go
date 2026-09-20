// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"

	"terraform-provider-zstack/zstack/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	_ resource.Resource                = &slbOfferingResource{}
	_ resource.ResourceWithConfigure   = &slbOfferingResource{}
	_ resource.ResourceWithImportState = &slbOfferingResource{}
)

type slbOfferingResource struct {
	client *client.ZSClient
}

type slbOfferingResourceModel struct {
	Uuid                  types.String `tfsdk:"uuid"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	CpuNum                types.Int64  `tfsdk:"cpu_num"`
	MemorySize            types.Int64  `tfsdk:"memory_size"`
	ManagementNetworkUuid types.String `tfsdk:"management_network_uuid"`
	ZoneUuid              types.String `tfsdk:"zone_uuid"`
	ImageUuid             types.String `tfsdk:"image_uuid"`
	Type                  types.String `tfsdk:"type"`
}

// Configure implements resource.ResourceWithConfigure.
func (r *slbOfferingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func SlbOfferingResource() resource.Resource {
	return &slbOfferingResource{}
}

// Create implements resource.Resource.
func (r *slbOfferingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan slbOfferingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Configuring ZStack client")
	offerParam := param.CreateSlbOfferingParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateSlbOfferingParamDetail{
			Name:                  plan.Name.ValueString(),
			CpuNum:                int(plan.CpuNum.ValueInt64()),
			MemorySize:            utils.MBToBytes(plan.MemorySize.ValueInt64()),
			ManagementNetworkUuid: plan.ManagementNetworkUuid.ValueString(),
			ZoneUuid:              plan.ZoneUuid.ValueString(),
			ImageUuid:             plan.ImageUuid.ValueString(),
			Type:                  stringPtr("SLB"),
		},
	}

	// BUG-055: Description is Optional+Computed with UseStateForUnknown — guard Unknown
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		offerParam.Params.Description = stringPtr(plan.Description.ValueString())
	}

	offering, err := r.client.CreateSlbOffering(offerParam)
	tflog.Debug(ctx, "Received SLB offering", map[string]interface{}{
		"offering": offering,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating SLB Offering",
			"Could not create SLB offering, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(offering.UUID)
	plan.Name = types.StringValue(offering.Name)
	plan.Description = types.StringValue(offering.Description)
	plan.CpuNum = types.Int64Value(int64(offering.CpuNum))
	plan.MemorySize = types.Int64Value(utils.BytesToMB(offering.MemorySize))
	plan.Type = types.StringValue(offering.Type)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete implements resource.Resource.
func (r *slbOfferingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state slbOfferingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteInstanceOffering(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting SLB Offering",
			"Could not delete SLB offering, unexpected error: "+err.Error(),
		)
		return
	}
}

// Metadata implements resource.Resource.
func (r *slbOfferingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_slb_offering"
}

// Read implements resource.Resource.
func (r *slbOfferingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state slbOfferingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	offering, err := findResourceByGet(r.client.GetSlbOffering, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading SLB Offering",
			"Could not read SLB offering UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	if offering.Type == "" {
		state.Type = types.StringValue("SLB")
	} else {
		state.Type = types.StringValue(offering.Type)
	}

	state.ZoneUuid = types.StringValue(offering.ZoneUuid)
	state.ManagementNetworkUuid = types.StringValue(offering.ManagementNetworkUuid)
	state.ImageUuid = types.StringValue(offering.ImageUuid)
	state.Uuid = types.StringValue(offering.UUID)
	state.Description = types.StringValue(offering.Description)
	state.Name = types.StringValue(offering.Name)
	state.CpuNum = types.Int64Value(int64(offering.CpuNum))
	state.MemorySize = types.Int64Value(utils.BytesToMB(offering.MemorySize))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements resource.Resource.
func (r *slbOfferingResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource allows you to manage SLB offerings in ZStack. " +
			"A SLB offering defines the configuration and resource settings for dedicated load balancer instances, such as CPU, memory, and management network. " +
			"You can define the offering's properties, such as its name, description, CPU and memory allocation, and the associated management network.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier (UUID) of the SLB offering. Automatically generated by ZStack.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the SLB offering. This is a mandatory field.",
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
				Description: "A description of the SLB offering, providing additional context or details about the configuration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cpu_num": schema.Int64Attribute{
				Required:    true,
				Description: "The number of CPUs allocated to the SLB offering. This is a mandatory field.",
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"memory_size": schema.Int64Attribute{
				Required:    true,
				Description: "The amount of memory allocated to the SLB offering. This is a mandatory field, in mebibytes (MiB)",
				Validators:  []validator.Int64{int64validator.AtLeast(1)},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"management_network_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the management network associated with the SLB offering. This is a mandatory field.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"zone_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the zone where the SLB offering is deployed. This is a mandatory field.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"image_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the image used by the SLB offering. This is a mandatory field.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of the SLB offering. Defaults to 'SLB' if not specified.",
				Validators: []validator.String{
					stringvalidator.OneOf("SLB"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *slbOfferingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"SLB Offering resource does not support updates. Please recreate the resource instead.",
	)
}

func (r *slbOfferingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
