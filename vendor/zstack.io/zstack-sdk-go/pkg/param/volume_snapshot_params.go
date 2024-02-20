package param

type VolumeSnapshotParam struct {
	BaseParam

	Params VolumeSnapshotDetailParam `json:"params"`
}

type VolumeSnapshotDetailParam struct {
	Name         string `json:"name" example:"chenjh-test-snapshot"`                         //快照名称
	Description  string `json:"description" example:"JUST a test VolumeSnapshot For chenjh"` //快照的详细描述(可选)
	ResourceUuid string `json:"resourceUuid" example:""`                                     //资源的Uuid(可选)
}

type UpdateVolumeSnapshotParam struct {
	BaseParam

	UpdateVolumeSnapshot UpdateVolumeSnapshotDetailParam `json:"updateVolumeSnapshot"`
}

type UpdateVolumeSnapshotDetailParam struct {
	Name        string `json:"name"`        //快照的新名称
	Description string `json:"description"` //快照的新详细描述
}
