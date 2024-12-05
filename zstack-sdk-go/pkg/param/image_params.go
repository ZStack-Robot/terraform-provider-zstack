// Copyright (c) ZStack.io, Inc.

package param

type Architecture string
type MediaType string
type ImageFormat string
type StateEvent string
type CandidateBackupStorageType string
type BootMode string

const (
	X86_64   Architecture = "x86_64"
	Aarch64  Architecture = "aarch64"
	Mips64el Architecture = "mips64el"

	RootVolumeTemplate MediaType = "RootVolumeTemplate"
	ISO                MediaType = "ISO"
	DataVolumeTemplate MediaType = "DataVolumeTemplate"

	Raw   ImageFormat = "raw"
	Qcow2 ImageFormat = "qcow2"
	Iso   ImageFormat = "iso"
	VMDK  ImageFormat = "vmdk"
	VHD   ImageFormat = "vhd"

	StateEventEnable  StateEvent = "enable"
	StateEventDisable StateEvent = "disable"

	CandidateBackupStorageTypeDefault         CandidateBackupStorageType = ""
	CandidateBackupStorageTypeVolumes         CandidateBackupStorageType = "volumes"
	CandidateBackupStorageTypeVolumeSnapshots CandidateBackupStorageType = "volume-snapshots"

	Legacy      BootMode = "Legacy"
	UEFI        BootMode = "UEFI"
	UEFIWITHCSM BootMode = "UEFI_WITH_CSM"

	SystemTagBootModeUEFI   = "bootMode::UEFI"
	SystemTagBootModeLegacy = "bootMode::Legacy"

	SystemTagApplianceTypeVRouter = "applianceType::vrouter"
)

type AddImageParam struct {
	BaseParam
	Params AddImageDetailParam `json:"params"`
}

type AddImageDetailParam struct {
	Name               string       `json:"name" example:"vm-image-1"`                                                                                        //镜像名称
	Description        string       `json:"description" example:"vm-image-1 desc"`                                                                            //详细描述
	Url                string       `json:"url" example:"http://172.20.20.132:8001/imagestore/download/image-d1f501b3887a6a084feb66d0a995215731f664e4.qcow2"` //被添加镜像的URL地址
	MediaType          MediaType    `json:"mediaType" example:"RootVolumeTemplate"`                                                                           //镜像的类型,RootVolumeTemplate,ISO,DataVolumeTemplate
	GuestOsType        string       `json:"guestOsType" example:"Windows 10"`                                                                                 //镜像对应客户机操作系统的类型
	System             bool         `json:"system" example:"false"`                                                                                           //                                                                                                           //是否系统镜像（如，云路由镜像）
	Format             ImageFormat  `json:"format" example:"raw"`                                                                                             //镜像的格式，比如：raw
	Platform           string       `json:"platform" example:"Windows"`                                                                                       //                                                                                                         //镜像的系统平台,Linux,Windows,WindowsVirtio,Other,Paravirtualization
	BackupStorageUuids []string     `json:"backupStorageUuids" example:"26684790e4734a0bbb506f40907f57da"`                                                    //指定添加镜像的镜像服务器UUID列表
	Type               string       `json:"type"`                                                                                                             //内部使用字段
	ResourceUuid       string       `json:"resourceUuid"`                                                                                                     //资源UUID。若指定，镜像会使用该字段值作为UUID。
	Architecture       Architecture `json:"architecture" example:"x86_64"`                                                                                    //x86_64,aarch64,mips64el
	TagUuids           []string     `json:"tagUuids"`                                                                                                         //标签UUID列表
	Virtio             bool         `json:"virtio"`
	SystemTags         []string     `json:"systemTags"`
	UserTags           []string     `json:"userTags"`
}

type UpdateImageParam struct {
	BaseParam
	UpdateImage UpdateImageDetailParam `json:"updateImage"`
}

