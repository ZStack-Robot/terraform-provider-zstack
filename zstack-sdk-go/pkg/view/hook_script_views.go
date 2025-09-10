// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package view

type HookScriptInventoryView struct {
	BaseTimeView
	UUID       string `json:"uuid"` // Resource UUID, unique identifier for the resource
	Name       string `json:"name"` // Resource name
	Type       string `json:"type"`
	HookScript string `json:"hookScript"`
}
