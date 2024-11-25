// Copyright (c) ZStack.io, Inc.

package view

import "time"

type LicenseInventoryView struct {
	UUID               string    `json:"uuid"`               // Resource UUID, uniquely identifies the resource
	User               string    `json:"user"`               // License owner username
	ProdInfo           string    `json:"prodInfo"`           // License product name
	CpuNum             int       `json:"cpuNum"`             // Licensed X86 CPU count
	HostNum            int       `json:"hostNum"`            // Licensed X86 server count
	VmNum              int       `json:"vmNum"`              // Licensed X86 VM count
	LicenseType        string    `json:"licenseType"`        // License type
	ExpiredDate        time.Time `json:"expiredDate"`        // License expiration date
	IssuedDate         time.Time `json:"issuedDate"`         // License issue date
	UploadDate         time.Time `json:"uploadDate"`         // License upload date
	ManagementNodeUuid string    `json:"managementNodeUuid"` // License-associated MN node UUID
	Expired            bool      `json:"expired"`            // Whether the license has expired
	LicenseRequest     string    `json:"licenseRequest"`     // License request code data
	AvailableHostNum   int       `json:"availableHostNum"`   // Available X86 server count
	AvailableCpuNum    int       `json:"availableCpuNum"`    // Available X86 CPU count
	AvailableVmNum     int       `json:"availableVmNum"`     // Available X86 VM count
	Source             string    `json:"source"`             // Source
}

type LicenseAddOnInventoryView struct {
	LicenseInventoryView

	Name    string   `json:"name"`    // License name
	Modules []string `json:"modules"` // Module information
}
