// Copyright (c) ZStack.io, Inc.

package view

type VolumeSnapshotGroupView struct {
	BaseInfoView
	BaseTimeView

	SnapshotCount      int                     `json:"snapshotCount"`  // Number of snapshots in the group
	VmInstanceUuid     string                  `json:"vmInstanceUuid"` // Cloud host UUID
	VolumeSnapshotRefs []VolumeSnapshotRefView `json:"volumeSnapshotRefs"`
}

type VolumeSnapshotRefView struct {
	BaseTimeView

	VolumeSnapshotUuid        string `json:"volumeSnapshotUuid"`        // Cloud disk snapshot UUID
	VolumeSnapshotGroupUuid   string `json:"volumeSnapshotGroupUuid"`   // Snapshot group UUID
	DeviceId                  int    `json:"deviceId"`                  // Mount sequence number of the cloud disk when the snapshot was taken
	SnapshotDeleted           bool   `json:"snapshotDeleted"`           // Whether the snapshot has been deleted
	VolumeUuid                string `json:"volumeUuid"`                // Cloud disk UUID
	VolumeName                string `json:"volumeName"`                // Name of the cloud disk
	VolumeType                string `json:"volumeType"`                // Type of the cloud disk
	VolumeSnapshotInstallPath string `json:"volumeSnapshotInstallPath"` // Installation path of the snapshot
	VolumeSnapshotName        string `json:"volumeSnapshotName"`        // Name of the snapshot
}

type VolumeSnapshotGroupAvailabilityView struct {
	UUID      string `json:"uuid"`      // Resource UUID, uniquely identifies the resource
	Available bool   `json:"available"` // Whether it can be restored
	Reason    string `json:"reason"`    // Reason for not being able to restore, empty if it can be restored
}
