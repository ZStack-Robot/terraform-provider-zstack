// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

func (cli *ZSClient) CreateVmInstanceScript(params param.CreateVmInstanceScriptParam) (*view.VmInstanceScriptInventoryView, error) {
	resp := view.VmInstanceScriptInventoryView{}
	if err := cli.Post("v1/scripts", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (cli *ZSClient) DeleteVmInstanceScrpt(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/scripts", uuid, string(deleteMode))
}

func (cli *ZSClient) GetVmInstanceScript(uuid string) (*view.VmInstanceScriptInventoryView, error) {
	var resp view.VmInstanceScriptInventoryView
	if err := cli.Get("v1/scripts", uuid, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (cli *ZSClient) QueryVmInstanceScript(params param.QueryParam) ([]view.VmInstanceScriptInventoryView, error) {
	var resp []view.VmInstanceScriptInventoryView
	return resp, cli.List("v1/scripts", &params, &resp)
}

func (cli *ZSClient) UpdateVmInstanceScript(uuid string, params param.UpdateVmInstanceScriptParam) error {
	return cli.Put("v1/scripts", uuid, params, nil)
}

func (cli *ZSClient) ExecuteVmInstanceScript(uuid string, params param.ExecuteVmInstanceScriptParam) (*view.VmInstanceScriptResultInventoryView, error) {
	var resp view.VmInstanceScriptResultInventoryView
	if err := cli.Put("v1/scripts", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (cli *ZSClient) GetVmInstanceScriptExecutedRecord(uuid string) (*view.VmInstanceScriptResultInventoryView, error) {
	var resp view.VmInstanceScriptResultInventoryView
	if err := cli.Get("v1/scripts/records", uuid, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (cli *ZSClient) QueryVmInstanceScriptExecutedRecord(params param.QueryParam) ([]view.VmInstanceScriptResultInventoryView, error) {
	var resp []view.VmInstanceScriptResultInventoryView
	return resp, cli.List("v1/scripts/records", &params, &resp)
}

func (cli *ZSClient) QueryGuestVmScriptExecutedRecordDetail(params param.QueryParam) ([]view.VmInstanceScriptResultDetailInventoryView, error) {
	var resp []view.VmInstanceScriptResultDetailInventoryView
	return resp, cli.List("v1/scripts/records/details", &params, &resp)
}
