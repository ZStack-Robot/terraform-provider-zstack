// Copyright (c) ZStack.io, Inc.

package view

type InstanceOfferingInventoryView struct {
	BaseInfoView
	BaseTimeView

	CpuNum            int    `json:"cpuNum"`            // Number of CPUs
	CpuSpeed          int    `json:"cpuSpeed"`          // CPU speed
	MemorySize        int64  `json:"memorySize"`        // Memory size
	Type              string `json:"type"`              // Type
	AllocatorStrategy string `json:"allocatorStrategy"` // Allocation strategy
	SortKey           int    `json:"sortKey"`
	State             string `json:"state"` // State (Enabled, Disabled)
}
