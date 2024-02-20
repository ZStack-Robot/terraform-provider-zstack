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
	Name        string  `json:"name"`        //镜像名称
	Description *string `json:"description"` //镜像的详细描述
}

type UpdateImageVirtioParam struct {
	BaseParam
	UpdateImage UpdateImageVirtioDetailParam `json:"updateImage"`
}
type UpdateImageVirtioDetailParam struct {
	Virtio bool   `json:"virtio"`
	UUID   string `json:"uuid"` //资源的UUID，唯一标示该资源
}

type UpdateImagePlatformParam struct {
	BaseParam
	UpdateImage UpdateImagePlatformDetailParam `json:"updateImage"`
}

type UpdateImagePlatformDetailParam struct {
	Platform    string `json:"platform"`    //平台
	GuestOsType string `json:"guestOsType"` //镜像对应的客户机操作系统类型
	UUID        string `json:"uuid"`        //资源的UUID，唯一标示该资源
}

type ExpungeImageParam struct {
	BaseParam
	BackupStorageUuids []string `json:"backupStorageUuids"`
}

type RecoverImageParam struct {
	BaseParam
	ImageUuid    string                   `json:"imageUuid"`    //镜像UUID
	RecoverImage RecoverImageDetailParams `json:"recoverImage"` //放backupStorageUuids
}

type RecoverImageDetailParams struct {
	BackupStorageUuids []string `json:"backupStorageUuids"` //指定添加镜像的镜像服务器UUID列表
}

type ChangeImageStateParam struct {
	BaseParam
	ImageUuid        string                      `json:"imageUuid"` //镜像UUID
	ChangeImageState ChangeImageStateDetailParam `json:"changeImageState"`
}

type ChangeImageStateDetailParam struct {
	StateEvent StateEvent `json:"stateEvent"`
}

type SyncImageSizeParam struct {
	BaseParam
	ImageUuid     string                   `json:"imageUuid"` //镜像UUID
	SyncImageSize SyncImageSizeDetailParam `json:"syncImageSize"`
}
type SyncImageSizeDetailParam struct {
}

type GetCandidateBackupStorageForCreatingImageParam struct {
	BaseParam
	CandidateBackupStorageType CandidateBackupStorageType `json:"candidateBackupStorageType"`
	VolumeUuid                 string                     `json:"volumeUuid" `         //云盘UUID，注意：volumeUuid 和 volumeSnapshotUuid 二选一
	VolumeSnapshotUuid         string                     `json:"volumeSnapshotUuid" ` //云盘快照UUID，注意：volumeUuid 和 volumeSnapshotUuid 二选一
}

type CreateRootVolumeTemplateFromRootVolumeParam struct {
	BaseParam
	RootVolumeUuid string                                            `json:"rootVolumeUuid"` //根云盘UUID
	Params         CreateRootVolumeTemplateFromRootVolumeDetailParam `json:"params"`         //结构体中的其他参数
}

type CreateRootVolumeTemplateFromRootVolumeDetailParam struct {
	Name               string   `json:"name"`                //名称
	RootVolumeUuid     string   `json:"rootVolumeUuid"`      //根云盘UUID
	Description        string   `json:"description"`         //详细描述
	GuestOsType        string   `json:"guestOsType" `        //根云盘镜像对应客户机操作系统类型
	BackupStorageUuids []string `json:"backupStorageUuids" ` //镜像服务器UUID列表
	Platform           string   `json:"platform" `           //镜像的系统平台,Linux,Windows,WindowsVirtio,Other,Paravirtualization
	System             bool     `json:"system"`              //是否系统根云盘镜像
	ResourceUuid       string   `json:"resourceUuid" `       //根云盘镜像UUID。若指定，根云盘镜像会使用该字段值作为UUID。
	Architecture       string   `json:"architecture"`        //x86_64,aarch64,mips64el
	TagUuids           []string `json:"tagUuids"`            //标签UUID列表
}

