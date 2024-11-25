// Copyright (c) ZStack.io, Inc.

package param

const (
	ResourceTypeVmInstanceVO          = "VmInstanceVO"
	ResourceTypeImageVO               = "ImageVO"
	ResourceTypeVolumeVo              = "VolumeVO"
	ResourceTypeVolumeSnapshotVO      = "VolumeSnapshotVO"
	ResourceTypeVolumeSnapshotGroupVO = "VolumeSnapshotGroupVO"
	ResourceTypeL3NetworkVO           = "L3NetworkVO"
)

type CreateTagParam struct {
	BaseParam

	Params CreateTagDetailParam `json:"params"`
}

type CreateTagDetailParam struct {
	ResourceType string `json:"resourceType"`
	ResourceUuid string `json:"resourceUuid"`
	Tag          string `json:"tag"`
}

type UpdateSystemTagParam struct {
	BaseParam

	UpdateSystemTag UpdateTagDetailParam `json:"updateSystemTag"`
}

type UpdateTagDetailParam struct {
	Tag string `json:"tag"`
}
