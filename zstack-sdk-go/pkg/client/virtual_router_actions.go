// Copyright (c) ZStack.io, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// QueryVirtualRouterVm Query VPC Virtual Router
func (cli *ZSClient) QueryVirtualRouterVm(params param.QueryParam) ([]view.VirtualRouterInventoryView, error) {
	resp := make([]view.VirtualRouterInventoryView, 0)
	return resp, cli.List("v1/vm-instances/appliances/virtual-routers", &params, &resp)
}

// GetVirtualRouterVm Query a specific VPC Virtual Router
func (cli *ZSClient) GetVirtualRouterVm(uuid string) (view.VirtualRouterInventoryView, error) {
	resp := view.VirtualRouterInventoryView{}
	return resp, cli.Get("v1/vm-instances/appliances/virtual-routers", uuid, nil, &resp)
}

// CreateVirtualRouterVM Create a Virtual Router VM instance
func (cli *ZSClient) CreateVirtualRouterInstance(params param.CreateVirtualRouterInstanceParam) (*view.VirtualRouterInventoryView, error) {
	resp := view.VirtualRouterInventoryView{}
	if err := cli.Post("v1/vpc/virtual-routers", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
