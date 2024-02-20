package param

type ReloadLicenseParam struct {
	BaseParam

	ReloadLicense ReloadLicenseDetailParam `json:"reloadLicense"`
}

type ReloadLicenseDetailParam struct {
	ManagementNodeUuids []string `json:"managementNodeUuids"` //管理节点UUID
}

type UpdateLicenseParam struct {
	BaseParam

	UpdateLicense UpdateLicenseDetailParam `json:"updateLicense"`
}

type UpdateLicenseDetailParam struct {
	License         string `json:"license"`         //进行过base64 encode的license内容，传递集群所有管理节点许可证文件，则可全部管理节点都更新
	AdditionSession string `json:"additionSession"` //额外的信息,是一个json字符串,可选参数
}
