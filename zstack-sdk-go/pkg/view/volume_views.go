// Copyright (c) ZStack.io, Inc.

package view

import "time"

type VolumeView struct {
	BaseInfoView
	BaseTimeView

	PrimaryStorageUUID string    `json:"primaryStorageUuid"` // Primary storage UUID
	VMInstanceUUID     string    `json:"vmInstanceUuid"`     // Cloud VM UUID
	LastVmInstanceUuid string    `json:"lastVmInstanceUuid"` // Last cloud VM UUID
	DiskOfferingUUID   string    `json:"diskOfferingUuid"`   // Cloud disk specification UUID
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
	LastDetachDate     time.Time `json:"lastDetachDate"` // Last detach time
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
	VolumeUuid                      string `json:"volumeUuid"`                      // Cloud disk UUID
	VolumeBandwidth                 int32  `json:"volumeBandwidth"`                 // Cloud disk bandwidth, default -1
	VolumeBandwidthRead             int32  `json:"volumeBandwidthRead"`             // Cloud disk read bandwidth, default -1
	VolumeBandwidthWrite            int32  `json:"volumeBandwidthWrite"`            // Cloud disk write bandwidth, default -1
	VolumeBandwidthUpthreshold      int32  `json:"volumeBandwidthUpthreshold"`      // Cloud disk bandwidth upper limit, default -1
	VolumeBandwidthReadUpthreshold  int32  `json:"volumeBandwidthReadUpthreshold"`  // Cloud disk read bandwidth upper limit, default -1
	VolumeBandwidthWriteUpthreshold int32  `json:"volumeBandwidthWriteUpthreshold"` // Cloud disk write bandwidth upper limit, default -1
	IopsTotal                       int32  `json:"iopsTotal"`                       // Total IOPS of the cloud disk
	IopsRead                        int32  `json:"iopsRead"`                        // Read IOPS of the cloud disk
	IopsWrite                       int32  `json:"iopsWrite"`                       // Write IOPS of the cloud disk
	IopsTotalUpthreshold            int32  `json:"iopsTotalUpthreshold"`            // Total IOPS upper limit of the cloud disk, -1 means no upper limit
	IopsReadUpthreshold             int32  `json:"iopsReadUpthreshold"`             // Read IOPS upper limit of the cloud disk, -1 means no upper limit
	IopsWriteUpthreshold            int32  `json:"iopsWriteUpthreshold"`            // Write IOPS upper limit of the cloud disk, -1 means no upper limit
}
