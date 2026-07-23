// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

const attachmentSystemTagPageSize = 500

var (
	_ resource.Resource                = &l2NetworkClusterAttachmentResource{}
	_ resource.ResourceWithConfigure   = &l2NetworkClusterAttachmentResource{}
	_ resource.ResourceWithImportState = &l2NetworkClusterAttachmentResource{}
)

type l2NetworkClusterAttachmentResource struct {
	client *client.ZSClient
}

type l2NetworkClusterAttachmentModel struct {
	ID             types.String `tfsdk:"id"`
	L2NetworkUuid  types.String `tfsdk:"l2_network_uuid"`
	ClusterUuid    types.String `tfsdk:"cluster_uuid"`
	L2ProviderType types.String `tfsdk:"l2_provider_type"`
	SystemTags     types.List   `tfsdk:"system_tags"`
}

func L2NetworkClusterAttachmentResource() resource.Resource {
	return &l2NetworkClusterAttachmentResource{}
}

func (r *l2NetworkClusterAttachmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *l2NetworkClusterAttachmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l2_network_cluster_attachment"
}

func (r *l2NetworkClusterAttachmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Attach a ZStack L2 network to a cluster. Destroying this resource detaches the L2 network from the cluster without deleting either resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Terraform resource ID in the format `l2_network_uuid:cluster_uuid`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"l2_network_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the L2 network to attach.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the cluster to attach the L2 network to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"l2_provider_type": schema.StringAttribute{
				Optional:    true,
				Description: "The L2 provider type sent at the top level of the attach request. VXLAN pool attachments commonly require this value.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"system_tags": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "System tags sent at the top level of the attach request. Use these for VXLAN VTEP CIDR tags when required.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *l2NetworkClusterAttachmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan l2NetworkClusterAttachmentModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddError("Client Not Configured", "The ZStack client was not properly configured.")
		return
	}

	l2NetworkUuid := plan.L2NetworkUuid.ValueString()
	clusterUuid := plan.ClusterUuid.ValueString()
	plan.ID = types.StringValue(l2NetworkClusterAttachmentID(l2NetworkUuid, clusterUuid))

	attached, err := r.isL2NetworkAttachedToCluster(l2NetworkUuid, clusterUuid)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating L2 Network Cluster Attachment",
			"Could not query L2 network cluster attachment: "+err.Error(),
		)
		return
	}

	if !attached {
		_, err = r.client.AttachL2NetworkToCluster(l2NetworkUuid, clusterUuid, l2NetworkClusterAttachParam(plan))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating L2 Network Cluster Attachment",
				fmt.Sprintf("Could not attach L2 network %s to cluster %s: %s", l2NetworkUuid, clusterUuid, err.Error()),
			)
			return
		}
	}

	tflog.Info(ctx, "L2 network cluster attachment created", map[string]any{
		"id":              plan.ID.ValueString(),
		"l2_network_uuid": l2NetworkUuid,
		"cluster_uuid":    clusterUuid,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *l2NetworkClusterAttachmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state l2NetworkClusterAttachmentModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	attached, err := r.isL2NetworkAttachedToCluster(state.L2NetworkUuid.ValueString(), state.ClusterUuid.ValueString())
	if err != nil {
		if isZStackNotFoundError(err) || errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading L2 Network Cluster Attachment",
			"Could not query L2 network cluster attachment: "+err.Error(),
		)
		return
	}

	if !attached {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(l2NetworkClusterAttachmentID(state.L2NetworkUuid.ValueString(), state.ClusterUuid.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *l2NetworkClusterAttachmentResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"This resource does not support updates. Any changes require replacement.",
	)
}

func (r *l2NetworkClusterAttachmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state l2NetworkClusterAttachmentModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddError("Client Not Configured", "The ZStack client was not properly configured.")
		return
	}

	l2NetworkUuid := state.L2NetworkUuid.ValueString()
	clusterUuid := state.ClusterUuid.ValueString()

	tflog.Debug(ctx, "Deleting L2 network cluster attachment", map[string]any{
		"id":              state.ID.ValueString(),
		"l2_network_uuid": l2NetworkUuid,
		"cluster_uuid":    clusterUuid,
	})

	err := r.client.DetachL2NetworkFromCluster(l2NetworkUuid, clusterUuid, param.DeleteModePermissive)
	if err != nil {
		attached, queryErr := r.isL2NetworkAttachedToCluster(l2NetworkUuid, clusterUuid)
		if queryErr != nil {
			if isZStackNotFoundError(queryErr) || errors.Is(queryErr, ErrResourceNotFound) {
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError(
				"Error deleting L2 Network Cluster Attachment",
				fmt.Sprintf("Detach failed (%s) and could not verify attachment status: %s", err.Error(), queryErr.Error()),
			)
			return
		}
		if attached {
			resp.Diagnostics.AddError(
				"Error deleting L2 Network Cluster Attachment",
				fmt.Sprintf("Could not detach L2 network %s from cluster %s: %s", l2NetworkUuid, clusterUuid, err.Error()),
			)
			return
		}
	}

	resp.State.RemoveResource(ctx)
}

func (r *l2NetworkClusterAttachmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	l2NetworkUuid, clusterUuid, err := parseL2NetworkClusterAttachmentID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected format: l2_network_uuid:cluster_uuid (e.g. abc123:def456).",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), l2NetworkClusterAttachmentID(l2NetworkUuid, clusterUuid))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("l2_network_uuid"), l2NetworkUuid)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cluster_uuid"), clusterUuid)...)
	if resp.Diagnostics.HasError() || r.client == nil {
		return
	}

	vxlanPool, err := findResourceByGet(r.client.GetL2VxlanNetworkPool, l2NetworkUuid)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return
		}
		resp.Diagnostics.AddError(
			"Error importing L2 Network Cluster Attachment",
			fmt.Sprintf("Could not query VXLAN pool %s: %s", l2NetworkUuid, err.Error()),
		)
		return
	}

	if vxlanPool.VSwitchType != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("l2_provider_type"), vxlanPool.VSwitchType)...)
	}

	systemTags, err := queryVxlanAttachmentSystemTags(ctx, r.client, l2NetworkUuid, clusterUuid)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing L2 Network Cluster Attachment",
			fmt.Sprintf("Could not query VXLAN attachment system tags: %s", err.Error()),
		)
		return
	}
	if len(systemTags) > 0 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("system_tags"), systemTags)...)
	}
}

