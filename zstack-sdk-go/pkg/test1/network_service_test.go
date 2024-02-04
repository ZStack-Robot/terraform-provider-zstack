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
	t.Log(provider)
}

func TestAttachNetworkServiceToL3Network(t *testing.T) {
	err := accountLoginCli.AttachNetworkServiceToL3Network("00025f3499ed43b998573dbe8225f142", param.AttachNetworkServiceToL3NetworkParam{
		BaseParam: param.BaseParam{},
		Params: param.AttachNetworkServiceToL3NetworkDetailParam{
			NetworkServices: map[string][]string{
				"590c129ef6dd451e914576d0aba74757": []string{
					"LoadBalancer",
				},
				"710a1f404ed5412595c0c4570cbde071": []string{"SecurityGroup"},
				"be4f28e3e9254526a4ad25617d6ccf59": []string{
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
