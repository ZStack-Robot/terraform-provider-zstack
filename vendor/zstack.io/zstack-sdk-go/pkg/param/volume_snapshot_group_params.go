package param

type VolumeSnapshotGroupParam struct {
	BaseParam

	Params VolumeSnapshotGroupDetailParam `json:"params"`
}

type VolumeSnapshotGroupDetailParam struct {
	RootVolumeUuid string `json:"rootVolumeUuid"` //根云盘UUID
	Name           string `json:"name"`           //资源名称
	Description    string `json:"description"`    //资源的详细描述(可选)
	ResourceUuid   string `json:"resourceUuid"`   //资源的Uuid(可选)
}

type UpdateVolumeSnapshotGroupParam struct {
	BaseParam

	UpdateVolumeSnapshotGroup UpdateVolumeSnapshotGroupDetailParam `json:"updateVolumeSnapshotGroup"`
}

type UpdateVolumeSnapshotGroupDetailParam struct {
	Name        string `json:"name"`        //资源名称
	Description string `json:"description"` //资源的详细描述
}
