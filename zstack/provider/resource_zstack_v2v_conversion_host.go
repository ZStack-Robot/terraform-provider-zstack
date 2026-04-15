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
	_ resource.Resource                = &v2vConversionHostResource{}
	_ resource.ResourceWithConfigure   = &v2vConversionHostResource{}
	_ resource.ResourceWithImportState = &v2vConversionHostResource{}
)

type v2vConversionHostResource struct {
	client *client.ZSClient
}

type v2vConversionHostModel struct {
	Uuid          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Type          types.String `tfsdk:"type"`
	HostUuid      types.String `tfsdk:"host_uuid"`
	StoragePath   types.String `tfsdk:"storage_path"`
	State         types.String `tfsdk:"state"`
	TotalSize     types.Int64  `tfsdk:"total_size"`
	AvailableSize types.Int64  `tfsdk:"available_size"`
	HostStatus    types.String `tfsdk:"host_status"`
	HostState     types.String `tfsdk:"host_state"`
}

func V2VConversionHostResource() resource.Resource { return &v2vConversionHostResource{} }

func (r *v2vConversionHostResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *v2vConversionHostResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_v2v_conversion_host"
}

func (r *v2vConversionHostResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage V2V conversion hosts in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the V2V conversion host.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the V2V conversion host.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the V2V conversion host.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the V2V conversion host.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"host_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The host UUID associated with this V2V conversion host.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"storage_path": schema.StringAttribute{
				Required:    true,
				Description: "The storage path of the V2V conversion host.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the V2V conversion host.",
			},
			"total_size": schema.Int64Attribute{
				Computed:    true,
				Description: "The total size in bytes.",
			},
			"available_size": schema.Int64Attribute{
				Computed:    true,
				Description: "The available size in bytes.",
			},
			"host_status": schema.StringAttribute{
				Computed:    true,
				Description: "The host status.",
			},
			"host_state": schema.StringAttribute{
				Computed:    true,
				Description: "The host state.",
			},
		},
	}
}

func (r *v2vConversionHostResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan v2vConversionHostModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.AddV2VConversionHostParam{
		BaseParam: param.BaseParam{},
		Params: param.AddV2VConversionHostParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			Type:        plan.Type.ValueString(),
			HostUuid:    plan.HostUuid.ValueString(),
			StoragePath: plan.StoragePath.ValueString(),
		},
	}

	item, err := r.client.AddV2VConversionHost(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating V2V Conversion Host",
			"Could not create v2v conversion host, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(item.UUID)
	plan.Name = types.StringValue(item.Name)
	plan.Description = stringValueOrNull(item.Description)
	plan.Type = types.StringValue(item.Type)
	plan.HostUuid = types.StringValue(item.HostUuid)
	plan.StoragePath = types.StringValue(item.StoragePath)
	plan.State = stringValueOrNull(item.State)
	plan.TotalSize = types.Int64Value(item.TotalSize)
	plan.AvailableSize = types.Int64Value(item.AvailableSize)
	plan.HostStatus = stringValueOrNull(item.HostStatus)
	plan.HostState = stringValueOrNull(item.HostState)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *v2vConversionHostResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state v2vConversionHostModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	item, err := findResourceByQuery(r.client.QueryV2VConversionHost, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading V2V Conversion Host",
			"Could not read V2V Conversion Host, unexpected error: "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(item.UUID)
	state.Name = types.StringValue(item.Name)
	state.Description = stringValueOrNull(item.Description)
	state.Type = types.StringValue(item.Type)
	state.HostUuid = types.StringValue(item.HostUuid)
	state.StoragePath = types.StringValue(item.StoragePath)
	state.State = stringValueOrNull(item.State)
	state.TotalSize = types.Int64Value(item.TotalSize)
	state.AvailableSize = types.Int64Value(item.AvailableSize)
	state.HostStatus = stringValueOrNull(item.HostStatus)
	state.HostState = stringValueOrNull(item.HostState)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (r *v2vConversionHostResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan v2vConversionHostModel
	var state v2vConversionHostModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateV2VConversionHostParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateV2VConversionHostParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			StoragePath: stringPtrOrNil(plan.StoragePath.ValueString()),
		},
	}

	item, err := r.client.UpdateV2VConversionHost(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating V2V Conversion Host",
			"Could not update v2v conversion host, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(item.UUID)
	plan.Name = types.StringValue(item.Name)
	plan.Description = stringValueOrNull(item.Description)
	plan.Type = types.StringValue(item.Type)
	plan.HostUuid = types.StringValue(item.HostUuid)
	plan.StoragePath = types.StringValue(item.StoragePath)
	plan.State = stringValueOrNull(item.State)
	plan.TotalSize = types.Int64Value(item.TotalSize)
	plan.AvailableSize = types.Int64Value(item.AvailableSize)
	plan.HostStatus = stringValueOrNull(item.HostStatus)
	plan.HostState = stringValueOrNull(item.HostState)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *v2vConversionHostResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state v2vConversionHostModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}


	err := r.client.DeleteV2VConversionHost(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting V2V Conversion Host",
			"Could not delete v2v conversion host, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *v2vConversionHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
