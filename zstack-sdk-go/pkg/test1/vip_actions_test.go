// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestCreateVip(t *testing.T) {
	vip, err := accountLoginCli.CreateVip(param.CreateVipParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVipDetailParam{
			Name:              "v1000",
			Description:       "v1000",
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
}

func TestDeleteVip(t *testing.T) {
	err := accountLoginCli.DeleteVip("7cc02dbe38cd49d09ae22a13470c011e", param.DeleteModePermissive)
	if err != nil {
		golog.Errorf("TestDeleteVip error %v ", err)
	}
}

func TestQueryVip(t *testing.T) {
	vip, err := accountLoginCli.QueryVip(param.NewQueryParam())
	if err != nil {
		golog.Errorf("TestQueryVip error %v ", err)
	}
	golog.Println(vip)
}

func TestGetVip(t *testing.T) {
	vip, err := accountLoginCli.GetVip("a341e00c910c496983452d1dd3ab5b2c")
	if err != nil {
		golog.Errorf("error %v ", err)
	}
	golog.Println(vip)
}
