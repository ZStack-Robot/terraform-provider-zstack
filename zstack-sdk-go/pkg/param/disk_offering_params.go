// Copyright (c) ZStack.io, Inc.

package param

type CreateDiskOfferingParam struct {
	BaseParam
	Params CreateDiskOfferingDetailParam `json:"params"`
}

type CreateDiskOfferingDetailParam struct {
	Name              string   `json:"name" validate:"required"`     // Resource name
	Description       *string  `json:"description"`                  // Detailed description of the resource
	DiskSize          int64    `json:"diskSize" validate:"required"` // disk size
	AllocatorStrategy *string  `json:"allocatorStrategy"`            // Allocation strategy
	SortKey           *int     `json:"sortKey"`                      // Sort key
	Type              *string  `json:"type"`                         // Type
	ResourceUuid      *string  `json:"resourceUuid"`                 // Resource UUID
	SystemTags        []string `json:"systemTags"`
	UserTags          []string `json:"userTags"`
}
