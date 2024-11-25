// Copyright (c) ZStack.io, Inc.

package param

type VolumeSnapshotGroupParam struct {
	BaseParam

	Params VolumeSnapshotGroupDetailParam `json:"params"`
}

type VolumeSnapshotGroupDetailParam struct {
	RootVolumeUuid string `json:"rootVolumeUuid"` // Root volume UUID
	Name           string `json:"name"`           // Resource name
	Description    string `json:"description"`    // Detailed description of the resource (optional)
	ResourceUuid   string `json:"resourceUuid"`   // Resource Uuid (optional)
}

type UpdateVolumeSnapshotGroupParam struct {
	BaseParam

	UpdateVolumeSnapshotGroup UpdateVolumeSnapshotGroupDetailParam `json:"updateVolumeSnapshotGroup"`
}

type UpdateVolumeSnapshotGroupDetailParam struct {
	Name        string `json:"name"`        // Resource name
	Description string `json:"description"` // Detailed description of the resource
}
