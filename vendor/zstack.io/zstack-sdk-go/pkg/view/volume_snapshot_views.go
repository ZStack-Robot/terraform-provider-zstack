// Copyright (c) HashiCorp, Inc.

package view

type VolumeSnapshotView struct {
	BaseInfoView
	BaseTimeView

	Type                      string `json:"type"`
	VolumeUUID                string `json:"volumeUuid"` //云盘UUID
	TreeUUID                  string `json:"treeUuid"`
	ParentUUID                string `json:"parentUuid"`
	PrimaryStorageUUID        string `json:"primaryStorageUuid"` //主存储UUID
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
	Size       int64 `json:"size"`       //快照容量
	ActualSize int64 `json:"actualSize"` //快照实际容量
	Success    bool  `json:"success"`
}

type VolumeSnapshotShrinkResultView struct {
	Result struct {
		OldSize   int64 `json:"oldSize"`
		Size      int64 `json:"size"`
		DeltaSize int64 `json:"deltaSize"`
	} `json:"shrinkResult"`
}
