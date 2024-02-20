package view

type BackupStorageInventoryView struct {
	BaseInfoView
	BaseTimeView
	Url               string   `json:"url" `
	TotalCapacity     int64    `json:"totalCapacity" `
	AvailableCapacity int64    `json:"availableCapacity" `
	Type              string   `json:"type"` //云主机类型 保留字段，无需指定。UserVm/ApplianceVm
	State             string   `json:"state" `
	Status            string   `json:"status" `
	AttachedZoneUuids []string `json:"attachedZoneUuids" `

	// imageStoreBackupStorage
	Hostname string `json:"hostname"`
	Username string `json:"username"`
	SshPort  int    `json:"sshPort"`

	// Ceph
	Fsid                  string    `json:"fsid"`
	Mons                  []CephMon `json:"mons"`
	PoolAvailableCapacity int64     `json:"poolAvailableCapacity"`
	PoolName              string    `json:"poolName"`
	PoolReplicatedSize    int       `json:"poolReplicatedSize"`
	PoolUsedCapacity      int64     `json:"poolUsedCapacity"`
}

type T struct {
	State         string `json:"state"`
	Status        string `json:"status"`
	TotalCapacity int64  `json:"totalCapacity"`
	Type          string `json:"type"`
	Url           string `json:"url"`
	Uuid          string `json:"uuid"`
}

type CephMon struct {
	BackupStorageUuid string `json:"backupStorageUuid"`
	CreateDate        string `json:"createDate"`
	Hostname          string `json:"hostname"`
	LastOpDate        string `json:"lastOpDate"`
	MonAddr           string `json:"monAddr"`
	MonPort           int    `json:"monPort"`
	MonUuid           string `json:"monUuid"`
	SshPassword       string `json:"sshPassword"`
	SshPort           int    `json:"sshPort"`
	SshUsername       string `json:"sshUsername"`
	Status            string `json:"status"`
}

type ExportImageFromBackupStorageResultView struct {
	ImageUrl     string `json:"imageUrl"`     //导出的镜像的URL
	ExportMd5Sum string `json:"exportMd5Sum"` //导出的镜像的MD5值
	Success      bool   `json:"success"`      //导出是否成功
	Error        string `json:"error"`        //导出失败的错误信息
}

type GcResultView struct {
	FreedSpaceInBytes int64 `json:"freedSpaceInBytes"`
}
