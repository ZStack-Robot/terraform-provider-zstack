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
	_ resource.Resource                = &cdpPolicyResource{}
	_ resource.ResourceWithConfigure   = &cdpPolicyResource{}
	_ resource.ResourceWithImportState = &cdpPolicyResource{}
)

type cdpPolicyResource struct {
	client *client.ZSClient
}

type cdpPolicyModel struct {
	Uuid                    types.String `tfsdk:"uuid"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	RetentionTimePerDay     types.Int64  `tfsdk:"retention_time_per_day"`
	HourlyRpSinceDay        types.Int64  `tfsdk:"hourly_rp_since_day"`
	DailyRpSinceDay         types.Int64  `tfsdk:"daily_rp_since_day"`
	ExpireTimeInDay         types.Int64  `tfsdk:"expire_time_in_day"`
	FullBackupIntervalInDay types.Int64  `tfsdk:"full_backup_interval_in_day"`
	RecoveryPointPerSecond  types.Int64  `tfsdk:"recovery_point_per_second"`
	State                   types.String `tfsdk:"state"`
}

func CdpPolicyResource() resource.Resource {
	return &cdpPolicyResource{}
}

func (r *cdpPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cdpPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdp_policy"
}

func (r *cdpPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage CDP policy in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the CDP policy.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the CDP policy.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the CDP policy.",
			},
			"retention_time_per_day": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Retention time per day.",
			},
			"hourly_rp_since_day": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Hourly recovery point since day.",
			},
			"daily_rp_since_day": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Daily recovery point since day.",
			},
			"expire_time_in_day": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Expire time in day.",
			},
			"full_backup_interval_in_day": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Full backup interval in day.",
			},
			"recovery_point_per_second": schema.Int64Attribute{
				Required:    true,
				Description: "Recovery point per second.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the CDP policy.",
			},
		},
	}
}

func (r *cdpPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cdpPolicyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateCdpPolicyParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateCdpPolicyParamDetail{
			Name:                   plan.Name.ValueString(),
			Description:            stringPtrOrNil(plan.Description.ValueString()),
			RecoveryPointPerSecond: int(plan.RecoveryPointPerSecond.ValueInt64()),
		},
	}

	if !plan.RetentionTimePerDay.IsNull() && !plan.RetentionTimePerDay.IsUnknown() {
		p.Params.RetentionTimePerDay = intPtr(int(plan.RetentionTimePerDay.ValueInt64()))
	}
	if !plan.HourlyRpSinceDay.IsNull() && !plan.HourlyRpSinceDay.IsUnknown() {
		p.Params.HourlyRpSinceDay = intPtr(int(plan.HourlyRpSinceDay.ValueInt64()))
	}
	if !plan.DailyRpSinceDay.IsNull() && !plan.DailyRpSinceDay.IsUnknown() {
		p.Params.DailyRpSinceDay = intPtr(int(plan.DailyRpSinceDay.ValueInt64()))
	}
	if !plan.ExpireTimeInDay.IsNull() && !plan.ExpireTimeInDay.IsUnknown() {
		p.Params.ExpireTimeInDay = intPtr(int(plan.ExpireTimeInDay.ValueInt64()))
	}
	if !plan.FullBackupIntervalInDay.IsNull() && !plan.FullBackupIntervalInDay.IsUnknown() {
		p.Params.FullBackupIntervalInDay = intPtr(int(plan.FullBackupIntervalInDay.ValueInt64()))
	}

	item, err := r.client.CreateCdpPolicy(p)
	if err != nil {
		resp.Diagnostics.AddError("Fail to create CDP policy", "Error "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(item.UUID)
	plan.Name = types.StringValue(item.Name)
	plan.Description = stringValueOrNull(item.Description)
	plan.RetentionTimePerDay = types.Int64Value(int64(item.RetentionTimePerDay))
	plan.HourlyRpSinceDay = types.Int64Value(int64(item.HourlyRpSinceDay))
	plan.DailyRpSinceDay = types.Int64Value(int64(item.DailyRpSinceDay))
	plan.ExpireTimeInDay = types.Int64Value(int64(item.ExpireTimeInDay))
	plan.FullBackupIntervalInDay = types.Int64Value(int64(item.FullBackupIntervalInDay))
	plan.RecoveryPointPerSecond = types.Int64Value(int64(item.RecoveryPointPerSecond))
	plan.State = stringValueOrNull(item.State)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *cdpPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cdpPolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	items, err := r.client.QueryCdpPolicy(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query CDP policies. It may have been deleted.: "+err.Error())
		state = cdpPolicyModel{Uuid: types.StringValue("")}
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
			state.RetentionTimePerDay = types.Int64Value(int64(item.RetentionTimePerDay))
			state.HourlyRpSinceDay = types.Int64Value(int64(item.HourlyRpSinceDay))
			state.DailyRpSinceDay = types.Int64Value(int64(item.DailyRpSinceDay))
			state.ExpireTimeInDay = types.Int64Value(int64(item.ExpireTimeInDay))
			state.FullBackupIntervalInDay = types.Int64Value(int64(item.FullBackupIntervalInDay))
			state.RecoveryPointPerSecond = types.Int64Value(int64(item.RecoveryPointPerSecond))
			state.State = stringValueOrNull(item.State)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "CDP policy not found. It might have been deleted outside of Terraform.")
		state = cdpPolicyModel{Uuid: types.StringValue("")}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *cdpPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan cdpPolicyModel
	var state cdpPolicyModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := param.UpdateCdpPolicyParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateCdpPolicyParamDetail{
			Name:                   plan.Name.ValueString(),
			Description:            stringPtrOrNil(plan.Description.ValueString()),
			RecoveryPointPerSecond: intPtr(int(plan.RecoveryPointPerSecond.ValueInt64())),
		},
	}

	if !plan.RetentionTimePerDay.IsNull() && !plan.RetentionTimePerDay.IsUnknown() {
		p.Params.RetentionTimePerDay = intPtr(int(plan.RetentionTimePerDay.ValueInt64()))
	}
	if !plan.HourlyRpSinceDay.IsNull() && !plan.HourlyRpSinceDay.IsUnknown() {
		p.Params.HourlyRpSinceDay = intPtr(int(plan.HourlyRpSinceDay.ValueInt64()))
	}
	if !plan.DailyRpSinceDay.IsNull() && !plan.DailyRpSinceDay.IsUnknown() {
		p.Params.DailyRpSinceDay = intPtr(int(plan.DailyRpSinceDay.ValueInt64()))
	}
	if !plan.ExpireTimeInDay.IsNull() && !plan.ExpireTimeInDay.IsUnknown() {
		p.Params.ExpireTimeInDay = intPtr(int(plan.ExpireTimeInDay.ValueInt64()))
	}
	if !plan.FullBackupIntervalInDay.IsNull() && !plan.FullBackupIntervalInDay.IsUnknown() {
		p.Params.FullBackupIntervalInDay = intPtr(int(plan.FullBackupIntervalInDay.ValueInt64()))
	}

	item, err := r.client.UpdateCdpPolicy(state.Uuid.ValueString(), p)
	if err != nil {
		resp.Diagnostics.AddError("Fail to update CDP policy", "Error "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(item.UUID)
	plan.Name = types.StringValue(item.Name)
	plan.Description = stringValueOrNull(item.Description)
	plan.RetentionTimePerDay = types.Int64Value(int64(item.RetentionTimePerDay))
	plan.HourlyRpSinceDay = types.Int64Value(int64(item.HourlyRpSinceDay))
	plan.DailyRpSinceDay = types.Int64Value(int64(item.DailyRpSinceDay))
	plan.ExpireTimeInDay = types.Int64Value(int64(item.ExpireTimeInDay))
	plan.FullBackupIntervalInDay = types.Int64Value(int64(item.FullBackupIntervalInDay))
	plan.RecoveryPointPerSecond = types.Int64Value(int64(item.RecoveryPointPerSecond))
	plan.State = stringValueOrNull(item.State)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *cdpPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cdpPolicyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "CDP policy UUID is empty, skipping delete.")
		return
	}

	if err := r.client.DeleteCdpPolicy(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError("Fail to delete CDP policy", err.Error())
		return
	}
}

func (r *cdpPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
