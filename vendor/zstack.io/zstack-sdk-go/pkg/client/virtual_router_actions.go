package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// QueryVirtualRouterVm 查询VPC路由器
func (cli *ZSClient) QueryVirtualRouterVm(params param.QueryParam) ([]view.VirtualRouterInventoryView, error) {
	resp := make([]view.VirtualRouterInventoryView, 0)
	return resp, cli.List("v1/vm-instances/appliances/virtual-routers", &params, &resp)
}

// GetVirtualRouterVm 查询VPC路由器
func (cli *ZSClient) GetVirtualRouterVm(uuid string) (view.VirtualRouterInventoryView, error) {
	resp := view.VirtualRouterInventoryView{}
	return resp, cli.Get("v1/vm-instances/appliances/virtual-routers", uuid, nil, &resp)
}
