// Copyright (c) ZStack.io, Inc.

package view

type GetCandidatePrimaryStoragesForCreatingVmView struct {
	DataVolumePrimaryStorages map[string][]PrimaryStorageInventoryView `json:"dataVolumePrimaryStorages"`
	RootVolumePrimaryStorages []PrimaryStorageInventoryView            `json:"rootVolumePrimaryStorages"`
}

type PrimaryStorageInventoryView struct {
	BaseInfoView
	BaseTimeView

	ZoneUuid                  string                                `json:"zoneUuid"` // Zone UUID
	Url                       string                                `json:"url"`
	TotalCapacity             int64                                 `json:"totalCapacity"`
	AvailableCapacity         int64                                 `json:"availableCapacity"`
	TotalPhysicalCapacity     int64                                 `json:"totalPhysicalCapacity"`
	AvailablePhysicalCapacity int64                                 `json:"availablePhysicalCapacity"`
	SystemUsedCapacity        int64                                 `json:"systemUsedCapacity"`
	Type                      string                                `json:"type"` // Cloud host type reserved field, no need to specify. UserVm/ApplianceVm
	State                     string                                `json:"state"`
	Status                    string                                `json:"status"`
	AttachedClusterUuids      []string                              `json:"attachedClusterUuids"`
	Mons                      []CephPrimaryStorageMonsView          `json:"mons"`
	Pools                     []CephPrimaryStoragePoolInventoryView `json:"pools"`
}

type CephPrimaryStorageMonsView struct {
	CreateDate         string `json:"createDate"`
	Hostname           string `json:"hostname"`
	LastOpDate         string `json:"lastOpDate"`
	MonAddr            string `json:"monAddr"`
	MonPort            int    `json:"monPort"`
	MonUuid            string `json:"monUuid"`
	PrimaryStorageUuid string `json:"primaryStorageUuid"`
	SshPassword        string `json:"-"`
	SshPort            int    `json:"sshPort"`
	SshUsername        string `json:"sshUsername"`
	Status             string `json:"status"`
}

type CephPrimaryStoragePoolInventoryView struct {
	BaseTimeView
	Uuid               string  `json:"uuid"`
	AvailableCapacity  int64   `json:"availableCapacity"`
	DiskUtilization    float64 `json:"diskUtilization"`
	PoolName           string  `json:"poolName"`
	PrimaryStorageUuid string  `json:"primaryStorageUuid"`
	ReplicatedSize     int     `json:"replicatedSize"`
	SecurityPolicy     string  `json:"securityPolicy"`
	TotalCapacity      int64   `json:"totalCapacity"`
	Type               string  `json:"type"`
	UsedCapacity       int64   `json:"usedCapacity"`

	AliasName   string `json:"aliasName"`
	Description string `json:"description"`
}
