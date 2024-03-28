// Copyright (c) HashiCorp, Inc.

package param

type InstanceType string
type InstanceStrategy string
type InstanceStopType string
type HA string
type ClockTrack string

const (
	UserVm      InstanceType = "UserVm"
	ApplianceVm InstanceType = "ApplianceVm"

	InstantStart  InstanceStrategy = "InstantStart"
	CreateStopped InstanceStrategy = "CreateStopped"

	Grace InstanceStopType = "grace" //优雅关机，需要云主机里安装了相关ACPI驱动
	Cold  InstanceStopType = "cold"  //冷关机，相当于直接断电

	NeverStop HA = "NeverStop" //开启高可用
	None      HA = "None"      //未开高可用

	Host  ClockTrack = "host"
	Guest ClockTrack = "guest"
)

type CreateVmInstanceParam struct {
	BaseParam
	Params CreateVmInstanceDetailParam `json:"params" `
}

type CreateVmInstanceDetailParam struct {
	Name                            string       `json:"name" `                            //云主机名称
	InstanceOfferingUUID            string       `json:"instanceOfferingUuid" `            //计算规格UUID 指定云主机的CPU、内存等参数。
	CpuNum                          int64        `json:"cpuNum"`                           //cpu 1
	MemorySize                      int64        `json:"memorySize"`                       //内存大小 1073741824
	ImageUUID                       string       `json:"imageUuid" `                       //镜像UUID 云主机的根云盘会从该字段指定的镜像创建。
	L3NetworkUuids                  []string     `json:"l3NetworkUuids" `                  //三层网络UUID列表 可指定一个或多个三层网络，云主机会在每个三层网络上创建一个网卡。
	Type                            InstanceType `json:"type" `                            //云主机类型 保留字段，无需指定。UserVm/ApplianceVm
	RootDiskOfferingUuid            string       `json:"rootDiskOfferingUuid" `            //根云盘规格UUID 如果imageUuid字段指定的镜像类型是ISO，该字段必须指定以确定需要创建的根云盘大小。如果镜像类型是非ISO，该字段无需指定。
	RootDiskSize                    *int64       `json:"rootDiskSize"`                     //跟云盘大小
	DataDiskOfferingUuids           []string     `json:"dataDiskOfferingUuids" `           //云盘规格UUID列表 可指定一个或多个云盘规格UUID（UUID可以重复）为云主机创建一个或多个数据云盘。
	DataDiskSizes                   []int64      `json:"dataDiskSizes"`                    //数据云盘大小
	ZoneUuid                        string       `json:"zoneUuid" `                        //区域UUID 若指定，云主机会在指定区域创建。
	ClusterUUID                     string       `json:"clusterUuid" `                     //集群UUID 若指定，云主机会在指定集群创建，该字段优先级高于zoneUuid。
	HostUuid                        string       `json:"hostUuid" `                        //物理机UUID 若指定，云主机会在指定物理机创建，该字段优先级高于zoneUuid和clusterUuid。
	PrimaryStorageUuidForRootVolume *string      `json:"primaryStorageUuidForRootVolume" ` //主存储UUID 若指定，云主机的根云盘会在指定主存储创建。
	Description                     string       `json:"description" `                     //云主机的详细描述
	DefaultL3NetworkUuid            string       `json:"defaultL3NetworkUuid" `            //默认三层网络UUID 当l3NetworkUuids指定了多个三层网络时，该字段指定提供默认路由的三层网络。若不指定，l3NetworkUuids的第一个网络被选为默认网络。
	ResourceUuid                    string       `json:"resourceUuid" `                    //资源UUID 若指定，云主机会使用该字段值作为UUID。

	TagUuids             []string         `json:"tagUuids" ` //标签UUID列表
	Strategy             InstanceStrategy `json:"strategy" ` //云主机创建策略 创建后立刻启动InstantStart 创建后不启动CreateStopped
	RootVolumeSystemTags []string         `json:"rootVolumeSystemTags"`
	DataVolumeSystemTags []string         `json:"dataVolumeSystemTags"`
}

type CreateVmFromVolumeParam struct {
	BaseParam
	Params CreateVmFromVolumeDetailParams `json:"params"`
}

type CreateVmFromVolumeDetailParams struct {
	Name                 string   `json:"name"`                 //云主机名称
	Description          string   `json:"description"`          //资源的详细描述
	InstanceOfferingUuid string   `json:"instanceOfferingUuid"` //计算规格UUID，注意：该参数与CPU数量、内存大小二选一
	CpuNum               int      `json:"cpuNum"`               //CPU数量/内存大小，注意：该参数与instanceOfferingUuid二选一
	MemorySize           int64    `json:"memorySize"`           //CPU数量/内存大小，注意：该参数与instanceOfferingUuid二选一
	L3NetworkUuids       []string `json:"l3NetworkUuids"`       //三层网络UUID列表 可指定一个或多个三层网络，云主机会在每个三层网络上创建一个网卡。
	Type                 string   `json:"type"`                 //云主机类型保留字段，无需指定。
	VolumeUuid           string   `json:"volumeUuid"`           //云盘UUID
	Platform             string   `json:"platform"`             //云盘系统平台
	ZoneUuid             string   `json:"zoneUuid"`             //区域UUID 若指定，云主机会在指定区域创建。
	ClusterUuid          string   `json:"clusterUuid"`          //集群UUID 若指定，云主机会在指定集群创建，该字段优先级高于zoneUuid
	HostUuid             string   `json:"hostUuid"`             //物理机UUID 若指定，云主机会在指定物理机创建，该字段优先级高于zoneUuid和clusterUuid
	PrimaryStorageUuid   string   `json:"primaryStorageUuid"`   //主存储UUID 若指定，云主机的根云盘会在指定主存储创建。
	DefaultL3NetworkUuid string   `json:"defaultL3NetworkUuid"` //默认三层网络UUID 当l3NetworkUuids指定了多个三层网络时，该字段指定提供默认路由的三层网络。若不指定，l3NetworkUuids的第一个网络被选为默认网络。
	Strategy             string   `json:"strategy"`             //云主机创建策略 1.创建后立刻启动 2.创建后不启动
	ResourceUuid         string   `json:"resourceUuid"`         //资源UUID 若指定，云主机会使用该字段值作为UUID。
	TagUuids             []string `json:"tagUuids"`             //标签UUID列表
}

