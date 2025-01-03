// Copyright (c) ZStack.io, Inc.

package view

type DiskOfferingInventoryView struct {
	BaseInfoView
	BaseTimeView

	DiskSize          int    `json:"diskSize"`          // Disk size
	Type              string `json:"type"`              // Type
	AllocatorStrategy string `json:"allocatorStrategy"` // Allocation strategy
	SortKey           int    `json:"sortKey"`
	State             string `json:"state"` // State (Enabled, Disabled)
}
