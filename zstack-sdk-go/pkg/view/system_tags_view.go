// Copyright (c) ZStack.io, Inc.

package view

type SystemTagView struct {
	BaseInfoView
	BaseTimeView

	Inherent     bool   `json:"inherent"`     // Internal system tag
	ResourceUuid string `json:"resourceUuid"` // User-specified resource UUID, if specified, the system will not randomly assign a UUID to the resource
	ResourceType string `json:"resourceType"` // When creating a tag, the user must specify the resource type associated with the tag
	Tag          string `json:"tag"`          // Tag string
}
