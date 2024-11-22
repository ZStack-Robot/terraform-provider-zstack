// Copyright (c) ZStack.io, Inc.

package client

import (
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// QueryL3Network Queries Layer 3 network
func (cli *ZSClient) QueryL3Network(params param.QueryParam) ([]view.L3NetworkInventoryView, error) {
	var network []view.L3NetworkInventoryView
	return network, cli.List("v1/l3-networks", &params, &network)
}

// PageL3network Paginated query for Layer 3 networks
func (cli *ZSClient) PageL3network(params param.QueryParam) ([]view.L3NetworkInventoryView, int, error) {
	var network []view.L3NetworkInventoryView
	total, err := cli.Page("v1/l3-networks", &params, &network)
	return network, total, err
}

// GetL3Network Queries a specific Layer 3 network
func (cli *ZSClient) GetL3Network(uuid string) (view.L3NetworkInventoryView, error) {
	var resp view.L3NetworkInventoryView
	return resp, cli.Get("v1/l3-networks", uuid, nil, &resp)
}

func (cli *ZSClient) CheckIpAvailability(l3NetworkUuid, ip string) (view.CheckIpView, error) {
	var resp view.CheckIpView

	return resp, cli.GetWithSpec("v1/l3-networks", fmt.Sprintf("%s/ip/%s/availability", l3NetworkUuid, ip), "", "", nil, &resp)
}

// GetFreeIp Queries free IPs
func (cli *ZSClient) GetFreeIp(l3NetworkUuid string, queryParam param.QueryParam) ([]view.FreeIpInventoryView, error) {
	var resp []view.FreeIpInventoryView
	return resp, cli.List(fmt.Sprintf("v1/l3-networks/%s/ip/free", l3NetworkUuid), &queryParam, &resp)
}

// GetIpAddressCapacity Retrieves IP address capacity for a given range
func (cli *ZSClient) GetIpAddressCapacity(ipRangeUuids string) (view.IpAddressCapacityView, error) {
	var resp view.IpAddressCapacityView
	queryParam := param.NewQueryParam()
	queryParam.Add("ipRangeUuids", ipRangeUuids)
	return resp, cli.GetWithSpec("v1/ip-capacity"+"?ipRangeUuids="+ipRangeUuids, "", "", "", nil, &resp)
}

// UpdateL3Network Updates a Layer 3 network
func (cli *ZSClient) UpdateL3Network(uuid string, params param.UpdateL3NetworkParam) (view.L3NetworkInventoryView, error) {
	var resp view.L3NetworkInventoryView
	return resp, cli.Put("v1/l3-networks", uuid, &params, &resp)
}

// DeleteL3Network Deletes a Layer 3 network
func (cli *ZSClient) DeleteL3Network(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/l3-networks", uuid, string(deleteMode))
}

// CreateL3Network Creates a Layer 3 network
func (cli *ZSClient) CreateL3Network(params param.CreateL3NetworkParam) (view.L3NetworkInventoryView, error) {
	var resp view.L3NetworkInventoryView
	return resp, cli.Post("v1/l3-networks", &params, &resp)
}

// AddDnsToL3Network Adds DNS to a Layer 3 network
func (cli *ZSClient) AddDnsToL3Network(l3NetworkUuid string, params param.AddDnsToL3NetworkParam) error {
	res := view.DnsInventoryView{}
	return cli.Post("v1/l3-networks/"+l3NetworkUuid+"/dns", &params, &res)
}

// AddIpRange Adds an IP address range
func (cli *ZSClient) AddIpRange(l3NetworkUuid string, params param.AddIpRangeParam) (view.IpRangeInventoryView, error) {
	var resp view.IpRangeInventoryView
	return resp, cli.Post("v1/l3-networks/"+l3NetworkUuid+"/ip-ranges", &params, &resp)
}

// AddIpv6Range Adds an IPv6 address range
func (cli *ZSClient) AddIpv6Range(l3NetworkUuid string, params param.AddIpv6RangeParam) (view.IpRangeInventoryView, error) {
	var resp view.IpRangeInventoryView
	return resp, cli.Post("v1/l3-networks/"+l3NetworkUuid+"/ipv6-ranges", &params, &resp)
}

// AddIpRangeByNetworkCidr Adds an IP address range via network CIDR
func (cli *ZSClient) AddIpRangeByNetworkCidr(l3NetworkUuid string, params param.AddIpRangeByNetworkCidrParam) (view.IpRangeInventoryView, error) {
	var resp view.IpRangeInventoryView
	return resp, cli.Post("v1/l3-networks/"+l3NetworkUuid+"/ip-ranges/by-cidr", &params, &resp)
}

func (cli *ZSClient) AddIpv6RangeByNetworkCidr(l3NetworkUuid string, params param.AddIpv6RangeByNetworkCidrParam) (view.IpRangeInventoryView, error) {
	var resp view.IpRangeInventoryView
	return resp, cli.Post("v1/l3-networks/"+l3NetworkUuid+"/ipv6-ranges/by-cidr", &params, &resp)
}

// GetL3NetworkDhcpIpAddress Retrieves the DHCP IP address for a network
func (cli *ZSClient) GetL3NetworkDhcpIpAddress(l3NetworkUuid string) (view.DhcpIpAddressView, error) {
	var resp view.DhcpIpAddressView
	return resp, cli.GetWithSpec("v1/l3-networks", l3NetworkUuid, "dhcp-ip", "", nil, &resp)
}

// GetL3NetworkMtu Retrieves the MTU value for a Layer 3 network
func (cli *ZSClient) GetL3NetworkMtu(l3NetworkUuid string) (view.MtuView, error) {
	var resp view.MtuView
	return resp, cli.GetWithSpec("v1/l3-networks", l3NetworkUuid, "mtu", "", nil, &resp)
}

// SetL3NetworkMtu Sets the MTU value for a Layer 3 network
func (cli *ZSClient) SetL3NetworkMtu(l3NetworkUuid string, mtu int64) error {
	return cli.Post("v1/l3-networks/"+l3NetworkUuid+"/mtu", &map[string]map[string]int64{
		"params": {
			"mtu": mtu,
		},
	}, nil)
}

// GetL3NetworkIpStatistic Retrieves IP address usage statistics for a Layer 3 network
func (cli *ZSClient) GetL3NetworkIpStatistic(l3NetworkUuid string) ([]view.IpStatisticView, error) {
	var resp []view.IpStatisticView
	return resp, cli.GetWithSpec("v1/l3-networks", l3NetworkUuid, "ip-statistic?limit=1000", "ipStatistics", nil, &resp)
}

// DeleteIpRange Deletes an IP address range
func (cli *ZSClient) DeleteIpRange(ipRangeUuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/l3-networks/ip-ranges", ipRangeUuid, string(deleteMode))
}

// RemoveDnsFromL3Network Removes DNS from a Layer 3 network
func (cli *ZSClient) RemoveDnsFromL3Network(l3NetworkUuid string, dns string) error {
	return cli.Delete("v1/l3-networks/"+l3NetworkUuid+"/dns", dns, "")
}

// QueryIpRange Queries IP address ranges
func (cli *ZSClient) QueryIpRange(queryParam param.QueryParam) ([]view.IpRangeInventoryView, error) {
	var resp []view.IpRangeInventoryView
	return resp, cli.List("v1/l3-networks/ip-ranges", &queryParam, &resp)
}

// GetIpRange Retrieves a specific IP address range
func (cli *ZSClient) GetIpRange(uuid string) (view.IpRangeInventoryView, error) {
	var resp view.IpRangeInventoryView
	return resp, cli.Get("v1/l3-networks/ip-ranges", uuid, nil, &resp)
}

// QueryIpAddress Queries IP addresses
func (cli *ZSClient) QueryIpAddress(queryParam param.QueryParam) ([]view.IpAddressInventoryView, error) {
	var resp []view.IpAddressInventoryView
	return resp, cli.List("v1/l3-networks/ip-address", &queryParam, &resp)
}

// AddReservedIpRange Adds a reserved IP address range
func (cli *ZSClient) AddReservedIpRange(l3NetworkUuid string, params param.AddReservedIpRangeParam) (view.ReservedIpRangeInventoryView, error) {
	var resp view.ReservedIpRangeInventoryView
	return resp, cli.Post("v1/l3-networks/"+l3NetworkUuid+"/reserved-ip-ranges", &params, &resp)
}

// DeleteReservedIpRange Deletes a reserved IP address range
func (cli *ZSClient) DeleteReservedIpRange(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/l3-networks/reserved-ip-ranges", uuid, string(deleteMode))
}
