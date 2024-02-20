package view

import "time"

type LicenseInventoryView struct {
	UUID               string    `json:"uuid"`               //资源的UUID，唯一标示该资源
	User               string    `json:"user"`               //许可证所属用户名
	ProdInfo           string    `json:"prodInfo"`           //许可证产品名称
	CpuNum             int       `json:"cpuNum"`             //许可证授权X86 CPU数量
	HostNum            int       `json:"hostNum"`            //许可证授权X86 服务器数量
	VmNum              int       `json:"vmNum"`              //许可证授权X86 VM数量
	LicenseType        string    `json:"licenseType"`        //许可证类型
	ExpiredDate        time.Time `json:"expiredDate"`        //许可证过期时间
	IssuedDate         time.Time `json:"issuedDate"`         //许可证申请时间
	UploadDate         time.Time `json:"uploadDate"`         //许可证上传时间
	ManagementNodeUuid string    `json:"managementNodeUuid"` //许可证所属MN节点UUID
	Expired            bool      `json:"expired"`            //许可证是否过期
	LicenseRequest     string    `json:"licenseRequest"`     //许可证申请码数据
	AvailableHostNum   int       `json:"availableHostNum"`   //可用X86 服务器数量
	AvailableCpuNum    int       `json:"availableCpuNum"`    //可用X86 CPU数量
	AvailableVmNum     int       `json:"availableVmNum"`     //可用X86 VM数量
	Source             string    `json:"source"`             // 来源
}

type LicenseAddOnInventoryView struct {
	LicenseInventoryView

	Name    string   `json:"name"`    // 许可证名称
	Modules []string `json:"modules"` // 模块信息
}
