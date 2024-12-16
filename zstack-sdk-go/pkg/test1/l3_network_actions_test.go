// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestZSClient_QueryL3Network(t *testing.T) {
	queryParam := param.NewQueryParam()
	//	queryParam.AddQ("category=Public")
	network, err := accountLoginCli.QueryL3Network(queryParam)
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", network, err)
	}
	golog.Infof("data: %v", network)
}

func TestZSClient_PageL3Network(t *testing.T) {
	queryParam := param.NewQueryParam()
	queryParam.Start(0).Limit(3)
	network, total, err := accountLoginCli.PageL3network(queryParam)
	if err != nil {
		golog.Errorf("TestZSClient_PageL3Network %v: %v", network, err)
	}
	golog.Infof("data: %v, total: %v", network, total)
}

func TestGetL3NetworkByName(t *testing.T) {
	queryParam := param.NewQueryParam()

	queryParam.AddQ("name=public")
	network, err := accountLoginCli.QueryL3Network(queryParam)
	if err != nil {
		golog.Errorf("TestGetL3NetworkByName %v: %v", network, err)
	}
	golog.Infof("data: %v", network)
}

func TestGetL3Network(t *testing.T) {
	network, err := accountLoginCli.GetL3Network("de7f26a7304d45aea9e9871a1ba7dbae")
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", network, err)
	}
	golog.Infof("data: %v", network)
}

func TestGetFreeIp(t *testing.T) {
	p := param.NewQueryParam()
	p.Limit(5)
	network, err := accountLoginCli.GetFreeIp("de7f26a7304d45aea9e9871a1ba7dbae", p)
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", network, err)
	}
	golog.Infof("data: %v", network)
}

func TestCheckIpAvailability(t *testing.T) {
	availability, err := accountLoginCli.CheckIpAvailability("001b44ad354a484d8fff6f41295904f8", "172.31.12.11")
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", availability, err)
	}
	golog.Infof("data: %v", availability)
}

func TestGetIpAddressCapacity(t *testing.T) {
	capacity, err := accountLoginCli.GetIpAddressCapacity("e3a869f887414e0599fc496eba47f454")
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", capacity, err)
	}
	golog.Infof("data: %v", capacity)
}

func TestUpdateL3Network(t *testing.T) {
	desc := "test"
	network, err := accountLoginCli.UpdateL3Network("088651c380af4d52b5791df02386b411", param.UpdateL3NetworkParam{
		BaseParam: param.BaseParam{},
		UpdateL3Network: param.UpdateL3NetworkDetailParam{
			Name:        "",
			Description: &desc,
		},
	})
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", network, err)
		return
	}
	golog.Infof("data: %v", network)
}

func TestDeleteL3Network(t *testing.T) {
	err := accountLoginCli.DeleteL3Network("088651c380af4d52b5791df02386b411", param.DeleteModePermissive)
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", nil, err)
		return
	}
}

func TestCreateL3Network(t *testing.T) {
	network, err := accountLoginCli.CreateL3Network(param.CreateL3NetworkParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateL3NetworkDetailParam{
			Name:          "l3-Test22",
			L2NetworkUuid: "0593971e63ea40a09b43f4f146819463",
			Category:      "Private",
			System:        false,
		},
	})
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", network, err)
		return
	}
	golog.Infof("data: %v", network)
}

func TestAddDnsToL3Network(t *testing.T) {
	err := accountLoginCli.AddDnsToL3Network("00025f3499ed43b998573dbe8225f142", param.AddDnsToL3NetworkParam{
		BaseParam: param.BaseParam{},
		Params: param.AddDnsToL3NetworkDetailParam{
			Dns: "5.5.5.6",
		},
	})
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", nil, err)
		return
	}

}

func TestAddIpRange(t *testing.T) {
	ipRange, err := accountLoginCli.AddIpRange("6a7c9dd9d6e449f992a59df8c102b3ba", param.AddIpRangeParam{
		BaseParam: param.BaseParam{},
		Params: param.AddIpRangeDetailParam{
			Name:        "192.168.1.100-92.168.1.200",
			StartIp:     "192.168.2.100",
			EndIp:       "192.168.2.200",
			Netmask:     "255.255.0.0",
			Gateway:     "192.168.100.1",
			IpRangeType: "",
		},
	})
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", ipRange, err)
		return
	}
	golog.Infof("data: %v", ipRange)
}

