// Copyright (c) ZStack.io, Inc.

package client

import (
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

func (cli *ZSClient) PageVmNic(params param.QueryParam) ([]view.VmNicInventoryView, int, error) {
	var resp []view.VmNicInventoryView
	total, err := cli.Page("v1/vm-instances/nics", &params, &resp)
	return resp, total, err
}

func (cli *ZSClient) QueryVmNic(params param.QueryParam) ([]view.VmNicInventoryView, error) {
	var resp []view.VmNicInventoryView
	return resp, cli.List("v1/vm-instances/nics", &params, &resp)
}

func (cli *ZSClient) GetVmNic(uuid string) (*view.VmNicInventoryView, error) {
	var resp view.VmNicInventoryView
	return &resp, cli.Get("v1/vm-instances/nics", uuid, nil, &resp)
}

func (cli *ZSClient) AttachL3NetworkToVm(l3NetworkUuid, vmInstanceUuid string, params param.AttachL3NetworkToVmParam) (*view.VmInstanceInventoryView, error) {
	var resp view.VmInstanceInventoryView
	return &resp, cli.Post(fmt.Sprintf("v1/vm-instances/%s/l3-networks/%s", vmInstanceUuid, l3NetworkUuid), params, &resp)
}

func (cli *ZSClient) UpdateVmNicMac(vmNicUuid string, params param.UpdateVmNicMacParam) (*view.VmNicInventoryView, error) {
	var resp view.VmNicInventoryView
	return &resp, cli.Put("v1/vm-instances/nics", vmNicUuid, params, &resp)
}

func (cli *ZSClient) SetVmStaticIp(vmInstanceUuid string, params param.SetVmStaticIpParam) error {
	return cli.Put("v1/vm-instances", vmInstanceUuid, params, nil)
}

func (cli *ZSClient) DeleteVmStaticIp(vmInstanceUuid string, params param.DeleteVmStaticIpParam) error {
	paramsStr := fmt.Sprintf("l3NetworkUuid=%s", params.Params.L3NetworkUuid)
	if params.Params.DeleteMode != "" {
		paramsStr += fmt.Sprintf("&deleteMode=%s", params.Params.DeleteMode)
	}
	return cli.DeleteWithSpec("v1/vm-instances", vmInstanceUuid, "static-ips", paramsStr, nil)
}

func (cli *ZSClient) ChangeVmNicNetwork(l3NetworkUuid, vmNicUuid string, params param.ChangeVmNicNetworkParam) (*view.VmNicInventoryView, error) {
	var resp view.VmNicInventoryView
	return &resp, cli.Post(fmt.Sprintf("v1/vm-instances/nics/%s/l3-networks/%s", vmNicUuid, l3NetworkUuid), params, &resp)
}

func (cli *ZSClient) DetachL3NetworkFromVm(vmNicUuid string) (*view.VmInstanceInventoryView, error) {
	var resp view.VmInstanceInventoryView
	return &resp, cli.DeleteWithSpec("v1/vm-instances/nics", vmNicUuid, "", "", &resp)
}
