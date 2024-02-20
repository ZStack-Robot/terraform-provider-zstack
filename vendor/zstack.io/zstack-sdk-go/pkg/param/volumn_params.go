package param

type VolumeState string

const (
	VolumeStateEnable  VolumeState = "enable"
	VolumeStateDisable VolumeState = "disable"
)

type CreateDataVolumeParam struct {
	BaseParam

	Params CreateDataVolumeDetailParam `json:"params"`
}

type CreateDataVolumeDetailParam struct {
	Name               string   `json:"name" example:"chenjh-DATA-TEST"`                     //云盘名称
	Description        string   `json:"description" example:"JUST a test Volume For chenjh"` //云盘的描述
	DiskOfferingUuid   string   `json:"diskOfferingUuid" example:""`                         //云盘规格UUID
	DiskSize           int64    `json:"diskSize" example:"1024"`                             //云盘大小
	PrimaryStorageUuid string   `json:"primaryStorageUuid" example:""`                       //主存储UUID
	ResourceUuid       string   `json:"resourceUuid" example:""`                             //资源UUID
	TagUuids           []string `json:"tagUuids" example:""`                                 //标签UUID列表
}

type CreateDataVolumeFromVolumeTemplateParam struct {
	BaseParam

	Params CreateDataVolumeFromVolumeTemplateDetailParam `json:"params"`
}

type CreateDataVolumeFromVolumeTemplateDetailParam struct {
	Name               string `json:"name"`               //云盘名称
	Description        string `json:"description"`        //云盘的详细描述
	PrimaryStorageUuid string `json:"primaryStorageUuid"` //主存储UUID
	HostUuid           string `json:"hostUuid"`           //物理机UUID
	ResourceUuid       string `json:"resourceUuid"`
}

type CreateDataVolumeFromVolumeSnapshotParam struct {
	BaseParam

	Params CreateDataVolumeFromVolumeSnapshotDetailParam `json:"params"`
}

type CreateDataVolumeFromVolumeSnapshotDetailParam struct {
	Name               string `json:"name"`               //云盘名称
	Description        string `json:"description"`        //云盘的详细描述
	VolumeSnapshotUuid string `json:"volumeSnapshotUuid"` //云盘快照UUID
	PrimaryStorageUuid string `json:"primaryStorageUuid"` //主存储UUID
	ResourceUuid       string `json:"resourceUuid"`       //资源的Uuid
}

type UpdateVolumeParam struct {
	BaseParam

	UpdateVolume UpdateVolumeDetailParam `json:"updateVolume"`
}

type UpdateVolumeDetailParam struct {
	Name        string  `json:"name"`        //资源名称
	Description *string `json:"description"` //资源的详细描述
}

type SetVolumeQoSParam struct {
	BaseParam

	SetVolumeQoS SetVolumeQoSDetailParam `json:"setVolumeQos"`
}

type SetVolumeQoSDetailParam struct {
	VolumeBandwidth int64  `json:"volumeBandwidth"` //云盘限速带宽
	Mode            string `json:"mode"`            //total read write
	ReadBandwidth   int64  `json:"readBandwidth"`
	WriteBandwidth  int64  `json:"writeBandwidth"`
	TotalBandwidth  int64  `json:"totalBandwidth"`
	ReadIOPS        int64  `json:"readIOPS"`
	WriteIOPS       int64  `json:"writeIOPS"`
	TotalIOPS       int64  `json:"totalIOPS"`
}

type ChangeVolumeStateParam struct {
	BaseParam

	ChangeVolumeState ChangeVolumeStateDetailParam `json:"changeVolumeState"`
}

type ChangeVolumeStateDetailParam struct {
	StateEvent VolumeState `json:"stateEvent"` //开启或关闭，取值范围：enable、disable
}
