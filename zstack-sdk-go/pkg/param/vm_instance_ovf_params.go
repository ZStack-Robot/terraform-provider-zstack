// Copyright (c) ZStack.io, Inc.

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
	XmlBase64 string `json:"xmlBase64"` // Base64 encoded OVF file content
}

type CreateVmInstanceFromOvfParam struct {
	BaseParam
	Params CreateVmInstanceFromOvfDetailParam `json:"params"`
}

type CreateVmInstanceFromOvfDetailParam struct {
	XmlBase64               string  `json:"xmlBase64"`               // Resource name
	JsonImageInfos          string  `json:"jsonImageInfos"`          // JSON string describing the relationship between disk ID and image file in OVF
	BackupStorageUuid       string  `json:"backupStorageUuid"`       // Backup storage UUID for storing uploaded image files
	JsonCreateVmParam       string  `json:"jsonCreateVmParam"`       // JSON string containing the message with VM creation parameters
	DeleteImageAfterSuccess bool    `json:"deleteImageAfterSuccess"` // Delete image file after successful deployment
	DeleteImageOnFail       bool    `json:"deleteImageOnFail"`       // Delete image file after deployment failure
	ResourceUuid            *string `json:"resourceUuid"`            // Resource UUID
}

type CreateVmFromOvfImageParam struct {
	OvfId string    `json:"ovfId"`
	Type  ImageType `json:"type"`
	Url   string    `json:"url"`
	Uuid  string    `json:"uuid"`
}
