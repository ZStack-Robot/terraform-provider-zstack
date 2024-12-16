// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"zstack.io/zstack-sdk-go/pkg/param"
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
