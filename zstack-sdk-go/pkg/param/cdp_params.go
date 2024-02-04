package param

type CdpPolicyState string
type CdpTaskType string
type CdpTaskStatus string
type CdpTaskState string

const (
	CdpPolicyEnabled  CdpPolicyState = "Enabled"
	CdpPolicyDisabled CdpPolicyState = "Disabled"

	CdpTaskTypeVM CdpTaskType = "VM"

	CdpTaskStatusCreated  CdpTaskStatus = "Created"
	CdpTaskStatusStarting CdpTaskStatus = "Starting"
	CdpTaskStatusRunning  CdpTaskStatus = "Running"
	CdpTaskStatusStopped  CdpTaskStatus = "Stopped"
	CdpTaskStatusUnknown  CdpTaskStatus = "Unknown"
	CdpTaskStatusFailed   CdpTaskStatus = "Failed"

	CdpTaskStateEnabled  CdpTaskState = "Enabled"
	CdpTaskStateDisabled CdpTaskState = "Disabled"
)

type CreateCdpPolicyParam struct {
	BaseParam
	Params CreateCdpPolicyDetailParam `json:"params"`
}

type CreateCdpPolicyDetailParam struct {
	Name                    string `json:"name"`
	Description             string `json:"description"`
	RecoveryPointPerSecond  int64  `json:"recoveryPointPerSecond"`  //恢复点间隔时间
	HourlyRpSinceDay        int64  `json:"hourlyRpSinceDay"`        //从哪天开始保留每小时一个恢复点
	DailyRpSinceDay         int64  `json:"dailyRpSinceDay"`         //从哪天开始保留每天一个恢复点
	ExpireTimeInDay         int64  `json:"expireTimeInDay"`         // 备份数据有效时间
	FullBackupIntervalInDay int64  `json:"fullBackupIntervalInDay"` //全量备份时间间隔
	ResourceUuid            string `json:"resourceUuid"`            //资源uuid
}

type UpdateCdpPolicyParam struct {
	BaseParam
	UpdateCdpPolicy UpdateCdpPolicyDetailParam `json:"updateCdpPolicy"`
}

type UpdateCdpPolicyDetailParam struct {
	Name                    *string `json:"name"`
	Description             *string `json:"description"`
	RetentionTimePerDay     *int64  `json:"retentionTimePerDay"`     // 恢复点保留时间
	RecoveryPointPerSecond  *int64  `json:"recoveryPointPerSecond"`  // 恢复点间隔时间
	HourlyRpSinceDay        *int64  `json:"hourlyRpSinceDay"`        // 从哪天开始保留每小时一个恢复点
	DailyRpSinceDay         *int64  `json:"dailyRpSinceDay"`         // 从哪天开始保留每天一个恢复点
	ExpireTimeInDay         *int64  `json:"expireTimeInDay"`         // 备份数据有效时间
	FullBackupIntervalInDay *int64  `json:"fullBackupIntervalInDay"` // 全量备份时间间隔
}

type CreateCdpTaskParam struct {
	BaseParam
	Params CreateCdpTaskDetailParam `json:"params"`
}

type CreateCdpTaskDetailParam struct {
	Name              string      `json:"name"`
	Description       string      `json:"description"`
	TaskType          CdpTaskType `json:"taskType"`          // CDP任务类型
	PolicyUuid        string      `json:"policyUuid"`        // 权限策略UUID
	BackupStorageUuid string      `json:"backupStorageUuid"` // 镜像存储UUID
	ResourceUuids     []string    `json:"resourceUuids"`     // 备份资源列表
	BackupBandwidth   int64       `json:"backupBandwidth"`   // 单个云盘备份速率
	MaxCapacity       int64       `json:"maxCapacity"`       // CDP任务规划容量
	MaxLatency        int64       `json:"maxLatency"`        // CDP任务RPO最大偏移量
}

