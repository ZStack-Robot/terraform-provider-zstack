// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"
	"terraform-provider-zstack/zstack/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

type gpuDeviceType string

type instanceResource struct {
	client *client.ZSClient
}

var (
	_ resource.Resource                = &instanceResource{}
	_ resource.ResourceWithConfigure   = &instanceResource{}
	_ resource.ResourceWithImportState = &instanceResource{}
)

var networkModelAttrTypes = map[string]attr.Type{
	"uuid":    types.StringType,
	"ip":      types.StringType,
	"netmask": types.StringType,
	"gateway": types.StringType,
}

const (
	mdevDevice gpuDeviceType = "mdevDevice"
	pciDevice  gpuDeviceType = "pciDevice"
)

type diskModel struct {
	VolumeUuid         types.String `tfsdk:"volume_uuid"`
	Size               types.Int64  `tfsdk:"size"`
	OfferingUuid       types.String `tfsdk:"offering_uuid"`
	VirtioSCSI         types.Bool   `tfsdk:"virtio_scsi"`
	PrimaryStorageUuid types.String `tfsdk:"primary_storage_uuid"`
	CephPoolName       types.String `tfsdk:"ceph_pool_name"`
}

type gpuModel struct {
	Uuid types.String `tfsdk:"uuid"`
	Type types.String `tfsdk:"type"`
}

type gpuSpecModel struct {
	Uuid   types.String `tfsdk:"uuid"`
	Type   types.String `tfsdk:"type"`
	Number types.Int64  `tfsdk:"number"`
}

type vmInstanceDataSourceModel struct {
	Uuid                 types.String `tfsdk:"uuid"`
	Name                 types.String `tfsdk:"name"`
	ImageUuid            types.String `tfsdk:"image_uuid"`
	NetworkInterfaces    types.List   `tfsdk:"network_interfaces"`
	RootDisk             types.Object `tfsdk:"root_disk"`
	DataDisks            types.List   `tfsdk:"data_disks"`
	ZoneUuid             types.String `tfsdk:"zone_uuid"`
	ClusterUuid          types.String `tfsdk:"cluster_uuid"`
	HostUuid             types.String `tfsdk:"host_uuid"`
	Description          types.String `tfsdk:"description"`
	InstanceOfferingUuid types.String `tfsdk:"instance_offering_uuid"`
	Strategy             types.String `tfsdk:"strategy"`
	MemorySize           types.Int64  `tfsdk:"memory_size"`
	CPUNum               types.Int64  `tfsdk:"cpu_num"`
	NeverStop            types.Bool   `tfsdk:"never_stop"`
	Marketplace          types.Bool   `tfsdk:"marketplace"`
	GPUDevices           types.List   `tfsdk:"gpu_devices"`
	GPUSpecs             types.Object `tfsdk:"gpu_device_specs"`
	UserData             types.String `tfsdk:"user_data"`
	VMNics               types.List   `tfsdk:"vm_nics"`
	Expunge              types.Bool   `tfsdk:"expunge"`
	HookScript           types.String `tfsdk:"hook_script"`
	Platform             types.String `tfsdk:"platform"`
	GuestOsType          types.String `tfsdk:"guest_os_type"`
	Architecture         types.String `tfsdk:"architecture"`
}

type NicsModel struct {
	Uuid    types.String `tfsdk:"uuid"`
	Ip      types.String `tfsdk:"ip"`
	Netmask types.String `tfsdk:"netmask"`
	Gateway types.String `tfsdk:"gateway"`
}

type NetworkInterfaceModel struct {
	L3NetworkUuid types.String `tfsdk:"l3_network_uuid"`
	DefaultL3     types.Bool   `tfsdk:"default_l3"`
	StaticIp      types.String `tfsdk:"static_ip"`
}

func InstanceResource() resource.Resource {
	return &instanceResource{}
}

