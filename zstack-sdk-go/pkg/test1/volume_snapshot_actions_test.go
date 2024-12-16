// # Copyright (c) ZStack.io, Inc.

package test

import (
	"fmt"
	"testing"

	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
)

const accountLogin, accessKeyAuth = "accountLogin", "accessKeyAuth"

var (
	snapshotID = map[string]string{
		accountLogin:  "489645ffd7b1424796de50031500b7fe",
		accessKeyAuth: "a30e07b6091e48ecb1deb50a861e27d2",
	}
	treeID = map[string]string{
		accountLogin:  "6afe7f35ddc14b5ab740f7c257a6fc62",
		accessKeyAuth: "25914e7134f84a7e8ae77ca28ce9b6d0",
	}
)

func TestZSClient_CreateVolumeSnapshot(t *testing.T) {
	params := param.VolumeSnapshotParam{
		Params: param.VolumeSnapshotDetailParam{
			Name: "chenjh-test-snapshot",
		},
	}
	type args struct {
		volumeUuid string
		params     param.VolumeSnapshotParam
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{volumeID[accountLogin], params}, false},
		{accessKeyAuth, accessKeyAuthCli, args{volumeID[accessKeyAuth], params}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.CreateVolumeSnapshot(tt.args.volumeUuid, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.CreateVolumeSnapshot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			snapshotID[tt.name] = got.UUID
			treeID[tt.name] = got.TreeUUID
			t.Logf("ZSClient.CreateVolumeSnapshot() = %v, tree = %v", got.UUID, got.TreeUUID)
		})
	}
}

func TestZSClient_QueryVolumeSnapshot(t *testing.T) {
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
			got, err := tt.cli.QueryVolumeSnapshot(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.QueryVolumeSnapshot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.QueryVolumeSnapshot() = %v", got)
		})
	}
}

func TestZSClient_PageVolumeSnapshot(t *testing.T) {
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
			got, got1, err := tt.cli.PageVolumeSnapshot(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.PageVolumeSnapshot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.PageVolumeSnapshot() got = %v, total %v", got, got1)
		})
	}
}

func TestZSClient_GetVolumeSnapshot(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{snapshotID[accountLogin]}, false},
		{accessKeyAuth, accessKeyAuthCli, args{snapshotID[accessKeyAuth]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.GetVolumeSnapshot(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.GetVolumeSnapshot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.GetVolumeSnapshot() = %v", got)
		})
	}
}

func TestZSClient_QueryVolumeSnapshotTree(t *testing.T) {
	tests := []struct {
		name    string
		cli     *client.ZSClient
		wantErr bool
	}{
		{accountLogin, accountLoginCli, false},
		{accessKeyAuth, accessKeyAuthCli, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := param.NewQueryParam()
			params.AddQ(fmt.Sprintf("volumeUuid=%s", volumeID[tt.name]))
			got, err := tt.cli.QueryVolumeSnapshotTree(params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.QueryVolumeSnapshotTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.QueryVolumeSnapshotTree() = %v", got)
		})
	}
}

func TestZSClient_GetVolumeSnapshotTree(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{treeID[accountLogin]}, false},
		{accessKeyAuth, accessKeyAuthCli, args{treeID[accessKeyAuth]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.GetVolumeSnapshotTree(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.GetVolumeSnapshotTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.GetVolumeSnapshotTree() = %v", got)
		})
	}
}

func TestZSClient_UpdateVolumeSnapshot(t *testing.T) {
	params := param.UpdateVolumeSnapshotParam{
		UpdateVolumeSnapshot: param.UpdateVolumeSnapshotDetailParam{
			Description: "Zero",
		},
	}
	type args struct {
		uuid   string
		params param.UpdateVolumeSnapshotParam
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{snapshotID[accountLogin], params}, false},
		{accessKeyAuth, accessKeyAuthCli, args{snapshotID[accessKeyAuth], params}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.UpdateVolumeSnapshot(tt.args.uuid, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.UpdateVolumeSnapshot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.UpdateVolumeSnapshot() = %v", got)
		})
	}
}

func TestZSClient_RevertVolumeFromSnapshot(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{snapshotID[accountLogin]}, false},
		{accessKeyAuth, accessKeyAuthCli, args{snapshotID[accessKeyAuth]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cli.RevertVolumeFromSnapshot(tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.RevertVolumeFromSnapshot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZSClient_GetVolumeSnapshotSize(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{snapshotID[accountLogin]}, false},
		{accessKeyAuth, accessKeyAuthCli, args{snapshotID[accessKeyAuth]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.GetVolumeSnapshotSize(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.GetVolumeSnapshotSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.GetVolumeSnapshotSize() = %v", got)
		})
	}
}

func TestZSClient_ShrinkVolumeSnapshot(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{snapshotID[accountLogin]}, true},
		{accessKeyAuth, accessKeyAuthCli, args{snapshotID[accessKeyAuth]}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.ShrinkVolumeSnapshot(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.ShrinkVolumeSnapshot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// TODO
			t.Logf("ZSClient.ShrinkVolumeSnapshot() = %v", got)
		})
	}
}

func TestZSClient_DeleteVolumeSnapshot(t *testing.T) {
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
		{accountLogin, accountLoginCli, args{snapshotID[accountLogin], param.DeleteModePermissive}, false},
		{accessKeyAuth, accessKeyAuthCli, args{snapshotID[accessKeyAuth], param.DeleteModePermissive}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cli.DeleteVolumeSnapshot(tt.args.uuid, tt.args.deleteMode); (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.DeleteVolumeSnapshot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