type UpdateImageDetailParam struct {
	Name        string  `json:"name"`        // Image name
	Description *string `json:"description"` // Detailed description of the image
}

type UpdateImageVirtioParam struct {
	BaseParam
	UpdateImage UpdateImageVirtioDetailParam `json:"updateImage"`
}

type UpdateImageVirtioDetailParam struct {
	Virtio bool   `json:"virtio"`
	UUID   string `json:"uuid"` // Resource UUID, uniquely identifies the resource
}

type UpdateImagePlatformParam struct {
	BaseParam
	UpdateImage UpdateImagePlatformDetailParam `json:"updateImage"`
}

type UpdateImagePlatformDetailParam struct {
	Platform    string `json:"platform"`    // Platform
	GuestOsType string `json:"guestOsType"` // Guest OS type corresponding to the image
	UUID        string `json:"uuid"`        // Resource UUID, uniquely identifies the resource
}

type ExpungeImageParam struct {
	BaseParam
	BackupStorageUuids []string `json:"backupStorageUuids"`
}

type RecoverImageParam struct {
	BaseParam
	ImageUuid    string                   `json:"imageUuid"`    // Image UUID
	RecoverImage RecoverImageDetailParams `json:"recoverImage"` // Backup storage UUIDs
}

type RecoverImageDetailParams struct {
	BackupStorageUuids []string `json:"backupStorageUuids"` // List of backup storage UUIDs to add the image to
}

type ChangeImageStateParam struct {
	BaseParam
	ImageUuid        string                      `json:"imageUuid"` // Image UUID
	ChangeImageState ChangeImageStateDetailParam `json:"changeImageState"`
}

type ChangeImageStateDetailParam struct {
	StateEvent StateEvent `json:"stateEvent"`
}

type SyncImageSizeParam struct {
	BaseParam
	ImageUuid     string                   `json:"imageUuid"` // Image UUID
	SyncImageSize SyncImageSizeDetailParam `json:"syncImageSize"`
}

type SyncImageSizeDetailParam struct {
}

type GetCandidateBackupStorageForCreatingImageParam struct {
	BaseParam
	CandidateBackupStorageType CandidateBackupStorageType `json:"candidateBackupStorageType"`
	VolumeUuid                 string                     `json:"volumeUuid"`         // Cloud volume UUID, note: either volumeUuid or volumeSnapshotUuid is required
	VolumeSnapshotUuid         string                     `json:"volumeSnapshotUuid"` // Cloud volume snapshot UUID, note: either volumeUuid or volumeSnapshotUuid is required
}

type CreateRootVolumeTemplateFromRootVolumeParam struct {
	BaseParam
	RootVolumeUuid string                                            `json:"rootVolumeUuid"` // Root cloud volume UUID
	Params         CreateRootVolumeTemplateFromRootVolumeDetailParam `json:"params"`         // Other parameters in the struct
}

type CreateRootVolumeTemplateFromRootVolumeDetailParam struct {
	Name               string   `json:"name"`               // Name
	RootVolumeUuid     string   `json:"rootVolumeUuid"`     // Root cloud volume UUID
	Description        string   `json:"description"`        // Detailed description
	GuestOsType        string   `json:"guestOsType"`        // Guest OS type corresponding to the root cloud volume image
	BackupStorageUuids []string `json:"backupStorageUuids"` // List of backup storage UUIDs
	Platform           string   `json:"platform"`           // Image system platform, Linux, Windows, WindowsVirtio, Other, Paravirtualization
	System             bool     `json:"system"`             // Whether it is a system root cloud volume image
	ResourceUuid       string   `json:"resourceUuid"`       // Root cloud volume image UUID. If specified, the root cloud volume image will use this value as its UUID.
	Architecture       string   `json:"architecture"`       // x86_64, aarch64, mips64el
	TagUuids           []string `json:"tagUuids"`           // List of tag UUIDs
}

