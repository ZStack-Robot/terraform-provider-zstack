// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &cephPrimaryStorageResource{}
	_ resource.ResourceWithConfigure   = &cephPrimaryStorageResource{}
	_ resource.ResourceWithImportState = &cephPrimaryStorageResource{}
)

type cephPrimaryStorageResource struct {
	client *client.ZSClient
}

type cephPrimaryStorageModel struct {
	Uuid               types.String `tfsdk:"uuid"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	ZoneUuid           types.String `tfsdk:"zone_uuid"`
	MonUrls            types.List   `tfsdk:"mon_urls"`
	RootVolumePoolName types.String `tfsdk:"root_volume_pool_name"`
	DataVolumePoolName types.String `tfsdk:"data_volume_pool_name"`
	ImageCachePoolName types.String `tfsdk:"image_cache_pool_name"`
	Url                types.String `tfsdk:"url"`
	Type               types.String `tfsdk:"type"`
	State              types.String `tfsdk:"state"`
	Status             types.String `tfsdk:"status"`
	TotalCapacity      types.Int64  `tfsdk:"total_capacity"`
	AvailableCapacity  types.Int64  `tfsdk:"available_capacity"`
}

func CephPrimaryStorageResource() resource.Resource {
	return &cephPrimaryStorageResource{}
}

func (r *cephPrimaryStorageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cephPrimaryStorageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ceph_primary_storage"
}

func (r *cephPrimaryStorageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZStack Ceph Primary Storage resource.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the Ceph primary storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Ceph primary storage.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the Ceph primary storage.",
			},
			"zone_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the zone where the Ceph primary storage will be created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mon_urls": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of Ceph monitor URLs.",
			},
			"root_volume_pool_name": schema.StringAttribute{
				Optional:    true,
				Description: "The Ceph pool name for root volumes.",
			},
			"data_volume_pool_name": schema.StringAttribute{
				Optional:    true,
				Description: "The Ceph pool name for data volumes.",
			},
			"image_cache_pool_name": schema.StringAttribute{
				Optional:    true,
				Description: "The Ceph pool name for image cache.",
			},
			"url": schema.StringAttribute{
				Computed:    true,
				Description: "The URL of the primary storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the primary storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the primary storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the primary storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"total_capacity": schema.Int64Attribute{
				Computed:    true,
				Description: "Total capacity in bytes.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"available_capacity": schema.Int64Attribute{
				Computed:    true,
				Description: "Available capacity in bytes.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *cephPrimaryStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cephPrimaryStorageModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	monUrls := listToStringSlice(plan.MonUrls)

	createParam := param.AddCephPrimaryStorageParam{
		BaseParam: param.BaseParam{},
		Params: param.AddCephPrimaryStorageParamDetail{
			Name:               plan.Name.ValueString(),
			ZoneUuid:           plan.ZoneUuid.ValueString(),
			MonUrls:            monUrls,
			RootVolumePoolName: stringPtrOrNil(plan.RootVolumePoolName.ValueString()),
			DataVolumePoolName: stringPtrOrNil(plan.DataVolumePoolName.ValueString()),
			ImageCachePoolName: stringPtrOrNil(plan.ImageCachePoolName.ValueString()),
			Description:        stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	tflog.Debug(ctx, "Creating Ceph primary storage", map[string]interface{}{
		"name": createParam.Params.Name,
	})

	storage, err := r.client.AddCephPrimaryStorage(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Ceph primary storage",
			"Could not create Ceph primary storage, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(storage.UUID)
	plan.Url = stringValueOrNull(storage.Url)
	plan.Type = stringValueOrNull(storage.Type)
	plan.State = stringValueOrNull(storage.State)
	plan.Status = stringValueOrNull(storage.Status)
	plan.TotalCapacity = types.Int64Value(storage.TotalCapacity)
	plan.AvailableCapacity = types.Int64Value(storage.AvailableCapacity)

	if storage.Description != "" {
		plan.Description = types.StringValue(storage.Description)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Created Ceph primary storage", map[string]interface{}{
		"uuid": storage.UUID,
	})
}

func (r *cephPrimaryStorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cephPrimaryStorageModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	queryParam.AddQ("uuid=" + state.Uuid.ValueString())

	tflog.Debug(ctx, "Reading Ceph primary storage", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})

	storages, err := r.client.QueryCephPrimaryStorage(&queryParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Ceph primary storage",
			"Could not read Ceph primary storage UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	if len(storages) == 0 {
		tflog.Warn(ctx, "Ceph primary storage not found, removing from state", map[string]interface{}{
			"uuid": state.Uuid.ValueString(),
		})
		state.Uuid = types.StringValue("")
		resp.State.Set(ctx, &state)
		return
	}

	storage := storages[0]

	state.Uuid = types.StringValue(storage.UUID)
	state.Name = types.StringValue(storage.Name)
	state.Description = stringValueOrNull(storage.Description)
	state.ZoneUuid = types.StringValue(storage.ZoneUuid)
	state.Url = stringValueOrNull(storage.Url)
	state.Type = stringValueOrNull(storage.Type)
	state.State = stringValueOrNull(storage.State)
	state.Status = stringValueOrNull(storage.Status)
	state.TotalCapacity = types.Int64Value(storage.TotalCapacity)
	state.AvailableCapacity = types.Int64Value(storage.AvailableCapacity)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *cephPrimaryStorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan cephPrimaryStorageModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Update called for Ceph primary storage (no-op)", map[string]interface{}{
		"uuid": plan.Uuid.ValueString(),
	})

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *cephPrimaryStorageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cephPrimaryStorageModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Ceph primary storage", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})

	err := r.client.DeletePrimaryStorage(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Ceph primary storage",
			"Could not delete Ceph primary storage UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "Deleted Ceph primary storage", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
}

func (r *cephPrimaryStorageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