type CreateRootVolumeTemplateFromVolumeSnapshotParam struct {
	BaseParam
	SnapshotUuid string                                                 `json:"snapshotUuid" ` //快照UUID
	Params       CreateRootVolumeTemplateFromVolumeSnapshotDetailParams `json:"params"`        //结构体中的其他参数
}
type CreateRootVolumeTemplateFromVolumeSnapshotDetailParams struct {
	Name               string   `json:"name" `               //名称
	Description        string   `json:"description" `        //详细描述
	GuestOsType        string   `json:"guestOsType"`         //根云盘镜像对应客户机操作系统类型
	BackupStorageUuids []string `json:"backupStorageUuids" ` //镜像服务器UUID列表
	Platform           string   `json:"platform" `           //镜像的系统平台,Linux,Windows,WindowsVirtio,Other,Paravirtualization
	System             bool     `json:"system" `             //是否系统根云盘镜像
	ResourceUuid       string   `json:"resourceUuid"`        //根云盘镜像UUID。若指定，根云盘镜像会使用该字段值作为UUID。
	Architecture       string   `json:"architecture"`        //x86_64,aarch64,mips64el
	TagUuids           []string `json:"tagUuids"`            //标签UUID列表
}

type CreateDataVolumeTemplateFromVolumeParam struct {
	BaseParam
	VolumeUuid string                                        `json:"volumeUuid" ` //快照UUID
	Params     CreateDataVolumeTemplateFromVolumeDetailParam `json:"params"`      //结构体中的其他参数

}

type CreateDataVolumeTemplateFromVolumeDetailParam struct {
	Name               string   `json:"name" `               //名称
	Description        string   `json:"description" `        //详细描述
	BackupStorageUuids []string `json:"backupStorageUuids" ` //镜像服务器UUID列表
	ResourceUuid       string   `json:"resourceUuid" `       //根云盘镜像UUID。若指定，根云盘镜像会使用该字段值作为UUID。
}

type CreateDataVolumeTemplateFromVolumeSnapshotParam struct {
	BaseParam
	SnapshotUuid string                                                `json:"snapshotUuid" ` //快照UUID
	Params       CreateDataVolumeTemplateFromVolumeSnapshotDetailParam `json:"params" `       //结构体中的其他参数
}

type CreateDataVolumeTemplateFromVolumeSnapshotDetailParam struct {
	Name               string   `json:"name" `               //名称
	Description        string   `json:"description" `        //详细描述
	BackupStorageUuids []string `json:"backupStorageUuids" ` //镜像服务器UUID列表
	ResourceUuid       string   `json:"resourceUuid" `       //根云盘镜像UUID。若指定，根云盘镜像会使用该字段值作为UUID。
	TagUuids           []string `json:"tagUuids" `           //标签UUID列表
}

type SetImageQgaParam struct {
	BaseParam
	Uuid        string                 `json:"uuid"`
	SetImageQga SetImageQgaDetailParam `json:"setImageQga"` //放enable
}

type SetImageQgaDetailParam struct {
	Enable bool `json:"enable" `
}

type SetImageBootModeRequest struct {
	BaseParam
	Uuid             string                 `json:"uuid" `
	SetImageBootMode SetImageBootModeParams `json:"setImageBootMode"` //放bootMode
}

type SetImageBootModeParams struct {
	BootMode BootMode `json:"bootMode"` //镜像启动模式,Legacy,UEFI,UEFI_WITH_CSM
}

type GetUploadImageJobDetailsParam struct {
	BaseParam
	ImageId string `json:"imageId" `
}

type UpdateImageArchitectureDetailParam struct {
	UUID         string       `json:"uuid"`         //资源的UUID，唯一标示该资源
	Architecture Architecture `json:"architecture"` //x86_64,aarch64,mips64el
}

type UpdateImageArchitectureParam struct {
	BaseParam
	UpdateImage UpdateImageArchitectureDetailParam `json:"updateImage"`
}
