// Copyright (c) ZStack.io, Inc.

package view

type ZoneView struct {
	BaseInfoView
	BaseTimeView

	State string `json:"state"`
	Type  string `json:"type"`
}
