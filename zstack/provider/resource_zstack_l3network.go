// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &l3networkResource{}
	_ resource.ResourceWithConfigure   = &l3networkResource{}
	_ resource.ResourceWithImportState = &l3networkResource{}
)

type l3networkResource struct {
	client *client.ZSClient
}

type l3networkResourceModel struct {
	Uuid          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	L2NetworkUuid types.String `tfsdk:"l2_network_uuid"`
	Type          types.String `tfsdk:"type"`
	Category      types.String `tfsdk:"category"`
	System        types.Bool   `tfsdk:"system"`
	DnsDomain     types.String `tfsdk:"dns_domain"`
	IpVersion     types.Int64  `tfsdk:"ip_version"`
	State         types.String `tfsdk:"state"`
	ZoneUuid      types.String `tfsdk:"zone_uuid"`
}

func L3NetworkResource() resource.Resource {
	return &l3networkResource{}
}

func (r *l3networkResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *l3networkResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_l3network"
}

func (r *l3networkResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage Layer 3 (L3) networks in ZStack. " +
			"An L3 network is a virtual network that spans across zones and provides IP addressing services to virtual machines.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the L3 network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the L3 network.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the L3 network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"l2_network_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the L2 network that this L3 network is built on.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of the L3 network (e.g., L3BasicNetwork).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"category": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The category of the L3 network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"system": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether this is a system L3 network.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"dns_domain": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The DNS domain for the L3 network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip_version": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The IP version for the L3 network (4 for IPv4, 6 for IPv6).",
				Validators: []validator.Int64{
					int64validator.OneOf(4, 6),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the L3 network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the zone that contains this L3 network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *l3networkResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan l3networkResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	var ipVersion *int
	if !plan.IpVersion.IsNull() {
		ipVersion = intPtr(int(plan.IpVersion.ValueInt64()))
	}

	system := false
	if !plan.System.IsNull() {
		system = plan.System.ValueBool()
	}

	p := param.CreateL3NetworkParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateL3NetworkParamDetail{
			Name:          plan.Name.ValueString(),
			Description:   stringPtrOrNil(plan.Description.ValueString()),
			Type:          stringPtrOrNil(plan.Type.ValueString()),
			L2NetworkUuid: plan.L2NetworkUuid.ValueString(),
			Category:      stringPtrOrNil(plan.Category.ValueString()),
			IpVersion:     ipVersion,
			System:        system,
			DnsDomain:     stringPtrOrNil(plan.DnsDomain.ValueString()),
		},
	}

	result, err := r.client.CreateL3Network(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating L3 Network",
			"Could not create L3 network, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = types.StringValue(result.Description)
	plan.L2NetworkUuid = types.StringValue(result.L2NetworkUuid)
	plan.Type = types.StringValue(result.Type)
	plan.Category = types.StringValue(result.Category)
	plan.System = types.BoolValue(result.System)
	plan.DnsDomain = types.StringValue(result.DnsDomain)
	plan.State = types.StringValue(result.State)
	plan.ZoneUuid = types.StringValue(result.ZoneUuid)
	if result.IpVersion > 0 {
		plan.IpVersion = types.Int64Value(int64(result.IpVersion))
	}

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *l3networkResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state l3networkResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	item, err := findResourceByQuery(r.client.QueryL3Network, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading L3 Network",
			"Could not read L3 Network, unexpected error: "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(item.UUID)
	state.Name = types.StringValue(item.Name)
	state.Description = types.StringValue(item.Description)
	state.L2NetworkUuid = types.StringValue(item.L2NetworkUuid)
	state.Type = types.StringValue(item.Type)
	state.Category = types.StringValue(item.Category)
	state.System = types.BoolValue(item.System)
	state.DnsDomain = types.StringValue(item.DnsDomain)
	state.State = types.StringValue(item.State)
	state.ZoneUuid = types.StringValue(item.ZoneUuid)
	if item.IpVersion > 0 {
		state.IpVersion = types.Int64Value(int64(item.IpVersion))
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *l3networkResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan l3networkResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var state l3networkResourceModel
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	var system *bool
	if !plan.System.IsNull() {
		systemVal := plan.System.ValueBool()
		system = &systemVal
	}

	p := param.UpdateL3NetworkParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateL3NetworkParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			DnsDomain:   stringPtrOrNil(plan.DnsDomain.ValueString()),
			Category:    stringPtrOrNil(plan.Category.ValueString()),
			System:      system,
		},
	}

	result, err := r.client.UpdateL3Network(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating L3 Network",
			"Could not update L3 network, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = types.StringValue(result.Description)
	plan.L2NetworkUuid = types.StringValue(result.L2NetworkUuid)
	plan.Type = types.StringValue(result.Type)
	plan.Category = types.StringValue(result.Category)
	plan.System = types.BoolValue(result.System)
	plan.DnsDomain = types.StringValue(result.DnsDomain)
	plan.State = types.StringValue(result.State)
	plan.ZoneUuid = types.StringValue(result.ZoneUuid)
	if result.IpVersion > 0 {
		plan.IpVersion = types.Int64Value(int64(result.IpVersion))
	}

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *l3networkResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state l3networkResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}


	err := r.client.DeleteL3Network(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError("Error deleting L3 Network", "Could not delete L3 network, unexpected error: "+err.Error())
		return
	}

}

func (r *l3networkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
