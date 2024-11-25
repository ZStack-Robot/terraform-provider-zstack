// Copyright (c) ZStack.io, Inc.

package param

type CreateInstanceOfferingParam struct {
	BaseParam
	Params CreateInstanceOfferingDetailParam `json:"params"`
}

type CreateInstanceOfferingDetailParam struct {
	Name              string   `json:"name" validate:"required"`       // Resource name
	Description       *string  `json:"description"`                    // Detailed description of the resource
	CpuNum            int      `json:"cpuNum" validate:"required"`     // Number of CPUs
	MemorySize        int64    `json:"memorySize" validate:"required"` // Memory size in bytes
	AllocatorStrategy *string  `json:"allocatorStrategy"`              // Allocation strategy
	SortKey           *int     `json:"sortKey"`                        // Sort key
	Type              *string  `json:"type"`                           // Type
	ResourceUuid      *string  `json:"resourceUuid"`                   // Resource UUID
	TagUuids          []string `json:"tagUuids"`
}
