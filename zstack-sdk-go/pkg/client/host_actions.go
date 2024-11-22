// Copyright (c) ZStack.io, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// QueryHost Query physical machine
func (cli *ZSClient) QueryHost(params param.QueryParam) ([]view.HostInventoryView, error) {
	var resp []view.HostInventoryView
	return resp, cli.List("v1/hosts", &params, &resp)
}

// PageHost Paging of physical machines
func (cli *ZSClient) PageHost(params param.QueryParam) ([]view.HostInventoryView, int, error) {
	var resp []view.HostInventoryView
	total, err := cli.Page("v1/hosts", &params, &resp)
	return resp, total, err
}

// GetHost Physical machine details
func (cli *ZSClient) GetHost(uuid string) (*view.HostInventoryView, error) {
	resp := view.HostInventoryView{}
	if err := cli.Get("v1/hosts", uuid, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateHost Update physical machine
func (cli *ZSClient) UpdateHost(uuid string, params param.UpdateHostParam) (*view.HostInventoryView, error) {
	resp := view.HostInventoryView{}
	if err := cli.Put("v1/hosts", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ChangeHostState Update physical machine enable state
func (cli *ZSClient) ChangeHostState(uuid string, params *param.ChangeHostStateParam) (*view.HostInventoryView, error) {
	resp := view.HostInventoryView{}
	if err := cli.Put("v1/hosts", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ReconnectHost Reconnect physical machine
func (cli *ZSClient) ReconnectHost(uuid string) (*view.HostInventoryView, error) {
	resp := view.HostInventoryView{}
	if err := cli.Put("v1/hosts", uuid, map[string]struct{}{
		"reconnectHost": {},
	}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AddKVMHost Add KVM physical machine
func (cli *ZSClient) AddKVMHost(params param.AddKVMHostParam) (*view.HostInventoryView, error) {
	resp := view.HostInventoryView{}
	if err := cli.Post("v1/hosts/kvm", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteHost Delete physical machine
func (cli *ZSClient) DeleteHost(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/hosts", uuid, string(deleteMode))
}

// QueryHostNetworkBonding Query physical machine bonding information
func (cli *ZSClient) QueryHostNetworkBonding(params param.QueryParam) ([]view.HostNetworkBondingInventoryView, error) {
	var resp []view.HostNetworkBondingInventoryView
	return resp, cli.List("v1/hosts/bondings", &params, &resp)
}

// QueryHostNetworkInterface Query physical machine network interface information
func (cli *ZSClient) QueryHostNetworkInterface(params param.QueryParam) ([]view.HostNetworkInterfaceInventoryView, error) {
	var resp []view.HostNetworkInterfaceInventoryView
	return resp, cli.List("v1/hosts/nics", &params, &resp)
}
