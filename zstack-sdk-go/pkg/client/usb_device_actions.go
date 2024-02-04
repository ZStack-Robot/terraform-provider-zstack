package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// PageUsbDevice 分页查询USB设备
func (cli *ZSClient) PageUsbDevice(params param.QueryParam) ([]view.UsbDeviceView, int, error) {
	usbs := []view.UsbDeviceView{}
	total, err := cli.Page("v1/usb-device/usb-devices", &params, &usbs)
	return usbs, total, err
}

// QueryUsbDevice 查询USB设备
func (cli *ZSClient) QueryUsbDevice(params param.QueryParam) ([]view.UsbDeviceView, error) {
	var usbs []view.UsbDeviceView
	return usbs, cli.List("v1/usb-device/usb-devices", &params, &usbs)
}

// GetUsbDevice 获取USB设备
func (cli *ZSClient) GetUsbDevice(uuid string) (view.UsbDeviceView, error) {
	var resp view.UsbDeviceView
	return resp, cli.Get("v1/usb-device/usb-devices", uuid, nil, &resp)
}

// UpdateUsbDevice 更新USB设备
func (cli *ZSClient) UpdateUsbDevice(uuid string, params param.UpdateUsbDeviceParam) (view.UsbDeviceView, error) {
	var resp view.UsbDeviceView
	return resp, cli.Put("v1/usb-device/usb-devices", uuid, &params, &resp)
}

// AttachUsbDeviceToVm 云主机加载物理机USB设备
func (cli *ZSClient) AttachUsbDeviceToVm(usbDeviceUuid string, params param.AttachUsbDeviceToVmParam) (view.UsbDeviceView, error) {
	var resp view.UsbDeviceView
	return resp, cli.Post("v1/usb-device/usb-devices/"+usbDeviceUuid+"/attach", &params, &resp)
}

// DetachUsbDeviceFromVm 将云主机挂载的USB设备卸载
func (cli *ZSClient) DetachUsbDeviceFromVm(usbDeviceUuid string, params param.DetachUsbDeviceFromVmParam) error {
	return cli.Post("v1/usb-device/usb-devices/"+usbDeviceUuid+"/detach", &params, nil)
}

// GetUsbDeviceCandidatesForAttachingVm 获取USB透传候选列表
func (cli *ZSClient) GetUsbDeviceCandidatesForAttachingVm(vmInstanceUuid string, attachType param.AttachType) ([]view.UsbDeviceView, error) {
	var usbs []view.UsbDeviceView
	url := ""
	if attachType != "" {
		url = string("?attachType=" + attachType)
	}
	return usbs, cli.Get("v1/vm-instances/", vmInstanceUuid+"/candidate-usb-devices"+url, nil, &usbs)
}
