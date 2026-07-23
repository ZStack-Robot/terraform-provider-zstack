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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

var (
	_ resource.Resource                = &l2VxlanNetworkPoolResource{}
	_ resource.ResourceWithConfigure   = &l2VxlanNetworkPoolResource{}
	_ resource.ResourceWithImportState = &l2VxlanNetworkPoolResource{}
)

type l2VxlanNetworkPoolResource struct {
	client *client.ZSClient
}

type l2VxlanNetworkPoolResourceModel struct {
	Uuid                 types.String `tfsdk:"uuid"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	ZoneUuid             types.String `tfsdk:"zone_uuid"`
	PhysicalInterface    types.String `tfsdk:"physical_interface"`
	Type                 types.String `tfsdk:"type"`
	VSwitchType          types.String `tfsdk:"vswitch_type"`
	VirtualNetworkId     types.Int64  `tfsdk:"virtual_network_id"`
	Isolated             types.Bool   `tfsdk:"isolated"`
	Pvlan                types.String `tfsdk:"pvlan"`
	ResourceUuid         types.String `tfsdk:"resource_uuid"`
	TagUuids             types.List   `tfsdk:"tag_uuids"`
	SystemTags           types.List   `tfsdk:"system_tags"`
	AttachedClusterUuids types.List   `tfsdk:"attached_cluster_uuids"`
}

func L2VxlanNetworkPoolResource() resource.Resource {
	return &l2VxlanNetworkPoolResource{}
}

func (r *l2VxlanNetworkPoolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *l2VxlanNetworkPoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l2vxlan_network_pool"
}

func (r *l2VxlanNetworkPoolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a ZStack L2 VXLAN network pool.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the L2 VXLAN network pool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the L2 VXLAN network pool.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the L2 VXLAN network pool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the zone where the VXLAN pool is created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"physical_interface": schema.StringAttribute{
				Required:    true,
				Description: "The physical network interface used by the VXLAN pool, for example `bond0`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The L2 network type reported by ZStack.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vswitch_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The virtual switch type, such as `LinuxBridge` or `OvsDpdk`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"virtual_network_id": schema.Int64Attribute{
				Computed:    true,
				Description: "The virtual network identifier reported for the pool.",
			},
			"isolated": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the VXLAN pool is isolated.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplaceIfConfigured(),
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"pvlan": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The private VLAN setting for the pool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A custom UUID requested when creating the VXLAN pool.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tag_uuids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Tag UUIDs attached while creating the VXLAN pool.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIfConfigured(),
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"system_tags": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "System tags sent while creating the VXLAN pool.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIfConfigured(),
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"attached_cluster_uuids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Cluster UUIDs currently attached to the VXLAN pool.",
			},
		},
	}
}

func (r *l2VxlanNetworkPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan l2VxlanNetworkPoolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createParam := l2VxlanNetworkPoolCreateParam(plan)

	pool, err := r.client.CreateL2VxlanNetworkPool(createParam)
	if err != nil {
		resp.Diagnostics.AddError("Error creating L2 VXLAN network pool", err.Error())
		return
	}

	state := l2VxlanNetworkPoolModelFromView(pool, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *l2VxlanNetworkPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state l2VxlanNetworkPoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pool, err := findResourceByGet(r.client.GetL2VxlanNetworkPool, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading L2 VXLAN network pool", err.Error())
		return
	}

	refreshed := l2VxlanNetworkPoolModelFromView(pool, state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &refreshed)...)
}

func (r *l2VxlanNetworkPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan l2VxlanNetworkPoolResourceModel
	var state l2VxlanNetworkPoolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Name.ValueString() != state.Name.ValueString() || plan.Description.ValueString() != state.Description.ValueString() {
		_, err := r.client.UpdateL2Network(state.Uuid.ValueString(), param.UpdateL2NetworkParam{
			BaseParam: param.BaseParam{},
			Params: param.UpdateL2NetworkParamDetail{
				Name:        plan.Name.ValueString(),
				Description: stringPtrOrNil(plan.Description.ValueString()),
			},
		})
		if err != nil {
			resp.Diagnostics.AddError("Error updating L2 VXLAN network pool", err.Error())
			return
		}
	}

	pool, err := findResourceByGet(r.client.GetL2VxlanNetworkPool, state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading L2 VXLAN network pool after update", err.Error())
		return
	}

	refreshed := l2VxlanNetworkPoolModelFromView(pool, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &refreshed)...)
}

func (r *l2VxlanNetworkPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state l2VxlanNetworkPoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteL2Network(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil && !isZStackNotFoundError(err) {
		resp.Diagnostics.AddError("Error deleting L2 VXLAN network pool", err.Error())
	}
}

func (r *l2VxlanNetworkPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func l2VxlanNetworkPoolModelFromView(pool *view.L2VxlanNetworkPoolInventoryView, prior l2VxlanNetworkPoolResourceModel) l2VxlanNetworkPoolResourceModel {
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

	return l2VxlanNetworkPoolResourceModel{
		Uuid:                 types.StringValue(pool.UUID),
		Name:                 types.StringValue(pool.Name),
		Description:          stringValueOrNull(pool.Description),
		ZoneUuid:             types.StringValue(pool.ZoneUuid),
		PhysicalInterface:    types.StringValue(pool.PhysicalInterface),
		Type:                 stringValueOrNull(pool.Type),
		VSwitchType:          stringValueOrNull(pool.VSwitchType),
		VirtualNetworkId:     types.Int64Value(int64(pool.VirtualNetworkId)),
		Isolated:             types.BoolValue(pool.Isolated),
		Pvlan:                stringValueOrNull(pool.Pvlan),
		ResourceUuid:         resourceUuid,
		TagUuids:             tagUuids,
		SystemTags:           systemTags,
		AttachedClusterUuids: stringSliceToList(pool.AttachedClusterUuids),
	}
}

func l2VxlanNetworkPoolCreateParam(plan l2VxlanNetworkPoolResourceModel) param.CreateL2VxlanNetworkPoolParam {
	createParam := param.CreateL2VxlanNetworkPoolParam{
		BaseParam: param.BaseParam{SystemTags: listToStringSlice(plan.SystemTags)},
		Params: param.CreateL2VxlanNetworkPoolParamDetail{
			Name:              plan.Name.ValueString(),
			ZoneUuid:          plan.ZoneUuid.ValueString(),
			PhysicalInterface: stringPtr(plan.PhysicalInterface.ValueString()),
			TagUuids:          listToStringSlice(plan.TagUuids),
		},
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createParam.Params.Description = stringPtr(plan.Description.ValueString())
	}
	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		createParam.Params.Type = stringPtr(plan.Type.ValueString())
	}
	if !plan.VSwitchType.IsNull() && !plan.VSwitchType.IsUnknown() {
		createParam.Params.VSwitchType = stringPtr(plan.VSwitchType.ValueString())
	}
	if !plan.Isolated.IsNull() && !plan.Isolated.IsUnknown() {
		createParam.Params.Isolated = boolPtr(plan.Isolated.ValueBool())
	}
	if !plan.Pvlan.IsNull() && !plan.Pvlan.IsUnknown() {
		createParam.Params.Pvlan = stringPtr(plan.Pvlan.ValueString())
	}
	if !plan.ResourceUuid.IsNull() && !plan.ResourceUuid.IsUnknown() {
		createParam.Params.ResourceUuid = stringPtr(plan.ResourceUuid.ValueString())
	}

	return createParam
}
