// Copyright (c) ZStack.io, Inc.

package view

type VirtualRouterOfferingInventoryView struct {
	BaseInfoView
	BaseTimeView

	CpuNum                int    `json:"cpuNum"`            // Number of CPUs
	CpuSpeed              int    `json:"cpuSpeed"`          // CPU speed
	MemorySize            int64  `json:"memorySize"`        // Memory size
	Type                  string `json:"type"`              // Type
	AllocatorStrategy     string `json:"allocatorStrategy"` // Allocation strategy
	SortKey               int    `json:"sortKey"`
	State                 string `json:"state"` // State (Enabled, Disabled)
	ManagementNetworkUuid string `json:"managementNetworkUuid"`
	PublicNetworkUuid     string `json:"publicNetworkUuid"`
	ZoneUuid              string `json:"zoneUuid"`
	ImageUuid             string `json:"imageUuid"`
	IsDefault             bool   `json:"isDefault"`
	ReservedMemorySize    string `json:"reservedMemorySize"`
}
