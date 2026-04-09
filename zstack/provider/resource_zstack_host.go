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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

var (
	_ resource.Resource                = &hostResource{}
	_ resource.ResourceWithConfigure   = &hostResource{}
	_ resource.ResourceWithImportState = &hostResource{}
)

type hostResource struct {
	client *client.ZSClient
}

type hostResourceModel struct {
	Uuid           types.String `tfsdk:"uuid"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	ManagementIp   types.String `tfsdk:"management_ip"`
	ClusterUuid    types.String `tfsdk:"cluster_uuid"`
	ZoneUuid       types.String `tfsdk:"zone_uuid"`
	HypervisorType types.String `tfsdk:"hypervisor_type"`
	State          types.String `tfsdk:"state"`
	Status         types.String `tfsdk:"status"`
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
	SshPort        types.Int64  `tfsdk:"ssh_port"`
	Architecture   types.String `tfsdk:"architecture"`
}

func HostResource() resource.Resource {
	return &hostResource{}
}

// Configure implements resource.ResourceWithConfigure.
func (r *hostResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Metadata implements resource.Resource.
func (r *hostResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

// Schema implements resource.Resource.
func (r *hostResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage ZStack KVM hosts. A host is a physical server that provides computing resources for running VM instances.",
		MarkdownDescription: "Manage ZStack KVM hosts. A host is a physical server that provides computing resources for running VM instances.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the host.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the host.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the host.",
			},
			"management_ip": schema.StringAttribute{
				Required:    true,
				Description: "The management IP address of the host.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the cluster this host belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"zone_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the zone this host belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"hypervisor_type": schema.StringAttribute{
				Computed:    true,
				Description: "The hypervisor type of the host (e.g., KVM).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The state of the host (Enabled, Disabled, PreMaintenance, Maintaining).",
				Validators: []validator.String{
					stringvalidator.OneOf("Enabled", "Disabled", "PreMaintenance", "Maintaining"),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the host (Connected, Disconnected, Connecting).",
			},
			"username": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The SSH username for connecting to the KVM host.",
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The SSH password for connecting to the KVM host.",
			},
			"ssh_port": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The SSH port for connecting to the KVM host (default 22).",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"architecture": schema.StringAttribute{
				Computed:    true,
				Description: "The CPU architecture of the host.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *hostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan hostResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating host", map[string]any{"name": plan.Name.ValueString()})

	sshPort := int(plan.SshPort.ValueInt64())
	if plan.SshPort.IsNull() || plan.SshPort.IsUnknown() {
		sshPort = 22
	}

	createParam := param.AddKVMHostParam{
		BaseParam: param.BaseParam{},
		Params: param.AddKVMHostParamDetail{
			Name:         plan.Name.ValueString(),
			ManagementIp: plan.ManagementIp.ValueString(),
			ClusterUuid:  plan.ClusterUuid.ValueString(),
			Username:     plan.Username.ValueString(),
			Password:     plan.Password.ValueString(),
			SshPort:      intPtr(sshPort),
		},
	}

	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		createParam.Params.Description = stringPtr(plan.Description.ValueString())
	}

	host, err := r.client.AddKVMHost(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Host",
			"Could not create host, unexpected error: "+err.Error(),
		)
		return
	}

	state := hostModelFromView(host, plan)

	// If the desired state is Disabled, change the host state after creation
	if !plan.State.IsNull() && !plan.State.IsUnknown() && strings.EqualFold(plan.State.ValueString(), "Disabled") {
		_, err := r.client.ChangeHostState(state.Uuid.ValueString(), param.ChangeHostStateParam{
			BaseParam: param.BaseParam{},
			Params: param.ChangeHostStateParamDetail{
				StateEvent: "disable",
			},
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error changing Host state",
				"Could not change host state to Disabled: "+err.Error(),
			)
			return
		}
		state.State = types.StringValue("Disabled")
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (r *hostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state hostResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	host, err := findResourceByGet(r.client.GetHost, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Host",
			"Could not read host, unexpected error: "+err.Error(),
		)
		return
	}

	// Preserve sensitive fields from prior state since API doesn't return them
	refreshedState := hostModelFromView(host, state)

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *hostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan hostResourceModel
	var state hostResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()

	// Update general host properties (name, description, managementIp)
	if plan.Name.ValueString() != state.Name.ValueString() ||
		plan.Description.ValueString() != state.Description.ValueString() ||
		plan.ManagementIp.ValueString() != state.ManagementIp.ValueString() {

		updateParam := param.UpdateHostParam{
			BaseParam: param.BaseParam{},
			Params: param.UpdateHostParamDetail{
				Name:        plan.Name.ValueString(),
				Description: stringPtrOrNil(plan.Description.ValueString()),
			},
		}

		if _, err := r.client.UpdateHost(uuid, updateParam); err != nil {
			resp.Diagnostics.AddError(
				"Error updating Host",
				"Could not update host, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Update KVM-specific properties (username, password, sshPort)
	if plan.Username.ValueString() != state.Username.ValueString() ||
		plan.Password.ValueString() != state.Password.ValueString() ||
		plan.SshPort.ValueInt64() != state.SshPort.ValueInt64() {

		sshPort := int(plan.SshPort.ValueInt64())
		if plan.SshPort.IsNull() || plan.SshPort.IsUnknown() {
			sshPort = 22
		}

		updateKVMParam := param.UpdateKVMHostParam{
			BaseParam: param.BaseParam{},
			Params: param.UpdateKVMHostParamDetail{
				Username: stringPtr(plan.Username.ValueString()),
				Password: stringPtr(plan.Password.ValueString()),
				SshPort:  intPtr(sshPort),
			},
		}

		if _, err := r.client.UpdateKVMHost(uuid, updateKVMParam); err != nil {
			resp.Diagnostics.AddError(
				"Error updating KVM Host",
				"Could not update KVM host, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Handle state changes (enable/disable/maintain)
	if !plan.State.IsNull() && !plan.State.IsUnknown() &&
		plan.State.ValueString() != state.State.ValueString() {

		var stateEvent string
		switch strings.ToLower(plan.State.ValueString()) {
		case "enabled":
			stateEvent = "enable"
		case "disabled":
			stateEvent = "disable"
		case "premaintenance", "maintaining":
			stateEvent = "maintain"
		default:
			resp.Diagnostics.AddError(
				"Error updating Host",
				fmt.Sprintf("Could not update host, unsupported state: %s", plan.State.ValueString()),
			)
			return
		}

		if _, err := r.client.ChangeHostState(uuid, param.ChangeHostStateParam{
			BaseParam: param.BaseParam{},
			Params: param.ChangeHostStateParamDetail{
				StateEvent: stateEvent,
			},
		}); err != nil {
			resp.Diagnostics.AddError(
				"Error changing Host state",
				"Could not change host state: "+err.Error(),
			)
			return
		}
	}

	// Read back the updated resource
	host, err := r.client.GetHost(uuid)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Host",
			"Could not read updated host: "+err.Error(),
		)
		return
	}

	refreshedState := hostModelFromView(host, plan)

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *hostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state hostResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid.IsNull() || state.Uuid.ValueString() == "" {
		tflog.Warn(ctx, "host uuid is empty, skip delete")
		return
	}

	if err := r.client.DeleteHost(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Host",
			"Could not delete host, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *hostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func hostModelFromView(h *view.HostInventoryView, plan hostResourceModel) hostResourceModel {
	model := hostResourceModel{
		Uuid:           types.StringValue(h.UUID),
		Name:           types.StringValue(h.Name),
		Description:    stringValueOrNull(h.Description),
		ManagementIp:   types.StringValue(h.ManagementIp),
		ClusterUuid:    types.StringValue(h.ClusterUuid),
		ZoneUuid:       types.StringValue(h.ZoneUuid),
		HypervisorType: types.StringValue(h.HypervisorType),
		State:          stringValueOrNull(h.State),
		Status:         stringValueOrNull(h.Status),
		Architecture:   stringValueOrNull(h.Architecture),
		// Preserve sensitive fields from plan/state since API doesn't return them
		Username: plan.Username,
		Password: plan.Password,
		SshPort:  plan.SshPort,
	}

	// Default SshPort to 22 if not set
	if model.SshPort.IsNull() || model.SshPort.IsUnknown() {
		model.SshPort = types.Int64Value(22)
	}

	return model
}
