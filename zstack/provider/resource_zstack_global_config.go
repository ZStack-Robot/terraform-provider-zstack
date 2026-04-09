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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

var (
	_ resource.Resource                = &globalConfigResource{}
	_ resource.ResourceWithConfigure   = &globalConfigResource{}
	_ resource.ResourceWithImportState = &globalConfigResource{}
)

type globalConfigResource struct {
	client *client.ZSClient
}

type globalConfigModel struct {
	Name         types.String `tfsdk:"name"`
	Category     types.String `tfsdk:"category"`
	Value        types.String `tfsdk:"value"`
	DefaultValue types.String `tfsdk:"default_value"`
	Description  types.String `tfsdk:"description"`
}

func GlobalConfigResource() resource.Resource {
	return &globalConfigResource{}
}

func (r *globalConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_global_config"
}

func (r *globalConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZStack Global Config resource.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the global config (acts as part of the ID).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"category": schema.StringAttribute{
				Description: "The category of the global config (acts as part of the ID).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				Description: "The value to set for this config.",
				Required:    true,
			},
			"default_value": schema.StringAttribute{
				Description: "The default value of the config.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the config.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *globalConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *globalConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan globalConfigModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateParam := param.UpdateGlobalConfigParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateGlobalConfigParamDetail{
			Value: stringPtr(plan.Value.ValueString()),
		},
	}

	result, err := r.client.UpdateGlobalConfig(plan.Category.ValueString(), plan.Name.ValueString(), updateParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Global Config",
			"Could not create global config, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Name = types.StringValue(result.Name)
	plan.Category = types.StringValue(result.Category)
	plan.Value = types.StringValue(result.Value)
	plan.DefaultValue = stringValueOrNull(result.DefaultValue)
	plan.Description = stringValueOrNull(result.Description)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Global config set", map[string]interface{}{
		"category": result.Category,
		"name":     result.Name,
	})
}

func (r *globalConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state globalConfigModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := findResourceByQuery(func(queryParam *param.QueryParam) ([]view.GlobalConfigInventoryView, error) {
		*queryParam = param.NewQueryParam()
		queryParam.AddQ("category=" + state.Category.ValueString())
		queryParam.AddQ("name=" + state.Name.ValueString())
		return r.client.QueryGlobalConfig(queryParam)
	}, state.Name.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			tflog.Warn(ctx, "Global config not found", map[string]interface{}{
				"category": state.Category.ValueString(),
				"name":     state.Name.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Global Config",
			"Could not read global config, unexpected error: "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(config.Name)
	state.Category = types.StringValue(config.Category)
	state.Value = types.StringValue(config.Value)
	state.DefaultValue = stringValueOrNull(config.DefaultValue)
	state.Description = stringValueOrNull(config.Description)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *globalConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan globalConfigModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateParam := param.UpdateGlobalConfigParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateGlobalConfigParamDetail{
			Value: stringPtr(plan.Value.ValueString()),
		},
	}

	result, err := r.client.UpdateGlobalConfig(plan.Category.ValueString(), plan.Name.ValueString(), updateParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Global Config",
			"Could not update global config, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Name = types.StringValue(result.Name)
	plan.Category = types.StringValue(result.Category)
	plan.Value = types.StringValue(result.Value)
	plan.DefaultValue = stringValueOrNull(result.DefaultValue)
	plan.Description = stringValueOrNull(result.Description)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Global config updated", map[string]interface{}{
		"category": result.Category,
		"name":     result.Name,
	})
}

func (r *globalConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state globalConfigModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Reset to default value by updating with the default value
	updateParam := param.UpdateGlobalConfigParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateGlobalConfigParamDetail{
			Value: stringPtr(state.DefaultValue.ValueString()),
		},
	}

	_, err := r.client.UpdateGlobalConfig(state.Category.ValueString(), state.Name.ValueString(), updateParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Global Config",
			"Could not delete global config, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "Global config reset to default", map[string]interface{}{
		"category": state.Category.ValueString(),
		"name":     state.Name.ValueString(),
	})
}

func (r *globalConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: category/name")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("category"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[1])...)
}