func (r *l2NetworkClusterAttachmentResource) isL2NetworkAttachedToCluster(l2NetworkUuid, clusterUuid string) (bool, error) {
	l2Network, err := findResourceByGet(r.client.GetL2Network, l2NetworkUuid)
	if err != nil {
		return false, err
	}

	for _, attachedClusterUuid := range l2Network.AttachedClusterUuids {
		if attachedClusterUuid == clusterUuid {
			return true, nil
		}
	}
	return false, nil
}

func l2NetworkClusterAttachmentID(l2NetworkUuid, clusterUuid string) string {
	return fmt.Sprintf("%s:%s", l2NetworkUuid, clusterUuid)
}

func parseL2NetworkClusterAttachmentID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("expected l2_network_uuid:cluster_uuid")
	}
	return parts[0], parts[1], nil
}

func queryVxlanAttachmentSystemTags(
	ctx context.Context,
	zstackClient *client.ZSClient,
	l2NetworkUuid, clusterUuid string,
) ([]string, error) {
	query := param.NewQueryParam()
	query.AddQ("resourceUuid=" + l2NetworkUuid)
	query.Limit(attachmentSystemTagPageSize).ReplyWithCount(true)

	expectedPrefix := fmt.Sprintf(
		"l2NetworkUuid::%s::clusterUuid::%s::cidr::{",
		l2NetworkUuid,
		clusterUuid,
	)
	offset := 0
	expectedTotal := -1
	seenUuids := make(map[string]struct{})
	matchedTags := make([]string, 0)

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		query.Start(offset)

		var page []view.SystemTagInventoryView
		total, err := zstackClient.ZSHttpClient.Page(ctx, "v1/system-tags", &query, &page)
		if err != nil {
			return nil, err
		}
		if total < 0 || total < offset+len(page) {
			return nil, fmt.Errorf("invalid system tag pagination total %d at offset %d with %d results", total, offset, len(page))
		}
		if expectedTotal == -1 {
			expectedTotal = total
		} else if total != expectedTotal {
			return nil, fmt.Errorf("system tag pagination total changed from %d to %d", expectedTotal, total)
		}

		for _, systemTag := range page {
			if systemTag.UUID == "" {
				return nil, fmt.Errorf("system tag query returned an item without a UUID")
			}
			if _, exists := seenUuids[systemTag.UUID]; exists {
				return nil, fmt.Errorf("system tag query returned duplicate UUID %s", systemTag.UUID)
			}
			seenUuids[systemTag.UUID] = struct{}{}
			if systemTag.ResourceUuid == l2NetworkUuid && strings.HasPrefix(systemTag.Tag, expectedPrefix) {
				matchedTags = append(matchedTags, systemTag.Tag)
			}
		}

		offset += len(page)
		if offset >= expectedTotal {
			break
		}
		if len(page) == 0 {
			return nil, fmt.Errorf("system tag pagination ended at %d of %d results", offset, expectedTotal)
		}
	}

	sort.Strings(matchedTags)
	return matchedTags, nil
}

func l2NetworkClusterAttachParam(plan l2NetworkClusterAttachmentModel) param.AttachL2NetworkToClusterParam {
	var l2ProviderType *string
	if !plan.L2ProviderType.IsNull() && !plan.L2ProviderType.IsUnknown() {
		l2ProviderType = stringPtr(plan.L2ProviderType.ValueString())
	}
	return attachL2NetworkToClusterParam(l2ProviderType, listToStringSlice(plan.SystemTags))
}

func attachL2NetworkToClusterParam(l2ProviderType *string, systemTags []string) param.AttachL2NetworkToClusterParam {
	return param.AttachL2NetworkToClusterParam{
		BaseParam: param.BaseParam{SystemTags: systemTags},
		Params: param.AttachL2NetworkToClusterParamDetail{
			L2ProviderType: l2ProviderType,
		},
	}
}
