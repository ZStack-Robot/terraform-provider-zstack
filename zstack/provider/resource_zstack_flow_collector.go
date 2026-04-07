// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &flowCollectorResource{}
	_ resource.ResourceWithConfigure   = &flowCollectorResource{}
	_ resource.ResourceWithImportState = &flowCollectorResource{}
)

type flowCollectorResource struct {
	client *client.ZSClient
}

type flowCollectorModel struct {
	Uuid          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	FlowMeterUuid types.String `tfsdk:"flow_meter_uuid"`
	Server        types.String `tfsdk:"server"`
	Port          types.Int64  `tfsdk:"port"`
}

func FlowCollectorResource() resource.Resource {
	return &flowCollectorResource{}
}

func (r *flowCollectorResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *flowCollectorResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_flow_collector"
}

func (r *flowCollectorResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Manages a flow collector in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the flow collector.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the flow collector.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the flow collector.",
			},
			"flow_meter_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the flow meter.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The server of the flow collector.",
			},
			"port": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The port of the flow collector.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *flowCollectorResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan flowCollectorModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateFlowCollectorParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateFlowCollectorParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			Server:      stringPtrOrNil(plan.Server.ValueString()),
		},
	}

	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		p.Params.Port = int64Ptr(plan.Port.ValueInt64())
	}

	flowCollector, err := r.client.CreateFlowCollector(plan.FlowMeterUuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError("Fail to create flow collector", "Error "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(flowCollector.UUID)
	plan.Name = stringValueOrNull(flowCollector.Name)
	plan.Description = stringValueOrNull(flowCollector.Description)
	plan.FlowMeterUuid = types.StringValue(flowCollector.FlowMeterUuid)
	plan.Server = stringValueOrNull(flowCollector.Server)
	plan.Port = types.Int64Value(flowCollector.Port)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *flowCollectorResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state flowCollectorModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	flowCollectors, err := r.client.QueryFlowCollector(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query flow collectors. It may have been deleted.: "+err.Error())
		state = flowCollectorModel{Uuid: types.StringValue("")}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false
	for _, flowCollector := range flowCollectors {
		if flowCollector.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(flowCollector.UUID)
			state.Name = stringValueOrNull(flowCollector.Name)
			state.Description = stringValueOrNull(flowCollector.Description)
			state.FlowMeterUuid = types.StringValue(flowCollector.FlowMeterUuid)
			state.Server = stringValueOrNull(flowCollector.Server)
			state.Port = types.Int64Value(flowCollector.Port)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Flow collector not found. It might have been deleted outside of Terraform.")
		state = flowCollectorModel{Uuid: types.StringValue("")}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *flowCollectorResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan flowCollectorModel
	var state flowCollectorModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateFlowCollectorParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateFlowCollectorParamDetail{
			Server: stringPtrOrNil(plan.Server.ValueString()),
		},
	}

	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		p.Params.Port = int64Ptr(plan.Port.ValueInt64())
	}

	flowCollector, err := r.client.UpdateFlowCollector(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError("Fail to update flow collector", "Error "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(flowCollector.UUID)
	plan.Name = stringValueOrNull(flowCollector.Name)
	plan.Description = stringValueOrNull(flowCollector.Description)
	plan.FlowMeterUuid = types.StringValue(flowCollector.FlowMeterUuid)
	plan.Server = stringValueOrNull(flowCollector.Server)
	plan.Port = types.Int64Value(flowCollector.Port)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *flowCollectorResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state flowCollectorModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Flow collector UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteFlowCollector(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("fail to delete flow collector", ""+err.Error())
		return
	}
}

func (r *flowCollectorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
