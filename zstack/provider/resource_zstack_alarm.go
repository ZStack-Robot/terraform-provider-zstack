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
	_ resource.Resource                = &alarmResource{}
	_ resource.ResourceWithConfigure   = &alarmResource{}
	_ resource.ResourceWithImportState = &alarmResource{}
)

type alarmResource struct {
	client *client.ZSClient
}

type alarmModel struct {
	Uuid               types.String  `tfsdk:"uuid"`
	Name               types.String  `tfsdk:"name"`
	Description        types.String  `tfsdk:"description"`
	ComparisonOperator types.String  `tfsdk:"comparison_operator"`
	Namespace          types.String  `tfsdk:"namespace"`
	MetricName         types.String  `tfsdk:"metric_name"`
	Threshold          types.Float64 `tfsdk:"threshold"`
	Period             types.Int64   `tfsdk:"period"`
	RepeatInterval     types.Int64   `tfsdk:"repeat_interval"`
	Status             types.String  `tfsdk:"status"`
	State              types.String  `tfsdk:"state"`
}

func AlarmResource() resource.Resource {
	return &alarmResource{}
}

func (r *alarmResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *alarmResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_alarm"
}

func (r *alarmResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage alarms in ZStack. " +
			"An alarm monitors metrics and triggers when thresholds are exceeded.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the alarm.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the alarm.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the alarm.",
			},
			"comparison_operator": schema.StringAttribute{
				Required:    true,
				Description: "The comparison operator (e.g., GreaterThanOrEqualTo, LessThan).",
				Validators: []validator.String{
					stringvalidator.OneOf("GreaterThanOrEqualTo", "GreaterThan", "LessThanOrEqualTo", "LessThan"),
				},
			},
			"namespace": schema.StringAttribute{
				Required:    true,
				Description: "The namespace of the metric.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"metric_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the metric to monitor.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"threshold": schema.Float64Attribute{
				Required:    true,
				Description: "The threshold value that triggers the alarm.",
			},
			"period": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The period in seconds over which the metric is evaluated.",
			},
			"repeat_interval": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The interval in seconds between repeated alarm notifications.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The current status of the alarm.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The current state of the alarm.",
			},
		},
	}
}

func (r *alarmResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan alarmModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateAlarmParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateAlarmParamDetail{
			Name:               plan.Name.ValueString(),
			ComparisonOperator: plan.ComparisonOperator.ValueString(),
			Namespace:          plan.Namespace.ValueString(),
			MetricName:         plan.MetricName.ValueString(),
			Threshold:          plan.Threshold.ValueFloat64(),
			Description:        stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	if !plan.Period.IsNull() && !plan.Period.IsUnknown() {
		val := int(plan.Period.ValueInt64())
		p.Params.Period = &val
	}

	if !plan.RepeatInterval.IsNull() && !plan.RepeatInterval.IsUnknown() {
		val := int(plan.RepeatInterval.ValueInt64())
		p.Params.RepeatInterval = &val
	}

	result, err := r.client.CreateAlarm(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Alarm",
			"Could not create alarm, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = stringValueOrNull(result.Description)
	plan.ComparisonOperator = types.StringValue(result.ComparisonOperator)
	plan.Namespace = types.StringValue(result.Namespace)
	plan.MetricName = types.StringValue(result.MetricName)
	plan.Threshold = types.Float64Value(result.Threshold)
	plan.Period = types.Int64Value(int64(result.Period))
	plan.RepeatInterval = types.Int64Value(int64(result.RepeatInterval))
	plan.Status = types.StringValue(result.Status)
	plan.State = types.StringValue(result.State)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *alarmResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state alarmModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	alarm, err := findResourceByQuery(r.client.QueryAlarm, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading Alarm",
			"Could not read alarm UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(alarm.UUID)
	state.Name = types.StringValue(alarm.Name)
	state.Description = stringValueOrNull(alarm.Description)
	state.ComparisonOperator = types.StringValue(alarm.ComparisonOperator)
	state.Namespace = types.StringValue(alarm.Namespace)
	state.MetricName = types.StringValue(alarm.MetricName)
	state.Threshold = types.Float64Value(alarm.Threshold)
	state.Period = types.Int64Value(int64(alarm.Period))
	state.RepeatInterval = types.Int64Value(int64(alarm.RepeatInterval))
	state.Status = types.StringValue(alarm.Status)
	state.State = types.StringValue(alarm.State)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *alarmResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan alarmModel
	var state alarmModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.UpdateAlarmParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateAlarmParamDetail{
			Name:               plan.Name.ValueString(),
			ComparisonOperator: stringPtrOrNil(plan.ComparisonOperator.ValueString()),
			Description:        stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	threshold := plan.Threshold.ValueFloat64()
	p.Params.Threshold = &threshold

	if !plan.Period.IsNull() && !plan.Period.IsUnknown() {
		val := int(plan.Period.ValueInt64())
		p.Params.Period = &val
	}

	if !plan.RepeatInterval.IsNull() && !plan.RepeatInterval.IsUnknown() {
		val := int(plan.RepeatInterval.ValueInt64())
		p.Params.RepeatInterval = &val
	}

	result, err := r.client.UpdateAlarm(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating Alarm",
			"Could not update alarm, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = stringValueOrNull(result.Description)
	plan.ComparisonOperator = types.StringValue(result.ComparisonOperator)
	plan.Namespace = types.StringValue(result.Namespace)
	plan.MetricName = types.StringValue(result.MetricName)
	plan.Threshold = types.Float64Value(result.Threshold)
	plan.Period = types.Int64Value(int64(result.Period))
	plan.RepeatInterval = types.Int64Value(int64(result.RepeatInterval))
	plan.Status = types.StringValue(result.Status)
	plan.State = types.StringValue(result.State)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *alarmResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state alarmModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Alarm UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteAlarm(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError("Error deleting Alarm", "Could not delete alarm, unexpected error: "+err.Error())
		return
	}
}

func (r *alarmResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
