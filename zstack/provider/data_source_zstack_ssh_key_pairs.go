// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"terraform-provider-zstack/zstack/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ datasource.DataSource              = &sshKeyPairDataSource{}
	_ datasource.DataSourceWithConfigure = &sshKeyPairDataSource{}
)

func ZStackSshKeyPairDataSource() datasource.DataSource {
	return &sshKeyPairDataSource{}
}

type sshKeyPairItem struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	PublicKey   types.String `tfsdk:"public_key"`
}

type sshKeyPairDataSourceModel struct {
	Name        types.String     `tfsdk:"name"`
	NamePattern types.String     `tfsdk:"name_pattern"`
	Filter      []Filter         `tfsdk:"filter"`
	SshKeyPairs []sshKeyPairItem `tfsdk:"ssh_key_pairs"`
}

type sshKeyPairDataSource struct {
	client *client.ZSClient
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *sshKeyPairDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Metadata implements datasource.DataSource.
func (d *sshKeyPairDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key_pairs"
}

// Read implements datasource.DataSource.
func (d *sshKeyPairDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sshKeyPairDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := param.NewQueryParam()

	if !state.Name.IsNull() {
		params.AddQ("name=" + state.Name.ValueString())
	} else if !state.NamePattern.IsNull() {
		params.AddQ("name~=" + state.NamePattern.ValueString())
	}

	sshKeyPairs, err := d.client.QuerySshKeyPair(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Query ZStack SSH Key Pairs",
			err.Error(),
		)
		return
	}

	filters := make(map[string][]string)
	for _, filter := range state.Filter {
		values := make([]string, 0, len(filter.Values.Elements()))
		diags := filter.Values.ElementsAs(ctx, &values, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		filters[filter.Name.ValueString()] = values
	}

	filterSshKeyPairs, filterDiags := utils.FilterResource(ctx, sshKeyPairs, filters, "ssh_key_pair")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, skp := range filterSshKeyPairs {
		state.SshKeyPairs = append(state.SshKeyPairs, sshKeyPairItem{
			Uuid:        types.StringValue(skp.UUID),
			Name:        types.StringValue(skp.Name),
			Description: types.StringValue(skp.Description),
			PublicKey:   types.StringValue(skp.PublicKey),
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Schema implements datasource.DataSource.
func (d *sshKeyPairDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query ZStack SSH Key Pairs by name, name pattern, or additional filters.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for querying an SSH key pair.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy matching SSH key pair names. Use % or _ like SQL.",
				Optional:    true,
			},
			"ssh_key_pairs": schema.ListNestedAttribute{
				Description: "List of matched SSH key pairs.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "UUID of the SSH key pair.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the SSH key pair.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the SSH key pair.",
							Computed:    true,
						},
						"public_key": schema.StringAttribute{
							Description: "The SSH public key content.",
							Computed:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter results by field values.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the field to filter by.",
							Required:    true,
						},
						"values": schema.SetAttribute{
							Description: "List of values to match. Treated as OR conditions.",
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}
