// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
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
	_ resource.Resource                = &flowMeterResource{}
	_ resource.ResourceWithConfigure   = &flowMeterResource{}
	_ resource.ResourceWithImportState = &flowMeterResource{}
)

type flowMeterResource struct {
	client *client.ZSClient
}

type flowMeterModel struct {
	Uuid             types.String `tfsdk:"uuid"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Type             types.String `tfsdk:"type"`
	Version          types.String `tfsdk:"version"`
	Sample           types.Int64  `tfsdk:"sample"`
	ExpireInterval   types.Int64  `tfsdk:"expire_interval"`
	Server           types.String `tfsdk:"server"`
	Port             types.Int64  `tfsdk:"port"`
	GenerateInterval types.Int64  `tfsdk:"generate_interval"`
}

func FlowMeterResource() resource.Resource {
	return &flowMeterResource{}
}

func (r *flowMeterResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *flowMeterResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_flow_meter"
}

func (r *flowMeterResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Manages a flow meter in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the flow meter.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the flow meter.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the flow meter.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the flow meter.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The version of the flow meter.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sample": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The sample interval.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"expire_interval": schema.Int64Attribute{
				Computed:    true,
				Description: "The expire interval.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"server": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The server of the flow meter.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"port": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The port of the flow meter.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"generate_interval": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The generate interval.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *flowMeterResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan flowMeterModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateFlowMeterParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateFlowMeterParamDetail{
			Version:     stringPtrOrNil(plan.Version.ValueString()),
			Type:        plan.Type.ValueString(),
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			Server:      stringPtrOrNil(plan.Server.ValueString()),
		},
	}

	if !plan.Sample.IsNull() && !plan.Sample.IsUnknown() {
		p.Params.Sample = intPtr(int(plan.Sample.ValueInt64()))
	}
	if !plan.GenerateInterval.IsNull() && !plan.GenerateInterval.IsUnknown() {
		p.Params.GenerateInterval = intPtr(int(plan.GenerateInterval.ValueInt64()))
	}
	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		p.Params.Port = int64Ptr(plan.Port.ValueInt64())
	}

	flowMeter, err := r.client.CreateFlowMeter(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Flow Meter",
			"Could not create flow meter, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(flowMeter.UUID)
	plan.Name = stringValueOrNull(flowMeter.Name)
	plan.Description = stringValueOrNull(flowMeter.Description)
	plan.Type = types.StringValue(flowMeter.Type)
	plan.Version = stringValueOrNull(flowMeter.Version)
	plan.Sample = types.Int64Value(flowMeter.Sample)
	plan.ExpireInterval = types.Int64Value(flowMeter.ExpireInterval)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *flowMeterResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state flowMeterModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	flowMeter, err := findResourceByQuery(r.client.QueryFlowMeter, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		tflog.Warn(ctx, "Unable to query flow meters. It may have been deleted.: "+err.Error())
		state = flowMeterModel{Uuid: types.StringValue("")}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	state.Uuid = types.StringValue(flowMeter.UUID)
	state.Name = stringValueOrNull(flowMeter.Name)
	state.Description = stringValueOrNull(flowMeter.Description)
	state.Type = types.StringValue(flowMeter.Type)
	state.Version = stringValueOrNull(flowMeter.Version)
	state.Sample = types.Int64Value(flowMeter.Sample)
	state.ExpireInterval = types.Int64Value(flowMeter.ExpireInterval)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *flowMeterResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan flowMeterModel
	var state flowMeterModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateFlowMeterParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateFlowMeterParamDetail{
			Version:     stringPtrOrNil(plan.Version.ValueString()),
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	if !plan.Sample.IsNull() && !plan.Sample.IsUnknown() {
		p.Params.Sample = int64Ptr(plan.Sample.ValueInt64())
	}
	if !plan.ExpireInterval.IsNull() && !plan.ExpireInterval.IsUnknown() {
		p.Params.ExpireInterval = int64Ptr(plan.ExpireInterval.ValueInt64())
	}

	flowMeter, err := r.client.UpdateFlowMeter(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating Flow Meter",
			"Could not update flow meter, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(flowMeter.UUID)
	plan.Name = stringValueOrNull(flowMeter.Name)
	plan.Description = stringValueOrNull(flowMeter.Description)
	plan.Type = types.StringValue(flowMeter.Type)
	plan.Version = stringValueOrNull(flowMeter.Version)
	plan.Sample = types.Int64Value(flowMeter.Sample)
	plan.ExpireInterval = types.Int64Value(flowMeter.ExpireInterval)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *flowMeterResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state flowMeterModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Flow meter UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteFlowMeter(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting Flow Meter",
			"Could not delete flow meter, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *flowMeterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
