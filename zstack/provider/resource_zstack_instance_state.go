// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

const (
	instanceStateRunning               = "Running"
	instanceStateStopped               = "Stopped"
	defaultInstanceStateStopType       = "grace"
	defaultInstanceStateTimeoutSeconds = int64(600)
)

var (
	_ resource.Resource                = &instanceStateResource{}
	_ resource.ResourceWithConfigure   = &instanceStateResource{}
	_ resource.ResourceWithImportState = &instanceStateResource{}
)

type instanceStateResource struct {
	client *client.ZSClient
}

type instanceStateModel struct {
	ID               types.String `tfsdk:"id"`
	VmInstanceUuid   types.String `tfsdk:"vm_instance_uuid"`
	State            types.String `tfsdk:"state"`
	StopType         types.String `tfsdk:"stop_type"`
	OperationTimeout types.Int64  `tfsdk:"operation_timeout"`
}

func InstanceStateResource() resource.Resource {
	return &instanceStateResource{}
}

func (r *instanceStateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *instanceStateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance_state"
}

func (r *instanceStateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the power state of an existing ZStack VM instance. This resource can start or stop a VM without owning the VM lifecycle. Destroying this resource only removes Terraform state management and does not change the VM power state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Same as vm_instance_uuid. Used for Terraform tracking.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vm_instance_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the ZStack VM instance whose power state is managed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Required:    true,
				Description: "Desired VM power state. Supported values are Running and Stopped.",
				Validators: []validator.String{
					stringvalidator.OneOf(instanceStateRunning, instanceStateStopped),
				},
			},
			"stop_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Stop mode used when state is Stopped. Currently only grace is supported.",
				Default:     stringdefault.StaticString(defaultInstanceStateStopType),
				Validators: []validator.String{
					stringvalidator.OneOf(defaultInstanceStateStopType),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"operation_timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Maximum number of seconds to wait for the VM to reach the desired state.",
				Default:     int64default.StaticInt64(defaultInstanceStateTimeoutSeconds),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *instanceStateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan instanceStateModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.reconcileInstanceState(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error managing VM Instance State",
			"Could not set VM instance state: "+err.Error(),
		)
		return
	}

	plan.ID = plan.VmInstanceUuid
	plan.State = types.StringValue(state)
	applyInstanceStateDefaults(&plan)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *instanceStateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state instanceStateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Error reading VM Instance State",
			"Could not read VM instance state, provider client is not properly configured.",
		)
		return
	}

	vm, err := findResourceByGet(r.client.GetVmInstance, state.VmInstanceUuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			tflog.Warn(ctx, "VM instance not found, removing instance state resource", map[string]any{
				"vm_instance_uuid": state.VmInstanceUuid.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading VM Instance State",
			"Could not read VM instance UUID "+state.VmInstanceUuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.ID = types.StringValue(vm.UUID)
	state.VmInstanceUuid = types.StringValue(vm.UUID)

	readState, normalized, err := instanceStateForRead(vm.State, state.State)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading VM Instance State",
			err.Error(),
		)
		return
	}
	if normalized {
		tflog.Info(ctx, "Normalizing non-terminal VM instance state for Terraform state", map[string]any{
			"vm_instance_uuid": state.VmInstanceUuid.ValueString(),
			"actual_state":     vm.State,
			"terraform_state":  readState,
		})
	}
	state.State = types.StringValue(readState)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *instanceStateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan instanceStateModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.reconcileInstanceState(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error managing VM Instance State",
			"Could not set VM instance state: "+err.Error(),
		)
		return
	}

	plan.ID = plan.VmInstanceUuid
	plan.State = types.StringValue(state)
	applyInstanceStateDefaults(&plan)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *instanceStateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state instanceStateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Removing VM instance state management from Terraform state only", map[string]any{
		"vm_instance_uuid": state.VmInstanceUuid.ValueString(),
	})
}

func (r *instanceStateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vm_instance_uuid"), types.StringValue(req.ID))...)
}

