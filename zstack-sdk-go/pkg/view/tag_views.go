// Copyright (c) ZStack.io, Inc.

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

type TagInventoryView struct {
	BaseInfoView
	BaseTimeView
	Color string `json:"color"`
	Type  string `json:"type"`
	// Note: For simple pattern tags, the value cannot be changed.
	// For withToken pattern tags, only the key values can be changed.
	// WithToken pattern format: name::{key1}::{key2}...::{keyN}
}

type TagInventory struct {
	Uuid         string `json:"uuid"`
	ResourceType string `json:"resourceType"`
	ResourceUuid string `json:"resourceUuid"`
	Tag          string `json:"tag"`
	Type         string `json:"type"`
	// Note: For simple pattern tags, the value cannot be changed.
	// For withToken pattern tags, only the key values can be changed.
	// WithToken pattern format: name::{key1}::{key2}...::{keyN}
	TagPatternUuid string         `json:"tagPatternUuid"`
	TagPattern     TagPatternView `json:"tagPattern"`
	CreateDate     string         `json:"createDate"`
	LastOpDate     string         `json:"lastOpDate"`
}

type TagPatternView struct {
	Uuid       string `json:"uuid"`       // Tag pattern UUID
	Name       string `json:"name"`       // Name of the tag pattern
	Value      string `json:"value"`      // Actual tag value (e.g. withToken::{xxx})
	Color      string `json:"color"`      // Color associated with the tag
	Type       string `json:"type"`       // Tag type: simple | withToken
	CreateDate string `json:"createDate"` // Creation timestamp
	LastOpDate string `json:"lastOpDate"` // Last update timestamp
}

type AttachTagToResourceResult struct {
	Success bool                            `json:"success"`
	Results []AttachTagToResourceResultItem `json:"results"`
}

type AttachTagToResourceResultItem struct {
	Inventory TagInventory `json:"inventory"`
	Success   bool         `json:"success"`
}
