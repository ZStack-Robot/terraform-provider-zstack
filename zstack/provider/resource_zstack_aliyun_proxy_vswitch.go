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

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &aliyunProxyVSwitchResource{}
	_ resource.ResourceWithConfigure   = &aliyunProxyVSwitchResource{}
	_ resource.ResourceWithImportState = &aliyunProxyVSwitchResource{}
)

// AliyunProxyVSwitchResource is a helper function to simplify the provider implementation.
func AliyunProxyVSwitchResource() resource.Resource {
	return &aliyunProxyVSwitchResource{}
}

// aliyunProxyVSwitchResource is the resource implementation.
type aliyunProxyVSwitchResource struct {
	client *client.ZSClient
}

// aliyunProxyVSwitchModel describes the resource data model.
type aliyunProxyVSwitchModel struct {
	UUID               types.String `tfsdk:"uuid"`
	Name               types.String `tfsdk:"name"`
	AliyunProxyVpcUuid types.String `tfsdk:"aliyun_proxy_vpc_uuid"`
	VpcL3NetworkUuid   types.String `tfsdk:"vpc_l3_network_uuid"`
	IsDefault          types.Bool   `tfsdk:"is_default"`
	Status             types.String `tfsdk:"status"`
}

// Metadata returns the resource type name.
func (r *aliyunProxyVSwitchResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aliyun_proxy_vswitch"
}

// Schema defines the schema for the resource.
func (r *aliyunProxyVSwitchResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Aliyun Proxy VSwitch resource",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The UUID of the Aliyun Proxy VSwitch",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the Aliyun Proxy VSwitch",
			},
			"aliyun_proxy_vpc_uuid": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The UUID of the Aliyun Proxy VPC",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpc_l3_network_uuid": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The UUID of the VPC L3 Network",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_default": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether this is the default VSwitch",
			},
			"status": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The status of the Aliyun Proxy VSwitch",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *aliyunProxyVSwitchResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource and sets the initial Terraform state.
func (r *aliyunProxyVSwitchResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aliyunProxyVSwitchModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the resource
	createParam := param.CreateAliyunProxyVSwitchParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateAliyunProxyVSwitchParamDetail{
			AliyunProxyVpcUuid: plan.AliyunProxyVpcUuid.ValueString(),
			VpcL3NetworkUuid:   plan.VpcL3NetworkUuid.ValueString(),
			IsDefault:          plan.IsDefault.ValueBool(),
		},
	}

	view, err := r.client.CreateAliyunProxyVSwitch(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Aliyun Proxy VSwitch",
			"Could not create Aliyun Proxy VSwitch, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response to model
	plan.UUID = types.StringValue(view.UUID)
	plan.Name = types.StringValue(view.Name)
	plan.Status = types.StringValue(view.Status)
	plan.IsDefault = types.BoolValue(view.IsDefault)

	tflog.Trace(ctx, "Created Aliyun Proxy VSwitch", map[string]any{"uuid": view.UUID})

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *aliyunProxyVSwitchResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state aliyunProxyVSwitchModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Query the resource
	view, err := r.client.GetAliyunProxyVSwitch(state.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Aliyun Proxy VSwitch",
			"Could not read Aliyun Proxy VSwitch, unexpected error: "+err.Error(),
		)
		return
	}

	if view == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Map response to model
	state.UUID = types.StringValue(view.UUID)
	state.Name = types.StringValue(view.Name)
	state.AliyunProxyVpcUuid = types.StringValue(view.AliyunProxyVpcUuid)
	state.VpcL3NetworkUuid = types.StringValue(view.VpcL3NetworkUuid)
	state.IsDefault = types.BoolValue(view.IsDefault)
	state.Status = types.StringValue(view.Status)

	tflog.Trace(ctx, "Read Aliyun Proxy VSwitch", map[string]any{"uuid": view.UUID})

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *aliyunProxyVSwitchResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan aliyunProxyVSwitchModel
	var state aliyunProxyVSwitchModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare update parameters
	updateParamDetail := param.UpdateAliyunProxyVSwitchParamDetail{
		IsDefault: boolPtr(plan.IsDefault.ValueBool()),
	}

	if !plan.Status.IsNull() {
		updateParamDetail.Status = stringPtr(plan.Status.ValueString())
	}

	updateParam := param.UpdateAliyunProxyVSwitchParam{
		BaseParam: param.BaseParam{},
		Params:    updateParamDetail,
	}

	view, err := r.client.UpdateAliyunProxyVSwitch(state.UUID.ValueString(), updateParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Aliyun Proxy VSwitch",
			"Could not update Aliyun Proxy VSwitch, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response to model
	plan.UUID = types.StringValue(view.UUID)
	plan.Name = types.StringValue(view.Name)
	plan.AliyunProxyVpcUuid = types.StringValue(view.AliyunProxyVpcUuid)
	plan.VpcL3NetworkUuid = types.StringValue(view.VpcL3NetworkUuid)
	plan.Status = types.StringValue(view.Status)
	plan.IsDefault = types.BoolValue(view.IsDefault)

	tflog.Trace(ctx, "Updated Aliyun Proxy VSwitch", map[string]any{"uuid": view.UUID})

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *aliyunProxyVSwitchResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state aliyunProxyVSwitchModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the resource
	err := r.client.DeleteAliyunProxyVSwitch(state.UUID.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Aliyun Proxy VSwitch",
			"Could not delete Aliyun Proxy VSwitch, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "Deleted Aliyun Proxy VSwitch", map[string]any{"uuid": state.UUID.ValueString()})
}

// ImportState imports a resource into Terraform state using UUID.
func (r *aliyunProxyVSwitchResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the UUID as the import identifier
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
