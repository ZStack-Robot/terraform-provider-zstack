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

const (
	maxVni            = 16777214
	vniRangePageSize  = 500
	vniRangeQueryPath = "v1/l2-networks/vxlan-pool/vni-range"
)

var (
	_ resource.Resource                   = &vniRangeResource{}
	_ resource.ResourceWithConfigure      = &vniRangeResource{}
	_ resource.ResourceWithImportState    = &vniRangeResource{}
	_ resource.ResourceWithModifyPlan     = &vniRangeResource{}
	_ resource.ResourceWithValidateConfig = &vniRangeResource{}
)

type vniRangeResource struct {
	client            *client.ZSClient
	vniRangePageQuery vniRangePageQuery
}

type vniRangePageQuery interface {
	PageVniRanges(context.Context, *param.QueryParam) ([]view.VniRangeInventoryView, int, error)
}

type zstackVniRangePageQuery struct {
	client *client.ZSClient
}

func (q zstackVniRangePageQuery) PageVniRanges(ctx context.Context, query *param.QueryParam) ([]view.VniRangeInventoryView, int, error) {
	var ranges []view.VniRangeInventoryView
	total, err := q.client.ZSHttpClient.Page(ctx, vniRangeQueryPath, query, &ranges)
	return ranges, total, err
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
	r.vniRangePageQuery = zstackVniRangePageQuery{client: configuredClient}
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
				Description: "The first VNI in the range (1-16777214 for software SDN).",
				Validators: []validator.Int64{
					int64validator.Between(1, maxVni),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"end_vni": schema.Int64Attribute{
				Required:    true,
				Description: "The last VNI in the range (1-16777214 for software SDN).",
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

func (r *vniRangeResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan vniRangeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() ||
		plan.StartVni.IsNull() || plan.StartVni.IsUnknown() ||
		plan.EndVni.IsNull() || plan.EndVni.IsUnknown() ||
		plan.PoolUuid.IsNull() || plan.PoolUuid.IsUnknown() {
		return
	}

	startVni := plan.StartVni.ValueInt64()
	endVni := plan.EndVni.ValueInt64()
	if startVni < 1 || startVni > maxVni || endVni < 1 || endVni > maxVni || startVni > endVni {
		return
	}

	currentUuid := ""
	if !req.State.Raw.IsNull() {
		var state vniRangeResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if !state.Uuid.IsNull() && !state.Uuid.IsUnknown() {
			currentUuid = state.Uuid.ValueString()
		}
	}

	if r.vniRangePageQuery == nil {
		resp.Diagnostics.AddError(
			"Unable to validate VNI range availability",
			"The ZStack client is not configured, so existing VNI ranges cannot be checked.",
		)
		return
	}

	conflict, err := r.findOverlappingVniRange(
		ctx,
		plan.PoolUuid.ValueString(),
		startVni,
		endVni,
		currentUuid,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to validate VNI range availability",
			fmt.Sprintf("Could not query all VNI ranges in pool %s: %s", plan.PoolUuid.ValueString(), err.Error()),
		)
		return
	}
	if conflict == nil {
		return
	}

	resp.Diagnostics.AddAttributeError(
		path.Root("start_vni"),
		"Overlapping VNI range",
		fmt.Sprintf(
			"Requested VNI range [%d, %d] overlaps existing range %q (%s) [%d, %d] in pool %s.",
			startVni,
			endVni,
			conflict.Name,
			conflict.UUID,
			conflict.StartVni,
			conflict.EndVni,
			plan.PoolUuid.ValueString(),
		),
	)
}

func (r *vniRangeResource) findOverlappingVniRange(
	ctx context.Context,
	poolUuid string,
	startVni, endVni int64,
	currentUuid string,
) (*view.VniRangeInventoryView, error) {
	ranges, err := r.queryAllVniRanges(ctx, poolUuid)
	if err != nil {
		return nil, err
	}

	return findVniRangeOverlap(ranges, startVni, endVni, currentUuid), nil
}

func findVniRangeOverlap(
	ranges []view.VniRangeInventoryView,
	startVni, endVni int64,
	currentUuid string,
) *view.VniRangeInventoryView {
	for i := range ranges {
		existing := &ranges[i]
		if existing.UUID == currentUuid {
			continue
		}
		if startVni <= int64(existing.EndVni) && endVni >= int64(existing.StartVni) {
			return existing
		}
	}
	return nil
}

func (r *vniRangeResource) queryAllVniRanges(ctx context.Context, poolUuid string) ([]view.VniRangeInventoryView, error) {
	return queryAllVniRanges(ctx, r.vniRangePageQuery, poolUuid)
}

func queryAllVniRanges(ctx context.Context, pageQuery vniRangePageQuery, poolUuid string) ([]view.VniRangeInventoryView, error) {
	if pageQuery == nil {
		return nil, errors.New("VNI range page query is not configured")
	}

	ranges := make([]view.VniRangeInventoryView, 0)
	seenUuids := make(map[string]struct{})
	expectedTotal := -1

	for start := 0; ; start = len(ranges) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		query := param.NewQueryParam()
		query.AddQ(fmt.Sprintf("l2NetworkUuid=%s", poolUuid))
		query.Start(start).Limit(vniRangePageSize).ReplyWithCount(true)

		page, total, err := pageQuery.PageVniRanges(ctx, &query)
		if err != nil {
			return nil, err
		}
		if total < 0 {
			return nil, fmt.Errorf("VNI range query returned invalid total %d", total)
		}
		if expectedTotal == -1 {
			expectedTotal = total
		} else if total != expectedTotal {
			return nil, fmt.Errorf("VNI range query total changed during pagination: expected %d, got %d", expectedTotal, total)
		}

		if len(page) == 0 {
			if len(ranges) != expectedTotal {
				return nil, fmt.Errorf("incomplete VNI range query: received %d of %d records", len(ranges), expectedTotal)
			}
			return ranges, nil
		}

		for _, existing := range page {
			if existing.L2NetworkUuid != poolUuid {
				return nil, fmt.Errorf(
					"VNI range query for pool %s returned range %s from pool %s",
					poolUuid,
					existing.UUID,
					existing.L2NetworkUuid,
				)
			}
			if existing.UUID != "" {
				if _, exists := seenUuids[existing.UUID]; exists {
					return nil, fmt.Errorf("VNI range query returned duplicate range %s", existing.UUID)
				}
				seenUuids[existing.UUID] = struct{}{}
			}
			ranges = append(ranges, existing)
		}

		if len(ranges) > expectedTotal {
			return nil, fmt.Errorf("VNI range query returned %d records, exceeding total %d", len(ranges), expectedTotal)
		}
		if len(ranges) == expectedTotal {
			return ranges, nil
		}
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
