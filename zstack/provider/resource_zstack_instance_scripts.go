// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
)

var (
	_ resource.Resource              = &scriptResource{}
	_ resource.ResourceWithConfigure = &scriptResource{}
)

type scriptResource struct {
	client *client.ZSClient
}

type scriptModel struct {
	Uuid          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	ScriptContent types.String `tfsdk:"script_content"`
	RenderParams  types.String `tfsdk:"render_params"`
	Platform      types.String `tfsdk:"platform"`
	ScriptType    types.String `tfsdk:"script_type"`
	ScriptTimeout types.Int64  `tfsdk:"script_timeout"`
}

func ScriptResource() resource.Resource {
	return &scriptResource{}
}

func (r *scriptResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *scriptResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_script"
}

func (r *scriptResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Manages a VM instance script in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the script.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the script.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the script.",
			},
			"script_content": schema.StringAttribute{
				Required:    true,
				Description: "The content of the script.",
			},
			"render_params": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "representing rendering parameters for the script.",
			},
			"platform": schema.StringAttribute{
				Optional:    true,
				Description: "The platform type for the script (e.g., 'Linux', 'Windows').",
			},
			"script_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the script (e.g., 'Shell', 'Python', 'PowerShell').",
				Validators: []validator.String{
					stringvalidator.OneOf("Shell", "Python", "PowerShell", "Perl", "Bat"),
				},
			},
			"script_timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The timeout for script execution in seconds.",
			},
		},
	}
}

