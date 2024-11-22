// Copyright (c) ZStack.io, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// PageUsbDevice Paginated query of USB devices
func (cli *ZSClient) PageUsbDevice(params param.QueryParam) ([]view.UsbDeviceView, int, error) {
	usbs := []view.UsbDeviceView{}
	total, err := cli.Page("v1/usb-device/usb-devices", &params, &usbs)
	return usbs, total, err
}

// QueryUsbDevice Query USB devices
func (cli *ZSClient) QueryUsbDevice(params param.QueryParam) ([]view.UsbDeviceView, error) {
	var usbs []view.UsbDeviceView
	return usbs, cli.List("v1/usb-device/usb-devices", &params, &usbs)
}

// GetUsbDevice Get a specific USB device
func (cli *ZSClient) GetUsbDevice(uuid string) (view.UsbDeviceView, error) {
	var resp view.UsbDeviceView
	return resp, cli.Get("v1/usb-device/usb-devices", uuid, nil, &resp)
}

// UpdateUsbDevice Update a USB device
func (cli *ZSClient) UpdateUsbDevice(uuid string, params param.UpdateUsbDeviceParam) (view.UsbDeviceView, error) {
	var resp view.UsbDeviceView
	return resp, cli.Put("v1/usb-device/usb-devices", uuid, &params, &resp)
}

// AttachUsbDeviceToVm Attach a physical USB device to a cloud VM
func (cli *ZSClient) AttachUsbDeviceToVm(usbDeviceUuid string, params param.AttachUsbDeviceToVmParam) (view.UsbDeviceView, error) {
	var resp view.UsbDeviceView
	return resp, cli.Post("v1/usb-device/usb-devices/"+usbDeviceUuid+"/attach", &params, &resp)
}

// DetachUsbDeviceFromVm Detach a USB device mounted on a cloud VM
func (cli *ZSClient) DetachUsbDeviceFromVm(usbDeviceUuid string, params param.DetachUsbDeviceFromVmParam) error {
	return cli.Post("v1/usb-device/usb-devices/"+usbDeviceUuid+"/detach", &params, nil)
}

// GetUsbDeviceCandidatesForAttachingVm Get the list of candidate USB devices for passthrough
func (cli *ZSClient) GetUsbDeviceCandidatesForAttachingVm(vmInstanceUuid string, attachType param.AttachType) ([]view.UsbDeviceView, error) {
	var usbs []view.UsbDeviceView
	url := ""
	if attachType != "" {
		url = string("?attachType=" + attachType)
	}
	return usbs, cli.Get("v1/vm-instances/", vmInstanceUuid+"/candidate-usb-devices"+url, nil, &usbs)
}
