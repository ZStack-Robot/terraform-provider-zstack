// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

var (
	_ resource.Resource                = &imageGroupResource{}
	_ resource.ResourceWithConfigure   = &imageGroupResource{}
	_ resource.ResourceWithImportState = &imageGroupResource{}
)

type imageGroupResource struct {
	client *client.ZSClient
}

type imageGroupResourceModel struct {
	Uuid                    types.String `tfsdk:"uuid"`
	RootVolumeTemplateUuid  types.String `tfsdk:"root_volume_template_uuid"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	DataVolumeTemplateUuids types.Set    `tfsdk:"data_volume_template_uuids"`
	ResourceUuid            types.String `tfsdk:"resource_uuid"`
	TagUuids                types.Set    `tfsdk:"tag_uuids"`
	SystemTags              types.Set    `tfsdk:"system_tags"`
	ImageCount              types.Int64  `tfsdk:"image_count"`
	Status                  types.String `tfsdk:"status"`
}

func ImageGroupResource() resource.Resource {
	return &imageGroupResource{}
}

func (r *imageGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *imageGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image_group"
}

func (r *imageGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a ZStack image group created from a root volume template and optional data volume templates.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the image group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"root_volume_template_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the root volume template used to create the image group.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the image group.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the image group.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"data_volume_template_uuids": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "The UUIDs of data volume templates included in the image group.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplaceIfConfigured(),
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A custom UUID requested when creating the image group.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tag_uuids": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Tag UUIDs attached while creating the image group.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplaceIfConfigured(),
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"system_tags": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "System tags sent while creating the image group.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplaceIfConfigured(),
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"image_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of images in the image group.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status reported by ZStack for the image group.",
			},
		},
	}
}

func (r *imageGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan imageGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	imageGroup, err := r.client.CreateImageGroupFromImage(
		plan.RootVolumeTemplateUuid.ValueString(),
		imageGroupCreateParam(plan),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating image group", err.Error())
		return
	}

	state := imageGroupModelFromView(imageGroup, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *imageGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state imageGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	imageGroup, err := findResourceByGet(r.client.GetImageGroup, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading image group", err.Error())
		return
	}

	rootUuid, dataUuids, err := imageGroupTemplates(r.client, imageGroup.UUID)
	if err != nil {
		resp.Diagnostics.AddError("Error resolving image group templates", err.Error())
		return
	}
	if state.RootVolumeTemplateUuid.IsNull() || state.RootVolumeTemplateUuid.IsUnknown() {
		state.RootVolumeTemplateUuid = types.StringValue(rootUuid)
	} else if state.RootVolumeTemplateUuid.ValueString() != rootUuid {
		resp.Diagnostics.AddError(
			"Invalid image group root template association",
			fmt.Sprintf(
				"Image group %s is associated with root volume template %s instead of the managed template %s.",
				imageGroup.UUID,
				rootUuid,
				state.RootVolumeTemplateUuid.ValueString(),
			),
		)
		return
	}
	state.DataVolumeTemplateUuids = stringSliceToSet(dataUuids)

	refreshed := imageGroupModelFromView(imageGroup, state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &refreshed)...)
}

func (r *imageGroupResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"ZStack does not provide an UpdateImageGroup API. Change a configured attribute to replace the image group.",
	)
}

func (r *imageGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state imageGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.ExpungeImageGroup(state.Uuid.ValueString()); err != nil && !isZStackNotFoundError(err) {
		resp.Diagnostics.AddError("Error expunging image group", err.Error())
	}
}

func (r *imageGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func imageGroupCreateParam(plan imageGroupResourceModel) param.CreateImageGroupFromImageParam {
	createParam := param.CreateImageGroupFromImageParam{
		BaseParam: param.BaseParam{
			SystemTags: setToStringSlice(plan.SystemTags),
		},
		Params: param.CreateImageGroupFromImageParamDetail{
			Name:                    plan.Name.ValueString(),
			DataVolumeTemplateUuids: setToStringSlice(plan.DataVolumeTemplateUuids),
			TagUuids:                setToStringSlice(plan.TagUuids),
		},
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createParam.Params.Description = stringPtr(plan.Description.ValueString())
	}
	if !plan.ResourceUuid.IsNull() && !plan.ResourceUuid.IsUnknown() {
		createParam.Params.ResourceUuid = stringPtr(plan.ResourceUuid.ValueString())
	}
	return createParam
}

func imageGroupModelFromView(imageGroup *view.ImageGroupInventoryView, prior imageGroupResourceModel) imageGroupResourceModel {
	dataUuids := knownSetOrEmpty(prior.DataVolumeTemplateUuids)
	tagUuids := knownSetOrEmpty(prior.TagUuids)
	systemTags := knownSetOrEmpty(prior.SystemTags)

	return imageGroupResourceModel{
		Uuid:                    types.StringValue(imageGroup.UUID),
		RootVolumeTemplateUuid:  prior.RootVolumeTemplateUuid,
		Name:                    types.StringValue(imageGroup.Name),
		Description:             stringValueOrNull(imageGroup.Description),
		DataVolumeTemplateUuids: dataUuids,
		ResourceUuid:            types.StringValue(imageGroup.UUID),
		TagUuids:                tagUuids,
		SystemTags:              systemTags,
		ImageCount:              types.Int64Value(int64(imageGroup.ImageCount)),
		Status:                  types.StringValue(imageGroup.Status),
	}
}

func imageGroupTemplates(cli *client.ZSClient, imageGroupUuid string) (string, []string, error) {
	query := param.NewQueryParam()
	query.AddQ("imageGroupUuid=" + imageGroupUuid)
	refs, err := cli.QueryImageGroupRef(&query)
	if err != nil {
		return "", nil, fmt.Errorf("query image group references: %w", err)
	}
	if len(refs) == 0 {
		return "", nil, fmt.Errorf("image group %s has no image references", imageGroupUuid)
	}

	var rootUuid string
	dataUuids := make([]string, 0, len(refs)-1)
	for _, ref := range refs {
		image, err := findResourceByGet(cli.GetImage, ref.ImageUuid)
		if err != nil {
			return "", nil, fmt.Errorf("read associated image %s: %w", ref.ImageUuid, err)
		}
		switch image.MediaType {
		case "RootVolumeTemplate":
			if rootUuid != "" && rootUuid != image.UUID {
				return "", nil, fmt.Errorf("image group %s has multiple root volume templates", imageGroupUuid)
			}
			rootUuid = image.UUID
		case "DataVolumeTemplate":
			dataUuids = append(dataUuids, image.UUID)
		default:
			return "", nil, fmt.Errorf(
				"image group %s contains image %s with unsupported media type %q",
				imageGroupUuid,
				image.UUID,
				image.MediaType,
			)
		}
	}
	if rootUuid == "" {
		return "", nil, fmt.Errorf("image group %s has no root volume template", imageGroupUuid)
	}
	return rootUuid, dataUuids, nil
}

func setToStringSlice(set types.Set) []string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}
	result := make([]string, 0, len(set.Elements()))
	for _, element := range set.Elements() {
		value, ok := element.(types.String)
		if ok && !value.IsNull() && !value.IsUnknown() && value.ValueString() != "" {
			result = append(result, value.ValueString())
		}
	}
	return result
}

func stringSliceToSet(values []string) types.Set {
	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		if value != "" {
			elements = append(elements, types.StringValue(value))
		}
	}
	return types.SetValueMust(types.StringType, elements)
}

func knownSetOrEmpty(value types.Set) types.Set {
	if value.IsNull() || value.IsUnknown() {
		return types.SetValueMust(types.StringType, []attr.Value{})
	}
	return value
}
