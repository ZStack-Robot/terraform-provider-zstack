// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VPC HA Group.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the VPC HA Group.",
			},
			"monitor_ips": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "List of IP addresses to monitor for HA.",
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
			"Fail to create VPC HA Group",
			"Error "+err.Error(),
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

	queryParam := param.NewQueryParam()
	vpcHaGroups, err := r.client.QueryVpcHaGroup(&queryParam)

	if err != nil {
		tflog.Warn(ctx, "Unable to query VPC HA Groups. It may have been deleted.: "+err.Error())
		state = vpcHaGroupModel{
			Uuid: types.StringValue(""),
		}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false

	for _, vpcHaGroup := range vpcHaGroups {
		if vpcHaGroup.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(vpcHaGroup.UUID)
			state.Name = types.StringValue(vpcHaGroup.Name)
			state.Description = stringValueOrNull(vpcHaGroup.Description)
			monitorIps := make([]string, len(vpcHaGroup.Monitors))
			for i, monitor := range vpcHaGroup.Monitors {
				monitorIps[i] = monitor.MonitorIp
			}
			state.MonitorIps = stringSliceToList(monitorIps)
			found = true
			break
		}
	}
	if !found {
		tflog.Warn(ctx, "VPC HA Group not found. It might have been deleted outside of Terraform.")
		state = vpcHaGroupModel{
			Uuid: types.StringValue(""),
		}
	}

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
			"Fail to update VPC HA Group",
			"Error "+err.Error(),
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

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "VPC HA Group UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteVpcHaGroup(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError("fail to delete VPC HA Group", ""+err.Error())
		return
	}
}

func (r *vpcHaGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