// Configure implements resource.ResourceWithConfigure.
func (r *instanceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. jiajian.chi@zstack.io", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata implements resource.Resource.
func (r *instanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

// Schema implements resource.Resource.
func (r *instanceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource allows you to manage virtual machine (VM) instances in ZStack. " +
			"A VM instance represents a virtualized compute resource that can be created, updated, and deleted. " +
			"You can define the VM's properties, such as its name, image, network configuration, disks, and GPU devices.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the VM instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VM instance.",
			},
			"network_interfaces": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Defines network interfaces attached to the VM. Each NIC corresponds to an L3 network, and optionally configures a static IP.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"l3_network_uuid": schema.StringAttribute{
							Required:    true,
							Description: "The UUID of the L3 network for this NIC.",
						},
						"default_l3": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							Description: "Whether this NIC is the default route NIC. " +
								"If omitted on every NIC, the first NIC is automatically chosen as the default. " +
								"After Create the server-resolved value is reflected back into state.",
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"static_ip": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Static IP address to assign. Optional — if omitted, the server picks one and reports it back. The format will be converted to system tag `staticIp::<l3_uuid>::<ip>`.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"vm_nics": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Required:    true,
							Description: "The UUID of the network.",
						},
						"ip": schema.StringAttribute{
							Computed:    true,
							Description: "The IP address assigned to the network.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"netmask": schema.StringAttribute{
							Computed:    true,
							Description: "The netmask of the network.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"gateway": schema.StringAttribute{
							Computed:    true,
							Description: "The gateway of the network.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
				Computed:    true,
				Description: "The IP address assigned to the VM instance.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_offering_uuid": schema.StringAttribute{
				Optional: true,
				Description: "The UUID of the instance offering used by the VM. Required if using instance offering uuid to create instances. " +
					"  Mutually exclusive with `cpu_num` and `memory_size`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"image_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the image used to create the VM instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"root_disk": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"volume_uuid": schema.StringAttribute{
						Computed:    true,
						Description: "The UUID of the root volume backing this disk (assigned by the server after the VM is created).",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"offering_uuid": schema.StringAttribute{
						Optional:    true,
						Description: "The UUID of the disk offering for the root disk.",
					},
					"size": schema.Int64Attribute{
						Optional:    true,
						Description: "The size of the root disk in gigabytes (GB).",
					},
					"primary_storage_uuid": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The UUID of the primary storage for the root disk. If not specified, the server picks one and returns it here.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"ceph_pool_name": schema.StringAttribute{
						Optional:    true,
						Description: "The Ceph pool name for the root disk.",
					},
					"virtio_scsi": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether the root disk uses Virtio-SCSI.",
					},
				},
				Optional:    true,
				Description: "The configuration for the root disk of the VM instance.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"data_disks": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"volume_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the data volume backing this disk (assigned by the server after the VM is created).",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"offering_uuid": schema.StringAttribute{
							Optional:    true,
							Description: "The UUID of the disk offering for the data disk.",
						},
						"size": schema.Int64Attribute{
							Optional:    true,
							Description: "The size of the data disk in gigabytes (GB).",
						},

						"primary_storage_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the primary storage for the data disk.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},

						"ceph_pool_name": schema.StringAttribute{
							Optional:    true,
							Description: "The Ceph pool name for the data disk.",
						},
						"virtio_scsi": schema.BoolAttribute{
							Optional:    true,
							Description: "Whether the data disk uses Virtio-SCSI.",
						},
					},
				},
				Optional:    true,
				Description: "The configuration for additional data disks.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"gpu_device_specs": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"uuid": schema.StringAttribute{
						Optional:    true,
						Description: "The UUID of the GPU device.",
					},
					"type": schema.StringAttribute{
						Optional:    true,
						Description: "The type of the GPU device. Must be one of: `mdevDevice` or `pciDevice`.",
						Validators: []validator.String{
							stringvalidator.OneOf("mdevDevice", "pciDevice"),
						},
					},
					"number": schema.Int64Attribute{
						Optional:    true,
						Description: "The number of GPU devices assigned.",
					},
				},
				Optional:    true,
				Description: "The GPU specifications for the VM instance.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"gpu_devices": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Optional:    true,
							Description: "The UUID of the GPU device.",
						},
						"type": schema.StringAttribute{
							Optional:    true,
							Description: "The type of the GPU device.  Must be one of: `mdevDevice` or `pciDevice`.",
							Validators: []validator.String{
								stringvalidator.OneOf("mdevDevice", "pciDevice"),
							},
						},
					},
				},
				Optional:    true,
				Description: "A list of GPU devices assigned to the VM instance.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"zone_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the zone where the VM instance is deployed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the cluster where the VM instance is deployed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"host_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the host where the VM instance is running.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description of the VM instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"memory_size": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The memory size allocated to the VM instance in megabytes (MB). When used together with `cpu_num`, the `instance_offering_uuid` is not required.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"cpu_num": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The number of CPUs allocated to the VM instance.  When used together with `memory_size`, the `instance_offering_uuid` is not required.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"strategy": schema.StringAttribute{
				Optional:    true,
				Description: "The deployment strategy for the VM instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_data": schema.StringAttribute{
				Optional:    true,
				Description: "User data injected into the VM instance at boot time.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"never_stop": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the VM instance should never stop automatically.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"expunge": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates if the instance should be expunged after deletion.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"marketplace": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates whether the VM instance is a marketplace instance.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"hook_script": schema.StringAttribute{
				Optional:    true,
				Description: "The uuid of hook script. Create Instance with custom xml Hook.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"platform": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "The platform of the guest OS (e.g. `Linux`, `Windows`, `Other`, `Paravirtualization`). " +
					"If unset the server inherits it from the image. " +
					"Updatable in place via the `UpdateVmInstance` API on a running cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"guest_os_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "The guest OS type / distribution (free-form string, e.g. `CentOS 7`, `Windows Server 2019`). " +
					"Server reports it back after Create. Updatable in place via `UpdateVmInstance`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"architecture": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "The CPU architecture of the guest (`x86_64`, `aarch64`, `mips64el`, etc.). " +
					"Inherited from the image when unset. Changing it requires the VM to be replaced.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *instanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vmInstanceDataSourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var rootDiskPlan diskModel
	var dataDisksPlan []diskModel

	var primaryStorageUuidForRootVolume *string
	hostUuid := ""
	clusterUuid := ""
	zoneUuid := ""
	//hook_script := ""
	var rootDiskSystemTags []string
	var dataDiskSizes []int64
	var dataDiskOfferingUuids []string
	var dataVolumeSystemTagsOnIndex []string
	var dataDiskSystemTags []string

	// SET ROOT DISK
	if !plan.RootDisk.IsNull() {
		diags = plan.RootDisk.As(ctx, &rootDiskPlan, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if rootDiskPlan.OfferingUuid.IsNull() && rootDiskPlan.Size.IsNull() {
			resp.Diagnostics.AddError(
				"Error creating VM Instance",
				"Could not create vm instance, rootDiskPlan offering_uuid and size cannot be null at the same time.",
			)
			return
		}
		err := isDiskParamValid(r, rootDiskPlan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating VM Instance",
				fmt.Sprintf("Could not create vm instance, invalid rootDiskPlan param: %v", err),
			)
			return
		}
		if !rootDiskPlan.PrimaryStorageUuid.IsNull() && rootDiskPlan.PrimaryStorageUuid.ValueString() != "" {
			primaryStorageUuidForRootVolume = rootDiskPlan.PrimaryStorageUuid.ValueStringPointer()
		}

		if !rootDiskPlan.CephPoolName.IsNull() && rootDiskPlan.CephPoolName.ValueString() != "" {
			rootDiskSystemTags = append(rootDiskSystemTags, fmt.Sprintf("ceph::rootPoolName::%s", rootDiskPlan.CephPoolName.ValueString()))
		}

		if !rootDiskPlan.Size.IsNull() && !rootDiskPlan.Size.IsUnknown() {
			rootDiskPlan.Size = types.Int64Value(utils.GBToBytes(rootDiskPlan.Size.ValueInt64()))
		}
	}

	// SET DATA DISK
	if !plan.DataDisks.IsNull() {
		plan.DataDisks.ElementsAs(ctx, &dataDisksPlan, false)

		for _, disk := range dataDisksPlan {
			if !disk.OfferingUuid.IsNull() {
				dataDiskOfferingUuids = append(dataDiskOfferingUuids, disk.OfferingUuid.ValueString())
			} else if !disk.Size.IsNull() && !disk.Size.IsUnknown() {
				dataDiskSizes = append(dataDiskSizes, utils.GBToBytes(disk.Size.ValueInt64()))
				if disk.VirtioSCSI.ValueBool() {
					dataVolumeSystemTagsOnIndex = append(dataVolumeSystemTagsOnIndex, "capability::virtio-scsi")
				}
			} else {
				resp.Diagnostics.AddError(
					"Error creating VM Instance",
					"Could not create vm instance, dataDisk offering_uuid and size cannot be null at the same time.",
				)
				return
			}
		}

		//only support one type data disk now
		if len(dataDisksPlan) > 0 {
			err := isDiskParamValid(r, dataDisksPlan[0])
			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating VM Instance",
					fmt.Sprintf("Could not create vm instance, invalid dataDisk param: %v", err),
				)
				return
			}
			if !dataDisksPlan[0].CephPoolName.IsNull() && dataDisksPlan[0].CephPoolName.ValueString() != "" {
				dataDiskSystemTags = append(dataDiskSystemTags, fmt.Sprintf("ceph::pool::%s", dataDisksPlan[0].CephPoolName.ValueString()))
			}
			dataDiskSystemTags = append(dataDiskSystemTags, dataVolumeSystemTagsOnIndex...)
		}
	}

	// SET NETWORK
	if plan.NetworkInterfaces.IsNull() || len(plan.NetworkInterfaces.Elements()) == 0 {
		resp.Diagnostics.AddError(
			"Error creating VM Instance",
			"Could not create vm instance, `network_interfaces` cannot be null or empty. At least one L3 network must be specified.",
		)
		return
	}

	var inputNics []NetworkInterfaceModel
	diags = plan.NetworkInterfaces.ElementsAs(ctx, &inputNics, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var l3NetworkUuids []string
	var defaultL3Uuid string
	var systemTags []string

	var createNics []NetworkInterfaceModel
	for _, nic := range inputNics {
		l3uuid := nic.L3NetworkUuid.ValueString()
		l3NetworkUuids = append(l3NetworkUuids, l3uuid)

		// Treat unset (null/unknown) the same as false for picking a default,
		// but record true matches in defaultL3Uuid.
		if !nic.DefaultL3.IsNull() && !nic.DefaultL3.IsUnknown() && nic.DefaultL3.ValueBool() {
			if defaultL3Uuid != "" && defaultL3Uuid != l3uuid {
				resp.Diagnostics.AddError(
					"Error creating VM Instance",
					"Could not create vm instance: more than one network_interfaces entry has default_l3=true. Exactly one (or zero) NIC may be marked default.",
				)
				return
			}
			defaultL3Uuid = l3uuid
		}

		var staticIp types.String
		if nic.StaticIp.IsNull() || nic.StaticIp.ValueString() == "" {
			staticIp = types.StringNull()
		} else {
			staticIp = nic.StaticIp
			systemTags = append(systemTags, fmt.Sprintf("staticIp::%s::%s", l3uuid, staticIp.ValueString()))
		}

		createNics = append(createNics, NetworkInterfaceModel{
			L3NetworkUuid: nic.L3NetworkUuid,
			DefaultL3:     nic.DefaultL3,
			StaticIp:      staticIp,
		})
	}

	// If no NIC was flagged as default, pick the first one. The server requires
	// a defaultL3NetworkUuid; this matches the "unset means first NIC is the
	// default" convention used by the ZStack UI.
	if defaultL3Uuid == "" && len(createNics) > 0 {
		defaultL3Uuid = createNics[0].L3NetworkUuid.ValueString()
	}

	//SET XML HOOK SCRIPT
	if !plan.HookScript.IsNull() && plan.HookScript.ValueString() != "" {
		systemTags = append(systemTags, fmt.Sprintf("xmlHook::%s", plan.HookScript.ValueString()))
	}

	// SET IMAGE
	image, err := r.client.GetImage(plan.ImageUuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating VM Instance",
			fmt.Sprintf("Could not create vm instance, failed to find image %s: %v", plan.ImageUuid.ValueString(), err),
		)
		return
	}

	if image.Status != "Ready" {
		resp.Diagnostics.AddError(
			"Error creating VM Instance",
			fmt.Sprintf("Could not create vm instance, image %s status is %s, not Ready.", plan.ImageUuid.ValueString(), image.State),
		)
		return
	}

	if image.State != "Enabled" {
		resp.Diagnostics.AddError(
			"Error creating VM Instance",
			fmt.Sprintf("Could not create vm instance, image %s state is %s, not Enabled.", plan.ImageUuid.ValueString(), image.State),
		)
		return
	}

	// SET HOST UUID
	if !plan.HostUuid.IsNull() && plan.HostUuid.ValueString() != "" {
		hostUuid = plan.HostUuid.ValueString()
	}

	// SET CLUSTER UUID
	if !plan.ClusterUuid.IsNull() && plan.ClusterUuid.ValueString() != "" {
		clusterUuid = plan.ClusterUuid.ValueString()
	}

	// SET ZONE UUID
	if !plan.ZoneUuid.IsNull() && plan.ZoneUuid.ValueString() != "" {
		zoneUuid = plan.ZoneUuid.ValueString()
	}

	// SET SYSTEM TAG
	//systemTags := []string{"resourceConfig::vm::vm.clock.track::guest", "cdroms::Empty::None::None"}
	if plan.Marketplace.ValueBool() {
		systemTags = append(systemTags, "marketplace::true")
	}
	if !plan.NeverStop.IsNull() && !plan.NeverStop.IsUnknown() && plan.NeverStop.ValueBool() {
		systemTags = append(systemTags, "ha::NeverStop")
	}

	if !plan.UserData.IsNull() && plan.UserData.ValueString() != "" {
		systemTags = append(systemTags, fmt.Sprintf("userdata::%s", plan.UserData.ValueString()))
	}

	// SET GPU
	// when gpu device is set, gpu spec is not work
	if !plan.GPUDevices.IsNull() {
		var gpuDevicesPlan []gpuModel
		plan.GPUDevices.ElementsAs(ctx, &gpuDevicesPlan, false)

		for _, gpuDevice := range gpuDevicesPlan {
			if gpuDevice.Type.ValueString() == string(mdevDevice) {
				systemTags = append(systemTags, fmt.Sprintf("mdevDevice::%s", gpuDevice.Uuid.ValueString()))
			} else if gpuDevice.Type.ValueString() == string(pciDevice) {
				systemTags = append(systemTags, fmt.Sprintf("pciDevice::%s", gpuDevice.Uuid.ValueString()))
			} else {
				resp.Diagnostics.AddError(
					"Error creating VM Instance",
					fmt.Sprintf("Could not create vm instance, gpu type %s is invalid.", gpuDevice.Type.ValueString()),
				)
				return
			}
		}
	} else if !plan.GPUSpecs.IsNull() {
		var gpuSpecPlan gpuSpecModel
		diags = plan.GPUSpecs.As(ctx, &gpuSpecPlan, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		number := 1
		if !gpuSpecPlan.Number.IsNull() && !gpuSpecPlan.Number.IsUnknown() {
			number = int(gpuSpecPlan.Number.ValueInt64())
		}

		if gpuSpecPlan.Type.ValueString() == string(mdevDevice) {
			systemTags = append(systemTags, fmt.Sprintf("mdevDeviceSpec::%s::%d", gpuSpecPlan.Uuid.ValueString(), number))
		} else if gpuSpecPlan.Type.ValueString() == string(pciDevice) {
			systemTags = append(systemTags, fmt.Sprintf("pciDeviceSpec::%s::%d", gpuSpecPlan.Uuid.ValueString(), number))
		} else {
			resp.Diagnostics.AddError(
				"Error creating VM Instance",
				fmt.Sprintf("Could not create vm instance, gpu type %s is invalid.", gpuSpecPlan.Type.ValueString()),
			)
			return
		}
	}

	//SET OTHER PARAM
	if !plan.Strategy.IsNull() {
		strategyValue := plan.Strategy.ValueString()
		if strategyValue != "InstantStart" && strategyValue != "CreateStopped" {
			resp.Diagnostics.AddError(
				"Error creating VM Instance",
				fmt.Sprintf("Could not create vm instance, strategy %s is invalid. Valid value is InstantStart or CreateStopped.", plan.Strategy.ValueString()),
			)
			return
		}
	}

	// Check if instance_offering_uuid is provided
	var memorySize int64
	var cpuNum int64
	if !plan.InstanceOfferingUuid.IsNull() && plan.InstanceOfferingUuid.ValueString() != "" {
		instanceOffering, err := r.client.GetInstanceOffering(plan.InstanceOfferingUuid.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating VM Instance",
				fmt.Sprintf("Could not create vm instance, failed to get instance offering %s: %v", plan.InstanceOfferingUuid.ValueString(), err),
			)
			return
		}
		memorySize = instanceOffering.MemorySize
		cpuNum = int64(instanceOffering.CpuNum)
	} else {
		if plan.MemorySize.IsNull() || plan.CPUNum.IsNull() {
			resp.Diagnostics.AddError(
				"Error creating VM Instance",
				"Could not create vm instance, memory_size and cpu_num must be provided if instance_offering_uuid is not set.",
			)
			return
		}

		memorySize = utils.MBToBytes(plan.MemorySize.ValueInt64())
		cpuNum = plan.CPUNum.ValueInt64()
	}

	createVmInstanceParam := param.CreateVmInstanceParam{
		BaseParam: param.BaseParam{
			SystemTags: systemTags,
			UserTags:   nil,
			RequestIp:  "",
		},
		Params: param.CreateVmInstanceParamDetail{
			Name:                            plan.Name.ValueString(),
			InstanceOfferingUuid:            stringPtrOrNil(plan.InstanceOfferingUuid.ValueString()),
			ImageUuid:                       stringPtr(plan.ImageUuid.ValueString()),
			L3NetworkUuids:                  l3NetworkUuids,
			Type:                            stringPtr("UserVm"),
			RootDiskOfferingUuid:            stringPtrOrNil(rootDiskPlan.OfferingUuid.ValueString()),
			RootDiskSize:                    rootDiskPlan.Size.ValueInt64Pointer(),
			PrimaryStorageUuidForRootVolume: primaryStorageUuidForRootVolume,
			DataDiskSizes:                   dataDiskSizes,
			DataDiskOfferingUuids:           dataDiskOfferingUuids,
			ZoneUuid:                        stringPtrOrNil(zoneUuid),
			ClusterUuid:                     stringPtrOrNil(clusterUuid),
			HostUuid:                        stringPtrOrNil(hostUuid),
			Description:                     stringPtrOrNil(plan.Description.ValueString()),
			DefaultL3NetworkUuid:            stringPtrOrNil(defaultL3Uuid),
			TagUuids:                        nil,
			Strategy:                        stringPtrOrNil(plan.Strategy.ValueString()),
			MemorySize:                      &memorySize,
			CpuNum:                          intPtr(int(cpuNum)),
			Platform:                        stringPtrOrNil(plan.Platform.ValueString()),
			GuestOsType:                     stringPtrOrNil(plan.GuestOsType.ValueString()),
			Architecture:                    stringPtrOrNil(plan.Architecture.ValueString()),
			RootVolumeSystemTags:            rootDiskSystemTags,
			DataVolumeSystemTags:            dataDiskSystemTags,
		},
	}

	instance, err := r.client.CreateVmInstance(createVmInstanceParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating VM Instance",
			fmt.Sprintf("Could not create vm instance, unexpected error: %v", err),
		)
		return
	}

	plan.Uuid = types.StringValue(instance.UUID)
	plan.Name = types.StringValue(instance.Name)
	plan.Description = stringValueOrNull(instance.Description)
	plan.MemorySize = types.Int64Value(utils.BytesToMB(instance.MemorySize))
	plan.CPUNum = types.Int64Value(int64(instance.CpuNum))
	plan.Platform = stringValueOrNull(instance.Platform)
	plan.GuestOsType = stringValueOrNull(instance.GuestOsType)
	plan.Architecture = stringValueOrNull(instance.Architecture)

	var updatedNics []NetworkInterfaceModel
	for _, nic := range createNics {
		var realIP string
		for _, vmNic := range instance.VmNics {
			if vmNic.L3NetworkUuid == nic.L3NetworkUuid.ValueString() {
				realIP = vmNic.Ip
				break
			}
		}

		staticIp := nic.StaticIp
		if staticIp.IsNull() || staticIp.ValueString() == "" {
			staticIp = types.StringValue(realIP)
		}

		// Reflect the server-resolved default_l3 back into state. The server
		// reports its choice via instance.DefaultL3NetworkUuid; we always echo
		// that (true/false) regardless of whether the user supplied an explicit
		// value, so plan and post-apply state agree.
		defaultL3 := types.BoolValue(instance.DefaultL3NetworkUuid != "" && nic.L3NetworkUuid.ValueString() == instance.DefaultL3NetworkUuid)

		updatedNics = append(updatedNics, NetworkInterfaceModel{
			L3NetworkUuid: nic.L3NetworkUuid,
			DefaultL3:     defaultL3,
			StaticIp:      staticIp,
		})
	}

	networkInterfaceAttrTypes := map[string]attr.Type{
		"l3_network_uuid": types.StringType,
		"default_l3":      types.BoolType,
		"static_ip":       types.StringType,
	}

	networkInterfacesList, diags := types.ListValueFrom(ctx,
		types.ObjectType{AttrTypes: networkInterfaceAttrTypes},
		updatedNics)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.NetworkInterfaces = networkInterfacesList

	var diskModelAttrTypes = map[string]attr.Type{
		"volume_uuid":          types.StringType,
		"offering_uuid":        types.StringType,
		"size":                 types.Int64Type,
		"primary_storage_uuid": types.StringType,
		"ceph_pool_name":       types.StringType,
		"virtio_scsi":          types.BoolType,
	}

	// Split volumes into root + data so we can populate the Computed
	// volume_uuid field on both root_disk and data_disks. The server returns
	// volumes in a single AllVolumes slice; the root volume is identified by
	// Type == "Root" (or by matching instance.RootVolumeUuid).
	var rootVolume *view.VolumeInventoryView
	var dataVolumes []view.VolumeInventoryView
	for i := range instance.AllVolumes {
		v := instance.AllVolumes[i]
		if rootVolume == nil && (v.Type == "Root" || v.UUID == instance.RootVolumeUuid) {
			rootVolume = &v
			continue
		}
		dataVolumes = append(dataVolumes, v)
	}

	// Populate root_disk in state. The schema makes root_disk Optional; if the
	// user did not declare it we leave it null. If they did, fill in the
	// Computed volume_uuid so Terraform's post-apply consistency check passes.
	if !plan.RootDisk.IsNull() {
		if rootVolume != nil {
			rootDiskPlan.VolumeUuid = types.StringValue(rootVolume.UUID)
			if rootDiskPlan.PrimaryStorageUuid.IsNull() || rootDiskPlan.PrimaryStorageUuid.ValueString() == "" {
				rootDiskPlan.PrimaryStorageUuid = stringValueOrNull(rootVolume.PrimaryStorageUuid)
			}
		} else {
			rootDiskPlan.VolumeUuid = types.StringNull()
		}
		rootDiskObj, rdDiags := types.ObjectValueFrom(ctx, diskModelAttrTypes, rootDiskPlan)
		resp.Diagnostics.Append(rdDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.RootDisk = rootDiskObj
	}

	if !plan.DataDisks.IsNull() {
		var dataDisksPlan []diskModel
		plan.DataDisks.ElementsAs(ctx, &dataDisksPlan, false)

		for i, disk := range dataVolumes {
			if i < len(dataDisksPlan) {
				dataDisksPlan[i].VolumeUuid = types.StringValue(disk.UUID)
				dataDisksPlan[i].PrimaryStorageUuid = types.StringValue(disk.PrimaryStorageUuid)
			}
		}

		dataDisksList, diags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: diskModelAttrTypes,
		}, dataDisksPlan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.DataDisks = dataDisksList

	}

	var vmNics []NicsModel
	for _, nic := range instance.VmNics {
		vmNics = append(vmNics, NicsModel{
			Uuid:    types.StringValue(nic.UUID),
			Ip:      types.StringValue(nic.Ip),
			Netmask: types.StringValue(nic.Netmask),
			Gateway: types.StringValue(nic.Gateway),
		})
	}

	plan.VMNics, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: networkModelAttrTypes}, vmNics)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

}

