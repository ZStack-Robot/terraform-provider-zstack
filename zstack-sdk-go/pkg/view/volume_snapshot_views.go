// Copyright (c) ZStack.io, Inc.

package view

type VolumeSnapshotView struct {
	BaseInfoView
	BaseTimeView

	Type                      string `json:"type"`
	VolumeUUID                string `json:"volumeUuid"`
	TreeUUID                  string `json:"treeUuid"`
	ParentUUID                string `json:"parentUuid"`
	PrimaryStorageUUID        string `json:"primaryStorageUuid"`
	PrimaryStorageInstallPath string `json:"primaryStorageInstallPath"`
	VolumeType                string `json:"volumeType"`
	Format                    string `json:"format"`
	Latest                    bool   `json:"latest"`
	Size                      int64  `json:"size"`
	State                     string `json:"state"`
	Status                    string `json:"status"`
	Distance                  int    `json:"distance"`
	GroupUuid                 string `json:"groupUuid"`
}

type VolumeSnapshotTreeView struct {
	BaseInfoView
	BaseTimeView

	Current    bool                       `json:"current"`
	Tree       VolumeSnapshotTreeNodeView `json:"tree"`
	Status     string                     `json:"status"`
	VolumeUUID string                     `json:"volumeUuid"`
}

type VolumeSnapshotTreeNodeView struct {
	Inventory VolumeSnapshotView           `json:"inventory"`
	Children  []VolumeSnapshotTreeNodeView `json:"children"`
}

type VolumeSnapshotSizeView struct {
	Size       int64 `json:"size"`
	ActualSize int64 `json:"actualSize"`
	Success    bool  `json:"success"`
}

type VolumeSnapshotShrinkResultView struct {
	Result struct {
		OldSize   int64 `json:"oldSize"`
		Size      int64 `json:"size"`
		DeltaSize int64 `json:"deltaSize"`
	} `json:"shrinkResult"`
}
