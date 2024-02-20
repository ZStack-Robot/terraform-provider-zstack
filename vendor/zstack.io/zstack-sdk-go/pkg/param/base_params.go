package param

type DeleteMode string

const (
	DeleteModePermissive DeleteMode = "Permissive"
	DeleteModeEnforcing  DeleteMode = "Enforcing"
)

type BaseParam struct {
	SystemTags []string `json:"systemTags,omitempty"` //系统标签
	UserTags   []string `json:"userTags,omitempty"`   //用户标签
	RequestIp  string   `json:"requestIp,omitempty"`  //请求IP
}

type HqlParam struct {
	OperationName string    `json:"operationName"` //请求名
	Query         string    `json:"query"`         //查询语句
	Variables     Variables `json:"variables"`     //语句对应参数
}

type Variables struct {
	Conditions      []Condition            `json:"conditions"`      //
	ExtraConditions []Condition            `json:"extraConditions"` //
	Input           map[string]interface{} `json:"input"`           //
	PageVar         `json:",inline,omitempty"`
	Type            string `json:"type"` //
}

type Condition struct {
	Key   string `json:"key"`   //
	Op    string `json:"op"`    //
	Value string `json:"value"` //
}

type PageVar struct {
	Start int `json:"start,omitempty"`
	Limit int `json:"limit,omitempty"`
}