func (r *scriptResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan scriptModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if plan.Name.IsNull() || plan.Name.IsUnknown() {
		response.Diagnostics.AddError("Missing Required Field", "The 'name' field must be provided.")
		return
	}
	if plan.ScriptContent.IsNull() || plan.ScriptContent.IsUnknown() {
		response.Diagnostics.AddError("Missing Required Field", "The 'script_content' field must be provided.")
		return
	}
	if plan.Platform.IsNull() || plan.Platform.IsUnknown() {
		response.Diagnostics.AddError("Missing Required Field", "The 'platform' field must be provided.")
		return
	}
	if plan.ScriptType.IsNull() || plan.ScriptType.IsUnknown() {
		response.Diagnostics.AddError("Missing Required Field", "The 'script_type' field must be provided.")
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	scriptTimeout := int64(300)
	if !plan.ScriptTimeout.IsNull() {
		scriptTimeout = plan.ScriptTimeout.ValueInt64()
	}
	renderParams := ""
	if !plan.RenderParams.IsNull() {
		renderParams = plan.RenderParams.ValueString()
	}

	description := ""
	if !plan.Description.IsNull() {
		description = plan.Description.ValueString()
	}

	ok, pattern := isScriptContentSafe(plan.ScriptContent.ValueString())
	if !ok {
		response.Diagnostics.AddError(
			"Dangerous script content detected",
			fmt.Sprintf("The script content contains dangerous commands and is not allowed. %s", pattern),
		)
		return
	}

	var systemTags []string

	Param := param.CreateVmInstanceScriptParam{
		BaseParam: param.BaseParam{
			SystemTags: systemTags,
		},
		Params: param.CreateVmInstanceScriptDetailParam{
			Name:          plan.Name.ValueString(),
			Description:   description,
			ScriptContent: plan.ScriptContent.ValueString(),
			Platform:      plan.Platform.ValueString(),
			ScriptType:    plan.ScriptType.ValueString(),
			ScriptTimeout: int(scriptTimeout),
			RenderParams:  renderParams,
		},
	}

	script, err := r.client.CreateVmInstanceScript(Param)
	if err != nil {
		response.Diagnostics.AddError(
			"Failed to create VM instance script",
			"Error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(script.UUID)
	plan.Name = types.StringValue(script.Name)
	plan.Description = types.StringValue(script.Description)
	plan.ScriptContent = types.StringValue(script.ScriptContent)
	plan.RenderParams = types.StringValue(script.RenderParams)
	plan.Platform = types.StringValue(script.Platform)
	plan.ScriptType = types.StringValue(script.ScriptType)
	plan.ScriptTimeout = types.Int64Value(int64(script.ScriptTimeout))

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *scriptResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state scriptModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddError("Client Not Configured", "The provider client is not properly configured.")
		return
	}

	scripts, err := r.client.GetVmInstanceScript(state.Uuid.ValueString())

	if err != nil {
		tflog.Warn(ctx, "Unable to retrieve VM instance script. It may have been deleted.", map[string]interface{}{
			"uuid":  state.Uuid.ValueString(),
			"error": err.Error(),
		})
		response.Diagnostics.Append(diags...)
		response.State.RemoveResource(ctx)
		return
	}

	state.Uuid = types.StringValue(scripts.UUID)
	state.Name = types.StringValue(scripts.Name)
	state.Description = types.StringValue(scripts.Description)
	state.ScriptContent = types.StringValue(scripts.ScriptContent)
	state.RenderParams = types.StringValue(scripts.RenderParams)
	state.Platform = types.StringValue(scripts.Platform)
	state.ScriptType = types.StringValue(scripts.ScriptType)
	state.ScriptTimeout = types.Int64Value(int64(scripts.ScriptTimeout))

	diags = response.State.Set(ctx, state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *scriptResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state scriptModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddError("Client Not Configured", "The provider client is not properly configured.")
		return
	}

	scriptTimeout := int64(300)
	if !plan.ScriptTimeout.IsNull() {
		scriptTimeout = plan.ScriptTimeout.ValueInt64()
	}

	renderParams := ""
	if !plan.RenderParams.IsNull() && !plan.RenderParams.IsUnknown() {
		renderParams = plan.RenderParams.ValueString()
	}
	plan.RenderParams = types.StringValue(renderParams)

	description := ""
	if !plan.Description.IsNull() {
		description = plan.Description.ValueString()
	}

	ok, pattern := isScriptContentSafe(plan.ScriptContent.ValueString())
	if !ok {
		response.Diagnostics.AddError(
			"Dangerous script content detected",
			fmt.Sprintf("The script content contains dangerous commands and is not allowed. %s", pattern),
		)
		return
	}

	var systemTags []string
	Param := param.UpdateVmInstanceScriptParam{
		BaseParam: param.BaseParam{
			SystemTags: systemTags,
		},
		Params: param.UpdateVmInstanceScriptDetailParam{
			Name:          plan.Name.ValueString(),
			Description:   description,
			ScriptContent: plan.ScriptContent.ValueString(),
			Platform:      plan.Platform.ValueString(),
			ScriptType:    plan.ScriptType.ValueString(),
			ScriptTimeout: int(scriptTimeout),
			RenderParams:  renderParams,
		},
	}

	tflog.Debug(ctx, "Updating VM instance script", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
		"name": plan.Name.ValueString(),
	})

	err := r.client.UpdateVmInstanceScript(state.Uuid.ValueString(), Param)
	if err != nil {
		response.Diagnostics.AddError(
			"Failed to update VM instance script",
			"Error: "+err.Error(),
		)
		return
	}

	plan.Uuid = state.Uuid

	tflog.Info(ctx, "VM instance script updated successfully", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)

}

func (r *scriptResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state scriptModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	err := r.client.DeleteVmInstanceScrpt(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError(
			"Failed to delete VM instance script",
			"Error: "+err.Error(),
		)
		return
	}

}

func (r *scriptResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing script by UUID", map[string]interface{}{
		"uuid": request.ID,
	})

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("uuid"), request.ID)...)
}

func isScriptContentSafe(script string) (bool, string) {
	dangerousPatterns := []string{
		`(?i)\brm\s+-rf\s+/`,      // rm -rf /
		`(?i)\brm\s+-[^\s]*f\s+/`, // any variant like rm -ef /
		`(?i):(){:|:&};:`,         // fork bomb
		`(?i)\bdd\s+if=`,          // dangerous disk overwrite
		`(?i)\bmkfs\.`,            // format disk
		`(?i)\bshutdown\b`,        // shutdown command
		`(?i)\breboot\b`,          // reboot
		`(?i)\bpoweroff\b`,
		`(?i)\bkill\s+-9\s+1\b`, // kill init
	}

	for _, pattern := range dangerousPatterns {
		match, _ := regexp.MatchString(pattern, script)
		if match {
			return false, pattern
		}
	}
	return true, ""
}
