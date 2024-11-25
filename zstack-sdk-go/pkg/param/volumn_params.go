// Copyright (c) ZStack.io, Inc.

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
	Name               string   `json:"name" example:"chenjh-DATA-TEST"`                     // Data volume name
	Description        string   `json:"description" example:"JUST a test Volume For chenjh"` // Description of the data volume
	DiskOfferingUuid   string   `json:"diskOfferingUuid" example:""`                         // Disk offering UUID
	DiskSize           int64    `json:"diskSize" example:"1024"`                             // Disk size
	PrimaryStorageUuid string   `json:"primaryStorageUuid" example:""`                       // Primary storage UUID
	ResourceUuid       string   `json:"resourceUuid" example:""`                             // Resource UUID
	TagUuids           []string `json:"tagUuids" example:""`                                 // List of tag UUIDs
}

type CreateDataVolumeFromVolumeTemplateParam struct {
	BaseParam

	Params CreateDataVolumeFromVolumeTemplateDetailParam `json:"params"`
}

type CreateDataVolumeFromVolumeTemplateDetailParam struct {
	Name               string `json:"name"`               // Data volume name
	Description        string `json:"description"`        // Detailed description of the data volume
	PrimaryStorageUuid string `json:"primaryStorageUuid"` // Primary storage UUID
	HostUuid           string `json:"hostUuid"`           // Host UUID
	ResourceUuid       string `json:"resourceUuid"`
}

type CreateDataVolumeFromVolumeSnapshotParam struct {
	BaseParam

	Params CreateDataVolumeFromVolumeSnapshotDetailParam `json:"params"`
}

type CreateDataVolumeFromVolumeSnapshotDetailParam struct {
	Name               string `json:"name"`               // Data volume name
	Description        string `json:"description"`        // Detailed description of the data volume
	VolumeSnapshotUuid string `json:"volumeSnapshotUuid"` // Volume snapshot UUID
	PrimaryStorageUuid string `json:"primaryStorageUuid"` // Primary storage UUID
	ResourceUuid       string `json:"resourceUuid"`       // Resource Uuid
}

type UpdateVolumeParam struct {
	BaseParam

	UpdateVolume UpdateVolumeDetailParam `json:"updateVolume"`
}

type UpdateVolumeDetailParam struct {
	Name        string  `json:"name"`        // Resource name
	Description *string `json:"description"` // Detailed description of the resource
}

type SetVolumeQoSParam struct {
	BaseParam

	SetVolumeQoS SetVolumeQoSDetailParam `json:"setVolumeQos"`
}

type SetVolumeQoSDetailParam struct {
	VolumeBandwidth int64  `json:"volumeBandwidth"` // Data volume speed limit bandwidth
	Mode            string `json:"mode"`            // total read write
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
	StateEvent VolumeState `json:"stateEvent"` // Enable or disable, valid values: enable, disable
}
