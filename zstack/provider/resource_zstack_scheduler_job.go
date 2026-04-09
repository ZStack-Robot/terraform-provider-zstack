// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &schedulerJobResource{}
	_ resource.ResourceWithConfigure   = &schedulerJobResource{}
	_ resource.ResourceWithImportState = &schedulerJobResource{}
)

type schedulerJobResource struct {
	client *client.ZSClient
}

type schedulerJobResourceModel struct {
	Uuid               types.String `tfsdk:"uuid"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	TargetResourceUuid types.String `tfsdk:"target_resource_uuid"`
	Type               types.String `tfsdk:"type"`
	State              types.String `tfsdk:"state"`
	JobData            types.String `tfsdk:"job_data"`
	JobClassName       types.String `tfsdk:"job_class_name"`
}

func SchedulerJobResource() resource.Resource {
	return &schedulerJobResource{}
}

func (r *schedulerJobResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *schedulerJobResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_scheduler_job"
}

func (r *schedulerJobResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Manages a scheduler job in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the scheduler job.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the scheduler job.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the scheduler job.",
			},
			"target_resource_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the target resource for the scheduler job.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the scheduler job.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the scheduler job.",
			},
			"job_data": schema.StringAttribute{
				Computed:    true,
				Description: "The job data of the scheduler job.",
			},
			"job_class_name": schema.StringAttribute{
				Computed:    true,
				Description: "The job class name of the scheduler job.",
			},
		},
	}
}

func (r *schedulerJobResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan schedulerJobResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	params := param.CreateSchedulerJobParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateSchedulerJobParamDetail{
			Name:               plan.Name.ValueString(),
			Description:        stringPtrOrNil(plan.Description.ValueString()),
			TargetResourceUuid: plan.TargetResourceUuid.ValueString(),
			Type:               plan.Type.ValueString(),
		},
	}

	job, err := r.client.CreateSchedulerJob(params)
	if err != nil {
		response.Diagnostics.AddError("Error creating scheduler job", err.Error())
		return
	}

	plan.Uuid = types.StringValue(job.UUID)
	plan.Name = types.StringValue(job.Name)
	plan.Description = stringValueOrNull(job.Description)
	plan.TargetResourceUuid = types.StringValue(job.TargetResourceUuid)
	plan.Type = types.StringValue(plan.Type.ValueString())
	plan.State = stringValueOrNull(job.State)
	plan.JobData = stringValueOrNull(job.JobData)
	plan.JobClassName = stringValueOrNull(job.JobClassName)

	diags = response.State.Set(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *schedulerJobResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state schedulerJobResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	job, err := findResourceByQuery(r.client.QuerySchedulerJob, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError("Error reading scheduler job", err.Error())
		return
	}
	state.Name = types.StringValue(job.Name)
	state.Description = stringValueOrNull(job.Description)
	state.TargetResourceUuid = types.StringValue(job.TargetResourceUuid)
	state.State = stringValueOrNull(job.State)
	state.JobData = stringValueOrNull(job.JobData)
	state.JobClassName = stringValueOrNull(job.JobClassName)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *schedulerJobResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state schedulerJobResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var plan schedulerJobResourceModel
	diags = request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	params := param.UpdateSchedulerJobParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateSchedulerJobParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	job, err := r.client.UpdateSchedulerJob(state.Uuid.ValueString(), params)
	if err != nil {
		response.Diagnostics.AddError("Error updating scheduler job", err.Error())
		return
	}

	plan.Uuid = types.StringValue(job.UUID)
	plan.Name = types.StringValue(job.Name)
	plan.Description = stringValueOrNull(job.Description)
	plan.TargetResourceUuid = types.StringValue(job.TargetResourceUuid)
	plan.State = stringValueOrNull(job.State)
	plan.JobData = stringValueOrNull(job.JobData)
	plan.JobClassName = stringValueOrNull(job.JobClassName)

	diags = response.State.Set(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *schedulerJobResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state schedulerJobResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSchedulerJob(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("Error deleting scheduler job", err.Error())
		return
	}
}

func (r *schedulerJobResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	parts := strings.SplitN(request.ID, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		response.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected format: <uuid>:<type> (e.g. abc123:VmInstanceBackupSchedulerJob). "+
				"The 'type' field is not returned by the API and must be supplied at import time.",
		)
		return
	}
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("uuid"), parts[0])...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("type"), parts[1])...)
}
