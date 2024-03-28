// Copyright (c) HashiCorp, Inc.

package view

import "time"

type VolumeView struct {
	BaseInfoView
	BaseTimeView

	PrimaryStorageUUID string    `json:"primaryStorageUuid"` //主存储UUID
	VMInstanceUUID     string    `json:"vmInstanceUuid"`     //云主机UUID
	LastVmInstanceUuid string    `json:"lastVmInstanceUuid"` //上一云主机UUID
	DiskOfferingUUID   string    `json:"diskOfferingUuid"`   //云盘规格UUID
	RootImageUUID      string    `json:"rootImageUuid"`
	InstallPath        string    `json:"installPath"`
	Type               string    `json:"type"`
	Format             string    `json:"format"`
	Size               int       `json:"size"`
	ActualSize         int       `json:"actualSize"`
	DeviceID           float32   `json:"deviceId"`
	State              string    `json:"state"`
	Status             string    `json:"status"`
	IsShareable        bool      `json:"isShareable"`
	LastDetachDate     time.Time `json:"lastDetachDate"` //最后一次卸载时间
}

type VolumeFormatView struct {
	Format                    string   `json:"format"`
	MasterHypervisorType      string   `json:"masterHypervisorType"`
	SupportingHypervisorTypes []string `json:"supportingHypervisorTypes"`
}

type VolumeCapabilitiesView struct {
	MigrationToOtherPrimaryStorage   bool `json:"MigrationToOtherPrimaryStorage"`
	MigrationInCurrentPrimaryStorage bool `json:"MigrationInCurrentPrimaryStorage"`
}

type VolumeQoSView struct {
	VolumeUuid                      string `json:"volumeUuid"`                      //云盘UUID
	VolumeBandwidth                 int32  `json:"volumeBandwidth"`                 //云盘带宽，默认-1
	VolumeBandwidthRead             int32  `json:"volumeBandwidthRead"`             //云盘读带宽，默认-1
	VolumeBandwidthWrite            int32  `json:"volumeBandwidthWrite"`            //云盘写带宽，默认-1
	VolumeBandwidthUpthreshold      int32  `json:"volumeBandwidthUpthreshold"`      //云盘带宽上限，默认-1
	VolumeBandwidthReadUpthreshold  int32  `json:"volumeBandwidthReadUpthreshold"`  //云盘读带宽上限，默认-1
	VolumeBandwidthWriteUpthreshold int32  `json:"volumeBandwidthWriteUpthreshold"` //云盘写带宽上限，默认-1
	IopsTotal                       int32  `json:"iopsTotal"`                       //云盘总IOPS
	IopsRead                        int32  `json:"iopsRead"`                        //云盘读取IOPS
	IopsWrite                       int32  `json:"iopsWrite"`                       //云盘写入IOPS
	IopsTotalUpthreshold            int32  `json:"iopsTotalUpthreshold"`            //云盘总IOPS上限，-1表示无上限限制
	IopsReadUpthreshold             int32  `json:"iopsReadUpthreshold"`             //云盘读取IOPS上限，-1表示无上限限制
	IopsWriteUpthreshold            int32  `json:"iopsWriteUpthreshold"`            //云盘写入IOPS上限，-1表示无上限限制
}
