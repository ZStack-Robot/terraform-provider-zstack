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
	_ resource.Resource                = &directoryResource{}
	_ resource.ResourceWithConfigure   = &directoryResource{}
	_ resource.ResourceWithImportState = &directoryResource{}
)

type directoryResource struct {
	client *client.ZSClient
}

type directoryModel struct {
	Uuid              types.String `tfsdk:"uuid"`
	Name              types.String `tfsdk:"name"`
	ParentUuid        types.String `tfsdk:"parent_uuid"`
	ZoneUuid          types.String `tfsdk:"zone_uuid"`
	Type              types.String `tfsdk:"type"`
	GroupName         types.String `tfsdk:"group_name"`
	RootDirectoryUuid types.String `tfsdk:"root_directory_uuid"`
}

func DirectoryResource() resource.Resource {
	return &directoryResource{}
}

func (r *directoryResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *directoryResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_directory"
}

func (r *directoryResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage directories in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the directory.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the directory.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"parent_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The parent directory UUID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"zone_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The zone UUID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the directory.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"group_name": schema.StringAttribute{
				Computed:    true,
				Description: "The group name of the directory.",
			},
			"root_directory_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The root directory UUID.",
			},
		},
	}
}

func (r *directoryResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan directoryModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.CreateDirectoryParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateDirectoryParamDetail{
			Name:       plan.Name.ValueString(),
			ParentUuid: stringPtrOrNil(plan.ParentUuid.ValueString()),
			ZoneUuid:   plan.ZoneUuid.ValueString(),
			Type:       plan.Type.ValueString(),
		},
	}

	directory, err := r.client.CreateDirectory(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Directory",
			"Could not create directory, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(directory.UUID)
	plan.Name = types.StringValue(directory.Name)
	plan.ParentUuid = stringValueOrNull(directory.ParentUuid)
	plan.ZoneUuid = types.StringValue(directory.ZoneUuid)
	plan.Type = types.StringValue(directory.Type)
	plan.GroupName = stringValueOrNull(directory.GroupName)
	plan.RootDirectoryUuid = stringValueOrNull(directory.RootDirectoryUuid)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *directoryResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state directoryModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	directory, err := findResourceByQuery(r.client.QueryDirectory, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading Directory",
			"Could not read directory UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(directory.UUID)
	state.Name = types.StringValue(directory.Name)
	state.ParentUuid = stringValueOrNull(directory.ParentUuid)
	state.ZoneUuid = types.StringValue(directory.ZoneUuid)
	state.Type = types.StringValue(directory.Type)
	state.GroupName = stringValueOrNull(directory.GroupName)
	state.RootDirectoryUuid = stringValueOrNull(directory.RootDirectoryUuid)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (r *directoryResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan directoryModel
	var state directoryModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateDirectoryParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateDirectoryParamDetail{
			Uuid: state.Uuid.ValueString(),
			Name: plan.Name.ValueString(),
		},
	}

	directory, err := r.client.UpdateDirectory(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating Directory",
			"Could not update directory, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(directory.UUID)
	plan.Name = types.StringValue(directory.Name)
	plan.ParentUuid = stringValueOrNull(directory.ParentUuid)
	plan.ZoneUuid = types.StringValue(directory.ZoneUuid)
	plan.Type = types.StringValue(directory.Type)
	plan.GroupName = stringValueOrNull(directory.GroupName)
	plan.RootDirectoryUuid = stringValueOrNull(directory.RootDirectoryUuid)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *directoryResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state directoryModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}


	err := r.client.DeleteDirectory(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting Directory",
			"Could not delete directory, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *directoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