// Read implements resource.Resource.
func (r *instanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vmInstanceDataSourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	vm, err := findResourceByGet(r.client.GetVmInstance, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			tflog.Warn(ctx, "vm not found, removing from state", map[string]interface{}{
				"uuid": state.Uuid.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading VM Instance",
			"Could not read VM instance UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(vm.UUID)
	state.Name = types.StringValue(vm.Name)
	state.Description = stringValueOrNull(vm.Description)
	state.ImageUuid = types.StringValue(vm.ImageUuid)
	state.MemorySize = types.Int64Value(utils.BytesToMB(vm.MemorySize))
	state.CPUNum = types.Int64Value(int64(vm.CpuNum))
	state.Platform = stringValueOrNull(vm.Platform)
	state.GuestOsType = stringValueOrNull(vm.GuestOsType)
	state.Architecture = stringValueOrNull(vm.Architecture)

	updatedDataDisks, err := syncInstanceDataDisksFromVM(ctx, state, vm)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading VM Instance",
			"Could not sync VM data disk state: "+err.Error(),
		)
		return
	}
	state.DataDisks = updatedDataDisks.DataDisks

	var vmNics []NicsModel
	for _, nic := range vm.VmNics {
		vmNics = append(vmNics, NicsModel{
			Uuid:    types.StringValue(nic.UUID),
			Ip:      types.StringValue(nic.Ip),
			Netmask: types.StringValue(nic.Netmask),
			Gateway: types.StringValue(nic.Gateway),
		})
	}

	state.VMNics, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: networkModelAttrTypes}, vmNics)
	resp.Diagnostics.Append(diags...)

	networkInterfaces := normalizeNetworkInterfacesFromVM(vm)
	state.NetworkInterfaces, _ = types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"l3_network_uuid": types.StringType,
			"default_l3":      types.BoolType,
			"static_ip":       types.StringType,
		},
	}, networkInterfaces)

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *instanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state vmInstanceDataSourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Error updating VM Instance",
			"Could not update vm instance, UUID is empty.")
		return
	}

	//"uuid of vm is empty, cannot upgrade vm."

	uuid := state.Uuid.ValueString()

	updateVmInstanceParam := param.UpdateVmInstanceParam{}
	updateVm := false

	if plan.Name.ValueString() != state.Name.ValueString() {
		updateVmInstanceParam.Params.Name = plan.Name.ValueString()
		updateVm = true
	}
	if plan.Description.ValueString() != state.Description.ValueString() {
		updateVmInstanceParam.Params.Description = plan.Description.ValueStringPointer()
		updateVm = true

	}
	if !plan.CPUNum.IsNull() && !plan.CPUNum.IsUnknown() && plan.CPUNum.ValueInt64() != state.CPUNum.ValueInt64() {
		updateVmInstanceParam.Params.CpuNum = utils.TfInt64ToIntPointer(plan.CPUNum)
		updateVm = true
	}
	if !plan.MemorySize.IsNull() && !plan.MemorySize.IsUnknown() && plan.MemorySize.ValueInt64() != state.MemorySize.ValueInt64() {
		memorySizeBytes := utils.MBToBytes(plan.MemorySize.ValueInt64())
		updateVmInstanceParam.Params.MemorySize = &memorySizeBytes
		updateVm = true
	}

	// platform / guest_os_type are online-modifiable via UpdateVmInstance.
	// architecture is RequiresReplace at the schema level, so it never reaches
	// here as a diff.
	if !plan.Platform.IsNull() && !plan.Platform.IsUnknown() && plan.Platform.ValueString() != state.Platform.ValueString() {
		v := plan.Platform.ValueString()
		updateVmInstanceParam.Params.Platform = &v
		updateVm = true
	}
	if !plan.GuestOsType.IsNull() && !plan.GuestOsType.IsUnknown() && plan.GuestOsType.ValueString() != state.GuestOsType.ValueString() {
		v := plan.GuestOsType.ValueString()
		updateVmInstanceParam.Params.GuestOsType = &v
		updateVm = true
	}

	if updateVm {
		if _, err := r.client.UpdateVmInstance(uuid, updateVmInstanceParam); err != nil {
			resp.Diagnostics.AddError(
				"Error updating VM Instance",
				"Could not update vm instance, unexpected error: "+err.Error())
			return
		}

		// Refresh from server to keep Update / Read state-construction in lockstep.
		vm, err := findResourceByGet(r.client.GetVmInstance, uuid)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating VM Instance",
				"Could not refresh vm instance after update: "+err.Error())
			return
		}

		plan, err = buildUpdatedStateFromVM(ctx, plan, vm)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating VM Instance",
				"Could not rebuild vm instance state after update: "+err.Error())
			return
		}

		diags := resp.State.Set(ctx, &plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	//return
}

