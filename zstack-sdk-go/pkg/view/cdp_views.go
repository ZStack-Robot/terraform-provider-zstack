// Copyright (c) ZStack.io, Inc.

package view

import "zstack.io/zstack-sdk-go/pkg/param"

type CdpPolicyInventoryView struct {
	BaseInfoView
	BaseTimeView
	RetentionTimePerDay     int64                `json:"retentionTimePerDay"`     // Retention time for recovery points
	RecoveryPointPerSecond  int64                `json:"recoveryPointPerSecond"`  // Interval time for recovery points
	State                   param.CdpPolicyState `json:"state"`                   // Policy state
	HourlyRpSinceDay        int64                `json:"hourlyRpSinceDay"`        // Start day to retain one recovery point per hour
	DailyRpSinceDay         int64                `json:"dailyRpSinceDay"`         // Start day to retain one recovery point per day
	ExpireTimeInDay         int64                `json:"expireTimeInDay"`         // Validity time for backup data
	FullBackupIntervalInDay int64                `json:"fullBackupIntervalInDay"` // Interval for full backup
}

type CdpTaskInventoryView struct {
	BaseInfoView
	BaseTimeView

	PolicyUuid        string                    `json:"policyUuid"`        // Policy UUID
	BackupStorageUuid string                    `json:"backupStorageUuid"` // Backup storage UUID
	BackupBandwidth   int64                     `json:"backupBandwidth"`   // Backup speed for a single volume
	MaxCapacity       int64                     `json:"maxCapacity"`       // Planned capacity for the CDP task
	UsedCapacity      int64                     `json:"usedCapacity"`      // Used capacity for CDP
	MaxLatency        int64                     `json:"maxLatency"`        // Maximum RPO offset for the CDP task
	LastLatency       int64                     `json:"lastLatency"`       // Last RPO offset for the CDP task
	Status            param.CdpTaskStatus       `json:"status"`
	State             param.CdpTaskState        `json:"state"`
	TaskType          param.CdpTaskType         `json:"taskType"`
	ResourceRefs      []CdpTaskResourceRefsView `json:"resourceRefs"` // Task resource list
}

type CdpTaskResourceRefsView struct {
	TaskUuid     string `json:"taskUuid"`     // CDP task UUID
	ResourceUuid string `json:"resourceUuid"` // Resource UUID
	ResourceType string `json:"resourceType"` // Task resource list
	BaseTimeView
}

type MountVmInstanceRecoveryPointView struct {
	ResourcePath string `json:"resourcePath"`
}

type UnmountVmInstanceRecoveryPointView MountVmInstanceRecoveryPointView
