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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &logServerResource{}
	_ resource.ResourceWithConfigure   = &logServerResource{}
	_ resource.ResourceWithImportState = &logServerResource{}
)

type logServerResource struct {
	client *client.ZSClient
}

type logServerModel struct {
	Uuid          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Category      types.String `tfsdk:"category"`
	Type          types.String `tfsdk:"type"`
	Level         types.String `tfsdk:"level"`
	Configuration types.String `tfsdk:"configuration"`
}

func LogServerResource() resource.Resource {
	return &logServerResource{}
}

func (r *logServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *logServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_server"
}

func (r *logServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage log server in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the log server.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the log server.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the log server.",
			},
			"category": schema.StringAttribute{
				Required:    true,
				Description: "The log category.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The log server type.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"level": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The log level.",
			},
			"configuration": schema.StringAttribute{
				Required:    true,
				Description: "The log server configuration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *logServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan logServerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.AddLogServerParam{
		BaseParam: param.BaseParam{},
		Params: param.AddLogServerParamDetail{
			Name:          plan.Name.ValueString(),
			Description:   stringPtrOrNil(plan.Description.ValueString()),
			Category:      plan.Category.ValueString(),
			Type:          plan.Type.ValueString(),
			Level:         stringPtrOrNil(plan.Level.ValueString()),
			Configuration: plan.Configuration.ValueString(),
		},
	}

	item, err := r.client.AddLogServer(p)
	if err != nil {
		resp.Diagnostics.AddError("Fail to add log server", "Error "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(item.UUID)
	plan.Name = types.StringValue(item.Name)
	plan.Description = stringValueOrNull(item.Description)
	plan.Category = stringValueOrNull(item.Category)
	plan.Type = stringValueOrNull(item.Type)
	plan.Level = stringValueOrNull(item.Level)
	plan.Configuration = stringValueOrNull(item.Configuration)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *logServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state logServerModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	items, err := r.client.QueryLogServer(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query log servers. It may have been deleted.: "+err.Error())
		state = logServerModel{Uuid: types.StringValue("")}
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
			state.Category = stringValueOrNull(item.Category)
			state.Type = stringValueOrNull(item.Type)
			state.Level = stringValueOrNull(item.Level)
			state.Configuration = stringValueOrNull(item.Configuration)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Log server not found. It might have been deleted outside of Terraform.")
		state = logServerModel{Uuid: types.StringValue("")}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *logServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan logServerModel
	var state logServerModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := param.UpdateLogServerParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateLogServerParamDetail{
			Uuid:        state.Uuid.ValueString(),
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	item, err := r.client.UpdateLogServer(p)
	if err != nil {
		resp.Diagnostics.AddError("Fail to update log server", "Error "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(item.UUID)
	plan.Name = types.StringValue(item.Name)
	plan.Description = stringValueOrNull(item.Description)
	plan.Category = stringValueOrNull(item.Category)
	plan.Type = stringValueOrNull(item.Type)
	plan.Level = stringValueOrNull(item.Level)
	plan.Configuration = stringValueOrNull(item.Configuration)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *logServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state logServerModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Log server UUID is empty, skipping delete.")
		return
	}

	if err := r.client.DeleteLogServer(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError("Fail to delete log server", err.Error())
		return
	}
}

func (r *logServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
