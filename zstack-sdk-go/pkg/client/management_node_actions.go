// Copyright (c) ZStack.io, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// QueryManagementNode Query management nodes
func (cli *ZSClient) QueryManagementNode(params param.QueryParam) ([]view.ManagementNodeInventoryView, error) {
	var resp []view.ManagementNodeInventoryView
	return resp, cli.List("v1/management-nodes", &params, &resp)
}

// GetVersion Retrieve the current version
func (cli *ZSClient) GetVersion() (string, error) {
	var resp string
	return resp, cli.PutWithRespKey("v1/management-nodes/actions", "", "version", map[string]struct{}{"getVersion": {}}, &resp)
}
