// Copyright (c) HashiCorp, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// QueryNetworkServiceProvider 查询网络服务模块
func (cli *ZSClient) QueryNetworkServiceProvider(params param.QueryParam) ([]view.NetworkServiceProviderInventoryView, error) {
	var resp []view.NetworkServiceProviderInventoryView
	return resp, cli.List("v1/network-services/providers", &params, &resp)
}

// AttachNetworkServiceToL3Network 挂载网络服务到三层网络
func (cli *ZSClient) AttachNetworkServiceToL3Network(l3NetworkUuid string, p param.AttachNetworkServiceToL3NetworkParam) error {
	return cli.Post("v1/l3-networks/"+l3NetworkUuid+"/network-services", p, nil)
}
