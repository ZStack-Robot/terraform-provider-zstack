package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

var (
	_ resource.Resource              = &imageResource{}
	_ resource.ResourceWithConfigure = &imageResource{}
)

type imageResource struct {
	client *client.ZSClient
}

type imageResourceModel struct {
	UUID        types.String `tfsdk:"uuid"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Url         types.String `tfsdk:"url"`
	//Mediatype         types.String `tfsdk:"mediatype"`
	Guestostype       types.String `tfsdk:"guestostype"`
	System            types.String `tfsdk:"system"`
	Platform          types.String `tfsdk:"platform"`
	Format            types.String `tfsdk:"format"`
	Backupstorageuuid types.String `tfsdk:"backupstorageuuid"`
	Architecture      types.String `tfsdk:"architecture"`
	Virtio            types.Bool   `tfsdk:"virtio"`
	Type              types.String `tfsdk:"type"`
}

// Configure implements resource.ResourceWithConfigure.
func (r *imageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func ImageResource() resource.Resource {
	return &imageResource{}
}

// Create implements resource.Resource.
func (r *imageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var imageplan imageResourceModel
	diags := req.Plan.Get(ctx, &imageplan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Configuring ZStack client")
	var imagesItem view.ImageView
	imageParam := param.AddImageParam{
		BaseParam: param.BaseParam{
			SystemTags: []string{"bootMode::Legacy"},
		},
		Params: param.AddImageDetailParam{
			Name:        imageplan.Name.ValueString(),
			Description: imageplan.Description.ValueString(),
			Url:         imageplan.Url.ValueString(),
			//MediaType:          param.RootVolumeTemplate,
			GuestOsType:        imageplan.Guestostype.ValueString(),
			System:             false,
			Format:             param.ImageFormat(imageplan.Format.ValueString()), // param.Qcow2,
			Platform:           imageplan.Platform.ValueString(),
			BackupStorageUuids: []string{imageplan.Backupstorageuuid.ValueString()}, //[]string{storage[0].UUID},
			Type:               imageplan.Type.ValueString(),
			ResourceUuid:       "",
			Architecture:       param.Architecture(imageplan.Architecture.ValueString()),
			Virtio:             imageplan.Virtio.ValueBool(),
		},
	}

	ctx = tflog.SetField(ctx, "url", imagesItem.Url)
	image, err := r.client.AddImage(imageParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Add image to ZStack Image storage"+image.Name, "Error "+err.Error(),
		)
		return
	}

	imageplan.UUID = types.StringValue(image.UUID)
	imageplan.Name = types.StringValue(image.Name)
	imageplan.Description = types.StringValue(image.Description)
	imageplan.Url = types.StringValue(image.Url)
	//imageplan.Mediatype = types.StringValue(image.MediaType)
	imageplan.Guestostype = types.StringValue(image.GuestOsType)
	imageplan.System = types.StringValue(image.System)
	imageplan.Format = types.StringValue(image.Format)
	imageplan.Platform = types.StringValue(image.Platform)
	imageplan.Type = types.StringValue(image.Type)
	imageplan.Architecture = types.StringValue(string(image.Architecture))
	imageplan.Virtio = types.BoolValue(image.Virtio)
	imageplan.Backupstorageuuid = types.StringValue(string(image.BackupStorageRefs[0].BackupStorageUuid))
	imageplan.LastUpdated = types.StringValue(image.LastOpDate.String())
	ctx = tflog.SetField(ctx, "url", imagesItem.Url)
	diags = resp.State.Set(ctx, imageplan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete implements resource.Resource.
func (r *imageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state imageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteImage(state.UUID.String(), param.DeleteModeEnforcing)

	if err != nil {
		resp.Diagnostics.AddError("", ""+err.Error())
		return
	}
}

// Metadata implements resource.Resource.
func (r *imageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

// Read implements resource.Resource.
func (r *imageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state imageResourceModel
	req.State.Schema.GetAttributes()

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	image, err := r.client.GetImage(state.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Geting ZStack Image uiid", "Could not read image uuid"+image.Name+": "+err.Error(),
		)
		return
	}

	state.UUID = types.StringValue(image.UUID)
	state.Name = types.StringValue(image.Name)
	state.Url = types.StringValue(image.Url)
	state.LastUpdated = types.StringValue(image.LastOpDate.GoString())
	state.Format = types.StringValue(image.Format)
	state.Description = types.StringValue(image.Description)
	state.Architecture = types.StringValue(string(image.Architecture))
	//state.Mediatype = types.StringValue(image.MediaType)
	state.Virtio = types.BoolValue(image.Virtio)

	//state.Description = types.StringValue(image.Description)

	diags = resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements resource.Resource.
func (r *imageResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"url": schema.StringAttribute{
				Required: true,
			},
			/*
				"mediatype": schema.StringAttribute{
					Optional: true,
				},
			*/
			"guestostype": schema.StringAttribute{
				Optional: true,
			},
			"system": schema.StringAttribute{
				Computed: true,
			},
			"platform": schema.StringAttribute{
				Optional: true,
			},
			"format": schema.StringAttribute{
				Optional: true,
			},
			"backupstorageuuid": schema.StringAttribute{
				Optional: true,
			},
			"architecture": schema.StringAttribute{
				Optional: true,
			},
			"type": schema.StringAttribute{
				Computed: true,
			},
			"virtio": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

// Update implements resource.Resource. //有问题，报PUT http://172.25.16.104:8080/zstack/v1/images {"updateImage":{"description":"test test","name":"C790123newname"}}:
// {"error":{"causes":null,"class":"","code":400,"data":{"fields":null,"id":""},"details":"the body doesn't contain action mapping to the URL[/v1/images]","request":{"body":{"updateImage":{"description":"test
// test","name":"C790123newname"}},"headers":{"Authorization":"*","Content-Length":"67","Content-Type":"application/json;
// charset=utf-8","User-Agent":"zstack-sdk-go/202206","X-Session-Id":"bacd061493b54ac9948b55bd6c31f5c9"},"method":"PUT","url":"http://172.25.16.104:8080/zstack/v1/images"}}}

func (r *imageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var imageplan imageResourceModel
	diags := req.Plan.Get(ctx, &imageplan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	desc := "test test"
	uuid := imageplan.UUID.ValueString()
	params := param.UpdateImageParam{
		UpdateImage: param.UpdateImageDetailParam{
			Name:        imageplan.Name.ValueString(), // "image4chenjtest",
			Description: &desc,                        //imageplan.Description.ValueStringPointer(),
		},
	}
	_, err := r.client.UpdateImage(uuid, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"", ""+err.Error(),
		)
		return
	}

	image, err := r.client.GetImage(imageplan.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("", ""+err.Error())
		return
	}
	imageplan.UUID = types.StringValue(image.UUID)
	imageplan.Name = types.StringValue(image.Name)
	imageplan.Description = types.StringValue(image.Description)
	imageplan.Url = types.StringValue(image.Url)
	//imageplan.Mediatype = types.StringValue(image.MediaType)
	imageplan.Guestostype = types.StringValue(image.GuestOsType)
	imageplan.System = types.StringValue(image.System)
	imageplan.Format = types.StringValue(image.Format)
	imageplan.Platform = types.StringValue(image.Platform)
	imageplan.Type = types.StringValue(image.Type)
	imageplan.Architecture = types.StringValue(string(image.Architecture))
	imageplan.Virtio = types.BoolValue(image.Virtio)
	imageplan.Backupstorageuuid = types.StringValue(string(image.BackupStorageRefs[0].BackupStorageUuid))
	imageplan.LastUpdated = types.StringValue(image.LastOpDate.String())

	diags = resp.State.Set(ctx, imageplan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