// Delete implements resource.Resource.
func (r *instanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vmInstanceDataSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	volumeUuids, err := instanceDataVolumeUUIDsFromState(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting VM Instance", "Could not determine VM data volume UUIDs from state: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "Deleting vm instance "+state.Uuid.String())

	//Delete existing vm instance
	err = r.client.DestroyVmInstance(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error destroying VM Instance", "Could not destroy vm instance, unexpected error: "+err.Error(),
		)
		return
	}

	//Delete vm data volume
	for _, uuid := range volumeUuids {
		err = r.client.DeleteDataVolume(uuid, param.DeleteModePermissive)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deleting Data Volume", "Could not delete data volume UUID "+uuid+": "+err.Error(),
			)
			return
		}
	}

	expunge := false
	if !state.Expunge.IsNull() && !state.Expunge.IsUnknown() {
		expunge = state.Expunge.ValueBool()
	}

	if expunge {
		tflog.Info(ctx, fmt.Sprintf("expunge instance %s", state.Uuid.ValueString()))
		//Expunge vm instance
		err = r.client.ExpungeVmInstance(state.Uuid.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error expunging VM Instance", "Could not expunge vm instance, unexpected error: "+err.Error(),
			)
			return
		}

		//Expunge vm data volume
		for _, uuid := range volumeUuids {
			err = r.client.ExpungeDataVolume(uuid)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error expunging Data Volume", "Could not expunge data volume UUID "+uuid+": "+err.Error(),
				)
				return
			}
		}
	}

}

func instanceDiskModelAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"volume_uuid":          types.StringType,
		"offering_uuid":        types.StringType,
		"size":                 types.Int64Type,
		"primary_storage_uuid": types.StringType,
		"ceph_pool_name":       types.StringType,
		"virtio_scsi":          types.BoolType,
	}
}

func syncInstanceDataDisksFromVM(ctx context.Context, state vmInstanceDataSourceModel, vm *view.VmInstanceInventoryView) (vmInstanceDataSourceModel, error) {
	if state.DataDisks.IsNull() || state.DataDisks.IsUnknown() || len(state.DataDisks.Elements()) == 0 {
		return state, nil
	}

	var dataDisks []diskModel
	if diags := state.DataDisks.ElementsAs(ctx, &dataDisks, false); diags.HasError() {
		return state, fmt.Errorf("failed decoding data_disks state: %v", diags)
	}

	dataVolumes := make([]view.VolumeInventoryView, 0, len(vm.AllVolumes))
	for _, volume := range vm.AllVolumes {
		if volume.Type == "Data" {
				dataVolumes = append(dataVolumes, volume)
		}
	}

	for i := range dataDisks {
		if i >= len(dataVolumes) {
			break
		}
		dataDisks[i].VolumeUuid = types.StringValue(dataVolumes[i].UUID)
		dataDisks[i].PrimaryStorageUuid = types.StringValue(dataVolumes[i].PrimaryStorageUuid)
	}

	listValue, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: instanceDiskModelAttrTypes()}, dataDisks)
	if diags.HasError() {
		return state, fmt.Errorf("failed encoding data_disks state: %v", diags)
	}
	state.DataDisks = listValue
	return state, nil
}

