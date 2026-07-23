// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

const maxVni = 16777215

var (
	_ resource.Resource                   = &vniRangeResource{}
	_ resource.ResourceWithConfigure      = &vniRangeResource{}
	_ resource.ResourceWithImportState    = &vniRangeResource{}
	_ resource.ResourceWithValidateConfig = &vniRangeResource{}
)

type vniRangeResource struct {
	client *client.ZSClient
}

type vniRangeResourceModel struct {
	Uuid         types.String `tfsdk:"uuid"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	StartVni     types.Int64  `tfsdk:"start_vni"`
	EndVni       types.Int64  `tfsdk:"end_vni"`
	PoolUuid     types.String `tfsdk:"pool_uuid"`
	ResourceUuid types.String `tfsdk:"resource_uuid"`
	TagUuids     types.List   `tfsdk:"tag_uuids"`
	SystemTags   types.List   `tfsdk:"system_tags"`
}

func VniRangeResource() resource.Resource {
	return &vniRangeResource{}
}

func (r *vniRangeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	configuredClient, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer.", req.ProviderData),
		)
		return
	}
	r.client = configuredClient
}

func (r *vniRangeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vni_range"
}

func (r *vniRangeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a VNI range in a ZStack L2 VXLAN network pool.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the VNI range.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VNI range.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the VNI range. Changing it recreates the range because the ZStack update API only accepts a name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"start_vni": schema.Int64Attribute{
				Required:    true,
				Description: "The first VNI in the range.",
				Validators: []validator.Int64{
					int64validator.Between(1, maxVni),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"end_vni": schema.Int64Attribute{
				Required:    true,
				Description: "The last VNI in the range.",
				Validators: []validator.Int64{
					int64validator.Between(1, maxVni),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"pool_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the L2 VXLAN network pool that owns the range.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A custom UUID requested when creating the VNI range.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tag_uuids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Tag UUIDs attached while creating the VNI range.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIfConfigured(),
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"system_tags": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "System tags sent while creating the VNI range.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIfConfigured(),
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *vniRangeResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config vniRangeResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() || config.StartVni.IsNull() || config.StartVni.IsUnknown() || config.EndVni.IsNull() || config.EndVni.IsUnknown() {
		return
	}

	if config.StartVni.ValueInt64() > config.EndVni.ValueInt64() {
		resp.Diagnostics.AddAttributeError(
			path.Root("start_vni"),
			"Invalid VNI range",
			"start_vni must be less than or equal to end_vni.",
		)
	}
}

func (r *vniRangeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vniRangeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.StartVni.ValueInt64() > plan.EndVni.ValueInt64() {
		resp.Diagnostics.AddError("Invalid VNI range", "start_vni must be less than or equal to end_vni.")
		return
	}

	createParam := vniRangeCreateParam(plan)

	vniRange, err := r.client.CreateVniRange(plan.PoolUuid.ValueString(), createParam)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VNI range", err.Error())
		return
	}

	state := vniRangeModelFromView(vniRange, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *vniRangeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vniRangeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vniRange, err := findResourceByGet(r.client.GetVniRange, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading VNI range", err.Error())
		return
	}

	refreshed := vniRangeModelFromView(vniRange, state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &refreshed)...)
}

func (r *vniRangeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan vniRangeResourceModel
	var state vniRangeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Name.ValueString() != state.Name.ValueString() {
		_, err := r.client.UpdateVniRange(state.Uuid.ValueString(), param.UpdateVniRangeParam{
			BaseParam: param.BaseParam{},
			Params: param.UpdateVniRangeParamDetail{
				Name: plan.Name.ValueString(),
			},
		})
		if err != nil {
			resp.Diagnostics.AddError("Error updating VNI range", err.Error())
			return
		}
	}

	vniRange, err := findResourceByGet(r.client.GetVniRange, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading VNI range after update", err.Error())
		return
	}

	refreshed := vniRangeModelFromView(vniRange, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &refreshed)...)
}

func (r *vniRangeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vniRangeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteVniRange(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil && !isZStackNotFoundError(err) {
		resp.Diagnostics.AddError("Error deleting VNI range", err.Error())
	}
}

func (r *vniRangeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func vniRangeModelFromView(vniRange *view.VniRangeInventoryView, prior vniRangeResourceModel) vniRangeResourceModel {
	resourceUuid := prior.ResourceUuid
	if resourceUuid.IsUnknown() {
		resourceUuid = types.StringNull()
	}
	tagUuids := prior.TagUuids
	if tagUuids.IsUnknown() {
		tagUuids = types.ListNull(types.StringType)
	}
	systemTags := prior.SystemTags
	if systemTags.IsUnknown() {
		systemTags = types.ListNull(types.StringType)
	}

	return vniRangeResourceModel{
		Uuid:         types.StringValue(vniRange.UUID),
		Name:         types.StringValue(vniRange.Name),
		Description:  stringValueOrNull(vniRange.Description),
		StartVni:     types.Int64Value(int64(vniRange.StartVni)),
		EndVni:       types.Int64Value(int64(vniRange.EndVni)),
		PoolUuid:     types.StringValue(vniRange.L2NetworkUuid),
		ResourceUuid: resourceUuid,
		TagUuids:     tagUuids,
		SystemTags:   systemTags,
	}
}

func vniRangeCreateParam(plan vniRangeResourceModel) param.CreateVniRangeParam {
	createParam := param.CreateVniRangeParam{
		BaseParam: param.BaseParam{SystemTags: listToStringSlice(plan.SystemTags)},
		Params: param.CreateVniRangeParamDetail{
			Name:     plan.Name.ValueString(),
			StartVni: int(plan.StartVni.ValueInt64()),
			EndVni:   int(plan.EndVni.ValueInt64()),
			TagUuids: listToStringSlice(plan.TagUuids),
		},
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createParam.Params.Description = stringPtr(plan.Description.ValueString())
	}
	if !plan.ResourceUuid.IsNull() && !plan.ResourceUuid.IsUnknown() {
		createParam.Params.ResourceUuid = stringPtr(plan.ResourceUuid.ValueString())
	}

	return createParam
}
