package view

type VipInventoryView struct {
	BaseInfoView
	BaseTimeView

	L3NetworkUUID      string         `json:"l3NetworkUuid"` //	三层网络UUID
	Ip                 string         `json:"ip"`
	State              string         `json:"state"`
	Gateway            string         `json:"gateway"`            //网关
	Netmask            string         `json:"netmask"`            //子网掩码
	PrefixLen          string         `json:"prefixLen"`          //掩码长度
	ServiceProvider    string         `json:"serviceProvider"`    //提供VIP服务的服务提供者
	PeerL3NetworkUuids string         `json:"peerL3NetworkUuids"` //提供VIP服务的L3网络UUID
	UseFor             string         `json:"useFor"`             //用途，例如：端口转发
	System             bool           `json:"system"`             //是否系统创建
	ServicesRefs       []ServicesRefs `json:"servicesRefs"`
}

type ServicesRefs struct {
	UUID        string `json:"uuid"`        //资源的UUID，唯一标示该资源
	ServiceType string `json:"serviceType"` //服务类型
	VipUuid     string `json:"vipUuid"`
	CreateDate  string `json:"createDate"` //创建时间
	LastOpDate  string `json:"lastOpDate"` //最后一次修改时间
}