func TestAddIpv6Range(t *testing.T) {
	ipRange, err := accountLoginCli.AddIpv6Range("18f983771c2e4ffe9521036a346b75e4", param.AddIpv6RangeParam{
		BaseParam: param.BaseParam{},
		Params: param.AddIpv6RangeDetailParam{
			Name:        "fd00::2000:3000-fd00::2000:3fff",
			StartIp:     "fd00::2000:3000",
			EndIp:       "fd00::2000:3fff",
			Gateway:     "fd00::1",
			PrefixLen:   64,
			AddressMode: "Stateful-DHCP",
		},
	})
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", ipRange, err)
		return
	}
	golog.Infof("data: %v", ipRange)
}

func TestAddIpRangeByNetworkCidr(t *testing.T) {
	cidr, err := accountLoginCli.AddIpRangeByNetworkCidr("6a7c9dd9d6e449f992a59df8c102b3ba", param.AddIpRangeByNetworkCidrParam{
		BaseParam: param.BaseParam{},
		Params: param.AddIpRangeByNetworkCidrDetailParam{
			Name:        "192.168.12.0/16",
			NetworkCidr: "192.168.15.0/16",
			Gateway:     "192.168.15.1",
			IpRangeType: "",
		},
	})
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", cidr, err)
		return
	}
	golog.Infof("data: %v", cidr)
}

func TestGetL3NetworkDhcpIpAddress(t *testing.T) {
	ip, err := accountLoginCli.GetL3NetworkDhcpIpAddress("bfd4b869faa54481bc2a4d94c4fa2b5a")
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", ip, err)
		return
	}
	golog.Infof("data: %v", ip)
}

func TestGetL3NetworkMtu(t *testing.T) {
	mtu, err := accountLoginCli.GetL3NetworkMtu("bfd4b869faa54481bc2a4d94c4fa2b5a")
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", mtu, err)
		return
	}
	golog.Infof("data: %v", mtu)
}

func TestSetL3NetworkMtu(t *testing.T) {
	err := accountLoginCli.SetL3NetworkMtu("bfd4b869faa54481bc2a4d94c4fa2b5a", 1700)
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", nil, err)
		return
	}
}
func TestGetL3NetworkIpStatistics(t *testing.T) {
	ipStatistics, err := accountLoginCli.GetL3NetworkIpStatistic("809dd092f9364ca4adb35cb8ae7b67b4")
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", ipStatistics, err)
		return
	}
	golog.Infof("len : %d    data: %v", len(ipStatistics), ipStatistics)
}

func TestDeleteIpRange(t *testing.T) {
	err := accountLoginCli.DeleteIpRange("26e189e698a84f6b983ccc872e8283c0", param.DeleteModePermissive)
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", nil, err)
		return
	}
}

func TestRemoveDnsFromL3Network(t *testing.T) {
	err := accountLoginCli.RemoveDnsFromL3Network("e565b8bc378441c79ed19418e611c492", "8.8.8.8")
	if err != nil {
		return
	}
}

func TestQueryIpRange(t *testing.T) {
	provider, err := accountLoginCli.QueryIpRange(param.NewQueryParam())
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", provider, err)
		return
	}
	golog.Infof("data: %v", provider)
}

// f674ae2fae194373966bfb04259740ae
func TestGetIpRange(t *testing.T) {
	ipRange, err := accountLoginCli.GetIpRange("6a97e9e8a9e045a4af526644c2c796f5")
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", ipRange, err)
		return
	}
	golog.Infof("data: %v", ipRange)
}

func TestQueryIpAddress(t *testing.T) {
	queryParam := param.NewQueryParam()
	queryParam.AddQ("ipRangeUuid=a55eddc580ac460d927d68f894894096")
	queryParam.AddQ("vmNicUuid not null")
	address, err := accountLoginCli.QueryIpAddress(queryParam)
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", address, err)
		return
	}
	golog.Infof("data: %v", address)

}
