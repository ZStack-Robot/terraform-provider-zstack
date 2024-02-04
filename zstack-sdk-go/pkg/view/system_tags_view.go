package view

type SystemTagView struct {
	BaseInfoView
	BaseTimeView

	Inherent     bool   `json:"inherent"`     //内部系统标签
	ResourceUuid string `json:"resourceUuid"` //用户指定的资源UUID，若指定，系统不会为该资源随机分配UUID
	ResourceType string `json:"resourceType"` //当创建一个标签时, 用户必须制定标签所关联的资源类型(resource type)
	Tag          string `json:"tag"`          //标签字符串
}
