// Copyright (c) ZStack.io, Inc.

package client

import (
	"fmt"
	"strings"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// QueryNetworkServiceProvider Query network service module
func (cli *ZSClient) QueryNetworkServiceProvider(params param.QueryParam) ([]view.NetworkServiceProviderInventoryView, error) {
	var resp []view.NetworkServiceProviderInventoryView
	return resp, cli.List("v1/network-services/providers", &params, &resp)
}

// AttachNetworkServiceToL3Network Attach network service to L3 network
func (cli *ZSClient) AttachNetworkServiceToL3Network(l3NetworkUuid string, p param.AttachNetworkServiceToL3NetworkParam) error {
	return cli.Post("v1/l3-networks/"+l3NetworkUuid+"/network-services", p, nil)
}

// QurySecurityGroup
func (cli *ZSClient) QuerySecurityGroup(params param.QueryParam) ([]view.SecurityGroupInventoryView, error) {
	var resp []view.SecurityGroupInventoryView
	return resp, cli.List("v1/security-groups", &params, &resp)
}

// GetSecurityGroup Get security group by UUID
func (cli *ZSClient) GetSecurityGroup(uuid string) ([]view.SecurityGroupInventoryView, error) {
	var resp []view.SecurityGroupInventoryView
	if err := cli.GetWithSpec("v1/security-groups", uuid, "", responseKeyInventories, nil, &resp); err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("security group with UUID %s not found", uuid)
	}
	return resp, nil
}

// AddVmNicToSecurityGroup Add VM NIC to security group  TODO
func (cli *ZSClient) AddVmNicToSecurityGroup(securityGroupUuid string, p param.AddVmNicToSecurityGroupParam) error {
	return cli.Post("v1/security-groups/"+securityGroupUuid+"/vm-instances/nics", p, nil)
}

// GetCandidateVmNicForSecurityGroup Get candidate VM NICs for security group
func (cli *ZSClient) GetCandidateVmNicForSecurityGroup(securityGroupUuid string) ([]view.VmNicInventoryView, error) {
	var resp []view.VmNicInventoryView
	if err := cli.GetWithSpec("v1/security-groups", securityGroupUuid, "/vm-instances/candidate-nics", responseKeyInventories, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteVmNicFromSecurityGroup Delete VM NIC from security group
func (cli *ZSClient) DeleteVmNicFromSecurityGroup(securityGroupUuid string, vmNicUuids []string) error {
	var uuidsStr []string
	for _, uuid := range vmNicUuids {
		uuidsStr = append(uuidsStr, fmt.Sprintf("vmNicUuids=%s", uuid))
	}
	uuidsQueryString := strings.Join(uuidsStr, "&")

	if err := cli.DeleteWithSpec("v1/security-groups", securityGroupUuid, "vm-instances/nics", uuidsQueryString, nil); err != nil {
		return err
	}
	return nil

}
