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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &diskOfferingResource{}
	_ resource.ResourceWithConfigure   = &diskOfferingResource{}
	_ resource.ResourceWithImportState = &diskOfferingResource{}
)

type diskOfferingResource struct {
	client *client.ZSClient
}

type diskOfferingResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Uuid        types.String `tfsdk:"uuid"`
	Description types.String `tfsdk:"description"`
	DiskSize    types.Int64  `tfsdk:"disk_size"`
	//AllocatorStrategy types.String `tfsdk:"allocator_strategy"` // Allocation strategy
}

// Configure implements resource.ResourceWithConfigure.
func (r *diskOfferingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func DiskOfferingResource() resource.Resource {
	return &diskOfferingResource{}
}

// Create implements resource.Resource.
func (r *diskOfferingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan diskOfferingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diskSizeBytes := utils.GBToBytes(plan.DiskSize.ValueInt64())
	offerParam := param.CreateDiskOfferingParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateDiskOfferingParamDetail{
			Name:        plan.Name.ValueString(),
			Description: plan.Description.ValueStringPointer(),
			DiskSize:    diskSizeBytes,
			//AllocatorStrategy: plan.AllocatorStrategy.ValueStringPointer(),
		},
	}

	disk_offer, err := r.client.CreateDiskOffering(offerParam)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Disk Offering",
			"Could not create disk offering, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(disk_offer.UUID)
	plan.Name = types.StringValue(disk_offer.Name)
	plan.Description = types.StringValue(disk_offer.Description)
	plan.DiskSize = types.Int64Value(utils.BytesToGB(int64(disk_offer.DiskSize)))
	//	plan.AllocatorStrategy = types.StringValue(disk_offer.AllocatorStrategy)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete implements resource.Resource.
func (r *diskOfferingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state diskOfferingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}


	err := r.client.DeleteDiskOffering(state.Uuid.ValueString(), param.DeleteModeEnforcing)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Disk Offering",
			"Could not delete disk offering, unexpected error: "+err.Error(),
		)
		return
	}
}

// Metadata implements resource.Resource.
func (r *diskOfferingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_disk_offer"
}

// Read implements resource.Resource.
func (r *diskOfferingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state diskOfferingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	disk_offer, err := findResourceByGet(r.client.GetDiskOffering, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Disk Offering",
			"Could not read disk offering UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(disk_offer.UUID)
	state.Description = types.StringValue(disk_offer.Description)
	state.Name = types.StringValue(disk_offer.Name)
	state.DiskSize = types.Int64Value(utils.BytesToGB(int64(disk_offer.DiskSize)))
	//state.AllocatorStrategy = types.StringValue(disk_offer.AllocatorStrategy)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements resource.Resource.
func (r *diskOfferingResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource allows you to manage disk offerings in ZStack. " +
			"An disk offering defines the configuration and resource settings for virtual machine  disks. " +
			"You can define the offering's properties, such as its name, description, disk_size.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the disk offering. Automatically generated by ZStack.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the disk offering. This is a mandatory field.",
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
				Description: "A description of the disk offering, providing additional context or details about the configuration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"disk_size": schema.Int64Attribute{
				Required:    true,
				Description: "The amount of disk size allocated to the disk offering. This is a mandatory field, in gigabytes (GB).",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			/*
				"allocator_strategy": schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Description: "The type of the allocator_strategy. ",
				},
			*/
		},
	}
}

func (r *diskOfferingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"Disk Offering resource does not support updates. Please recreate the resource instead.",
	)
}

func (r *diskOfferingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
