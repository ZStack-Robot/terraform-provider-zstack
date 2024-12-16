// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestCreateEip(t *testing.T) {
	vip, err := accountLoginCli.CreateVip(param.CreateVipParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVipDetailParam{
			Name:              "v5",
			Description:       "v3",
			L3NetworkUUID:     "d8d0e446fedb42048b84d0eb3b01af66",
			IpRangeUUID:       "",
			AllocatorStrategy: "",
			RequiredIp:        "",
			ResourceUuid:      "",
		},
	})
	if err != nil {
		golog.Errorf("TestCreateVip error %v ", err)
	}
	golog.Println(vip)

	eip, err := accountLoginCli.CreateEip(param.CreateEipParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateEipDetailParam{
			Name:         "e3",
			Description:  "e3",
			VipUuid:      vip.UUID,
			VmNicUuid:    "",
			ResourceUuid: "",
		},
	})
	if err != nil {
		golog.Errorf("error %v ", err)
	}

	golog.Println(eip)
}

func TestDeleteEip(t *testing.T) {
	err := accountLoginCli.DeleteEip("df4569ca89f148f9b871f245dd17aa5b", param.DeleteModeEnforcing)
	if err != nil {
		golog.Errorf("error %v ", err)
	}
}

func TestQueryEip(t *testing.T) {
	eip, err := accountLoginCli.QueryEip(param.NewQueryParam())
	if err != nil {
		golog.Errorf("error %v ", err)
	}
	golog.Println(eip)
}

func TestGetEip(t *testing.T) {
	eip, err := accountLoginCli.GetEip("92af54c988a74260b3c024a8203e9989")
	if err != nil {
		golog.Errorf("error %v ", err)
	}
	golog.Println(eip)
}

func TestUpdateEip(t *testing.T) {
	eip, err := accountLoginCli.UpdateEip(param.UpdateEipParam{
		BaseParam: param.BaseParam{},
		UUID:      "51106acfe3ad4b67a97e9e3f02743593",
		UpdateEip: param.UpdateEipDetailParam{
			Name:        "f2",
			Description: "f2",
		},
	})
	if err != nil {
		golog.Errorf("error %v ", err)
	}
	golog.Println(eip)
}

func TestChangeEipState(t *testing.T) {
	eip, err := accountLoginCli.ChangeEipState(param.ChangeEipStateParam{
		BaseParam: param.BaseParam{},
		UUID:      "92af54c988a74260b3c024a8203e9989",
		ChangeEipState: param.ChangeEipStateDetailParam{
			StateEvent: param.StateEventEnable,
		},
	})
	if err != nil {
		golog.Errorf("error %v ", err)
	}
	golog.Println(eip)
}

func TestGetEipAttachableVmNics(t *testing.T) {
	nics, err := accountLoginCli.GetEipAttachableVmNics(param.GetEipAttachableVmNicsParam{
		BaseParam: param.BaseParam{},
		EipUuid:   "c60c19dd73f44215b6ff7587f5942dd2",
		//	VipUuid:   "66fee56a6b234d7192722ce4c8866d0d",
	})
	if err != nil {
		golog.Errorf("error %v ", err)
	}
	golog.Println(nics)
}

//0b306af4a07747de9ff47278dafa5157

func TestGetVmNicAttachableEips(t *testing.T) {
	eips, err := accountLoginCli.GetVmNicAttachableEips(param.GetVmNicAttachableEipsParam{
		BaseParam: param.BaseParam{},
		VmNicUuid: "92af54c988a74260b3c024a8203e9989",
		IpVersion: 4,
		Limit:     10,
		Start:     0,
	})
	if err != nil {
		golog.Errorf("error %v ", err)
	}
	golog.Println(eips)
}

func TestAttachEip(t *testing.T) {
	err := accountLoginCli.AttachEip("92af54c988a74260b3c024a8203e9989", "62ed6e2886ed498b8fdf3405e957b8d8")
	if err != nil {
		golog.Errorf("error %v ", err)
	}
}

func TestDetachEip(t *testing.T) {
	err := accountLoginCli.DetachEip("92af54c988a74260b3c024a8203e9989")
	if err != nil {
		golog.Errorf("error %v ", err)
	}
}
