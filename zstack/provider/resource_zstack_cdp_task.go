// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &cdpTaskResource{}
	_ resource.ResourceWithConfigure   = &cdpTaskResource{}
	_ resource.ResourceWithImportState = &cdpTaskResource{}
)

type cdpTaskResource struct {
	client *client.ZSClient
}

type cdpTaskModel struct {
	Uuid              types.String `tfsdk:"uuid"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	TaskType          types.String `tfsdk:"task_type"`
	PolicyUuid        types.String `tfsdk:"policy_uuid"`
	BackupStorageUuid types.String `tfsdk:"backup_storage_uuid"`
	ResourceUuids     types.List   `tfsdk:"resource_uuids"`
	BackupBandwidth   types.Int64  `tfsdk:"backup_bandwidth"`
	MaxCapacity       types.Int64  `tfsdk:"max_capacity"`
	MaxLatency        types.Int64  `tfsdk:"max_latency"`
	Status            types.String `tfsdk:"status"`
	State             types.String `tfsdk:"state"`
}

func CdpTaskResource() resource.Resource {
	return &cdpTaskResource{}
}

func (r *cdpTaskResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cli, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", req.ProviderData),
		)
		return
	}

	r.client = cli
}

func (r *cdpTaskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdp_task"
}

func (r *cdpTaskResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage CDP tasks in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the CDP task.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the CDP task.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the CDP task.",
			},
			"task_type": schema.StringAttribute{
				Required:    true,
				Description: "The task type.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The policy UUID for the CDP task.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"backup_storage_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The backup storage UUID for the CDP task.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_uuids": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "The resource UUID list for the CDP task.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"backup_bandwidth": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Backup bandwidth.",
			},
			"max_capacity": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Max capacity.",
			},
			"max_latency": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Max latency.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the CDP task.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the CDP task.",
			},
		},
	}
}

func (r *cdpTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cdpTaskModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateCdpTaskParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateCdpTaskParamDetail{
			Name:              plan.Name.ValueString(),
			Description:       stringPtrOrNil(plan.Description.ValueString()),
			TaskType:          plan.TaskType.ValueString(),
			PolicyUuid:        plan.PolicyUuid.ValueString(),
			BackupStorageUuid: plan.BackupStorageUuid.ValueString(),
			ResourceUuids:     listToStringSlice(plan.ResourceUuids),
		},
	}

	if !plan.BackupBandwidth.IsNull() && !plan.BackupBandwidth.IsUnknown() {
		p.Params.BackupBandwidth = int64Ptr(plan.BackupBandwidth.ValueInt64())
	}
	if !plan.MaxCapacity.IsNull() && !plan.MaxCapacity.IsUnknown() {
		p.Params.MaxCapacity = int64Ptr(plan.MaxCapacity.ValueInt64())
	}
	if !plan.MaxLatency.IsNull() && !plan.MaxLatency.IsUnknown() {
		p.Params.MaxLatency = int64Ptr(plan.MaxLatency.ValueInt64())
	}

	item, err := r.client.CreateCdpTask(p)
	if err != nil {
		resp.Diagnostics.AddError("Fail to create CDP task", "Error "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(item.UUID)
	plan.Name = types.StringValue(item.Name)
	plan.Description = stringValueOrNull(item.Description)
	plan.TaskType = stringValueOrNull(item.TaskType)
	plan.PolicyUuid = stringValueOrNull(item.PolicyUuid)
	plan.BackupStorageUuid = stringValueOrNull(item.BackupStorageUuid)
	plan.BackupBandwidth = types.Int64Value(item.BackupBandwidth)
	plan.MaxCapacity = types.Int64Value(item.MaxCapacity)
	plan.MaxLatency = types.Int64Value(item.MaxLatency)
	plan.Status = stringValueOrNull(item.Status)
	plan.State = stringValueOrNull(item.State)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *cdpTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cdpTaskModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	items, err := r.client.QueryCdpTask(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query CDP tasks. It may have been deleted.: "+err.Error())
		state = cdpTaskModel{Uuid: types.StringValue("")}
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		return
	}

	found := false
	for _, item := range items {
		if item.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(item.UUID)
			state.Name = types.StringValue(item.Name)
			state.Description = stringValueOrNull(item.Description)
			state.TaskType = stringValueOrNull(item.TaskType)
			state.PolicyUuid = stringValueOrNull(item.PolicyUuid)
			state.BackupStorageUuid = stringValueOrNull(item.BackupStorageUuid)
			state.BackupBandwidth = types.Int64Value(item.BackupBandwidth)
			state.MaxCapacity = types.Int64Value(item.MaxCapacity)
			state.MaxLatency = types.Int64Value(item.MaxLatency)
			state.Status = stringValueOrNull(item.Status)
			state.State = stringValueOrNull(item.State)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "CDP task not found. It might have been deleted outside of Terraform.")
		state = cdpTaskModel{Uuid: types.StringValue("")}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *cdpTaskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan cdpTaskModel
	var state cdpTaskModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := param.UpdateCdpTaskParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateCdpTaskParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	if !plan.BackupBandwidth.IsNull() && !plan.BackupBandwidth.IsUnknown() {
		p.Params.BackupBandwidth = int64Ptr(plan.BackupBandwidth.ValueInt64())
	}
	if !plan.MaxCapacity.IsNull() && !plan.MaxCapacity.IsUnknown() {
		p.Params.MaxCapacity = int64Ptr(plan.MaxCapacity.ValueInt64())
	}
	if !plan.MaxLatency.IsNull() && !plan.MaxLatency.IsUnknown() {
		p.Params.MaxLatency = int64Ptr(plan.MaxLatency.ValueInt64())
	}

	item, err := r.client.UpdateCdpTask(state.Uuid.ValueString(), p)
	if err != nil {
		resp.Diagnostics.AddError("Fail to update CDP task", "Error "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(item.UUID)
	plan.Name = types.StringValue(item.Name)
	plan.Description = stringValueOrNull(item.Description)
	plan.TaskType = stringValueOrNull(item.TaskType)
	plan.PolicyUuid = stringValueOrNull(item.PolicyUuid)
	plan.BackupStorageUuid = stringValueOrNull(item.BackupStorageUuid)
	plan.BackupBandwidth = types.Int64Value(item.BackupBandwidth)
	plan.MaxCapacity = types.Int64Value(item.MaxCapacity)
	plan.MaxLatency = types.Int64Value(item.MaxLatency)
	plan.Status = stringValueOrNull(item.Status)
	plan.State = stringValueOrNull(item.State)
	plan.ResourceUuids = state.ResourceUuids

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *cdpTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cdpTaskModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "CDP task UUID is empty, skipping delete.")
		return
	}

	if err := r.client.DeleteCdpTask(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError("Fail to delete CDP task", err.Error())
		return
	}
}

func (r *cdpTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
