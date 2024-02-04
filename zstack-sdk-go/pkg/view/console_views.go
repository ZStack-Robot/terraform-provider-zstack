package view

type ConsoleInventoryView struct {
	Scheme   string `json:"scheme" bson:"scheme"`     //访问协议类型
	Hostname string `json:"hostname" bson:"hostname"` //宿主机名称
	Port     int    `json:"port" bson:"port"`         //端口
	Token    string `json:"token" bson:"token"`       //口令
}
