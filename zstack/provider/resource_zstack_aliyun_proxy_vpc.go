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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

var (
	_ resource.Resource                = &aliyunProxyVpcResource{}
	_ resource.ResourceWithConfigure   = &aliyunProxyVpcResource{}
	_ resource.ResourceWithImportState = &aliyunProxyVpcResource{}
)

type aliyunProxyVpcResource struct {
	client *client.ZSClient
}

type aliyunProxyVpcModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CidrBlock   types.String `tfsdk:"cidr_block"`
	VRouterUuid types.String `tfsdk:"vrouter_uuid"`
	IsDefault   types.Bool   `tfsdk:"is_default"`
	VpcName     types.String `tfsdk:"vpc_name"`
	Status      types.String `tfsdk:"status"`
}

func AliyunProxyVpcResource() resource.Resource {
	return &aliyunProxyVpcResource{}
}

// Metadata implements resource.Resource.
func (r *aliyunProxyVpcResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aliyun_proxy_vpc"
}

// Schema implements resource.Resource.
func (r *aliyunProxyVpcResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage ZStack Aliyun Proxy VPC resources.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The UUID of the Aliyun Proxy VPC.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the Aliyun Proxy VPC.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the Aliyun Proxy VPC.",
			},
			"cidr_block": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The CIDR block of the Aliyun Proxy VPC.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vrouter_uuid": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The UUID of the Virtual Router.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_default": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether this is the default Aliyun Proxy VPC.",
			},
			"vpc_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the VPC.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the Aliyun Proxy VPC.",
			},
		},
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *aliyunProxyVpcResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create implements resource.Resource.
func (r *aliyunProxyVpcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aliyunProxyVpcModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating Aliyun Proxy VPC: %s", plan.Name.ValueString()))

	createParam := param.CreateAliyunProxyVpcParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateAliyunProxyVpcParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			CidrBlock:   plan.CidrBlock.ValueString(),
			VRouterUuid: plan.VRouterUuid.ValueString(),
			IsDefault:   plan.IsDefault.ValueBool(),
		},
	}

	result, err := r.client.CreateAliyunProxyVpc(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Aliyun Proxy VPC",
			"Could not create aliyun proxy vpc, unexpected error: "+err.Error(),
		)
		return
	}

	state := aliyunProxyVpcModelFromView(result)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (r *aliyunProxyVpcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state aliyunProxyVpcModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading Aliyun Proxy VPC: %s", state.Uuid.ValueString()))

	result, err := findResourceByGet(r.client.GetAliyunProxyVpc, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Aliyun Proxy VPC",
			"Could not read aliyun proxy vpc, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState := aliyunProxyVpcModelFromView(result)
	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *aliyunProxyVpcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan aliyunProxyVpcModel
	var state aliyunProxyVpcModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Updating Aliyun Proxy VPC: %s", uuid))

	updateParam := param.UpdateAliyunProxyVpcParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateAliyunProxyVpcParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			IsDefault:   boolPtr(plan.IsDefault.ValueBool()),
		},
	}

	result, err := r.client.UpdateAliyunProxyVpc(uuid, updateParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Aliyun Proxy VPC",
			"Could not update aliyun proxy vpc, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState := aliyunProxyVpcModelFromView(result)
	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *aliyunProxyVpcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state aliyunProxyVpcModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Deleting Aliyun Proxy VPC: %s", uuid))

	if err := r.client.DeleteAliyunProxyVpc(uuid, param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Aliyun Proxy VPC",
			"Could not delete aliyun proxy vpc, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *aliyunProxyVpcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func aliyunProxyVpcModelFromView(v *view.AliyunProxyVpcInventoryView) aliyunProxyVpcModel {
	return aliyunProxyVpcModel{
		Uuid:        types.StringValue(v.UUID),
		Name:        types.StringValue(v.Name),
		Description: stringValueOrNull(v.Description),
		CidrBlock:   types.StringValue(v.CidrBlock),
		VRouterUuid: types.StringValue(v.VRouterUuid),
		IsDefault:   types.BoolValue(v.IsDefault),
		VpcName:     types.StringValue(v.VpcName),
		Status:      types.StringValue(v.Status),
	}
}