type CloneVmInstanceParam struct {
	BaseParam
	CloneVmInstance CloneVmInstanceDetailParam `json:"cloneVmInstance"`
}

type CloneVmInstanceDetailParam struct {
	Names                           []string         `json:"names"`    //云主机名称
	Strategy                        InstanceStrategy `json:"strategy"` //策略 克隆后立刻启动InstantStart 克隆后不启动JustCreate
	Full                            *bool            `json:"full"`     //是否克隆已挂载数据盘
	PrimaryStorageUuidForRootVolume *string          `json:"primaryStorageUuidForRootVolume" `
	PrimaryStorageUuidForDataVolume *string          `json:"primaryStorageUuidForDataVolume" `
	RootVolumeSystemTags            []string         `json:"rootVolumeSystemTags" `
	DataVolumeSystemTags            []string         `json:"dataVolumeSystemTags" `
}

type StartVmInstanceParam struct {
	BaseParam
	StartVmInstance StartVmInstanceDetailParam `json:"startVmInstance"` //可传hostUuid
}

type StartVmInstanceDetailParam struct {
	HostUuid string `json:"hostUuid"` //物理机UUID
}

type StopVmInstanceParam struct {
	BaseParam
	StopVmInstance StopVmInstanceDetailParam `json:"stopVmInstance"` //需要存uuid和 type
}

type StopVmInstanceDetailParam struct {
	Type   InstanceStopType `json:"type"`   //默认为grace：优雅关机；cold：冷关机（关闭电源）
	StopHA bool             `json:"stopHa"` //彻底关闭HA云主机
}

type UpdateVmInstanceParam struct {
	BaseParam
	UpdateVmInstance UpdateVmInstanceDetailParam `json:"updateVmInstance"`
}

type UpdateVmInstanceDetailParam struct {
	Name                 string  `json:"name"`        //云主机名称
	Description          *string `json:"description"` //云主机的详细描述
	State                string  `json:"state"`
	DefaultL3NetworkUuid string  `json:"defaultL3NetworkUuid"` //默认三层网络UUID 当l3NetworkUuids指定了多个三层网络时，该字段指定提供默认路由的三层网络。若不指定，l3NetworkUuids的第一个网络被选为默认网络。
	Platform             string  `json:"platform"`             //云盘系统平台
	CpuNum               *int    `json:"cpuNum"`               //CPU数目
	MemorySize           *int64  `json:"memorySize"`           //CPU数量/内存大小，注意：该参数与instanceOfferingUuid二选一
	GuestOsType          string  `json:"guestOsType"`
}

type UpdateVmInstanceQgaParam struct {
	BaseParam
	SetVmQga SetVmQgaParam `json:"setVmQga"`
	UUID     string        `json:"UUID"`
}

type SetVmQgaParam struct {
	Enable bool `json:"enable"`
}

type SetVmBootModeParam struct {
	BaseParam
	SetVmBootMode SetVmBootModeDetailParam `json:"setVmBootMode"`
}

type SetVmBootModeDetailParam struct {
	BootMode BootMode `json:"bootMode"` //启动模式 Legacy,UEFI,UEFI_WITH_CSM
}

type UpdateVmInstanceSshKeyParam struct {
	UUID        string         `json:"uuid"`
	SetVmSshKey SetSshKeyParam `json:"setVmSshKey"`
}
type SetSshKeyParam struct {
	SshKey string `json:"SshKey"`
}

type UpdateVmInstanceChangePwdParam struct {
	UUID             string                `json:"uuid"`
	ChangeVmPassword ChangeVmPasswordParam `json:"changeVmPassword"`
}
type ChangeVmPasswordParam struct {
	Password string `json:"password"`
	Account  string `json:"account"`
}

type UpdateVmInstanceClockTrackParam struct {
	BaseParam
	SetVmClockTrack UpdateVmInstanceClockTrackDetailParam `json:"setVmClockTrack"`
}

type UpdateVmInstanceClockTrackDetailParam struct {
	Track             ClockTrack `json:"track"`             //时钟同步方式，可选值：guest, host
	SyncAfterVMResume bool       `json:"syncAfterVMResume"` //是否在云主机恢复时同步时钟
	IntervalInSeconds float64    `json:"intervalInSeconds"` //时钟同步间隔，单位：秒0  60 600 1800 3600 7200 21600 43200 86400
}
type UpdateVmCdRomParam struct {
	BaseParam
	UpdateVmCdRom UpdateVmCdRomDetailParam `json:"updateVmCdRom"`
}
type UpdateVmCdRomDetailParam struct {
	Name string `json:"name"`
}

type CreateVmCdRomParam struct {
	BaseParam
	Params CreateVmCdRomDetailParam `json:"params"`
}

type CreateVmCdRomDetailParam struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	VmInstanceUuid string `json:"vmInstanceUuid"`
	IsoUuid        string `json:"isoUuid"`
	ResourceUuid   string `json:"resourceUuid"`
}
