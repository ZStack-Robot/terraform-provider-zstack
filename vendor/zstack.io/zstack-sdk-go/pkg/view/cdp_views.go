package view

import "zstack.io/zstack-sdk-go/pkg/param"

type CdpPolicyInventoryView struct {
	BaseInfoView
	BaseTimeView
	RetentionTimePerDay     int64                `json:"retentionTimePerDay"`     //恢复点保留时间
	RecoveryPointPerSecond  int64                `json:"recoveryPointPerSecond"`  //恢复点间隔时间
	State                   param.CdpPolicyState `json:"state"`                   //策略状态
	HourlyRpSinceDay        int64                `json:"hourlyRpSinceDay"`        // 从哪天开始保留每小时一个恢复点
	DailyRpSinceDay         int64                `json:"dailyRpSinceDay"`         // 从哪天开始保留每天一个恢复点
	ExpireTimeInDay         int64                `json:"expireTimeInDay"`         // 备份数据有效时间
	FullBackupIntervalInDay int64                `json:"fullBackupIntervalInDay"` // 全量备份时间间隔
}

type CdpTaskInventoryView struct {
	BaseInfoView
	BaseTimeView

	PolicyUuid        string                    `json:"policyUuid"`        // 权限策略UUID
	BackupStorageUuid string                    `json:"backupStorageUuid"` // 镜像存储UUID
	BackupBandwidth   int64                     `json:"backupBandwidth"`   // 单个云盘备份速率
	MaxCapacity       int64                     `json:"maxCapacity"`       // CDP任务规划容量
	UsedCapacity      int64                     `json:"usedCapacity"`      //CDP已用容量
	MaxLatency        int64                     `json:"maxLatency"`        // CDP任务RPO最大偏移量
	LastLatency       int64                     `json:"lastLatency"`       // CDP任务RPO最后偏移量
	Status            param.CdpTaskStatus       `json:"status"`
	State             param.CdpTaskState        `json:"state"`
	TaskType          param.CdpTaskType         `json:"taskType"`
	ResourceRefs      []CdpTaskResourceRefsView `json:"resourceRefs"` // 任务资源清单
}

type CdpTaskResourceRefsView struct {
	TaskUuid     string `json:"taskUuid"`     // CDP任务UUID
	ResourceUuid string `json:"resourceUuid"` // 资源UUID
	ResourceType string `json:"resourceType"` // 任务资源清单
	BaseTimeView
}

type MountVmInstanceRecoveryPointView struct {
	ResourcePath string `json:"resourcePath"`
}

type UnmountVmInstanceRecoveryPointView MountVmInstanceRecoveryPointView
