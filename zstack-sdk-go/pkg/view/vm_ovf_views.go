// Copyright (c) ZStack.io, Inc.

package view

type OvfInfo struct {
	VmName     string        `json:"vmName"` // Cloud host name
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
	Index         int    `json:"index"`         // Disk index
	DiskId        string `json:"diskId"`        // Disk ID
	FileRef       string `json:"fileRef"`       // File reference name
	FileName      string `json:"fileName"`      // Image file name
	Format        string `json:"format"`        // Image file format
	PopulatedSize int64  `json:"populatedSize"` // Image file size
	Capacity      int64  `json:"capacity"`      // Disk capacity, unit Byte
}

type OvfNetwork struct {
	Name string `json:"name"` // Network name
}

type OvfCpuInfo struct {
	InstanceId     string `json:"instanceId"`     // Hardware ID
	Quantity       int    `json:"quantity"`       // Number of CPU cores
	CoresPerSocket int    `json:"coresPerSocket"` // Cores per CPU
}

type OvfMemoryInfo struct {
	InstanceId string `json:"instanceId"` // Hardware ID
	Quantity   int64  `json:"quantity"`   // Memory capacity, unit Byte
}

type OvfOsInfo struct {
	Id          int    `json:"id"`          // Operating system ID
	Version     string `json:"version"`     // Operating system version
	OsType      string `json:"osType"`      // Operating system type
	Description string `json:"description"` // Operating system description
}

type OvfSystemInfo struct {
	VirtualSystemType string `json:"virtualSystemType"` // Hardware system type
	FirmwareType      string `json:"firmwareType"`      // Firmware type
}

type OvfNic struct {
	NicName        string `json:"nicName"`        // Network name
	NicModel       string `json:"nicModel"`       // NIC model
	NetworkName    string `json:"networkName"`    // NIC name
	AutoAllocation bool   `json:"autoAllocation"` // Whether to auto-allocate
}

type OvfCdDriver struct {
	AutoAllocation bool   `json:"autoAllocation"` // Whether to auto-allocate
	DriverType     string `json:"driverType"`     // CD-ROM controller type
	SubType        string `json:"subType"`        // Subtype
	Name           string `json:"name"`           // CD-ROM name
}

type OvfVolume struct {
	Name       string `json:"name"`       // Disk name
	DiskId     string `json:"diskId"`     // Disk ID
	DriverType string `json:"driverType"` // Disk driver type
}
