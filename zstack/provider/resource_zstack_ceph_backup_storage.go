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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &cephBackupStorageResource{}
	_ resource.ResourceWithConfigure   = &cephBackupStorageResource{}
	_ resource.ResourceWithImportState = &cephBackupStorageResource{}
)

type cephBackupStorageResource struct {
	client *client.ZSClient
}

type cephBackupStorageModel struct {
	Uuid              types.String `tfsdk:"uuid"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	MonUrls           types.List   `tfsdk:"mon_urls"`
	PoolName          types.String `tfsdk:"pool_name"`
	Url               types.String `tfsdk:"url"`
	Type              types.String `tfsdk:"type"`
	State             types.String `tfsdk:"state"`
	Status            types.String `tfsdk:"status"`
	TotalCapacity     types.Int64  `tfsdk:"total_capacity"`
	AvailableCapacity types.Int64  `tfsdk:"available_capacity"`
}

func CephBackupStorageResource() resource.Resource {
	return &cephBackupStorageResource{}
}

func (r *cephBackupStorageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cephBackupStorageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ceph_backup_storage"
}

func (r *cephBackupStorageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZStack Ceph Backup Storage resource.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the Ceph backup storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Ceph backup storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the Ceph backup storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mon_urls": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of Ceph monitor URLs.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"pool_name": schema.StringAttribute{
				Optional:    true,
				Description: "The Ceph pool name for backup storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Computed:    true,
				Description: "The URL of the backup storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the backup storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the backup storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the backup storage.",
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

func (r *cephBackupStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cephBackupStorageModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	monUrls := listToStringSlice(plan.MonUrls)

	createParam := param.AddCephBackupStorageParam{
		BaseParam: param.BaseParam{},
		Params: param.AddCephBackupStorageParamDetail{
			Name:        plan.Name.ValueString(),
			MonUrls:     monUrls,
			PoolName:    stringPtrOrNil(plan.PoolName.ValueString()),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	tflog.Debug(ctx, "Creating Ceph backup storage", map[string]interface{}{
		"name": createParam.Params.Name,
	})

	storage, err := r.client.AddCephBackupStorage(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Ceph Backup Storage",
			"Could not create ceph backup storage, unexpected error: "+err.Error(),
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

	tflog.Debug(ctx, "Created Ceph backup storage", map[string]interface{}{
		"uuid": storage.UUID,
	})
}

func (r *cephBackupStorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cephBackupStorageModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Ceph backup storage", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})

	storage, err := findResourceByQuery(r.client.QueryCephBackupStorage, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Ceph Backup Storage",
			"Could not read ceph backup storage UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(storage.UUID)
	state.Name = types.StringValue(storage.Name)
	state.Description = stringValueOrNull(storage.Description)
	state.Url = stringValueOrNull(storage.Url)
	state.Type = stringValueOrNull(storage.Type)
	state.State = stringValueOrNull(storage.State)
	state.Status = stringValueOrNull(storage.Status)
	state.TotalCapacity = types.Int64Value(storage.TotalCapacity)
	state.AvailableCapacity = types.Int64Value(storage.AvailableCapacity)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *cephBackupStorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan cephBackupStorageModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Update called for Ceph backup storage (no-op)", map[string]interface{}{
		"uuid": plan.Uuid.ValueString(),
	})

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *cephBackupStorageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cephBackupStorageModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Ceph backup storage", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})

	err := r.client.DeleteBackupStorage(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Ceph Backup Storage",
			"Could not delete ceph backup storage UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "Deleted Ceph backup storage", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
}

func (r *cephBackupStorageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