type UpdateCdpTaskParam struct {
	BaseParam
	UpdateCdpTask UpdateCdpTaskDetailParam `json:"updateCdpTask"`
}
type UpdateCdpTaskDetailParam struct {
	Name            *string `json:"name"`
	Description     *string `json:"description"`
	BackupBandwidth *int64  `json:"backupBandwidth"`
	MaxCapacity     *int64  `json:"maxCapacity"`
	MaxLatency      *int64  `json:"maxLatency"`
}

type MountVmInstanceRecoveryPointParam struct {
	BaseParam
	Params MountVmInstanceRecoveryPointDetailParam `json:"params"`
}

type MountVmInstanceRecoveryPointDetailParam struct {
	VmUuid  string `json:"vmUuid"`
	GroupId int64  `json:"groupId"`
	Https   bool   `json:"https"`
}

type UnmountVmInstanceRecoveryPointParam MountVmInstanceRecoveryPointParam

type CreateVmFromCdpBackupParam struct {
	BaseParam
	CreateVmFromCdpBackup CreateVmFromCdpBackupDetailParam `json:"createVmFromCdpBackup"`
}
type CreateVmFromCdpBackupDetailParam struct {
	Name                            string   `json:"name" validate:"required"`
	GroupId                         int64    `json:"groupId" validate:"required"`
	CdpTaskUuid                     string   `json:"cdpTaskUuid" validate:"required"`
	InstanceOfferingUuid            string   `json:"instanceOfferingUuid" validate:"required"`
	DefaultL3NetworkUuid            string   `json:"defaultL3NetworkUuid" `
	L3NetworkUuids                  []string `json:"l3NetworkUuids" validate:"required"`
	Type                            *string  `json:"type" `
	ZoneUuid                        *string  `json:"zoneUuid" `
	ClusterUuid                     *string  `json:"clusterUuid" `
	HostUuid                        *string  `json:"hostUuid" `
	PrimaryStorageUuidForRootVolume *string  `json:"primaryStorageUuidForRootVolume" `
	PrimaryStorageUuidForDataVolume *string  `json:"primaryStorageUuidForDataVolume" `
	RecoverBandwidth                *int64   `json:"recoverBandwidth" `
	Description                     *string  `json:"description" `
	RootVolumeSystemTags            []string `json:"rootVolumeSystemTags" `
	DataVolumeSystemTags            []string `json:"dataVolumeSystemTags" `
	ResourceUuid                    *string  `json:"resourceUuid" `
	TagUuids                        []string `json:"tagUuids" `
}
type CreateVmFromCdpBackupJobData struct {
	BaseParam
	CreateVmFromCdpBackupDetailParam
}

type RevertVmFromCdpBackupParam struct {
	BaseParam
	RevertVmFromCdpBackup RevertVmFromCdpBackupDetailParam `json:"revertVmFromCdpBackup"`
}
type RevertVmFromCdpBackupDetailParam struct {
	BackupStorageUuid               string   `json:"backupStorageUuid" validate:"required"`
	GroupId                         int64    `json:"groupId" validate:"required"`
	PrimaryStorageUuidForRootVolume *string  `json:"primaryStorageUuidForRootVolume" `
	PrimaryStorageUuidForDataVolume *string  `json:"primaryStorageUuidForDataVolume" `
	RootVolumeSystemTags            []string `json:"rootVolumeSystemTags" `
	DataVolumeSystemTags            []string `json:"dataVolumeSystemTags" `
	HostUuid                        *string  `json:"hostUuid" `
	UseExistingVolume               *bool    `json:"useExistingVolume" `
	RecoverBandwidth                *int64   `json:"recoverBandwidth" `
}
type RevertVmFromCdpBackupJobData struct {
	BaseParam
	RevertVmFromCdpBackupDetailParam
	VmInstanceUuid string `json:"vmInstanceUuid"`
	StopVm         bool   `json:"stopVm"`
}