func normalizeNetworkInterfacesFromVM(vm *view.VmInstanceInventoryView) []NetworkInterfaceModel {
	networkInterfaces := make([]NetworkInterfaceModel, 0, len(vm.VmNics))
	for _, nic := range vm.VmNics {
		networkInterfaces = append(networkInterfaces, NetworkInterfaceModel{
			L3NetworkUuid: types.StringValue(nic.L3NetworkUuid),
			DefaultL3:     types.BoolValue(vm.DefaultL3NetworkUuid != "" && nic.L3NetworkUuid == vm.DefaultL3NetworkUuid),
			StaticIp:      types.StringValue(nic.Ip),
		})
	}
	return networkInterfaces
}

func buildUpdatedStateFromVM(ctx context.Context, current vmInstanceDataSourceModel, vm *view.VmInstanceInventoryView) (vmInstanceDataSourceModel, error) {
	updated := current
	updated.Uuid = types.StringValue(vm.UUID)
	updated.Name = types.StringValue(vm.Name)
	updated.Description = stringValueOrNull(vm.Description)
	updated.ImageUuid = types.StringValue(vm.ImageUuid)
	updated.MemorySize = types.Int64Value(utils.BytesToMB(vm.MemorySize))
	updated.CPUNum = types.Int64Value(int64(vm.CpuNum))
	updated.Platform = stringValueOrNull(vm.Platform)
	updated.GuestOsType = stringValueOrNull(vm.GuestOsType)
	updated.Architecture = stringValueOrNull(vm.Architecture)

	var vmNics []NicsModel
	for _, nic := range vm.VmNics {
		vmNics = append(vmNics, NicsModel{
			Uuid:    types.StringValue(nic.UUID),
			Ip:      types.StringValue(nic.Ip),
			Netmask: types.StringValue(nic.Netmask),
			Gateway: types.StringValue(nic.Gateway),
		})
	}
	vmNicsValue, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: networkModelAttrTypes}, vmNics)
	if diags.HasError() {
		return current, fmt.Errorf("failed encoding vm_nics state: %v", diags)
	}
	updated.VMNics = vmNicsValue

	networkInterfaces := normalizeNetworkInterfacesFromVM(vm)
	networkInterfacesValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"l3_network_uuid": types.StringType,
			"default_l3":      types.BoolType,
			"static_ip":       types.StringType,
		},
	}, networkInterfaces)
	if diags.HasError() {
		return current, fmt.Errorf("failed encoding network_interfaces state: %v", diags)
	}
	updated.NetworkInterfaces = networkInterfacesValue

	updatedDataDisks, err := syncInstanceDataDisksFromVM(ctx, updated, vm)
	if err != nil {
		return current, err
	}

	return updatedDataDisks, nil
}

