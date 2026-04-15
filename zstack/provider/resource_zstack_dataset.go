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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &datasetResource{}
	_ resource.ResourceWithConfigure   = &datasetResource{}
	_ resource.ResourceWithImportState = &datasetResource{}
)

type datasetResource struct {
	client *client.ZSClient
}

type datasetModel struct {
	Uuid            types.String `tfsdk:"uuid"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Url             types.String `tfsdk:"url"`
	ModelCenterUuid types.String `tfsdk:"model_center_uuid"`
	Token           types.String `tfsdk:"token"`
	System          types.Bool   `tfsdk:"system"`
	InstallPath     types.String `tfsdk:"install_path"`
	Size            types.Int64  `tfsdk:"size"`
}

func DatasetResource() resource.Resource {
	return &datasetResource{}
}

func (r *datasetResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *datasetResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_dataset"
}

func (r *datasetResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage datasets in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"url": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"model_center_uuid": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"system": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"install_path": schema.StringAttribute{
				Computed: true,
			},
			"size": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func (r *datasetResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan datasetModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.CreateDatasetParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateDatasetParamDetail{
			Name:            plan.Name.ValueString(),
			Description:     stringPtrOrNil(plan.Description.ValueString()),
			Url:             plan.Url.ValueString(),
			ModelCenterUuid: plan.ModelCenterUuid.ValueString(),
		},
	}

	if !plan.Token.IsNull() && !plan.Token.IsUnknown() {
		p.Params.Token = stringPtrOrNil(plan.Token.ValueString())
	}

	if !plan.System.IsNull() && !plan.System.IsUnknown() {
		p.Params.System = boolPtr(plan.System.ValueBool())
	}

	item, err := r.client.CreateDataset(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Dataset",
			"Could not create dataset, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(item.UUID)
	plan.Name = types.StringValue(item.Name)
	plan.Description = stringValueOrNull(item.Description)
	plan.Url = types.StringValue(item.Url)
	plan.ModelCenterUuid = types.StringValue(item.ModelCenterUuid)
	plan.System = types.BoolValue(item.System)
	plan.InstallPath = stringValueOrNull(item.InstallPath)
	plan.Size = types.Int64Value(item.Size)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *datasetResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state datasetModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	item, err := findResourceByQuery(r.client.QueryDataset, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading Dataset",
			"Could not read dataset UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(item.UUID)
	state.Name = types.StringValue(item.Name)
	state.Description = stringValueOrNull(item.Description)
	state.Url = types.StringValue(item.Url)
	state.ModelCenterUuid = types.StringValue(item.ModelCenterUuid)
	state.System = types.BoolValue(item.System)
	state.InstallPath = stringValueOrNull(item.InstallPath)
	state.Size = types.Int64Value(item.Size)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (r *datasetResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan datasetModel
	var state datasetModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateDatasetParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateDatasetParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	item, err := r.client.UpdateDataset(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating Dataset",
			"Could not update dataset, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(item.UUID)
	plan.Name = types.StringValue(item.Name)
	plan.Description = stringValueOrNull(item.Description)
	plan.Url = types.StringValue(item.Url)
	plan.ModelCenterUuid = types.StringValue(item.ModelCenterUuid)
	plan.System = types.BoolValue(item.System)
	plan.InstallPath = stringValueOrNull(item.InstallPath)
	plan.Size = types.Int64Value(item.Size)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *datasetResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state datasetModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}


	err := r.client.DeleteDataset(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting Dataset",
			"Could not delete dataset, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *datasetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
