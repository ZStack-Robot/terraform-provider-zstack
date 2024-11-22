// Copyright (c) ZStack.io, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// RequestConsoleAccess Request console access URL
func (cli *ZSClient) RequestConsoleAccess(params param.RequestConsoleAccessParam) (view.ConsoleInventoryView, error) {
	var resp view.ConsoleInventoryView
	return resp, cli.Post("v1/consoles", &params, &resp)
}