func instanceDataVolumeUUIDsFromState(ctx context.Context, state vmInstanceDataSourceModel) ([]string, error) {
	if state.DataDisks.IsNull() || state.DataDisks.IsUnknown() || len(state.DataDisks.Elements()) == 0 {
		return nil, nil
	}

	var dataDisks []diskModel
	if diags := state.DataDisks.ElementsAs(ctx, &dataDisks, false); diags.HasError() {
		return nil, fmt.Errorf("failed decoding data_disks state: %v", diags)
	}

	volumeUUIDs := make([]string, 0, len(dataDisks))
	for _, disk := range dataDisks {
		if disk.VolumeUuid.IsNull() || disk.VolumeUuid.IsUnknown() || disk.VolumeUuid.ValueString() == "" {
			continue
		}
		volumeUUIDs = append(volumeUUIDs, disk.VolumeUuid.ValueString())
	}
	return volumeUUIDs, nil
}

func isDiskParamValid(r *instanceResource, model diskModel) error {
	if model.PrimaryStorageUuid.IsNull() || model.PrimaryStorageUuid.ValueString() == "" {
		return nil
	}

	dataDiskPrimaryStorageUuid := model.PrimaryStorageUuid.ValueString()

	qparam := param.NewQueryParam()
	qparam.AddQ("uuid=" + dataDiskPrimaryStorageUuid)
	qparam.AddQ("state=Enabled")
	qparam.Limit(1)
	primaryStorages, err := r.client.QueryPrimaryStorage(&qparam)
	if err != nil {
		return fmt.Errorf("failed to get primary storage %s, err: %v", dataDiskPrimaryStorageUuid, err)
	}

	if len(primaryStorages) == 0 {
		return fmt.Errorf("unable to find primary storage %s, err: %v", dataDiskPrimaryStorageUuid, err)
	}

	return nil
}

func (r *instanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

