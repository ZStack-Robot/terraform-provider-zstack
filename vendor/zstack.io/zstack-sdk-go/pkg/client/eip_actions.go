// Copyright (c) HashiCorp, Inc.

package client

import (
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// CreateEip 创建弹性IP
func (cli *ZSClient) CreateEip(params param.CreateEipParam) (view.EipInventoryView, error) {
	var resp view.EipInventoryView
	return resp, cli.Post("v1/eips", params, &resp)
}

// DeleteEip 删除弹性IP
func (cli *ZSClient) DeleteEip(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/eips", uuid, string(deleteMode))
}

// PageEip 分页
func (cli *ZSClient) PageEip(params param.QueryParam) ([]view.EipInventoryView, int, error) {
	var eips []view.EipInventoryView
	total, err := cli.Page("v1/eips", &params, &eips)
	return eips, total, err
}

// QueryEip 查询弹性IP
func (cli *ZSClient) QueryEip(params param.QueryParam) ([]view.EipInventoryView, error) {
	resp := make([]view.EipInventoryView, 0)
	return resp, cli.List("v1/eips", &params, &resp)
}

// GetEip 查询弹性IP
func (cli *ZSClient) GetEip(uuid string) (view.EipInventoryView, error) {
	var resp view.EipInventoryView
	return resp, cli.Get("v1/eips", uuid, nil, &resp)
}

// UpdateEip 更新弹性IP
func (cli *ZSClient) UpdateEip(params param.UpdateEipParam) (view.EipInventoryView, error) {
	var resp view.EipInventoryView
	return resp, cli.Put("v1/eips", params.UUID, params, &resp)
}

// ChangeEipState 更改虚拟IP启用状态
func (cli *ZSClient) ChangeEipState(params param.ChangeEipStateParam) (view.EipInventoryView, error) {
	var resp view.EipInventoryView
	return resp, cli.Put("v1/eips", params.UUID, params, &resp)
}

// GetEipAttachableVmNics 获取可绑定指定弹性IP的云主机网卡
func (cli *ZSClient) GetEipAttachableVmNics(params param.GetEipAttachableVmNicsParam) ([]view.VmNicInventoryView, error) {
	resp := make([]view.VmNicInventoryView, 0)
	return resp, cli.GetWithSpec("v1/eips", params.EipUuid, "vm-instances/candidate-nics", responseKeyInventories, nil, &resp)
}

// GetVmNicAttachableEips 获取vmNic可绑定的EIp
func (cli *ZSClient) GetVmNicAttachableEips(params param.GetVmNicAttachableEipsParam) ([]view.EipInventoryView, error) {
	resp := make([]view.EipInventoryView, 0)
	return resp, cli.GetWithSpec("v1/vm-instances/nics", params.VmNicUuid, "candidate-eips", responseKeyInventories, nil, &resp)
}

// AttachEip 绑定弹性IP
func (cli *ZSClient) AttachEip(eipUuid, vmNicUuid string) error {
	return cli.PutWithSpec("v1/eips", eipUuid, fmt.Sprintf("vm-instances/nics/%s", vmNicUuid), "", map[string]string{}, nil)
}

// DetachEip 解绑弹性IP
func (cli *ZSClient) DetachEip(eipUuid string) error {
	return cli.Delete("v1/eips", fmt.Sprintf("%s/vm-instances/nics", eipUuid), "")
}
