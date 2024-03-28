// Copyright (c) HashiCorp, Inc.

package view

type VolumeSnapshotGroupView struct {
	BaseInfoView
	BaseTimeView

	SnapshotCount      int                     `json:"snapshotCount"`  //组内快照数量
	VmInstanceUuid     string                  `json:"vmInstanceUuid"` //云主机UUID
	VolumeSnapshotRefs []VolumeSnapshotRefView `json:"volumeSnapshotRefs"`
}

type VolumeSnapshotRefView struct {
	BaseTimeView

	VolumeSnapshotUuid        string `json:"volumeSnapshotUuid"`        //云盘快照UUID
	VolumeSnapshotGroupUuid   string `json:"volumeSnapshotGroupUuid"`   //快照组UUID
	DeviceId                  int    `json:"deviceId"`                  //打快照时云盘的加载序号
	SnapshotDeleted           bool   `json:"snapshotDeleted"`           //快照是否已经被删除
	VolumeUuid                string `json:"volumeUuid"`                //云盘UUID
	VolumeName                string `json:"volumeName"`                //云盘的名字
	VolumeType                string `json:"volumeType"`                //云盘的类型
	VolumeSnapshotInstallPath string `json:"volumeSnapshotInstallPath"` //快照的安装路径
	VolumeSnapshotName        string `json:"volumeSnapshotName"`        //快照的名字
}

type VolumeSnapshotGroupAvailabilityView struct {
	UUID      string `json:"uuid"`      //资源的UUID，唯一标示该资源
	Available bool   `json:"available"` //是否可以恢复
	Reason    string `json:"reason"`    //不可恢复的理由，如可恢复则为空
}
