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
)

var (
	_ resource.Resource                = &vipResource{}
	_ resource.ResourceWithConfigure   = &vipResource{}
	_ resource.ResourceWithImportState = &vipResource{}
)

type vipResource struct {
	client *client.ZSClient
}

type vipModel struct {
	Uuid          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	L3NetworkUuid types.String `tfsdk:"l3_network_uuid"`
	IpRangeUuid   types.String `tfsdk:"ip_range_uuid"`
	VIP           types.String `tfsdk:"vip"`
}

func VipResource() resource.Resource {
	return &vipResource{}
}

func (r *vipResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*client.ZSClient)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", request.ProviderData),
		)

		return
	}
	r.client = client
}

func (r *vipResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vip"
}

func (r *vipResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage Virtual IP (VIP) addresses in ZStack. " +
			"A VIP is a dedicated IP address that can be used for various network services, such as load balancing or high availability. " +
			"You can define the VIP's properties, such as its name, description, associated L3 network, and the IP range it belongs to.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the VIP network service.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VIP network service.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the VIP network service.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"l3_network_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the L3 network associated with this VIP network service.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ip_range_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The type of IP range. Possible values depend on the ZStack configuration (e.g., 'Normal' or 'Reserved').",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vip": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "create vip ip address  for this VPC network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *vipResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan vipModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if plan.Description.IsNull() {
		plan.Description = types.StringValue("")
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateVipParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVipParamDetail{
			Name:          plan.Name.ValueString(),
			Description:   stringPtr(plan.Description.ValueString()),
			L3NetworkUuid: plan.L3NetworkUuid.ValueString(),
		},
	}

	if !plan.IpRangeUuid.IsNull() && plan.IpRangeUuid.ValueString() != "" {
		p.Params.IpRangeUuid = stringPtr(plan.IpRangeUuid.ValueString())
	}

	if !plan.VIP.IsNull() && plan.VIP.ValueString() != "" {
		p.Params.RequiredIp = stringPtr(plan.VIP.ValueString())
	}

	vip, err := r.client.CreateVip(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating VIP",
			"Could not create vip, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(vip.UUID)
	plan.Name = types.StringValue(vip.Name)
	plan.Description = types.StringValue(vip.Description)
	plan.L3NetworkUuid = types.StringValue(vip.L3NetworkUuid)
	plan.VIP = types.StringValue(vip.Ip)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *vipResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state vipModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	vip, err := findResourceByQuery(r.client.QueryVip, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		tflog.Warn(ctx, "Unable to query VIPs. It may have been deleted.: "+err.Error())
		state = vipModel{
			Uuid: types.StringValue(""),
		}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	state.Uuid = types.StringValue(vip.UUID)
	state.Name = types.StringValue(vip.Name)
	state.Description = types.StringValue(vip.Description)
	state.L3NetworkUuid = types.StringValue(vip.L3NetworkUuid)
	state.VIP = types.StringValue(vip.Ip)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *vipResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	response.Diagnostics.AddError(
		"Update not supported",
		"VIP resource does not support updates. Please recreate the resource instead.",
	)
}

func (r *vipResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state vipModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "reserved ip UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteVip(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting VIP",
			"Could not delete vip, unexpected error: "+err.Error(),
		)
		return
	}

}

func (r *vipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
