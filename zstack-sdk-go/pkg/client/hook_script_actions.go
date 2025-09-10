// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// QueryVmUserDefinedXmlHookScript queries user custom xml hook scripts
func (cli *ZSClient) QueryVmUserDefinedXmlHookScript(params param.QueryParam) ([]view.HookScriptInventoryView, error) {
	var views []view.HookScriptInventoryView
	return views, cli.List("v1/vm-instances/xml-hook-script", &params, &views)
}

// GetVmUserDefinedXmlHookScript retrieves a specific user custom xml hook script by UUID
func (cli *ZSClient) GetVmUserDefinedXmlHookScript(uuid string) (*view.HookScriptInventoryView, error) {
	var resp view.HookScriptInventoryView
	if err := cli.Get("v1/vm-instances/xml-hook-script", uuid, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
