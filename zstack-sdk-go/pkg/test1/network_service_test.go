// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/util/jsonutils"
)

func TestQueryNetworkServiceProvider(t *testing.T) {
	provider, err := accountLoginCli.QueryNetworkServiceProvider(param.NewQueryParam())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("begin------------------------")
	t.Log(provider[0].Name)
	t.Log(provider[0].Uuid)
	t.Log(provider[0].Type)
	t.Log(provider[0].NetworkServiceTypes)
	t.Log(provider[0].AttachedL2NetworkUuids)
	t.Log("------------------------")
	t.Log(provider[1].Name)
	t.Log(provider[1].Uuid)
	t.Log(provider[1].NetworkServiceTypes)
	t.Log(provider[1].AttachedL2NetworkUuids)
	t.Log("------------------------")
	t.Log(provider[2].Name)
	t.Log(provider[2].Uuid)
	t.Log(provider[2].NetworkServiceTypes)
	t.Log(provider[2].AttachedL2NetworkUuids)
}

func TestAttachNetworkServiceToL3Network(t *testing.T) {
	err := accountLoginCli.AttachNetworkServiceToL3Network("00025f3499ed43b998573dbe8225f142", param.AttachNetworkServiceToL3NetworkParam{
		BaseParam: param.BaseParam{},
		Params: param.AttachNetworkServiceToL3NetworkDetailParam{
			NetworkServices: map[string][]string{
				"590c129ef6dd451e914576d0aba74757": {"LoadBalancer"},  // Simplified
				"710a1f404ed5412595c0c4570cbde071": {"SecurityGroup"}, // Simplified
				"22de5a6792bf4de1835b92125ac3c419": { // Simplified
					"VipQos",
					"DNS",
					"HostRoute",
					"Userdata",
					"Eip",
					"DHCP",
				},
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestQuerySecurityGroup(t *testing.T) {
	data, err := accountLoginCli.QuerySecurityGroup(param.NewQueryParam())
	if err != nil {
		t.Error(err)
		return
	}

	golog.Info(jsonutils.Marshal(data))

}

func TestGetSecurityGroup(t *testing.T) {
	data, err := accountLoginCli.GetSecurityGroup("f450b20497c34397977091bc1c8f87f9")
	if err != nil {
		t.Error(err)
		return
	}
	golog.Info(jsonutils.Marshal(data))
}

func TestGetCandidateVmNicForSecurityGroup(t *testing.T) {
	data, err := accountLoginCli.GetCandidateVmNicForSecurityGroup("f450b20497c34397977091bc1c8f87f9")
	if err != nil {
		t.Error(err)
		return
	}
	golog.Info(jsonutils.Marshal(data))
}

func TestAddVmNicToSecurityGroup(t *testing.T) {
	err := accountLoginCli.AddVmNicToSecurityGroup("f450b20497c34397977091bc1c8f87f9", param.AddVmNicToSecurityGroupParam{
		BaseParam: param.BaseParam{},
		Params: param.AddVmNicToSecurityGroupDetailParam{
			VmNicUuids: []string{"20ff9a2ba9ca4209a361c1ee52ff1b0f", "a8aa88c413704717b138190832864b54"},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDeleteVmNicFromSecurityGroup(t *testing.T) {
	err := accountLoginCli.DeleteVmNicFromSecurityGroup("f450b20497c34397977091bc1c8f87f9", []string{"20ff9a2ba9ca4209a361c1ee52ff1b0f", "a8aa88c413704717b138190832864b54"})
	if err != nil {
		t.Error(err)
		return
	}
}
