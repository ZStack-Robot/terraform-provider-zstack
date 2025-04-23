// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

var (
	_ resource.Resource              = &scriptExecutionResource{}
	_ resource.ResourceWithConfigure = &scriptExecutionResource{}
)

type scriptExecutionResource struct {
	client *client.ZSClient
}

type scriptExecutionModel struct {
	Uuid           types.String `tfsdk:"uuid"`
	ScriptUuid     types.String `tfsdk:"script_uuid"`
	InstanceUuid   types.String `tfsdk:"instance_uuid"`
	ScriptTimeout  types.Int64  `tfsdk:"script_timeout"`
	RecordName     types.String `tfsdk:"record_name"`
	Status         types.String `tfsdk:"status"`
	Executor       types.String `tfsdk:"executor"`
	ExecutionCount types.Int64  `tfsdk:"execution_count"`
	Version        types.Int32  `tfsdk:"version"`
}

func ScriptExecutionResource() resource.Resource {
	return &scriptExecutionResource{}
}

func (r *scriptExecutionResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *scriptExecutionResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_script_execution"
}

func (r *scriptExecutionResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Executes a VM instance script in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the script execution record.",
			},
			"script_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The content of the script.",
			},
			"instance_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the VM instance to execute the script on.",
			},
			"script_timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The timeout for script execution in seconds.",
				Default:     int64default.StaticInt64(300),
			},
			"record_name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the execution record.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the script execution (Running, Succeeded, Failed).",
			},
			"executor": schema.StringAttribute{
				Computed:    true,
				Description: "The executor of the script.",
			},
			"version": schema.Int32Attribute{
				Computed:    true,
				Description: "The version of the script.",
			},
			"execution_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The execution count of the script.",
			},
		},
	}
}

func (r *scriptExecutionResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan scriptExecutionModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	scriptTimeout := int(plan.ScriptTimeout.ValueInt64())
	if plan.ScriptTimeout.IsNull() {
		scriptTimeout = 300
	}

	var systemTags []string
	executeParam := param.ExecuteVmInstanceScriptParam{
		BaseParam: param.BaseParam{
			SystemTags: systemTags,
		},
		Params: param.ExecuteVmInstanceScriptDetailParam{
			VmInstanceUuids: []string{plan.InstanceUuid.ValueString()},
			ScriptTimeout:   scriptTimeout,
		},
	}

	tflog.Debug(ctx, "Executing VM instance script", map[string]interface{}{
		"script_uuid":      plan.ScriptUuid.ValueString(),
		"vm_instance_uuid": plan.InstanceUuid.ValueString(),
	})

	scriptExecuteResult, err := r.client.ExecuteVmInstanceScript(plan.ScriptUuid.ValueString(), executeParam)
	if err != nil {
		response.Diagnostics.AddError(
			"Failed to execute VM instance script",
			"Error: "+err.Error(),
		)
		return
	}

	recordUuid := scriptExecuteResult.UUID
	maxWaitTime := time.Duration(scriptTimeout) * time.Second
	pollingInterval := 10 * time.Second
	startTime := time.Now()

	tflog.Debug(ctx, "Waiting for script execution to complete", map[string]interface{}{
		"record_uuid":   recordUuid,
		"max_wait_time": maxWaitTime,
	})

	var record *view.VmInstanceScriptResultInventoryView
	for {
		if time.Since(startTime) > maxWaitTime {
			response.Diagnostics.AddError(
				"Script execution timed out",
				fmt.Sprintf("The script execution did not complete within %s seconds.", maxWaitTime),
			)
			return
		}

		record, err = r.client.GetVmInstanceScriptExecutedRecord(recordUuid)
		if err != nil {
			response.Diagnostics.AddError(
				"Failed to retrieve script execution record",
				"Error: "+err.Error(),
			)
			return
		}

		if record.Status == "Succeed" || record.Status == "Failed" || record.Status == "Exception" || record.Status == "Succeeded" {
			tflog.Info(ctx, "Script execution completed", map[string]interface{}{
				"record_uuid":  recordUuid,
				"status":       record.Status,
				"elapsed_time": time.Since(startTime),
			})
			break
		}

		tflog.Debug(ctx, "Script execution in progress", map[string]interface{}{
			"status": record.Status,
		})

		time.Sleep(pollingInterval)
	}

	plan.Uuid = types.StringValue(record.UUID)
	plan.RecordName = types.StringValue(record.RecordName)
	plan.Status = types.StringValue(record.Status)
	plan.Executor = types.StringValue(record.Executor)
	plan.ScriptUuid = types.StringValue(record.ScriptUUID)
	plan.ExecutionCount = types.Int64Value(int64(record.ExecutionCount))
	plan.Version = types.Int32Value(int32(record.Version))

	if record.Status == "Failed" {
		response.Diagnostics.AddWarning(
			"Script Execution Failed",
			fmt.Sprintf("The script execution failed. Please verify the script content and the status of the target VM. Execution record UUID: %s", record.UUID),
		)
	}

	tflog.Info(ctx, "VM instance script executed successfully", map[string]interface{}{
		"id":     record.UUID,
		"status": record.Status,
	})

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *scriptExecutionResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state scriptExecutionModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddError("Client Not Configured", "The provider client is not properly configured.")
		return
	}

	record, err := r.client.GetVmInstanceScriptExecutedRecord(state.Uuid.ValueString())

	if err != nil {
		tflog.Warn(ctx, "Unable to retrieve VM instance script execution record. It may have been deleted.", map[string]interface{}{
			"id":    state.Uuid.ValueString(),
			"error": err.Error(),
		})
		response.State.RemoveResource(ctx)
		return
	}

	state.RecordName = types.StringValue(record.RecordName)
	state.Status = types.StringValue(record.Status)
	state.Executor = types.StringValue(record.Executor)
	state.ScriptUuid = types.StringValue(record.ScriptUUID)
	state.ExecutionCount = types.Int64Value(int64(record.ExecutionCount))
	state.Version = types.Int32Value(int32(record.Version))

	diags = response.State.Set(ctx, state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *scriptExecutionResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	response.Diagnostics.AddWarning(
		"Update Not Supported",
		"Updating script executions is not supported. Terraform will not perform any update operations on this resource.",
	)
}

func (r *scriptExecutionResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state scriptExecutionModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Removing script execution from Terraform state only", map[string]interface{}{
		"uuid": request.State.GetAttribute(ctx, path.Root("uuid"), nil),
	})
}
