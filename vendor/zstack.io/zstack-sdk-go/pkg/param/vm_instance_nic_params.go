package param

type AttachL3NetworkToVmParam struct {
	BaseParam
	Params AttachL3NetworkToVmDetailParam `json:"params"`
}

type AttachL3NetworkToVmDetailParam struct {
	StaticIp string `json:"staticIp"` //指定分配给云主机的IP地址
}

type UpdateVmNicMacParam struct {
	BaseParam
	UpdateVmNicMac UpdateVmNicMacDetailParam `json:"updateVmNicMac"`
}

type UpdateVmNicMacDetailParam struct {
	Mac string `json:"mac"` //mac地址
}

type SetVmStaticIpParam struct {
	BaseParam
	SetVmStaticIp SetVmStaticIpDetailParam `json:"setVmStaticIp"`
}

type SetVmStaticIpDetailParam struct {
	L3NetworkUuid string `json:"l3NetworkUuid"` //三层网络UUID
	Ip            string `json:"ip"`            //指定IP地址
	Ip6           string `json:"ip6"`           //指定IPv6地址
}

type DeleteVmStaticIpParam struct {
	BaseParam
	Params DeleteVmStaticIpDetailParam `json:"params"`
}

type DeleteVmStaticIpDetailParam struct {
	L3NetworkUuid string     `json:"l3NetworkUuid"` //三层网络UUID
	DeleteMode    DeleteMode `json:"deleteMode"`
}

type ChangeVmNicNetworkParam struct {
	BaseParam
	Params ChangeVmNicNetworkDetailParam `json:"params"`
}

type ChangeVmNicNetworkDetailParam struct {
	DestL3NetworkUuid string `json:"destL3NetworkUuid"` //指定三层网络UUID
}
