// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"terraform-provider-zstack/zstack/utils"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
)

var (
	_ resource.Resource              = &instanceOfferingResource{}
	_ resource.ResourceWithConfigure = &instanceOfferingResource{}
)

type instanceOfferingResource struct {
	client *client.ZSClient
}

type instanceOfferingResourceModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CpuNum      types.Int64  `tfsdk:"cpu_num"`
	MemorySize  types.Int64  `tfsdk:"memory_size"`
	Type        types.String `tfsdk:"type"`
}

// Configure implements resource.ResourceWithConfigure.
func (r *instanceOfferingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func InstanceOfferingResource() resource.Resource {
	return &instanceOfferingResource{}
}

// Create implements resource.Resource.
func (r *instanceOfferingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan instanceOfferingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	offerType := "UserVm"
	tflog.Info(ctx, "Configuring ZStack client")
	offerParam := param.CreateInstanceOfferingParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateInstanceOfferingDetailParam{
			Name:        plan.Name.ValueString(),
			Description: plan.Description.ValueStringPointer(),
			CpuNum:      int(plan.CpuNum.ValueInt64()),
			MemorySize:  utils.MBToBytes(plan.MemorySize.ValueInt64()),
			Type:        &offerType,
		},
	}

	instance_offer, err := r.client.CreateInstanceOffering(&offerParam)

	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Add instance offering to ZStack", "Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(instance_offer.UUID)
	plan.Name = types.StringValue(instance_offer.Name)
	plan.Description = types.StringValue(instance_offer.Description)
	plan.CpuNum = types.Int64Value(int64(instance_offer.CpuNum))
	plan.MemorySize = types.Int64Value(utils.BytesToMB(instance_offer.MemorySize))
	plan.Type = types.StringValue(instance_offer.Type)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete implements resource.Resource.
func (r *instanceOfferingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state instanceOfferingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "instance offering uuid is empty, so nothing to delete, skip it")
		return
	}

	err := r.client.DeleteInstanceOffering(state.Uuid.ValueString(), param.DeleteModeEnforcing)

	if err != nil {
		resp.Diagnostics.AddError("fail to delete instance offering ", ""+err.Error())
		return
	}
}

// Metadata implements resource.Resource.
func (r *instanceOfferingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance_offer"
}

// Read implements resource.Resource.
func (r *instanceOfferingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state instanceOfferingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	instance_offer, err := r.client.GetInstanceOffering(state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting virtual router image uuid", "Could not read instance offer uuid: "+err.Error(),
		)
		return
	}

	state.Type = types.StringValue(instance_offer.Type)
	if instance_offer.Type == "" {
		state.Type = types.StringValue("UserVm")
	} else {
		state.Type = types.StringValue(instance_offer.Type)
	}

	state.Uuid = types.StringValue(instance_offer.UUID)
	state.Description = types.StringValue(instance_offer.Description)
	state.Name = types.StringValue(instance_offer.Name)
	state.CpuNum = types.Int64Value(int64(instance_offer.CpuNum))
	state.MemorySize = types.Int64Value(utils.BytesToMB(instance_offer.MemorySize))
	state.Type = types.StringValue(instance_offer.Type)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements resource.Resource.
func (r *instanceOfferingResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource allows you to manage instance offerings in ZStack. " +
			"An instance offering defines the configuration and resource settings for virtual machine instances, such as CPU and memory. " +
			"You can define the offering's properties, such as its name, description, CPU and memory allocation.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the instance offering. Automatically generated by ZStack.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the instance offering. This is a mandatory field.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description of the instance offering, providing additional context or details about the configuration.",
			},
			"cpu_num": schema.Int64Attribute{
				Required:    true,
				Description: "The number of CPUs allocated to the instance offering. This is a mandatory field.",
			},
			"memory_size": schema.Int64Attribute{
				Required:    true,
				Description: "The amount of memory (in megabytes, MB) allocated to the instance offering. This is a mandatory field.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of the instance offering. Defaults to 'UserVm' if not specified.",
			},
		},
	}
}

func (r *instanceOfferingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}
