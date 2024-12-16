// # Copyright (c) ZStack.io, Inc.

package test

import (
	"fmt"
	"testing"

	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestZSClient_QueryGlobalConfig(t *testing.T) {
	params := param.NewQueryParam()
	params.AddQ("category?=ha,vm")
	params.AddQ("name?=self.fencer.strategy,numa,vm.clock.track,vm.clock.sync.interval.in.seconds")
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
		{accessKeyAuth, accessKeyAuthCli, args{params}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.QueryGlobalConfig(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.QueryGlobalConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.QueryGlobalConfig() = %v", got)
		})
	}
}

func TestZSClient_QueryResourceConfig(t *testing.T) {
	params := param.NewQueryParam()
	params.AddQ("name?=numa,vm.clock.sync.interval.in.seconds,vm.clock.track")
	params.AddQ("category=vm")
	params.AddQ(fmt.Sprintf("resourceUuid=%s", "3cf85c1f312a4706a4718623d1b1b04e"))
	params.AddQ(fmt.Sprintf("resourceType=%s", param.ResourceTypeVmInstanceVO))

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
		{accessKeyAuth, accessKeyAuthCli, args{params}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.QueryResourceConfig(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.QueryResourceConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.QueryResourceConfig() = %v", got)
		})
	}
}

func TestGetResourceConfig(t *testing.T) {
	config, err := accountLoginCli.GetResourceConfig("8f6924be22b54a3db872b17bc7fa08c6", "vm", "vm.clock.track")
	if err != nil {
		t.Errorf("GetResourceConfig() error = %v", err)
		return
	}
	t.Log(config)
}

func TestUpdateResourceConfig(t *testing.T) {
	config, err := accountLoginCli.UpdateResourceConfig("vm", "vm.clock.track", "8f6924be22b54a3db872b17bc7fa08c6", param.UpdateResourceConfigParam{
		UpdateResourceConfig: param.UpdateConfigDetailParam{
			Value: "host",
		},
	})
	if err != nil {
		t.Errorf("UpdateResourceConfig() error = %v", err)
		return
	}
	t.Log(config)
}

func TestZSClient_UpdateGlobalConfig(t *testing.T) {
	type args struct {
		category string
		name     string
		params   param.UpdateGlobalConfigParam
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{name: accountLogin,
			cli:     accountLoginCli,
			args:    args{name: "self.fencer.strategy", category: "ha", params: param.UpdateGlobalConfigParam{UpdateGlobalConfig: param.UpdateConfigDetailParam{Value: string(client.Force)}}},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.UpdateGlobalConfig(tt.args.category, tt.args.name, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateGlobalConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}
