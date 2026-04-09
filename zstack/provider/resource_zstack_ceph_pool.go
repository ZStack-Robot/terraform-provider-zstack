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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &cephPoolResource{}
	_ resource.ResourceWithConfigure   = &cephPoolResource{}
	_ resource.ResourceWithImportState = &cephPoolResource{}
)

type cephPoolResource struct {
	client *client.ZSClient
}

type cephPoolModel struct {
	Uuid               types.String `tfsdk:"uuid"`
	PoolName           types.String `tfsdk:"pool_name"`
	AliasName          types.String `tfsdk:"alias_name"`
	Description        types.String `tfsdk:"description"`
	Type               types.String `tfsdk:"type"`
	IsCreate           types.Bool   `tfsdk:"is_create"`
	PrimaryStorageUuid types.String `tfsdk:"primary_storage_uuid"`
}

func CephPoolResource() resource.Resource {
	return &cephPoolResource{}
}

func (r *cephPoolResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *cephPoolResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ceph_pool"
}

func (r *cephPoolResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage Ceph primary storage pools in ZStack. " +
			"A Ceph pool stores data for volumes or images in Ceph primary storage.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the Ceph pool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pool_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Ceph pool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"alias_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "An alias name for the Ceph pool.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the Ceph pool.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the pool (e.g., Data, Root).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("Data", "Root"),
				},
			},
			"is_create": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to create the pool on the Ceph cluster (true) or just register an existing pool (false). Only used during creation.",
			},
			"primary_storage_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the Ceph primary storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *cephPoolResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan cephPoolModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.AddCephPrimaryStoragePoolParam{
		BaseParam: param.BaseParam{},
		Params: param.AddCephPrimaryStoragePoolParamDetail{
			PoolName:    plan.PoolName.ValueString(),
			AliasName:   stringPtrOrNil(plan.AliasName.ValueString()),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			Type:        plan.Type.ValueString(),
		},
	}

	if !plan.IsCreate.IsNull() && !plan.IsCreate.IsUnknown() {
		p.Params.IsCreate = plan.IsCreate.ValueBool()
	}

	result, err := r.client.AddCephPrimaryStoragePool(plan.PrimaryStorageUuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Ceph Pool",
			"Could not create ceph pool, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.PoolName = types.StringValue(result.PoolName)
	plan.AliasName = stringValueOrNull(result.AliasName)
	plan.Description = stringValueOrNull(result.Description)
	plan.Type = types.StringValue(result.Type)
	plan.PrimaryStorageUuid = types.StringValue(result.PrimaryStorageUuid)
	plan.IsCreate = types.BoolNull()

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *cephPoolResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state cephPoolModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	cephPool, err := findResourceByQuery(r.client.QueryCephPrimaryStoragePool, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading Ceph Pool",
			"Could not read ceph pool UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(cephPool.UUID)
	state.PoolName = types.StringValue(cephPool.PoolName)
	state.AliasName = stringValueOrNull(cephPool.AliasName)
	state.Description = stringValueOrNull(cephPool.Description)
	state.Type = types.StringValue(cephPool.Type)
	state.PrimaryStorageUuid = types.StringValue(cephPool.PrimaryStorageUuid)
	state.IsCreate = types.BoolNull()

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *cephPoolResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan cephPoolModel
	var state cephPoolModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.UpdateCephPrimaryStoragePoolParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateCephPrimaryStoragePoolParamDetail{
			AliasName:   stringPtrOrNil(plan.AliasName.ValueString()),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	result, err := r.client.UpdateCephPrimaryStoragePool(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating Ceph Pool",
			"Could not update ceph pool, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.PoolName = types.StringValue(result.PoolName)
	plan.AliasName = stringValueOrNull(result.AliasName)
	plan.Description = stringValueOrNull(result.Description)
	plan.Type = types.StringValue(result.Type)
	plan.PrimaryStorageUuid = types.StringValue(result.PrimaryStorageUuid)
	plan.IsCreate = types.BoolNull()

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *cephPoolResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state cephPoolModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Ceph primary storage pool UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteCephPrimaryStoragePool(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError("Error deleting Ceph Pool", "Could not delete ceph pool, unexpected error: "+err.Error())
		return
	}
}

func (r *cephPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
