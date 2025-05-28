// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"zstack.io/zstack-sdk-go/pkg/client"
)

var (
	_ resource.Resource              = &tagAttachmentResource{}
	_ resource.ResourceWithConfigure = &tagAttachmentResource{}
)

type tagAttachmentResource struct {
	client *client.ZSClient
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
	r.client = client
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
			},
			"resource_uuids": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "The uuid of the resource.",
			},
			"tokens": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "If attach type is 'withToken', you must set this. For simple attach, leave it empty or omit.",
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

	attachType := []string{}
	if !plan.Tokens.IsNull() && len(plan.Tokens.Elements()) > 0 {
		tokenMap := make(map[string]interface{})
		for k, v := range plan.Tokens.Elements() {
			tokenMap[k] = v.(types.String).ValueString()
		}
		tokenJson, err := json.Marshal(tokenMap)
		if err != nil {
			response.Diagnostics.AddError("Error marshaling tokens", err.Error())
			return
		}
		attachType = []string{"withToken", string(tokenJson)}
	}

	_, err := r.client.AttachTagToResource(plan.TagUuid.ValueString(), resourceUuids, attachType...)
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
	result, err := r.client.GetUserTag(tagUuid)
	if err != nil {
		response.Diagnostics.AddError(
			"Unable to Read Tag Binding",
			fmt.Sprintf("Error retrieving tag binding for tag UUID %s: %s", tagUuid, err),
		)
		return
	}
	if len(result) == 0 {
		response.State.RemoveResource(ctx)
		return
	}

	var resourceUuids []string
	for _, item := range result {
		resourceUuids = append(resourceUuids, item.ResourceUuid)
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
	response.Diagnostics.AddWarning("Update not supported", "Please recreate the resource if you need changes")
}

func (r *tagAttachmentResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state tagAttachmentModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	resourceUuids := make([]string, 0, len(state.ResourceUuids.Elements()))
	for _, v := range state.ResourceUuids.Elements() {
		resourceUuids = append(resourceUuids, v.(types.String).ValueString())
	}

	err := r.client.DetachTagFromResource(state.TagUuid.ValueString(), resourceUuids)
	if err != nil {
		response.Diagnostics.AddError("Error detaching tag", err.Error())
		return
	}

}
