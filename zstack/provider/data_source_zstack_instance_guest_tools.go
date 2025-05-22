// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"zstack.io/zstack-sdk-go/pkg/client"
)

var (
	_ datasource.DataSource              = &guestToolsDataSource{}
	_ datasource.DataSourceWithConfigure = &guestToolsDataSource{}
)

func ZStackGuestToolsDataSource() datasource.DataSource {
	return &guestToolsDataSource{}
}

type guestToolsDataSource struct {
	client *client.ZSClient
}

type guestToolsDataSourceModel struct {
	InstanceUuid types.String `tfsdk:"instance_uuid"`
	Version      types.String `tfsdk:"version"`
	Status       types.String `tfsdk:"status"`
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *guestToolsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *guestToolsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_guest_tools"
}

func (d *guestToolsDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a  instance guest tools info.",
		Attributes: map[string]schema.Attribute{
			"version": schema.StringAttribute{
				Description: "Version of the guest tools installed on the instance.",
				Computed:    true,
			},
			"instance_uuid": schema.StringAttribute{
				Description: "UUID of the instance to fetch guest tools info for.",
				Required:    true,
			},
			"status": schema.StringAttribute{
				Description: "Current status of the guest tools (e.g. Running, NotInstalled).",
				Computed:    true,
			},
		},
	}
}

func (d *guestToolsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state guestToolsDataSourceModel
	//var state clusterModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.InstanceUuid.IsNull() || state.InstanceUuid.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"'instance_uuid' must be provided.",
		)
	}

	guest_tools, err := d.client.GetVmGuestToolsInfo(state.InstanceUuid.ValueString())

	if err != nil {
		resp.Diagnostics.AddWarning(
			"Unable to read guest tools info",
			fmt.Sprintf("Failed to read guest tools info for instance UUID: %s, error: %s",
				state.InstanceUuid.ValueString(), err.Error()),
		)
		return
	}

	if guest_tools == nil {
		resp.Diagnostics.AddError(
			"Unable to Read Guest Tools Info",
			fmt.Sprintf("Guest tools not found for instance UUID: %s", state.InstanceUuid.ValueString()),
		)
		return
	}

	state.Version = types.StringValue(guest_tools.Version)
	state.Status = types.StringValue(guest_tools.Status)

	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
