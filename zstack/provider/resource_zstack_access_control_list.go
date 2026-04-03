// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &accessControlListResource{}
	_ resource.ResourceWithConfigure   = &accessControlListResource{}
	_ resource.ResourceWithImportState = &accessControlListResource{}
)

type accessControlListResource struct {
	client *client.ZSClient
}

type accessControlListModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	IpVersion   types.Int64  `tfsdk:"ip_version"`
}

func AccessControlListResource() resource.Resource {
	return &accessControlListResource{}
}

func (r *accessControlListResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *accessControlListResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_access_control_list"
}

func (r *accessControlListResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage Access Control Lists (ACLs) in ZStack. " +
			"An ACL is used to control access to network resources based on IP addresses.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the Access Control List.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Access Control List.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the Access Control List.",
			},
			"ip_version": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The IP version (4 for IPv4, 6 for IPv6).",
			},
		},
	}
}

func (r *accessControlListResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan accessControlListModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateAccessControlListParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateAccessControlListParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			IpVersion:   intPtrFromInt64OrNil(plan.IpVersion),
		},
	}

	acl, err := r.client.CreateAccessControlList(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to create Access Control List",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(acl.UUID)
	plan.Name = types.StringValue(acl.Name)
	plan.Description = stringValueOrNull(acl.Description)
	plan.IpVersion = types.Int64Value(int64(acl.IpVersion))

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *accessControlListResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state accessControlListModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	acls, err := r.client.QueryAccessControlList(&queryParam)

	if err != nil {
		tflog.Warn(ctx, "Unable to query Access Control Lists. It may have been deleted.: "+err.Error())
		state = accessControlListModel{
			Uuid: types.StringValue(""),
		}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false

	for _, acl := range acls {
		if acl.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(acl.UUID)
			state.Name = types.StringValue(acl.Name)
			state.Description = stringValueOrNull(acl.Description)
			state.IpVersion = types.Int64Value(int64(acl.IpVersion))
			found = true
			break
		}
	}
	if !found {
		tflog.Warn(ctx, "Access Control List not found. It might have been deleted outside of Terraform.")
		state = accessControlListModel{
			Uuid: types.StringValue(""),
		}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *accessControlListResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan accessControlListModel
	var state accessControlListModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateAccessControlListParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateAccessControlListParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	acl, err := r.client.UpdateAccessControlList(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to update Access Control List",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(acl.UUID)
	plan.Name = types.StringValue(acl.Name)
	plan.Description = stringValueOrNull(acl.Description)
	plan.IpVersion = types.Int64Value(int64(acl.IpVersion))

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *accessControlListResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state accessControlListModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Access Control List UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteAccessControlList(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError("fail to delete Access Control List", ""+err.Error())
		return
	}
}

func (r *accessControlListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