func (r *instanceStateResource) reconcileInstanceState(ctx context.Context, plan instanceStateModel) (string, error) {
	if r.client == nil {
		return "", fmt.Errorf("provider client is not properly configured")
	}

	uuid := plan.VmInstanceUuid.ValueString()
	desiredState := plan.State.ValueString()

	vm, err := findResourceByGet(r.client.GetVmInstance, uuid)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return "", fmt.Errorf("VM instance %s not found", uuid)
		}
		return "", fmt.Errorf("read VM instance %s: %w", uuid, err)
	}

	if vm.State == desiredState {
		return vm.State, nil
	}

	switch desiredState {
	case instanceStateRunning:
		if vm.State != "Starting" {
			tflog.Info(ctx, "Starting VM instance", map[string]any{"vm_instance_uuid": uuid})
			if _, err := r.client.StartVmInstance(uuid, param.StartVmInstanceParam{
				Params: param.StartVmInstanceParamDetail{},
			}); err != nil {
				return "", fmt.Errorf("start VM instance %s: %w", uuid, err)
			}
		}
	case instanceStateStopped:
		if vm.State != "Stopping" {
			stopType := instanceStateStopType(plan)
			tflog.Info(ctx, "Stopping VM instance", map[string]any{
				"vm_instance_uuid": uuid,
				"stop_type":        stopType,
			})
			if _, err := r.client.StopVmInstance(uuid, param.StopVmInstanceParam{
				Params: param.StopVmInstanceParamDetail{
					Type: stringPtr(stopType),
				},
			}); err != nil {
				return "", fmt.Errorf("stop VM instance %s: %w", uuid, err)
			}
		}
	default:
		return "", fmt.Errorf("unsupported state %q", desiredState)
	}

	state, err := r.waitForInstanceState(ctx, uuid, desiredState, instanceStateTimeout(plan))
	if err != nil {
		return "", err
	}

	return state, nil
}

func (r *instanceStateResource) waitForInstanceState(ctx context.Context, uuid string, desiredState string, timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	lastState := ""

	for {
		vm, err := findResourceByGet(r.client.GetVmInstance, uuid)
		if err != nil {
			return "", fmt.Errorf("read VM instance %s while waiting for state %s: %w", uuid, desiredState, err)
		}

		lastState = vm.State
		if lastState == desiredState {
			return lastState, nil
		}

		remaining := time.Until(deadline)
		if remaining <= 0 {
			return "", fmt.Errorf("timed out waiting for VM instance %s to reach state %s; last state was %s", uuid, desiredState, lastState)
		}

		sleepFor := 5 * time.Second
		if remaining < sleepFor {
			sleepFor = remaining
		}

		timer := time.NewTimer(sleepFor)
		select {
		case <-ctx.Done():
			timer.Stop()
			return "", ctx.Err()
		case <-timer.C:
		}
	}
}

func instanceStateStopType(model instanceStateModel) string {
	if model.StopType.IsNull() || model.StopType.IsUnknown() || model.StopType.ValueString() == "" {
		return defaultInstanceStateStopType
	}
	return model.StopType.ValueString()
}

func instanceStateTimeout(model instanceStateModel) time.Duration {
	if model.OperationTimeout.IsNull() || model.OperationTimeout.IsUnknown() || model.OperationTimeout.ValueInt64() <= 0 {
		return time.Duration(defaultInstanceStateTimeoutSeconds) * time.Second
	}
	return time.Duration(model.OperationTimeout.ValueInt64()) * time.Second
}

func instanceStateForRead(actualState string, previousState types.String) (string, bool, error) {
	switch actualState {
	case instanceStateRunning, instanceStateStopped:
		return actualState, false, nil
	case "Starting":
		return instanceStateRunning, true, nil
	case "Stopping":
		return instanceStateStopped, true, nil
	}

	if !previousState.IsNull() && !previousState.IsUnknown() {
		switch previousState.ValueString() {
		case instanceStateRunning, instanceStateStopped:
			return previousState.ValueString(), true, nil
		}
	}

	return "", false, fmt.Errorf("VM instance is in unsupported state %q. Retry when it reaches Running or Stopped", actualState)
}

func applyInstanceStateDefaults(model *instanceStateModel) {
	if model.StopType.IsNull() || model.StopType.IsUnknown() || model.StopType.ValueString() == "" {
		model.StopType = types.StringValue(defaultInstanceStateStopType)
	}
	if model.OperationTimeout.IsNull() || model.OperationTimeout.IsUnknown() || model.OperationTimeout.ValueInt64() <= 0 {
		model.OperationTimeout = types.Int64Value(defaultInstanceStateTimeoutSeconds)
	}
}
