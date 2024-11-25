// Copyright (c) ZStack.io, Inc.

package view

type UsbDeviceView struct {
	BaseInfoView
	BaseTimeView
	AttachType     string `json:"attachType"`
	BusNum         string `json:"busNum"`
	DevNum         string `json:"devNum"`
	HostUuid       string `json:"hostUuid"`
	IManufacturer  string `json:"iManufacturer"`
	IProduct       string `json:"iProduct"`
	ISerial        string `json:"iSerial"`
	IdProduct      string `json:"idProduct"`
	IdVendor       string `json:"idVendor"`
	State          string `json:"state"`
	UsbVersion     string `json:"usbVersion"`
	VmInstanceUuid string `json:"vmInstanceUuid"`
}
