// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestQueryVirtualRouterVm(t *testing.T) {
	vm, err := accountLoginCli.QueryVirtualRouterVm(param.NewQueryParam())
	if err != nil {
		golog.Errorf("TestQueryVirtualRouterVm error %v", err)
		return
	}
	golog.Println(vm)
}

func TestGetVirtualRouterVm(t *testing.T) {
	vm, err := accountLoginCli.GetVirtualRouterVm("")
	if err != nil {
		golog.Errorf("TestGetVirtualRouterVm error %v", err)
		return
	}
	golog.Println(vm)
}

func TestCreateVirtualRouterOffering(t *testing.T) {
	vroffer, err := accountLoginCli.CreateVirtualRouterOffering(param.CreateVirtualRouterOfferingParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVirtualRouterOfferingDetailParam{
			Name:                  "vrOffering",
			Description:           "vrOffering",
			ZoneUuid:              "d29f4847a99f4dea83bc446c8fe6e64c",
			ManagementNetworkUuid: "50e8c0d69681447fbe347c8dae2b1bef",
			ImageUuid:             "93005c8a2a314a489635eca8c30794d4",
			//PublicNetworkUuid:     "",
			//IsDefault:  true,
			CpuNum:     1,
			MemorySize: 1073741824, // Byte 1GB = 1024MB = 1024*1024KB = 1024*1024*1024B = 1073741824Byte
			Type:       "VirtualRouter",
		},
	})
	if err != nil {
		t.Errorf("TestCreateVirtualRouterOffering %v", err)
	}
	golog.Println(vroffer)
}

func TestCreateVirtualRouterInstance(t *testing.T) {
	vrInstance, err := accountLoginCli.CreateVirtualRouterInstance(param.CreateVirtualRouterInstanceParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVirtualRouterInstanceDetailParam{
			Name:                      "vrInstance",
			Description:               "testing",
			VirtualRouterOfferingUuid: "43ad9df1194243bea1259fc276615a7f",
		},
	})
	if err != nil {
		t.Errorf("TestCreateVirtualRouterInstance %v", err)
	}
	golog.Println(vrInstance)
}
