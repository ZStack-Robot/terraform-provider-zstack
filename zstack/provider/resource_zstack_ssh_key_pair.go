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
	_ resource.Resource                = &sshKeyPairResource{}
	_ resource.ResourceWithConfigure   = &sshKeyPairResource{}
	_ resource.ResourceWithImportState = &sshKeyPairResource{}
)

type sshKeyPairResource struct {
	client *client.ZSClient
}

type sshKeyPairResourceModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	PublicKey   types.String `tfsdk:"public_key"`
}

func SshKeyPairResource() resource.Resource {
	return &sshKeyPairResource{}
}

// Configure implements resource.ResourceWithConfigure.
func (r *sshKeyPairResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Metadata implements resource.Resource.
func (r *sshKeyPairResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key_pair"
}

// Schema implements resource.Resource.
func (r *sshKeyPairResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource allows you to manage SSH key pairs in ZStack. " +
			"An SSH key pair stores a public key that can be attached to VM instances for key-based authentication. " +
			"Only importing existing public keys is supported (not generating new key pairs).",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier (UUID) of the SSH key pair.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the SSH key pair.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description of the SSH key pair.",
			},
			"public_key": schema.StringAttribute{
				Required:    true,
				Description: "The SSH public key content.",
			},
		},
	}
}

// Create implements resource.Resource.
func (r *sshKeyPairResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sshKeyPairResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating SSH key pair")

	createParam := param.CreateSshKeyPairParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateSshKeyPairParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			PublicKey:   plan.PublicKey.ValueString(),
		},
	}

	sshKeyPair, err := r.client.CreateSshKeyPair(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not create SSH key pair in ZStack", "Error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(sshKeyPair.UUID)
	plan.Name = types.StringValue(sshKeyPair.Name)
	plan.Description = types.StringValue(sshKeyPair.Description)
	plan.PublicKey = types.StringValue(sshKeyPair.PublicKey)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (r *sshKeyPairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sshKeyPairResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sshKeyPair, err := r.client.GetSshKeyPair(state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading SSH key pair", "Could not read SSH key pair UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(sshKeyPair.UUID)
	state.Name = types.StringValue(sshKeyPair.Name)
	state.Description = types.StringValue(sshKeyPair.Description)
	state.PublicKey = types.StringValue(sshKeyPair.PublicKey)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *sshKeyPairResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sshKeyPairResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state sshKeyPairResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateParam := param.UpdateSshKeyPairParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateSshKeyPairParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	sshKeyPair, err := r.client.UpdateSshKeyPair(state.Uuid.ValueString(), updateParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not update SSH key pair", "Error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(sshKeyPair.UUID)
	plan.Name = types.StringValue(sshKeyPair.Name)
	plan.Description = types.StringValue(sshKeyPair.Description)
	plan.PublicKey = types.StringValue(sshKeyPair.PublicKey)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *sshKeyPairResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sshKeyPairResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "SSH key pair uuid is empty, so nothing to delete, skip it")
		return
	}

	err := r.client.DeleteSshKeyPair(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete SSH key pair", err.Error())
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *sshKeyPairResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
