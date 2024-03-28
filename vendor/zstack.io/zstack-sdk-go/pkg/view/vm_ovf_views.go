// Copyright (c) HashiCorp, Inc.

package view

type OvfInfo struct {
	VmName     string        `json:"vmName"` //云主机名称
	Disks      []OvfDisk     `json:"disks"`
	Networks   []OvfNetwork  `json:"networks"`
	Cpu        OvfCpuInfo    `json:"cpu"`
	Memory     OvfMemoryInfo `json:"memory"`
	Os         OvfOsInfo     `json:"os"`
	SystemInfo OvfSystemInfo `json:"systemInfo"`
	Nics       []OvfNic      `json:"nics"`
	CdDrivers  []OvfCdDriver `json:"cdDrivers"`
	Volumes    []OvfVolume   `json:"volumes"`
}
type OvfDisk struct {
	Index         int    `json:"index"`         //磁盘序号
	DiskId        string `json:"diskId"`        //磁盘ID
	FileRef       string `json:"fileRef"`       //文件引用名称
	FileName      string `json:"fileName"`      //镜像文件名
	Format        string `json:"format"`        //镜像文件格式
	PopulatedSize int64  `json:"populatedSize"` //镜像文件大小
	Capacity      int64  `json:"capacity"`      //磁盘容量，单位Byte
}
type OvfNetwork struct {
	Name string `json:"name"` //网络名称
}
type OvfCpuInfo struct {
	InstanceId     string `json:"instanceId"`     //硬件ID
	Quantity       int    `json:"quantity"`       //CPU内核数量
	CoresPerSocket int    `json:"coresPerSocket"` //每CPU内核数
}
type OvfMemoryInfo struct {
	InstanceId string `json:"instanceId"` //硬件ID
	Quantity   int64  `json:"quantity"`   //内存容量，单位Byte
}
type OvfOsInfo struct {
	Id          int    `json:"id"`          //操作系统ID
	Version     string `json:"version"`     //操作系统版本
	OsType      string `json:"osType"`      //操作系统类型
	Description string `json:"description"` //操作系统描述
}
type OvfSystemInfo struct {
	VirtualSystemType string `json:"virtualSystemType"` //硬件系统类型
	FirmwareType      string `json:"firmwareType"`      //固件类型
}
type OvfNic struct {
	NicName        string `json:"nicName"`        //网络名称
	NicModel       string `json:"nicModel"`       //网卡型号
	NetworkName    string `json:"networkName"`    //网卡名称
	AutoAllocation bool   `json:"autoAllocation"` //是否自动分配
}
type OvfCdDriver struct {
	AutoAllocation bool   `json:"autoAllocation"` //是否自动分配
	DriverType     string `json:"driverType"`     //光驱控制器类型
	SubType        string `json:"subType"`        //子类型
	Name           string `json:"name"`           //光驱名称
}
type OvfVolume struct {
	Name       string `json:"name"`       //磁盘名称
	DiskId     string `json:"diskId"`     //磁盘ID
	DriverType string `json:"driverType"` //磁盘驱动器类型
}
