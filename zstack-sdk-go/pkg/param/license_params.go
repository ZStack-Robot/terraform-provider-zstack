// Copyright (c) ZStack.io, Inc.

package param

type ReloadLicenseParam struct {
	BaseParam
	ReloadLicense ReloadLicenseDetailParam `json:"reloadLicense"`
}

type ReloadLicenseDetailParam struct {
	ManagementNodeUuids []string `json:"managementNodeUuids"` // Management node UUIDs
}

type UpdateLicenseParam struct {
	BaseParam
	UpdateLicense UpdateLicenseDetailParam `json:"updateLicense"`
}

type UpdateLicenseDetailParam struct {
	License         string `json:"license"`         // Base64 encoded license content, passing the license file for all management nodes will update all management nodes
	AdditionSession string `json:"additionSession"` // Additional information, a JSON string, optional parameter
}
