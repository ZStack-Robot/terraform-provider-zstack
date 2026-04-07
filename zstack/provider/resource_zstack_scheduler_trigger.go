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
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &schedulerTriggerResource{}
	_ resource.ResourceWithConfigure   = &schedulerTriggerResource{}
	_ resource.ResourceWithImportState = &schedulerTriggerResource{}
)

type schedulerTriggerResource struct {
	client *client.ZSClient
}

type schedulerTriggerResourceModel struct {
	Uuid              types.String `tfsdk:"uuid"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	SchedulerType     types.String `tfsdk:"scheduler_type"`
	SchedulerInterval types.Int64  `tfsdk:"scheduler_interval"`
	RepeatCount       types.Int64  `tfsdk:"repeat_count"`
	StartTime         types.Int64  `tfsdk:"start_time"`
	Cron              types.String `tfsdk:"cron"`
	StopTime          types.Int64  `tfsdk:"stop_time"`
}

func SchedulerTriggerResource() resource.Resource {
	return &schedulerTriggerResource{}
}

func (r *schedulerTriggerResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*client.ZSClient)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", request.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *schedulerTriggerResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_scheduler_trigger"
}

func (r *schedulerTriggerResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Manages a scheduler trigger in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the scheduler trigger.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the scheduler trigger.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the scheduler trigger.",
			},
			"scheduler_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the scheduler trigger (e.g. 'simple' or 'cron').",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"scheduler_interval": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The scheduler interval in seconds.",
			},
			"repeat_count": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The number of times to repeat the trigger.",
			},
			"start_time": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The start time of the trigger as Unix timestamp.",
			},
			"cron": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The cron expression for cron-type triggers.",
			},
			"stop_time": schema.Int64Attribute{
				Computed:    true,
				Description: "The stop time of the trigger as Unix timestamp.",
			},
		},
	}
}

func (r *schedulerTriggerResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan schedulerTriggerResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var startTime *int64
	if !plan.StartTime.IsNull() && !plan.StartTime.IsUnknown() {
		val := plan.StartTime.ValueInt64()
		startTime = &val
	}

	params := param.CreateSchedulerTriggerParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateSchedulerTriggerParamDetail{
			Name:              plan.Name.ValueString(),
			Description:       stringPtrOrNil(plan.Description.ValueString()),
			SchedulerInterval: intPtrFromInt64OrNil(plan.SchedulerInterval),
			RepeatCount:       intPtrFromInt64OrNil(plan.RepeatCount),
			StartTime:         startTime,
			SchedulerType:     plan.SchedulerType.ValueString(),
			Cron:              stringPtrOrNil(plan.Cron.ValueString()),
		},
	}

	trigger, err := r.client.CreateSchedulerTrigger(params)
	if err != nil {
		response.Diagnostics.AddError("Error creating scheduler trigger", err.Error())
		return
	}

	plan.Uuid = types.StringValue(trigger.UUID)
	plan.Name = types.StringValue(trigger.Name)
	plan.Description = stringValueOrNull(trigger.Description)
	plan.SchedulerType = types.StringValue(trigger.SchedulerType)
	plan.SchedulerInterval = types.Int64Value(int64(trigger.SchedulerInterval))
	plan.RepeatCount = types.Int64Value(int64(trigger.RepeatCount))
	if !trigger.StartTime.IsZero() {
		plan.StartTime = types.Int64Value(trigger.StartTime.Unix())
	} else {
		plan.StartTime = types.Int64Null()
	}
	plan.Cron = stringValueOrNull(trigger.Cron)
	if !trigger.StopTime.IsZero() {
		plan.StopTime = types.Int64Value(trigger.StopTime.Unix())
	} else {
		plan.StopTime = types.Int64Null()
	}

	diags = response.State.Set(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *schedulerTriggerResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state schedulerTriggerResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	queryParam.AddQ("uuid=" + state.Uuid.ValueString())

	triggers, err := r.client.QuerySchedulerTrigger(&queryParam)
	if err != nil {
		response.Diagnostics.AddError("Error reading scheduler trigger", err.Error())
		return
	}

	if len(triggers) == 0 {
		response.State.RemoveResource(ctx)
		return
	}

	trigger := triggers[0]
	state.Name = types.StringValue(trigger.Name)
	state.Description = stringValueOrNull(trigger.Description)
	state.SchedulerType = types.StringValue(trigger.SchedulerType)
	state.SchedulerInterval = types.Int64Value(int64(trigger.SchedulerInterval))
	state.RepeatCount = types.Int64Value(int64(trigger.RepeatCount))
	if !trigger.StartTime.IsZero() {
		state.StartTime = types.Int64Value(trigger.StartTime.Unix())
	} else {
		state.StartTime = types.Int64Null()
	}
	state.Cron = stringValueOrNull(trigger.Cron)
	if !trigger.StopTime.IsZero() {
		state.StopTime = types.Int64Value(trigger.StopTime.Unix())
	} else {
		state.StopTime = types.Int64Null()
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *schedulerTriggerResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state schedulerTriggerResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var plan schedulerTriggerResourceModel
	diags = request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	params := param.UpdateSchedulerTriggerParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateSchedulerTriggerParamDetail{
			Name:              plan.Name.ValueString(),
			Description:       stringPtrOrNil(plan.Description.ValueString()),
			SchedulerInterval: intPtrFromInt64OrNil(plan.SchedulerInterval),
			RepeatCount:       intPtrFromInt64OrNil(plan.RepeatCount),
			Cron:              stringPtrOrNil(plan.Cron.ValueString()),
		},
	}

	trigger, err := r.client.UpdateSchedulerTrigger(state.Uuid.ValueString(), params)
	if err != nil {
		response.Diagnostics.AddError("Error updating scheduler trigger", err.Error())
		return
	}

	plan.Uuid = types.StringValue(trigger.UUID)
	plan.Name = types.StringValue(trigger.Name)
	plan.Description = stringValueOrNull(trigger.Description)
	plan.SchedulerType = types.StringValue(trigger.SchedulerType)
	plan.SchedulerInterval = types.Int64Value(int64(trigger.SchedulerInterval))
	plan.RepeatCount = types.Int64Value(int64(trigger.RepeatCount))
	if !trigger.StartTime.IsZero() {
		plan.StartTime = types.Int64Value(trigger.StartTime.Unix())
	} else {
		plan.StartTime = types.Int64Null()
	}
	plan.Cron = stringValueOrNull(trigger.Cron)
	if !trigger.StopTime.IsZero() {
		plan.StopTime = types.Int64Value(trigger.StopTime.Unix())
	} else {
		plan.StopTime = types.Int64Null()
	}

	diags = response.State.Set(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *schedulerTriggerResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state schedulerTriggerResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSchedulerTrigger(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("Error deleting scheduler trigger", err.Error())
		return
	}
}

func (r *schedulerTriggerResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), request, response)
}
