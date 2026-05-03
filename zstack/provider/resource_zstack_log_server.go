// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	Uuid           types.String `tfsdk:"uuid"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Category       types.String `tfsdk:"category"`
	Type           types.String `tfsdk:"type"`
	Level          types.String `tfsdk:"level"`
	Configuration  types.String `tfsdk:"configuration"`
	AppenderType   types.String `tfsdk:"appender_type"`
	AppenderConfig types.Map    `tfsdk:"appender_configuration"`
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the log server.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the log server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"category": schema.StringAttribute{
				Required:    true,
				Description: "The log category.",
				Validators: []validator.String{
					stringvalidator.OneOf("ManagementNodeLog", "PlatformOperationLog"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The log server type.",
				Validators: []validator.String{
					stringvalidator.OneOf("Log4j2", "FluentBit"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"level": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The log level.",
				Validators: []validator.String{
					stringvalidator.OneOf("OFF", "FATAL", "ERROR", "WARN", "INFO", "DEBUG", "TRACE", "ALL"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"configuration": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Raw log server configuration JSON. Must be a JSON object containing appenderType and configuration. Mutually exclusive with appender_type/appender_configuration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"appender_type": schema.StringAttribute{
				Optional:    true,
				Description: "Appender type used to build configuration, for example Syslog. Must be set with appender_configuration and omitted when configuration is set.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"appender_configuration": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Appender configuration key/value map used to build the nested configuration JSON. For Syslog, typical keys are hostname, port, protocol, and facility.",
				Validators: []validator.Map{
					mapvalidator.SizeAtLeast(1),
				},
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
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

	configuration, ok := buildLogServerConfiguration(ctx, plan, &resp.Diagnostics)
	if !ok {
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
			Configuration: configuration,
		},
	}

	item, err := r.client.AddLogServer(p)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Log Server", "Could not create log server, unexpected error: "+err.Error())
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

	item, err := findResourceByQuery(r.client.QueryLogServer, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Log Server",
			"Could not read log server, unexpected error: "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(item.UUID)
	state.Name = types.StringValue(item.Name)
	state.Description = stringValueOrNull(item.Description)
	state.Category = stringValueOrNull(item.Category)
	state.Type = stringValueOrNull(item.Type)
	state.Level = stringValueOrNull(item.Level)
	state.Configuration = stringValueOrNull(item.Configuration)

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
		resp.Diagnostics.AddError("Error updating Log Server", "Could not update log server UUID "+state.Uuid.ValueString()+": "+err.Error())
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

	if err := r.client.DeleteLogServer(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError("Error deleting Log Server", "Could not delete log server UUID "+state.Uuid.ValueString()+": "+err.Error())
		return
	}
}

func (r *logServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

type logServerConfiguration struct {
	AppenderType  string            `json:"appenderType"`
	Configuration map[string]string `json:"configuration"`
}

type rawLogServerConfiguration struct {
	AppenderType  string          `json:"appenderType"`
	Configuration json.RawMessage `json:"configuration"`
}

func buildLogServerConfiguration(ctx context.Context, plan logServerModel, diags *diag.Diagnostics) (string, bool) {
	rawSet := !plan.Configuration.IsNull() && !plan.Configuration.IsUnknown() && plan.Configuration.ValueString() != ""
	appenderTypeSet := !plan.AppenderType.IsNull() && !plan.AppenderType.IsUnknown() && plan.AppenderType.ValueString() != ""
	appenderConfigSet := !plan.AppenderConfig.IsNull() && !plan.AppenderConfig.IsUnknown()

	if rawSet && (appenderTypeSet || appenderConfigSet) {
		diags.AddAttributeError(
			path.Root("configuration"),
			"Conflicting log server configuration",
			"Use either raw configuration or appender_type/appender_configuration, not both.",
		)
		return "", false
	}
	if rawSet {
		if !isValidRawLogServerConfiguration(plan.Configuration.ValueString()) {
			diags.AddAttributeError(
				path.Root("configuration"),
				"Invalid log server configuration",
				"configuration must be a JSON object containing appenderType and configuration. Example: {\"appenderType\":\"Syslog\",\"configuration\":{\"hostname\":\"192.168.0.11\",\"port\":\"514\",\"protocol\":\"UDP\",\"facility\":\"LOCAL5\"}}.",
			)
			return "", false
		}
		return plan.Configuration.ValueString(), true
	}
	if appenderTypeSet != appenderConfigSet {
		diags.AddAttributeError(
			path.Root("appender_type"),
			"Incomplete log server appender configuration",
			"appender_type and appender_configuration must be set together when raw configuration is omitted.",
		)
		return "", false
	}
	if !appenderTypeSet {
		diags.AddAttributeError(
			path.Root("configuration"),
			"Missing log server configuration",
			"Set either raw configuration or appender_type/appender_configuration.",
		)
		return "", false
	}

	var appenderConfiguration map[string]string
	diags.Append(plan.AppenderConfig.ElementsAs(ctx, &appenderConfiguration, false)...)
	if diags.HasError() {
		return "", false
	}

	payload, err := json.Marshal(logServerConfiguration{
		AppenderType:  plan.AppenderType.ValueString(),
		Configuration: appenderConfiguration,
	})
	if err != nil {
		diags.AddAttributeError(
			path.Root("appender_configuration"),
			"Invalid log server appender configuration",
			"Could not encode appender_configuration as JSON: "+err.Error(),
		)
		return "", false
	}

	return string(payload), true
}

func isValidRawLogServerConfiguration(configuration string) bool {
	var payload rawLogServerConfiguration
	if err := json.Unmarshal([]byte(configuration), &payload); err != nil {
		return false
	}
	if payload.AppenderType == "" || len(payload.Configuration) == 0 {
		return false
	}

	var config map[string]json.RawMessage
	if err := json.Unmarshal(payload.Configuration, &config); err != nil {
		return false
	}
	return config != nil
}
