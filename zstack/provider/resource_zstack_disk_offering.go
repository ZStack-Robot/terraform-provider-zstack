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
	_ resource.Resource              = &diskOfferingResource{}
	_ resource.ResourceWithConfigure = &diskOfferingResource{}
)

type diskOfferingResource struct {
	client *client.ZSClient
}

type diskOfferingResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Uuid        types.String `tfsdk:"uuid"`
	Description types.String `tfsdk:"description"`
	DiskSize    types.Int64  `tfsdk:"disk_size"`
	//	Type              types.String `tfsdk:"type"`               // Type
	AllocatorStrategy types.String `tfsdk:"allocator_strategy"` // Allocation strategy
	// SortKey           types.Int32  `tfsdk:"sort_key"`
	// State             types.String `tfsdk:"state"` // State (Enabled, Disabled)
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

	tflog.Info(ctx, "Configuring ZStack client")
	offerParam := param.CreateDiskOfferingParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateDiskOfferingDetailParam{
			Name:        plan.Name.ValueString(),
			Description: plan.Description.ValueStringPointer(),
			DiskSize:    plan.DiskSize.ValueInt64(),
			//	Type:              plan.Type.ValueStringPointer(),
			AllocatorStrategy: plan.AllocatorStrategy.ValueStringPointer(),
		},
	}

	disk_offer, err := r.client.CreateDiskOffering(&offerParam)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Add disk offering to ZStack", "Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(disk_offer.UUID)
	plan.Name = types.StringValue(disk_offer.Name)
	plan.Description = types.StringValue(disk_offer.Description)
	plan.DiskSize = types.Int64Value(int64(disk_offer.DiskSize))
	//plan.Type = types.StringValue(disk_offer.Type)
	plan.AllocatorStrategy = types.StringValue(disk_offer.AllocatorStrategy)
	//plan.State = types.StringValue(disk_offer.State)

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

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "disk offering uuid is empty, so nothing to delete, skip it")
		return
	}

	err := r.client.DeleteDiskOffering(state.Uuid.ValueString(), param.DeleteModeEnforcing)

	if err != nil {
		resp.Diagnostics.AddError("fail to delete disk offering ", ""+err.Error())
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
	disk_offer, err := r.client.GetDiskOffering(state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting disk offering uuid", "Could not read disk offer uuid: "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(disk_offer.UUID)
	state.Description = types.StringValue(disk_offer.Description)
	state.Name = types.StringValue(disk_offer.Name)
	state.DiskSize = types.Int64Value(int64(disk_offer.DiskSize))
	//state.State = types.StringValue(disk_offer.State)
	state.AllocatorStrategy = types.StringValue(disk_offer.AllocatorStrategy)
	//state.Type = types.StringValue(disk_offer.Type)

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
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the disk offering. This is a mandatory field.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description of the disk offering, providing additional context or details about the configuration.",
			},
			"disk_size": schema.Int64Attribute{
				Required:    true,
				Description: "The amount of disk size (in bytes) allocated to the disk offering. This is a mandatory field.",
			},
			"allocator_strategy": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of the allocator_strategy. ",
			},
		},
	}
}

func (r *diskOfferingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}
