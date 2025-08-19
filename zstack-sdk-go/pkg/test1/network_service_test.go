// # Copyright (c) ZStack.io, Inc.

package test

import (
	"fmt"
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
	data, err := accountLoginCli.GetSecurityGroup("20a8f0ebe92840af8ad0383aacfa4022")
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

func TestCreateSecurityGroupTypeBridge(t *testing.T) {
	data, err := accountLoginCli.CreateSecurityGroup(param.CreateSecurityGroupParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateSecurityGroupDetailParam{
			Name:        "test-security-bridge",
			Description: "This is a test security group",
			IpVersion:   4,
			VSwitchType: "LinuxBridge",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	golog.Info(jsonutils.Marshal(data))
}

func TestCreateSecurityGroupTypeOvnDpdk(t *testing.T) {
	data, err := accountLoginCli.CreateSecurityGroup(param.CreateSecurityGroupParam{
		BaseParam: param.BaseParam{
			SystemTags: []string{"SdnControllerUuid::65589889039944b5a2efeb2ed4d67594"},
		},
		Params: param.CreateSecurityGroupDetailParam{
			Name:        "test-security-group-ovn-dpdk",
			Description: "This is a test security group",
			IpVersion:   4,
			VSwitchType: "OvnDpdk",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	golog.Info(jsonutils.Marshal(data))
}

func TestDeleteSecurityGroup(t *testing.T) {
	err := accountLoginCli.DeleteSecurityGroup("29625f3dc9614b2aba987b1473ef3cc6", param.DeleteModePermissive)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestChangeSecurityGroupState(t *testing.T) {
	state, err := accountLoginCli.ChangeSecurityGroupState(param.ChangeSecurityGroupStateParam{
		BaseParam:         param.BaseParam{},
		SecurityGroupUuid: "a6c0eec6940d48cea18363fe9321fca1",
		ChangeImageState: param.ChangeSecurityGroupStateDetailParam{
			StateEvent: param.StateEventDisable, // or StateEventEnable
		},
	})
	if err != nil {
		golog.Errorf("TestZSClient_UpdateSecurityGroupState error:%v", err)
	}
	fmt.Println(state)
}

func TestAddSecurityGroupRule(t *testing.T) {
	sgUuid := "2aa25c537ec94495b060850e2ffec762"

	params := param.AddSecurityGroupRuleParam{
		BaseParam: param.BaseParam{},
		Params: param.AddSecurityGroupRuleDetailParam{
			Rules: []param.AddSecurityGroupRule{
				{
					RuleType:     "Ingress",
					State:        "Enabled",
					Description:  "Allow HTTP traffic5",
					IpVersion:    4,
					Protocol:     "TCP",
					SrcIpRange:   "10.5.1.200-10.5.1.210",
					DstPortRange: "80,443,8080-9090",
					Action:       "ACCEPT",
				},
				{
					RuleType:     "Ingress",
					State:        "Enabled",
					Description:  "Allow outbound traffic5",
					IpVersion:    4,
					Protocol:     "TCP",
					SrcIpRange:   "13.13.15.19",
					DstPortRange: "80,443",
					Action:       "DROP",
				},
			},
			Priority: 2,
		},
	}

	resp, err := accountLoginCli.AddSecurityGroupRule(sgUuid, params)
	if err != nil {
		t.Error(err)
		return
	}
	golog.Info(jsonutils.Marshal(resp))
}

func TestGetSecurityGroupRules(t *testing.T) {
	ruleUuid := "a2b24c9acb9045ebbc4923a53b2d2b52"

	rules, err := accountLoginCli.GetSecurityGroupRule(ruleUuid)
	if err != nil {
		t.Error(err)
		return
	}
	golog.Info(jsonutils.Marshal(rules))
}

func TestQuerySecurityGroupRules(t *testing.T) {
	params := param.NewQueryParam()

	rules, err := accountLoginCli.QuerySecurityGroupRule(params)
	if err != nil {
		t.Error(err)
		return
	}
	golog.Info(jsonutils.Marshal(rules))
}

func TestDeleteSecurityGroupRule(t *testing.T) {
	ruleUuid := "cbdc3c4c9c1e4e2faa3607f53b6e4217"

	err := accountLoginCli.DeleteSecurityGroupRule(ruleUuid)
	if err != nil {
		t.Error(err)
		return
	}
	golog.Info("Security group rule deleted successfully")
}

func TestUpdateSecurityGroupRule(t *testing.T) {
	ruleUuid := "6ed5510bef2d43f688ac04c9c35647a6"

	params := param.UpdateSecurityGroupRuleParam{
		BaseParam: param.BaseParam{},
		ChangeSecurityGroupRule: param.UpdateSecurityGroupRuleDetailParam{
			Description:  "Updated rule description",
			State:        "Enabled", // or "Enabled"
			Priority:     8,
			Protocol:     "TCP",
			SrcIpRange:   "13.13.13.111",
			Action:       "DROP", // or "ACCEPT"
			DstPortRange: "80,443",
		},
	}
	resp, err := accountLoginCli.UpdateSecurityGroupRule(ruleUuid, params)
	if err != nil {
		t.Error(err)
		return
	}
	golog.Info(jsonutils.Marshal(resp))
}