type CreateRootVolumeTemplateFromVolumeSnapshotParam struct {
	BaseParam
	SnapshotUuid string                                                 `json:"snapshotUuid"` // Snapshot UUID
	Params       CreateRootVolumeTemplateFromVolumeSnapshotDetailParams `json:"params"`       // Other parameters in the struct
}

type CreateRootVolumeTemplateFromVolumeSnapshotDetailParams struct {
	Name               string   `json:"name"`               // Name
	Description        string   `json:"description"`        // Detailed description
	GuestOsType        string   `json:"guestOsType"`        // Guest OS type corresponding to the root cloud volume image
	BackupStorageUuids []string `json:"backupStorageUuids"` // List of backup storage UUIDs
	Platform           string   `json:"platform"`           // Image system platform, Linux, Windows, WindowsVirtio, Other, Paravirtualization
	System             bool     `json:"system"`             // Whether it is a system root cloud volume image
	ResourceUuid       string   `json:"resourceUuid"`       // Root cloud volume image UUID. If specified, the root cloud volume image will use this value as its UUID.
	Architecture       string   `json:"architecture"`       // x86_64, aarch64, mips64el
	TagUuids           []string `json:"tagUuids"`           // List of tag UUIDs
}

type CreateDataVolumeTemplateFromVolumeParam struct {
	BaseParam
	VolumeUuid string                                        `json:"volumeUuid"` // Snapshot UUID
	Params     CreateDataVolumeTemplateFromVolumeDetailParam `json:"params"`     // Other parameters in the struct
}

type CreateDataVolumeTemplateFromVolumeDetailParam struct {
	Name               string   `json:"name"`               // Name
	Description        string   `json:"description"`        // Detailed description
	BackupStorageUuids []string `json:"backupStorageUuids"` // List of backup storage UUIDs
	ResourceUuid       string   `json:"resourceUuid"`       // Root cloud volume image UUID. If specified, the root cloud volume image will use this value as its UUID.
}

type CreateDataVolumeTemplateFromVolumeSnapshotParam struct {
	BaseParam
	SnapshotUuid string                                                `json:"snapshotUuid"` // Snapshot UUID
	Params       CreateDataVolumeTemplateFromVolumeSnapshotDetailParam `json:"params"`       // Other parameters in the struct
}

type CreateDataVolumeTemplateFromVolumeSnapshotDetailParam struct {
	Name               string   `json:"name"`               // Name
	Description        string   `json:"description"`        // Detailed description
	BackupStorageUuids []string `json:"backupStorageUuids"` // List of backup storage UUIDs
	ResourceUuid       string   `json:"resourceUuid"`       // Root cloud volume image UUID. If specified, the root cloud volume image will use this value as its UUID.
	TagUuids           []string `json:"tagUuids"`           // List of tag UUIDs
}

type SetImageQgaParam struct {
	BaseParam
	Uuid        string                 `json:"uuid"`
	SetImageQga SetImageQgaDetailParam `json:"setImageQga"` // Enable
}

type SetImageQgaDetailParam struct {
	Enable bool `json:"enable"`
}

type SetImageBootModeRequest struct {
	BaseParam
	Uuid             string                 `json:"uuid"`
	SetImageBootMode SetImageBootModeParams `json:"setImageBootMode"` // Boot mode
}

type SetImageBootModeParams struct {
	BootMode BootMode `json:"bootMode"` // Image boot mode, Legacy, UEFI, UEFI_WITH_CSM
}

type GetUploadImageJobDetailsParam struct {
	BaseParam
	ImageId string `json:"imageId"`
}

type UpdateImageArchitectureDetailParam struct {
	UUID         string       `json:"uuid"`         // Resource UUID, uniquely identifies the resource
	Architecture Architecture `json:"architecture"` // x86_64, aarch64, mips64el
}

type UpdateImageArchitectureParam struct {
	BaseParam
	UpdateImage UpdateImageArchitectureDetailParam `json:"updateImage"`
}
