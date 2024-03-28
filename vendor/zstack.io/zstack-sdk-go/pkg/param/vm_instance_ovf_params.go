// Copyright (c) HashiCorp, Inc.

package param

type ImageType string

const (
	Download  ImageType = "Download"
	Upload    ImageType = "Upload"
	ImageUuid ImageType = "ImageUuid"
)

type ParseOvfParam struct {
	BaseParam
	Params ParseOvfDetailParam `json:"params"`
}

type ParseOvfDetailParam struct {
	XmlBase64 string `json:"xmlBase64"` //Base64编码的OVF文件内容
}

type CreateVmInstanceFromOvfParam struct {
	BaseParam
	Params CreateVmInstanceFromOvfDetailParam `json:"params"`
}

type CreateVmInstanceFromOvfDetailParam struct {
	XmlBase64               string  `json:"xmlBase64"`               //资源名称
	JsonImageInfos          string  `json:"jsonImageInfos"`          //描述OVF中disk ID与镜像文件对应关系的JSON字符串
	BackupStorageUuid       string  `json:"backupStorageUuid"`       //用于存储上传镜像文件的镜像存储UUID
	JsonCreateVmParam       string  `json:"jsonCreateVmParam"`       //包含云主机创建参数的消息的JSON字符串
	DeleteImageAfterSuccess bool    `json:"deleteImageAfterSuccess"` //部署完成后删除镜像文件
	DeleteImageOnFail       bool    `json:"deleteImageOnFail"`       //部署失败后删除镜像文件
	ResourceUuid            *string `json:"resourceUuid"`            //资源UUID
}

type CreateVmFromOvfImageParam struct {
	OvfId string    `json:"ovfId"`
	Type  ImageType `json:"type"`
	Url   string    `json:"url"`
	Uuid  string    `json:"uuid"`
}
