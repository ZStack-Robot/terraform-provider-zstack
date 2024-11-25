// Copyright (c) ZStack.io, Inc.

package param

type VolumeSnapshotParam struct {
	BaseParam

	Params VolumeSnapshotDetailParam `json:"params"`
}

type VolumeSnapshotDetailParam struct {
	Name         string `json:"name" example:"chenjh-test-snapshot"`                         // Snapshot name
	Description  string `json:"description" example:"JUST a test VolumeSnapshot For chenjh"` // Detailed description of the snapshot (optional)
	ResourceUuid string `json:"resourceUuid" example:""`                                     // Resource Uuid (optional)
}

type UpdateVolumeSnapshotParam struct {
	BaseParam

	UpdateVolumeSnapshot UpdateVolumeSnapshotDetailParam `json:"updateVolumeSnapshot"`
}

type UpdateVolumeSnapshotDetailParam struct {
	Name        string `json:"name"`        // New name for the snapshot
	Description string `json:"description"` // New detailed description for the snapshot
}
