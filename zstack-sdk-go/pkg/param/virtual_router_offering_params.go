// Copyright (c) ZStack.io, Inc.

package param

type CreateVirtualRouterOfferingParam struct {
	BaseParam
	Params CreateVirtualRouterOfferingDetailParam `json:"params"`
}

type CreateVirtualRouterOfferingDetailParam struct {
	Name                  string   `json:"name" validate:"required"` // Resource name
	Description           string   `json:"description"`              // Detailed description of the resource
	ZoneUuid              string   `json:"zoneUuid" validate:"required"`
	ManagementNetworkUuid string   `json:"managementNetworkUuid" validate:"required"`
	ImageUuid             string   `json:"imageUuid" validate:"required"`
	PublicNetworkUuid     string   `json:"publicNetworkUuid"`
	IsDefault             bool     `json:"isDefault"`
	CpuNum                int      `json:"cpuNum" validate:"required"`     // Number of CPUs
	MemorySize            int64    `json:"memorySize" validate:"required"` // Memory size in bytes
	AllocatorStrategy     string   `json:"allocatorStrategy"`              // Allocation strategy
	SortKey               int      `json:"sortKey"`                        // Sort key
	Type                  string   `json:"type" validate:"required"`       // Type
	ResourceUuid          string   `json:"resourceUuid"`                   // Resource UUID
	TagUuids              []string `json:"tagUuids"`
}
