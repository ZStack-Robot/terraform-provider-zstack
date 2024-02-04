package param

type CreateInstanceOfferingParam struct {
	BaseParam
	Params CreateInstanceOfferingDetailParam `json:"params"`
}
type CreateInstanceOfferingDetailParam struct {
	Name              string   `json:"name" validate:"required"`       //资源名称
	Description       *string  `json:"description" `                   //资源的详细描述
	CpuNum            int      `json:"cpuNum" validate:"required"`     //CPU数目
	MemorySize        int64    `json:"memorySize" validate:"required"` //内存大小, 单位Byte
	AllocatorStrategy *string  `json:"allocatorStrategy" `             //分配策略
	SortKey           *int     `json:"sortKey" `                       //排序键
	Type              *string  `json:"type" `                          //类型
	ResourceUuid      *string  `json:"resourceUuid" `                  //资源UUID
	TagUuids          []string `json:"tagUuids" `
}
