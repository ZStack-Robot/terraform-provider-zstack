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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &vpcHaGroupResource{}
	_ resource.ResourceWithConfigure   = &vpcHaGroupResource{}
	_ resource.ResourceWithImportState = &vpcHaGroupResource{}
)

type vpcHaGroupResource struct {
	client *client.ZSClient
}

type vpcHaGroupModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	MonitorIps  types.List   `tfsdk:"monitor_ips"`
}

func VpcHaGroupResource() resource.Resource {
	return &vpcHaGroupResource{}
}

func (r *vpcHaGroupResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *vpcHaGroupResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpc_ha_group"
}

func (r *vpcHaGroupResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage VPC HA Groups in ZStack. " +
			"A VPC HA Group provides high availability for VPC virtual routers.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the VPC HA Group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VPC HA Group.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the VPC HA Group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"monitor_ips": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "List of IP addresses to monitor for HA.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
					listplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *vpcHaGroupResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan vpcHaGroupModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateVpcHaGroupParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVpcHaGroupParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			MonitorIps:  listToStringSlice(plan.MonitorIps),
		},
	}

	vpcHaGroup, err := r.client.CreateVpcHaGroup(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating VPC HA Group",
			"Could not create vpc ha group, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(vpcHaGroup.UUID)
	plan.Name = types.StringValue(vpcHaGroup.Name)
	plan.Description = stringValueOrNull(vpcHaGroup.Description)
	monitorIps := make([]string, len(vpcHaGroup.Monitors))
	for i, monitor := range vpcHaGroup.Monitors {
		monitorIps[i] = monitor.MonitorIp
	}
	plan.MonitorIps = stringSliceToList(monitorIps)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vpcHaGroupResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state vpcHaGroupModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	vpcHaGroup, err := findResourceByQuery(r.client.QueryVpcHaGroup, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading VPC HA Group",
			"Could not read VPC HA Group, unexpected error: "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(vpcHaGroup.UUID)
	state.Name = types.StringValue(vpcHaGroup.Name)
	state.Description = stringValueOrNull(vpcHaGroup.Description)
	monitorIps := make([]string, len(vpcHaGroup.Monitors))
	for i, monitor := range vpcHaGroup.Monitors {
		monitorIps[i] = monitor.MonitorIp
	}
	state.MonitorIps = stringSliceToList(monitorIps)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vpcHaGroupResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan vpcHaGroupModel
	var state vpcHaGroupModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateVpcHaGroupParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateVpcHaGroupParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	vpcHaGroup, err := r.client.UpdateVpcHaGroup(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating VPC HA Group",
			"Could not update vpc ha group, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(vpcHaGroup.UUID)
	plan.Name = types.StringValue(vpcHaGroup.Name)
	plan.Description = stringValueOrNull(vpcHaGroup.Description)
	monitorIps := make([]string, len(vpcHaGroup.Monitors))
	for i, monitor := range vpcHaGroup.Monitors {
		monitorIps[i] = monitor.MonitorIp
	}
	plan.MonitorIps = stringSliceToList(monitorIps)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vpcHaGroupResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state vpcHaGroupModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}


	err := r.client.DeleteVpcHaGroup(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting VPC HA Group",
			"Could not delete vpc ha group, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *vpcHaGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
