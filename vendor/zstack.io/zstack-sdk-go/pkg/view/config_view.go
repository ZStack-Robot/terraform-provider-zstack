package view

type GlobalConfigView struct {
	BaseInfoView

	Category     string `json:"category"`
	DefaultValue string `json:"defaultValue"`
	Value        string `json:"value"`
}

type ResourceConfigView struct {
	BaseInfoView
	BaseTimeView

	ResourceUuid string `json:"resourceUuid"` //配置对应的资源UUID
	ResourceType string `json:"resourceType"` //配置对应的资源类型
	Category     string `json:"category"`     //配置类别
	Value        string `json:"value"`        //配置的值
}
