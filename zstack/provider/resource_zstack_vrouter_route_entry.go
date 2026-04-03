// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &vrouterRouteEntryResource{}
	_ resource.ResourceWithConfigure   = &vrouterRouteEntryResource{}
	_ resource.ResourceWithImportState = &vrouterRouteEntryResource{}
)

type vrouterRouteEntryResource struct {
	client *client.ZSClient
}

type vrouterRouteEntryResourceModel struct {
	Uuid           types.String `tfsdk:"uuid"`
	RouteTableUuid types.String `tfsdk:"route_table_uuid"`
	Description    types.String `tfsdk:"description"`
	Type           types.String `tfsdk:"type"`
	Destination    types.String `tfsdk:"destination"`
	Target         types.String `tfsdk:"target"`
	Distance       types.Int64  `tfsdk:"distance"`
}

func VRouterRouteEntryResource() resource.Resource {
	return &vrouterRouteEntryResource{}
}

func (r *vrouterRouteEntryResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *vrouterRouteEntryResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vrouter_route_entry"
}

func (r *vrouterRouteEntryResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage VRouter route entries in ZStack. " +
			"A VRouter route entry defines a routing rule within a route table.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the VRouter route entry.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"route_table_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the route table.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the VRouter route entry.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of the route entry.",
			},
			"destination": schema.StringAttribute{
				Required:    true,
				Description: "The destination CIDR of the route entry.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The target of the route entry.",
			},
			"distance": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The distance (priority) of the route entry.",
			},
		},
	}
}

func (r *vrouterRouteEntryResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan vrouterRouteEntryResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.AddVRouterRouteEntryParam{
		BaseParam: param.BaseParam{},
		Params: param.AddVRouterRouteEntryParamDetail{
			Description: stringPtrOrNil(plan.Description.ValueString()),
			Type:        stringPtrOrNil(plan.Type.ValueString()),
			Destination: plan.Destination.ValueString(),
			Target:      stringPtrOrNil(plan.Target.ValueString()),
			Distance:    intPtrFromInt64OrNil(plan.Distance),
		},
	}

	vrouterRouteEntry, err := r.client.AddVRouterRouteEntry(plan.RouteTableUuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to add VRouter route entry",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(vrouterRouteEntry.UUID)
	plan.RouteTableUuid = types.StringValue(vrouterRouteEntry.RouteTableUuid)
	plan.Description = types.StringValue(vrouterRouteEntry.Description)
	plan.Type = types.StringValue(vrouterRouteEntry.Type)
	plan.Destination = types.StringValue(vrouterRouteEntry.Destination)
	plan.Target = types.StringValue(vrouterRouteEntry.Target)
	plan.Distance = types.Int64Value(int64(vrouterRouteEntry.Distance))

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vrouterRouteEntryResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state vrouterRouteEntryResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	vrouterRouteEntries, err := r.client.QueryVRouterRouteEntry(&queryParam)

	if err != nil {
		tflog.Warn(ctx, "Unable to query VRouter route entries. It may have been deleted.: "+err.Error())
		state = vrouterRouteEntryResourceModel{
			Uuid: types.StringValue(""),
		}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false

	for _, vrouterRouteEntry := range vrouterRouteEntries {
		if vrouterRouteEntry.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(vrouterRouteEntry.UUID)
			state.RouteTableUuid = types.StringValue(vrouterRouteEntry.RouteTableUuid)
			state.Description = types.StringValue(vrouterRouteEntry.Description)
			state.Type = types.StringValue(vrouterRouteEntry.Type)
			state.Destination = types.StringValue(vrouterRouteEntry.Destination)
			state.Target = types.StringValue(vrouterRouteEntry.Target)
			state.Distance = types.Int64Value(int64(vrouterRouteEntry.Distance))
			found = true
			break
		}
	}
	if !found {
		tflog.Warn(ctx, "VRouter route entry not found. It might have been deleted outside of Terraform.")
		state = vrouterRouteEntryResourceModel{
			Uuid: types.StringValue(""),
		}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vrouterRouteEntryResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (r *vrouterRouteEntryResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state vrouterRouteEntryResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "VRouter route entry UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteVRouterRouteEntry(state.RouteTableUuid.ValueString(), state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError("fail to delete VRouter route entry", ""+err.Error())
		return
	}
}

func (r *vrouterRouteEntryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func intPtrFromInt64OrNil(v types.Int64) *int {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	val := int(v.ValueInt64())
	return &val
}
