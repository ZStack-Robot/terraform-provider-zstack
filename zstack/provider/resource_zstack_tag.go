// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &tagResource{}
	_ resource.ResourceWithConfigure   = &tagResource{}
	_ resource.ResourceWithImportState = &tagResource{}
)

type tagResource struct {
	client *client.ZSClient
}

type tagResourceModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Value       types.String `tfsdk:"value"`
	Description types.String `tfsdk:"description"`
	Color       types.String `tfsdk:"color"`
	Type        types.String `tfsdk:"type"` //  validValues = {"simple", "withToken"}
}

func TagResource() resource.Resource {
	return &tagResource{}
}

func (r *tagResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *tagResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_tag"
}

func (r *tagResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage tags in ZStack. Tags can be used to categorize and manage resources within ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the tag.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the tag.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the tag.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"value": schema.StringAttribute{
				Required:    true,
				Description: "The value of the tag.",
			},
			"color": schema.StringAttribute{
				Optional:    true,
				Description: "The color associated with the tag, used for visual categorization.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "The type of the tag. Valid values are 'simple' and 'withToken'.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("simple", "withToken"),
				},
			},
		},
	}
}

func (r *tagResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan tagResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	params := param.CreateTagParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateTagParamDetail{
			Name:        plan.Name.ValueString(),
			Value:       plan.Value.ValueString(),
			Description: stringPtr(plan.Description.ValueString()),
			Color:       stringPtr(plan.Color.ValueString()),
			Type:        stringPtr(plan.Type.ValueString()),
		},
	}

	tag, err := r.client.CreateTag(params)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Tag",
			"Could not create tag, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(tag.UUID)
	plan.Name = types.StringValue(tag.Name)
	plan.Description = stringValueOrNull(tag.Description)
	plan.Color = stringValueOrNull(tag.Color)
	plan.Type = stringValueOrNull(tag.Type)

	diags = response.State.Set(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}
func (r *tagResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state tagResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tag, err := findResourceByGet(r.client.GetTag, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading Tag",
			"Could not read tag UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(tag.Name)
	state.Value = stringValueOrNull(tag.Value)
	state.Color = stringValueOrNull(tag.Color)
	state.Type = stringValueOrNull(tag.Type)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}
func (r *tagResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state tagResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var plan tagResourceModel
	diags = request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	//uuid := ("uuid".ValueString())
	if state.Uuid.ValueString() == "" {
		response.Diagnostics.AddError(
			"Error updating Tag",
			"Could not update tag without a UUID.",
		)
		return
	}

	if state.Type.ValueString() == "withToken" {
		oldValue := state.Value.ValueString()
		newValue := plan.Value.ValueString()

		oldParts := strings.SplitN(oldValue, "::", 2)
		newParts := strings.SplitN(newValue, "::", 2)

		if len(oldParts) != 2 || len(newParts) != 2 {
			response.Diagnostics.AddError(
				"Error updating Tag",
				fmt.Sprintf("Could not update tag, expected value format 'token::key'. Got old='%s' new='%s'", oldValue, newValue),
			)
			return
		}

		if oldParts[0] != newParts[0] {
			response.Diagnostics.AddError(
				"Error updating Tag",
				fmt.Sprintf("Could not update tag because token part '%s' cannot be modified to '%s'. Only the key part after '::' can be changed.", oldParts[0], newParts[0]),
			)
			return
		}
	}

	params := param.UpdateTagParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateTagParamDetail{
			Name:        plan.Name.ValueString(),
			Value:       stringPtr(plan.Value.ValueString()),
			Description: stringPtr(plan.Description.ValueString()),
			Color:       stringPtr(plan.Color.ValueString()),
		},
	}

	tag, err := r.client.UpdateTag(state.Uuid.ValueString(), params)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating Tag",
			"Could not update tag, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(tag.UUID)
	plan.Name = types.StringValue(tag.Name)
	plan.Description = stringValueOrNull(tag.Description)
	plan.Color = stringValueOrNull(tag.Color)
	plan.Type = stringValueOrNull(tag.Type)

	diags = response.State.Set(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}
func (r *tagResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state tagResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid.ValueString() == "" {
		response.Diagnostics.AddWarning("Tag UUID is empty", "Nothing to delete, skipping.")
		return
	}

	err := r.client.DeleteTag(state.Uuid.ValueString(), param.DeleteModeEnforcing)
	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting Tag",
			"Could not delete tag, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *tagResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	tagUUID := request.ID

	if tagUUID == "" {
		response.Diagnostics.AddError("Missing Tag UUID", "Import requires a valid tag UUID as the ID.")
		return
	}

	tag, err := r.client.GetTag(tagUUID)
	if err != nil {
		response.Diagnostics.AddError(
			"Tag Not Found",
			fmt.Sprintf("Could not find tag with UUID '%s': '%v'", tagUUID, err),
		)
		return
	}

	response.State.Set(ctx, &tagResourceModel{
		Uuid:        types.StringValue(tag.UUID),
		Name:        types.StringValue(tag.Name),
		Value:       stringValueOrNull(tag.Value),
		Description: stringValueOrNull(tag.Description),
		Color:       stringValueOrNull(tag.Color),
		Type:        stringValueOrNull(tag.Type),
	})
}
