// Copyright (c) HashiCorp, Inc.

package view

import "time"

type VmNicInventoryView struct {
	UUID           string   `json:"uuid"`           //资源的UUID，唯一标示该资源
	VMInstanceUUID string   `json:"vmInstanceUuid"` //	云主机UUID
	L3NetworkUUID  string   `json:"l3NetworkUuid"`  //	三层网络UUID
	IP             string   `json:"ip"`             //ip地址
	Mac            string   `json:"mac"`            //mac地址
	HypervisorType string   `json:"hypervisorType"` //虚拟化类型
	Netmask        string   `json:"netmask"`        //子网掩码
	Gateway        string   `json:"gateway"`        //网关
	MetaData       string   `json:"metaData"`       //内部使用的保留域，元数据
	IpVersion      int      `json:"ipVersion"`      //IP地址版本
	DeviceID       int      `json:"deviceId"`       //设备ID 标识网卡在客户操作系统（guest operating system）以太网设备中顺序的整形数字。例如， 0通常代表eth0，1通常代表eth1。
	DriverType     string   `json:"driverType"`     //网卡型号
	Type           string   `json:"type"`           //网卡类型
	CreateDate     string   `json:"createDate"`     //创建时间
	LastOpDate     string   `json:"lastOpDate"`     //最后一次修改时间
	InternalName   string   `json:"internalName"`   //
	UsedIps        []UsedIp `json:"usedIps"`
}

type UsedIp struct {
	Uuid          string    `json:"uuid"`          //资源的UUID，唯一标示该资源
	IpRangeUuid   string    `json:"ipRangeUuid"`   //IP段UUID
	L3NetworkUuid string    `json:"l3NetworkUuid"` //三层网络UUID
	IpVersion     int       `json:"ipVersion"`     //IP协议号
	Ip            string    `json:"ip"`            //IP地址
	Netmask       string    `json:"netmask"`       //网络掩码
	Gateway       string    `json:"gateway"`       //网关地址
	UsedFor       string    `json:"usedFor"`       //
	IpInLong      int64     `json:"ipInLong"`      //
	VmNicUuid     string    `json:"vmNicUuid"`     //云主机网卡UUID
	CreateDate    time.Time `json:"createDate"`    //创建时间
	LastOpDate    time.Time `json:"lastOpDate"`    //最后一次修改时间
}

func GetIpFromUsedIps(usedIps []UsedIp) (ip string, ip6 string) {
	for _, usedIp := range usedIps {
		if usedIp.IpVersion == 4 {
			ip = usedIp.Ip
		}
		if usedIp.IpVersion == 6 {
			ip6 = usedIp.Ip
		}
	}
	return
}

type NicSimpleView struct {
	Ip        string    `json:"ip"`        //
	IpVersion string    `json:"ipVersion"` //
	Uuid      string    `json:"uuid"`      //
	VmNicUuid string    `json:"vmNicUuid"` //
	VmNic     VmNicView `json:"vmNic"`     //
}
type VmNicView struct {
	InternalName string `json:"internalName"` //
	Mac          string `json:"mac"`
	Uuid         string `json:"uuid"`
}
