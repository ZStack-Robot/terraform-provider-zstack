// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestQueryL2Network(t *testing.T) {
	network, err := accountLoginCli.QueryL2Network(param.NewQueryParam())
	if err != nil {
		golog.Errorf("TestQueryL2Network  %v", err)
		return
	}
	golog.Println(network)
}

func TestPageL2Network(t *testing.T) {
	queryParam := param.NewQueryParam()
	queryParam.Start(0).Limit(3)

	network, total, err := accountLoginCli.PageL2Network(queryParam)
	if err != nil {
		golog.Errorf("TestPageL2Network  %v", err)
		return
	}
	golog.Println(network, total)
}

func TestGetL2Network(t *testing.T) {
	network, err := accountLoginCli.GetL2Network("6de83607f46544e497c84c7eb085b498")
	if err != nil {
		golog.Errorf("TestGetL2Network  %v", err)
		return
	}
	golog.Println(network)
}

func TestUpdateL2Network(t *testing.T) {
	network, err := accountLoginCli.UpdateL2Network("8f065d55690a418aa7c37192c1b45b1e", param.UpdateL2NetworkParam{
		UpdateL2Network: param.UpdateL2NetworkDetailParam{
			Name:        "L2Network-1",
			Description: nil,
		},
	})
	if err != nil {
		golog.Errorf("TestUpdateL2Network  %v", err)
		return
	}
	golog.Println(network)
}

func TestDeleteL2Network(t *testing.T) {
	err := accountLoginCli.DeleteL2Network("4f3f73a3ebdc4f09b401f06eed583890", param.DeleteModePermissive)
	if err != nil {
		golog.Errorf("TestDeleteL2Network  %v", err)
		return
	}
}

func TestCreateL2NoVlanNetwork(t *testing.T) {
	network, err := accountLoginCli.CreateL2NoVlanNetwork(param.CreateL2NoVlanNetworkParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateL2NoVlanNetworkDetailParam{
			Name:              "L2Network-2",
			Description:       "sdffaddcsf",
			ZoneUuid:          "b37fa249602947d5960d51adb1b35f6c",
			PhysicalInterface: "eth2",
		},
	})
	if err != nil {
		golog.Errorf("TestCreateL2NoVlanNetwork  %v", err)
		return
	}
	golog.Println(network)
}

func TestCreateL2VlanNetwork(t *testing.T) {
	network, err := accountLoginCli.CreateL2VlanNetwork(param.CreateL2VlanNetworkParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateL2VlanNetworkDetailParam{
			Vlan:              4,
			Name:              "CreateL2VlanNetwork",
			Description:       "sdffaddcsfccc",
			ZoneUuid:          "b37fa249602947d5960d51adb1b35f6c",
			PhysicalInterface: "eth4",
		},
	})
	if err != nil {
		golog.Errorf("TestCreateL2VlanNetwork  %v", err)
		return
	}
	golog.Println(network)
}

func TestGetClusterHostNetworkFacts(t *testing.T) {
	facts, err := accountLoginCli.GetClusterHostNetworkFacts("9c2a62406143446dadd7ffbb300a424a")
	if err != nil {
		golog.Errorf("TestGetClusterHostNetworkFacts  %v", err)
		return
	}
	golog.Println(facts)
}

func TestAttachL2NetworkToCluster(t *testing.T) {
	cluster, err := accountLoginCli.AttachL2NetworkToCluster("9c2a62406143446dadd7ffbb300a424a", "18b2cf10a6a14634afdb11cd6e195f26")
	if err != nil {
		golog.Errorf("TestAttachL2NetworkToCluster  %v", err)
		return
	}
	golog.Println(cluster)
}
