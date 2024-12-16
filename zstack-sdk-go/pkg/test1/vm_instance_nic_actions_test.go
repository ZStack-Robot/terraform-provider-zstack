// # Copyright (c) ZStack.io, Inc.

package test

import (
	"fmt"
	"reflect"
	"testing"

	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

func TestZSClient_PageVmNic(t *testing.T) {
	type args struct {
		params param.QueryParam
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		want    []view.VmNicInventoryView
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.cli.PageVmNic(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.PageVmNic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZSClient.PageVmNic() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ZSClient.PageVmNic() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestZSClient_QueryVmNic(t *testing.T) {
	params := param.NewQueryParam()
	filterParams := param.NewQueryParam()
	filterParams.AddQ(fmt.Sprintf("vmInstanceUuid=%s", "3c430ff7da8440479920032a72f88ca0"))
	type args struct {
		params param.QueryParam
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{params}, false},
		{accountLogin + "Filter", accountLoginCli, args{filterParams}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.QueryVmNic(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.QueryVmNic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.QueryVmNic() = %v", got)
			t.Logf("======================================")
			for _, r := range got {
				t.Logf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s", r.UUID, r.VMInstanceUUID, r.L3NetworkUUID, r.IP, r.Mac, r.HypervisorType, r.Netmask, r.Gateway)
			}
			t.Logf("======================================")
		})
	}
}

func TestZSClient_GetVmNic(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		want    *view.VmNicInventoryView
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.GetVmNic(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.GetVmNic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZSClient.GetVmNic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZSClient_AttachL3NetworkToVm(t *testing.T) {
	type args struct {
		l3NetworkUuid  string
		vmInstanceUuid string
		params         param.AttachL3NetworkToVmParam
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{
			l3NetworkUuid:  "ff830e26aca7442c96672129ca3c54ff",
			vmInstanceUuid: "3c430ff7da8440479920032a72f88ca0",
			params: param.AttachL3NetworkToVmParam{
				Params: param.AttachL3NetworkToVmDetailParam{
					StaticIp: "192.168.1.209",
				},
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.AttachL3NetworkToVm(tt.args.l3NetworkUuid, tt.args.vmInstanceUuid, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.AttachL3NetworkToVm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.AttachL3NetworkToVm() = %v", got)
		})
	}
}

func TestZSClient_UpdateVmNicMac(t *testing.T) {
	type args struct {
		vmNicUuid string
		params    param.UpdateVmNicMacParam
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{
			vmNicUuid: "b906b4d3da0f416eab56df47aaa117a5",
			params: param.UpdateVmNicMacParam{
				UpdateVmNicMac: param.UpdateVmNicMacDetailParam{
					Mac: "FA:D2:E4:AA:7E:03",
				},
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.UpdateVmNicMac(tt.args.vmNicUuid, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.UpdateVmNicMac() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.UpdateVmNicMac() = %v", got)
		})
	}
}

func TestZSClient_SetVmStaticIp(t *testing.T) {
	type args struct {
		vmInstanceUuid string
		params         param.SetVmStaticIpParam
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{
			vmInstanceUuid: "3c430ff7da8440479920032a72f88ca0",
			params: param.SetVmStaticIpParam{
				SetVmStaticIp: param.SetVmStaticIpDetailParam{
					L3NetworkUuid: "ff830e26aca7442c96672129ca3c54ff",
					Ip:            "192.168.1.3",
				},
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cli.SetVmStaticIp(tt.args.vmInstanceUuid, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.SetVmStaticIp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.SetVmStaticIp()")
		})
	}
}

func TestZSClient_ChangeVmNicNetwork(t *testing.T) {
	type args struct {
		l3NetworkUuid string
		vmNicUuid     string
		params        param.ChangeVmNicNetworkParam
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{
			l3NetworkUuid: "06a2aeee77e1492cb8a22e9f5ed6c737",
			vmNicUuid:     "28e7de430f3d41b38891e6306ee6f4c0",
			params: param.ChangeVmNicNetworkParam{
				Params: param.ChangeVmNicNetworkDetailParam{
					DestL3NetworkUuid: "ff830e26aca7442c96672129ca3c54ff",
				},
				BaseParam: param.BaseParam{
					SystemTags: []string{"staticIp::ff830e26aca7442c96672129ca3c54ff::192.168.1.3"},
				},
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.ChangeVmNicNetwork(tt.args.l3NetworkUuid, tt.args.vmNicUuid, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.ChangeVmNicNetwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.ChangeVmNicNetwork() = %v", got)
		})
	}
}

func TestZSClient_DetachL3NetworkFromVm(t *testing.T) {
	type args struct {
		vmNicUuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{"b906b4d3da0f416eab56df47aaa117a5"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.DetachL3NetworkFromVm(tt.args.vmNicUuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.DetachL3NetworkFromVm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.DetachL3NetworkFromVm() = %v", got)
		})
	}
}
