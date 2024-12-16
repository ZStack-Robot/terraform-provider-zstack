// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
)

var (
	groupID = map[string]string{
		accountLogin:  "4232a5a7e5d2465f8ae6578b0dceb7b6",
		accessKeyAuth: "615d1911fb3247f395adcbd55e8b9ad4",
	}
)

func TestZSClient_CreateVolumeSnapshotGroup(t *testing.T) {
	params := param.VolumeSnapshotGroupParam{
		Params: param.VolumeSnapshotGroupDetailParam{
			RootVolumeUuid: rootVolumeID[accountLogin],
			Name:           "chenjh-test-snapshot-group",
		},
	}
	paramsAK := param.VolumeSnapshotGroupParam{
		Params: param.VolumeSnapshotGroupDetailParam{
			RootVolumeUuid: rootVolumeID[accessKeyAuth],
			Name:           "chenjh-test-snapshot-group",
		},
	}
	type args struct {
		params param.VolumeSnapshotGroupParam
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{params}, false},
		{accessKeyAuth, accessKeyAuthCli, args{paramsAK}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.CreateVolumeSnapshotGroup(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.CreateVolumeSnapshotGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			groupID[tt.name] = got.UUID
			t.Logf("ZSClient.CreateVolumeSnapshotGroup() = %v", groupID)
		})
	}
}

func TestZSClient_UpdateVolumeSnapshotGroup(t *testing.T) {
	params := param.UpdateVolumeSnapshotGroupParam{
		UpdateVolumeSnapshotGroup: param.UpdateVolumeSnapshotGroupDetailParam{
			Description: "Zero",
		},
	}
	type args struct {
		uuid   string
		params param.UpdateVolumeSnapshotGroupParam
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{groupID[accountLogin], params}, false},
		{accessKeyAuth, accessKeyAuthCli, args{groupID[accessKeyAuth], params}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.UpdateVolumeSnapshotGroup(tt.args.uuid, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.UpdateVolumeSnapshotGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.UpdateVolumeSnapshotGroup() = %v", got)
		})
	}
}

func TestZSClient_QueryVolumeSnapshotGroup(t *testing.T) {
	params := param.NewQueryParam()
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
			got, err := tt.cli.QueryVolumeSnapshotGroup(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.QueryVolumeSnapshotGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.QueryVolumeSnapshotGroup() = %v", got)
		})
	}
}

func TestZSClient_PageVolumeSnapshotGroup(t *testing.T) {
	params := param.NewQueryParam()
	params.Start(0).Limit(3)
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
			got, got1, err := tt.cli.PageVolumeSnapshotGroup(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.PageVolumeSnapshotGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.PageVolumeSnapshotGroup() got = %v, total %v", got, got1)
		})
	}
}

func TestZSClient_GetVolumeSnapshotGroup(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{groupID[accountLogin]}, false},
		{accessKeyAuth, accessKeyAuthCli, args{groupID[accessKeyAuth]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.GetVolumeSnapshotGroup(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.GetVolumeSnapshotGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.GetVolumeSnapshotGroup() = %v", got)
		})
	}
}
func TestZSClient_CheckVolumeSnapshotGroupAvailability(t *testing.T) {
	type args struct {
		uuids []string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{[]string{groupID[accountLogin]}}, false},
		{accessKeyAuth, accessKeyAuthCli, args{[]string{groupID[accessKeyAuth]}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.CheckVolumeSnapshotGroupAvailability(tt.args.uuids)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.CheckVolumeSnapshotGroupAvailability() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.CheckVolumeSnapshotGroupAvailability() = %v", got)
		})
	}
}

func TestZSClient_RevertVmFromSnapshotGroup(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{groupID[accountLogin]}, false},
		{accessKeyAuth, accessKeyAuthCli, args{groupID[accessKeyAuth]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cli.RevertVmFromSnapshotGroup(tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.RevertVmFromSnapshotGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZSClient_UngroupVolumeSnapshotGroup(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{groupID[accountLogin]}, false},
		{accessKeyAuth, accessKeyAuthCli, args{groupID[accessKeyAuth]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cli.UngroupVolumeSnapshotGroup(tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.UngroupVolumeSnapshotGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZSClient_DeleteVolumeSnapshotGroup(t *testing.T) {
	type args struct {
		uuid       string
		deleteMode param.DeleteMode
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{groupID[accountLogin], param.DeleteModePermissive}, false},
		{accessKeyAuth, accessKeyAuthCli, args{groupID[accessKeyAuth], param.DeleteModePermissive}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cli.DeleteVolumeSnapshotGroup(tt.args.uuid, tt.args.deleteMode); (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.DeleteVolumeSnapshotGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
