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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &certificateResource{}
	_ resource.ResourceWithConfigure   = &certificateResource{}
	_ resource.ResourceWithImportState = &certificateResource{}
)

type certificateResource struct {
	client *client.ZSClient
}

type certificateModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Certificate types.String `tfsdk:"certificate"`
	Description types.String `tfsdk:"description"`
}

func CertificateResource() resource.Resource {
	return &certificateResource{}
}

func (r *certificateResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *certificateResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_certificate"
}

func (r *certificateResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage Certificates in ZStack. " +
			"A Certificate is used for SSL/TLS termination in load balancers and other services.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the Certificate.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Certificate.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"certificate": schema.StringAttribute{
				Required:    true,
				Description: "The PEM-encoded certificate content.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the Certificate.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *certificateResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan certificateModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateCertificateParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateCertificateParamDetail{
			Name:        plan.Name.ValueString(),
			Certificate: plan.Certificate.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	certificate, err := r.client.CreateCertificate(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Certificate",
			"Could not create certificate, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(certificate.UUID)
	plan.Name = types.StringValue(certificate.Name)
	plan.Certificate = types.StringValue(certificate.Certificate)
	plan.Description = stringValueOrNull(certificate.Description)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *certificateResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state certificateModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	certificate, err := findResourceByQuery(r.client.QueryCertificate, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading Certificate",
			"Could not read certificate UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(certificate.UUID)
	state.Name = types.StringValue(certificate.Name)
	state.Certificate = types.StringValue(certificate.Certificate)
	state.Description = stringValueOrNull(certificate.Description)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *certificateResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan certificateModel
	var state certificateModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateCertificateParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateCertificateParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	certificate, err := r.client.UpdateCertificate(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating Certificate",
			"Could not update certificate, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(certificate.UUID)
	plan.Name = types.StringValue(certificate.Name)
	plan.Certificate = types.StringValue(certificate.Certificate)
	plan.Description = stringValueOrNull(certificate.Description)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *certificateResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state certificateModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}


	err := r.client.DeleteCertificate(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError("Error deleting Certificate", "Could not delete certificate, unexpected error: "+err.Error())
		return
	}
}

func (r *certificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
