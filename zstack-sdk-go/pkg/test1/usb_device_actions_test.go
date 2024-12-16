// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestZSClient_QueryUsbDevice(t *testing.T) {
	device, err := accountLoginCli.QueryUsbDevice(param.NewQueryParam())
	if err != nil {
		t.Errorf("TestZSClient_QueryUsbDevice %v", err)
		return
	}
	t.Log(device)
}

func TestGetUsbDevice(t *testing.T) {
	device, err := accountLoginCli.GetUsbDevice("3c9a51222e4f454ebfc39784d8c06485")
	if err != nil {
		t.Errorf("TestGetUsbDevice %v", err)
		return
	}
	t.Log(device)
}

func TestUpdateUsbDevice(t *testing.T) {
	name := "USB-002-003"
	device, err := accountLoginCli.UpdateUsbDevice("3c9a51222e4f454ebfc39784d8c06485", param.UpdateUsbDeviceParam{
		BaseParam: param.BaseParam{},
		UpdateUsbDevice: param.UpdateUsbDeviceDetailParam{
			Name: &name,
		},
	})
	if err != nil {
		t.Errorf("TestUpdateUsbDevice %v", err)
		return
	}
	t.Log(device)
}

func TestAttachUsbDeviceToVm(t *testing.T) {
	vm, err := accountLoginCli.AttachUsbDeviceToVm("3c9a51222e4f454ebfc39784d8c06485", param.AttachUsbDeviceToVmParam{
		BaseParam: param.BaseParam{},
		Params: param.AttachUsbDeviceToVmDetailParam{
			VmInstanceUuid: "e3eca0bf47354e9ab52fd23baf93d85a",
			AttachType:     "",
		},
	})
	if err != nil {
		t.Errorf("TestAttachUsbDeviceToVm %v", err)
		return
	}
	t.Log(vm)
}

func TestDetachUsbDeviceFromVm(t *testing.T) {
	err := accountLoginCli.DetachUsbDeviceFromVm("3c9a51222e4f454ebfc39784d8c06485", param.DetachUsbDeviceFromVmParam{
		BaseParam: param.BaseParam{},
	})
	if err != nil {
		t.Errorf("TestDetachUsbDeviceFromVm %v", err)
		return
	}
}
