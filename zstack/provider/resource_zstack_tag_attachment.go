// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

// tagClient defines the minimal interface needed by tagAttachmentResource.
// This interface enables testing with spy/fake clients.
type tagClient interface {
	AttachTagToResources(tagUuid string, params param.AttachTagToResourcesParam) (*view.AttachTagToResourcesEventView, error)
	DetachTagFromResources(uuid string, deleteMode param.DeleteMode, resourceUuids []string) error
	QueryUserTag(params *param.QueryParam) ([]view.UserTagInventoryView, error)
}

// zstackTagClient adapts *client.ZSClient to the tagClient interface.
type zstackTagClient struct {
	*client.ZSClient
}

func (c *zstackTagClient) DetachTagFromResources(uuid string, deleteMode param.DeleteMode, resourceUuids []string) error {
	return c.ZSClient.DetachTagFromResources(uuid, deleteMode)
}

func (c *zstackTagClient) AttachTagToResources(tagUuid string, params param.AttachTagToResourcesParam) (*view.AttachTagToResourcesEventView, error) {
	return c.ZSClient.AttachTagToResources(tagUuid, params)
}

func (c *zstackTagClient) QueryUserTag(params *param.QueryParam) ([]view.UserTagInventoryView, error) {
	return c.ZSClient.QueryUserTag(params)
}

var (
	_ resource.Resource              = &tagAttachmentResource{}
	_ resource.ResourceWithConfigure = &tagAttachmentResource{}
)

type tagAttachmentResource struct {
	client tagClient
}

type tagAttachmentModel struct {
	//ID            types.String `tfsdk:"id"`
	TagUuid       types.String `tfsdk:"tag_uuid"`
	ResourceUuids types.List   `tfsdk:"resource_uuids"`
	Tokens        types.Map    `tfsdk:"tokens"`
}

func TagAttachmentResource() resource.Resource {
	return &tagAttachmentResource{}
}

func (r *tagAttachmentResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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
	r.client = &zstackTagClient{ZSClient: client}
}

func (r *tagAttachmentResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_tag_attachment"
}

func (r *tagAttachmentResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "ZStack Tag Attach to  Resource",
		Attributes: map[string]schema.Attribute{
			"tag_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the tag.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_uuids": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "The uuid of the resource.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"tokens": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "If attach type is 'withToken', you must set this. For simple attach, leave it empty or omit.",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *tagAttachmentResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan tagAttachmentModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	resourceUuids := make([]string, 0, len(plan.ResourceUuids.Elements()))
	for _, v := range plan.ResourceUuids.Elements() {
		resourceUuids = append(resourceUuids, v.(types.String).ValueString())
	}

	attachParams := param.AttachTagToResourcesParam{
		BaseParam: param.BaseParam{},
		Params: param.AttachTagToResourcesParamDetail{
			ResourceUuids: resourceUuids,
		},
	}

	if !plan.Tokens.IsNull() && len(plan.Tokens.Elements()) > 0 {
		tokenMap := make(map[string]string)
		for k, v := range plan.Tokens.Elements() {
			tokenMap[k] = v.(types.String).ValueString()
		}
		attachParams.Params.Tokens = tokenMap
	}

	_, err := r.client.AttachTagToResources(plan.TagUuid.ValueString(), attachParams)
	if err != nil {
		response.Diagnostics.AddError("Error attaching tag", err.Error())
		return
	}

	//id := fmt.Sprintf("%s|%s", plan.TagUuid.ValueString(), strings.Join(resourceUuids, ","))
	response.State.Set(ctx, &tagAttachmentModel{
		TagUuid:       plan.TagUuid,
		ResourceUuids: plan.ResourceUuids,
		Tokens:        plan.Tokens,
	})

}

func (r *tagAttachmentResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {

	var state tagAttachmentModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	tagUuid := state.TagUuid.ValueString()
	if tagUuid == "" {
		response.Diagnostics.AddError("Tag UUID is empty", "Cannot read tag attachment without a valid tag UUID.")
		return
	}
	desiredResourceUuids := extractResourceUuids(state.ResourceUuids)
	q := param.NewQueryParam()
	q.AddQ("tagPatternUuid=" + tagUuid)
	results, err := r.client.QueryUserTag(&q)
	if err != nil {
		response.Diagnostics.AddError(
			"Unable to Read Tag Binding",
			fmt.Sprintf("Error retrieving tag binding for tag UUID %s: %s", tagUuid, err),
		)
		return
	}

	desired := make(map[string]struct{}, len(desiredResourceUuids))
	for _, uuid := range desiredResourceUuids {
		desired[uuid] = struct{}{}
	}

	resourceUuids := make([]string, 0, len(desiredResourceUuids))
	for _, result := range results {
		if result.UUID == "" || result.ResourceUuid == "" {
			continue
		}
		if _, ok := desired[result.ResourceUuid]; ok {
			resourceUuids = append(resourceUuids, result.ResourceUuid)
		}
	}

	if len(resourceUuids) == 0 {
		response.State.RemoveResource(ctx)
		return
	}

	// Update state.ResourceUuids
	var resourceUuidsValues []attr.Value
	for _, uuid := range resourceUuids {
		resourceUuidsValues = append(resourceUuidsValues, types.StringValue(uuid))
	}

	state.ResourceUuids = types.ListValueMust(types.StringType, resourceUuidsValues)
	state.TagUuid = types.StringValue(tagUuid)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)

}

func (r *tagAttachmentResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	response.Diagnostics.AddError(
		"Update not supported",
		"Tag Attachment resource does not support updates. Please recreate the resource instead.",
	)
}

func extractResourceUuids(listValue types.List) []string {
	resourceUuids := make([]string, 0, len(listValue.Elements()))
	for _, v := range listValue.Elements() {
		resourceUuids = append(resourceUuids, v.(types.String).ValueString())
	}
	return resourceUuids
}

func (r *tagAttachmentResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state tagAttachmentModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	resourceUuids := extractResourceUuids(state.ResourceUuids)

	err := r.client.DetachTagFromResources(state.TagUuid.ValueString(), param.DeleteModePermissive, resourceUuids)
	if err != nil {
		response.Diagnostics.AddError("Error detaching tag", err.Error())
		return
	}

}
