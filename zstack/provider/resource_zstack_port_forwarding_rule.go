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
	_ resource.Resource                = &portForwardingRuleResource{}
	_ resource.ResourceWithConfigure   = &portForwardingRuleResource{}
	_ resource.ResourceWithImportState = &portForwardingRuleResource{}
)

type portForwardingRuleResource struct {
	client *client.ZSClient
}

type portForwardingRuleResourceModel struct {
	Uuid             types.String `tfsdk:"uuid"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	VipUuid          types.String `tfsdk:"vip_uuid"`
	VipPortStart     types.Int64  `tfsdk:"vip_port_start"`
	VipPortEnd       types.Int64  `tfsdk:"vip_port_end"`
	PrivatePortStart types.Int64  `tfsdk:"private_port_start"`
	PrivatePortEnd   types.Int64  `tfsdk:"private_port_end"`
	ProtocolType     types.String `tfsdk:"protocol_type"`
	VmNicUuid        types.String `tfsdk:"vm_nic_uuid"`
	AllowedCidr      types.String `tfsdk:"allowed_cidr"`
	VipIp            types.String `tfsdk:"vip_ip"`
	GuestIp          types.String `tfsdk:"guest_ip"`
	State            types.String `tfsdk:"state"`
}

func PortForwardingRuleResource() resource.Resource {
	return &portForwardingRuleResource{}
}

// Configure implements resource.ResourceWithConfigure.
func (r *portForwardingRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *portForwardingRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_forwarding_rule"
}

// Schema implements resource.Resource.
func (r *portForwardingRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage ZStack port forwarding rules. Port forwarding rules map VIP ports to VM instance private ports, supporting TCP and UDP protocols.",
		MarkdownDescription: "Manage ZStack port forwarding rules. Port forwarding rules map VIP ports to VM instance private ports, supporting TCP and UDP protocols.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the port forwarding rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the port forwarding rule.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the port forwarding rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vip_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the VIP used for port forwarding.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vip_port_start": schema.Int64Attribute{
				Required:    true,
				Description: "The start port on the VIP side.",
			},
			"vip_port_end": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The end port on the VIP side. Defaults to the same as vip_port_start.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"private_port_start": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The start port on the private (VM) side. Defaults to the same as vip_port_start.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"private_port_end": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The end port on the private (VM) side. Defaults to the same as vip_port_end.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"protocol_type": schema.StringAttribute{
				Required:    true,
				Description: "The protocol type: TCP or UDP.",
				Validators: []validator.String{
					stringvalidator.OneOf("TCP", "UDP"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vm_nic_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The UUID of the VM NIC to which the port forwarding rule is attached.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_cidr": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The CIDR block allowed to access this port forwarding rule (e.g., 0.0.0.0/0).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vip_ip": schema.StringAttribute{
				Computed:    true,
				Description: "The IP address of the VIP.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"guest_ip": schema.StringAttribute{
				Computed:    true,
				Description: "The IP address of the guest VM.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the port forwarding rule (Enabled, Disabled).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *portForwardingRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan portForwardingRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating port forwarding rule", map[string]any{"name": plan.Name.ValueString()})

	createParam := param.CreatePortForwardingRuleParam{
		BaseParam: param.BaseParam{},
		Params: param.CreatePortForwardingRuleParamDetail{
			Name:         plan.Name.ValueString(),
			VipUuid:      plan.VipUuid.ValueString(),
			VipPortStart: int(plan.VipPortStart.ValueInt64()),
			ProtocolType: plan.ProtocolType.ValueString(),
		},
	}

	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		createParam.Params.Description = stringPtr(plan.Description.ValueString())
	}
	if !plan.VipPortEnd.IsNull() {
		createParam.Params.VipPortEnd = intPtr(int(plan.VipPortEnd.ValueInt64()))
	}
	if !plan.PrivatePortStart.IsNull() {
		createParam.Params.PrivatePortStart = intPtr(int(plan.PrivatePortStart.ValueInt64()))
	}
	if !plan.PrivatePortEnd.IsNull() {
		createParam.Params.PrivatePortEnd = intPtr(int(plan.PrivatePortEnd.ValueInt64()))
	}
	if !plan.VmNicUuid.IsNull() && plan.VmNicUuid.ValueString() != "" {
		createParam.Params.VmNicUuid = stringPtr(plan.VmNicUuid.ValueString())
	}
	if !plan.AllowedCidr.IsNull() && plan.AllowedCidr.ValueString() != "" {
		createParam.Params.AllowedCidr = stringPtr(plan.AllowedCidr.ValueString())
	}

	rule, err := r.client.CreatePortForwardingRule(createParam)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Port Forwarding Rule", "Could not create port forwarding rule, unexpected error: "+err.Error())
		return
	}

	state := portForwardingRuleModelFromView(rule)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (r *portForwardingRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state portForwardingRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := findResourceByGet(r.client.GetPortForwardingRule, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Port Forwarding Rule", "Could not read port forwarding rule UUID "+state.Uuid.ValueString()+": "+err.Error())
		return
	}

	refreshedState := portForwardingRuleModelFromView(rule)

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *portForwardingRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan portForwardingRuleResourceModel
	var state portForwardingRuleResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()

	// Update name/description if changed
	if plan.Name.ValueString() != state.Name.ValueString() || plan.Description.ValueString() != state.Description.ValueString() {
		updateParam := param.UpdatePortForwardingRuleParam{
			BaseParam: param.BaseParam{},
			Params: param.UpdatePortForwardingRuleParamDetail{
				Name:        plan.Name.ValueString(),
				Description: stringPtrOrNil(plan.Description.ValueString()),
			},
		}

		if _, err := r.client.UpdatePortForwardingRule(uuid, updateParam); err != nil {
			resp.Diagnostics.AddError("Error updating Port Forwarding Rule", "Could not update port forwarding rule UUID "+uuid+": "+err.Error())
			return
		}
	}

	// Reconcile VM NIC attachment
	currentNic := state.VmNicUuid.ValueString()
	desiredNic := plan.VmNicUuid.ValueString()

	if currentNic != desiredNic {
		// Detach from current NIC if attached
		if currentNic != "" {
			if err := r.client.DetachPortForwardingRule(uuid, param.DeleteModePermissive); err != nil {
				resp.Diagnostics.AddError("Error detaching Port Forwarding Rule from VM NIC", "Could not detach port forwarding rule from VM NIC, unexpected error: "+err.Error())
				return
			}
		}

		// Attach to new NIC if specified
		if desiredNic != "" {
			if err := r.attachToVmNic(uuid, desiredNic); err != nil {
				resp.Diagnostics.AddError("Error attaching Port Forwarding Rule to VM NIC", "Could not attach port forwarding rule to VM NIC, unexpected error: "+err.Error())
				return
			}
		}
	}

	// Read back the updated resource
	rule, err := r.client.GetPortForwardingRule(uuid)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Port Forwarding Rule", "Could not read port forwarding rule after update UUID "+uuid+": "+err.Error())
		return
	}

	refreshedState := portForwardingRuleModelFromView(rule)

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *portForwardingRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state portForwardingRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid.IsNull() || state.Uuid.ValueString() == "" {
		tflog.Warn(ctx, "port forwarding rule uuid is empty, skip delete")
		return
	}

	if err := r.client.DeletePortForwardingRule(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError("Error deleting Port Forwarding Rule", "Could not delete port forwarding rule UUID "+state.Uuid.ValueString()+": "+err.Error())
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *portForwardingRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *portForwardingRuleResource) attachToVmNic(ruleUuid, vmNicUuid string) error {
	attachParam := param.AttachPortForwardingRuleParam{
		BaseParam: param.BaseParam{},
		Params:    param.AttachPortForwardingRuleParamDetail{},
	}

	if _, err := r.client.AttachPortForwardingRule(ruleUuid, vmNicUuid, attachParam); err != nil {
		return err
	}
	return nil
}

func portForwardingRuleModelFromView(rule *view.PortForwardingRuleInventoryView) portForwardingRuleResourceModel {
	return portForwardingRuleResourceModel{
		Uuid:             types.StringValue(rule.UUID),
		Name:             types.StringValue(rule.Name),
		Description:      stringValueOrNull(rule.Description),
		VipUuid:          types.StringValue(rule.VipUuid),
		VipPortStart:     types.Int64Value(int64(rule.VipPortStart)),
		VipPortEnd:       types.Int64Value(int64(rule.VipPortEnd)),
		PrivatePortStart: types.Int64Value(int64(rule.PrivatePortStart)),
		PrivatePortEnd:   types.Int64Value(int64(rule.PrivatePortEnd)),
		ProtocolType:     types.StringValue(rule.ProtocolType),
		VmNicUuid:        stringValueOrNull(rule.VmNicUuid),
		AllowedCidr:      stringValueOrNull(rule.AllowedCidr),
		VipIp:            stringValueOrNull(rule.VipIp),
		GuestIp:          stringValueOrNull(rule.GuestIp),
		State:            stringValueOrNull(rule.State),
	}
}
