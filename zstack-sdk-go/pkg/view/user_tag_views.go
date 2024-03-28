// Copyright (c) HashiCorp, Inc.

package view

type UserTagInventoryView struct {
	Uuid         string `json:"uuid"`
	ResourceType string `json:"resourceType"`
	ResourceUuid string `json:"resourceUuid"`
	Tag          string `json:"tag"`
	Type         string `json:"type"`
	CreateDate   string `json:"createDate"`
	LastOpDate   string `json:"lastOpDate"`
}
