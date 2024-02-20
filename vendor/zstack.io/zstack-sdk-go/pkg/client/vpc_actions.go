package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// QueryVpcRouter 查询VPC路由器
func (cli *ZSClient) QueryVpcRouter(params param.QueryParam) ([]view.VpcRouterVmInventoryView, error) {
	resp := make([]view.VpcRouterVmInventoryView, 0)
	return resp, cli.List("v1/vpc/virtual-routers", &params, &resp)
}

// GetVpcRouter 查询VPC路由器
func (cli *ZSClient) GetVpcRouter(uuid string) (view.VpcRouterVmInventoryView, error) {
	var resp view.VpcRouterVmInventoryView
	return resp, cli.Get("v1/vpc/virtual-routers", uuid, nil, &resp)
}
