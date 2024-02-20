package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// CreateVip 创建虚拟IP
func (cli *ZSClient) CreateVip(params param.CreateVipParam) (view.VipInventoryView, error) {
	var resp view.VipInventoryView
	return resp, cli.Post("v1/vips", params, &resp)
}

// DeleteVip 删除虚拟IP
func (cli *ZSClient) DeleteVip(uuid string, mode param.DeleteMode) error {
	return cli.Delete("v1/vips", uuid, string(mode))
}

// QueryVip 查询虚拟IP
func (cli *ZSClient) QueryVip(params param.QueryParam) ([]view.VipInventoryView, error) {
	var resp []view.VipInventoryView
	return resp, cli.List("v1/vips", &params, &resp)
}

// GetVip 查询虚拟IP
func (cli *ZSClient) GetVip(uuid string) (view.VipInventoryView, error) {
	var resp view.VipInventoryView
	return resp, cli.Get("v1/vips", uuid, nil, &resp)
}

// UpdateVip 更新虚拟IP
func (cli *ZSClient) UpdateVip(params param.UpdateVipParam) (view.VipInventoryView, error) {
	var resp view.VipInventoryView
	return resp, cli.Put("v1/vips", params.UUID, params, &resp)
}
