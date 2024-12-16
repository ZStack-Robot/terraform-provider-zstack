// # Copyright (c) ZStack.io, Inc.

package test

import (
	"fmt"
	"testing"

	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
)

var (
	volumeID = map[string]string{
		accountLogin:  "6ca3f6e0b7af45e9bc6ad301b0e72042",
		accessKeyAuth: "67d2c8fb2bac4736a49d102c1d725248",
	}
	imageID = map[string]string{
		accountLogin:  "968e87334a12422fbe78c8b72bcfab68",
		accessKeyAuth: "8d88bf390a3543efb11dfda6afebc655",
	}
	hostID = map[string]string{
		accountLogin:  "43a562cb71744784b41d5d3663eb620f",
		accessKeyAuth: "b0de6e34be6042faa34069babcb64878",
	}
	primaryStorageID = map[string]string{
		accountLogin:  "ace08e7a30c14609b5a92e5114f19e82",
		accessKeyAuth: "dd2ae6841a054ce2b582545db9e7f787",
	}
	vmID = map[string]string{
		accountLogin:  "22f6836626bb4683b3d5ccf5bd9e0ae0",
		accessKeyAuth: "69f1c9d494414042860d355d386d91ba",
	}
	rootVolumeID = map[string]string{
		accountLogin:  "0eb9776b41184a108f53b4fd9b11acfa",
		accessKeyAuth: "aba61ccd0ea5426188cc05a61ffe1581",
	}
)

func TestZSClient_CreateDataVolume(t *testing.T) {
	params := param.CreateDataVolumeParam{
		BaseParam: param.BaseParam{
			UserTags: []string{"userID=10", "parentID=1"},
		},
		Params: param.CreateDataVolumeDetailParam{
			Name:        "chenjh-DATA-TEST2",
			Description: "JUST a test Volume For chenjh",
			DiskSize:    10240,
		},
	}
	type args struct {
		params param.CreateDataVolumeParam
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
			r, err := tt.cli.CreateDataVolume(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.CreateDataVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			volumeID[tt.name] = r.UUID
			t.Logf("ZSClient.CreateDataVolume() = %v", r)
			t.Logf("======================================")
			t.Logf("%s\t%s\t%s\t%s\t%s\t%d\t%d\t%s", r.UUID, r.Name, r.Type, r.State, r.Status, r.ActualSize, r.Size, r.Description)
			t.Logf("======================================")
		})
	}
}

func TestZSClient_ChangeVolumeState(t *testing.T) {
	params := param.ChangeVolumeStateParam{
		ChangeVolumeState: param.ChangeVolumeStateDetailParam{
			StateEvent: param.VolumeStateEnable,
		},
	}
	type args struct {
		uuid   string
		params param.ChangeVolumeStateParam
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
			r, err := tt.cli.ChangeVolumeState(tt.args.uuid, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.ChangeVolumeState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.ChangeVolumeState() = %v", r)
			t.Logf("======================================")
			t.Logf("%s\t%s\t%s\t%s\t%s\t%d\t%d\t%s", r.UUID, r.Name, r.Type, r.State, r.Status, r.ActualSize, r.Size, r.Description)
			t.Logf("======================================")
		})
	}
}

func TestZSClient_CreateDataVolumeFromVolumeTemplate(t *testing.T) {
	params := param.CreateDataVolumeFromVolumeTemplateParam{
		Params: param.CreateDataVolumeFromVolumeTemplateDetailParam{
			Name:               "chenjh-DATA-Image-TEST",
			Description:        "JUST a test Volume For chenjh",
			PrimaryStorageUuid: primaryStorageID[accountLogin],
			HostUuid:           hostID[accountLogin],
		},
	}
	type args struct {
		imageUuid string
		params    param.CreateDataVolumeFromVolumeTemplateParam
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{imageID[accountLogin], params}, false},
		{accessKeyAuth, accessKeyAuthCli, args{imageID[accessKeyAuth], params}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.CreateDataVolumeFromVolumeTemplate(tt.args.imageUuid, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.CreateDataVolumeFromVolumeTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.CreateDataVolumeFromVolumeTemplate() = %v", got)
		})
	}
}

func TestZSClient_CreateDataVolumeFromVolumeSnapshot(t *testing.T) {
	params := param.CreateDataVolumeFromVolumeSnapshotParam{
		Params: param.CreateDataVolumeFromVolumeSnapshotDetailParam{
			Name:               "chenjh-DATA-Snapshot-TEST",
			Description:        "JUST a test Volume For chenjh",
			VolumeSnapshotUuid: snapshotID[accountLogin],
		},
	}
	type args struct {
		volumeSnapshotUuid string
		params             param.CreateDataVolumeFromVolumeSnapshotParam
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
			got, err := tt.cli.CreateDataVolumeFromVolumeSnapshot(tt.args.volumeSnapshotUuid, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.CreateDataVolumeFromVolumeSnapshot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.CreateDataVolumeFromVolumeSnapshot() = %v", got)
		})
	}
}

func TestZSClient_QueryVolume(t *testing.T) {
	// params := param.NewQueryParam()
	filterParams := param.NewQueryParam()
	// filterParams.AddQ(fmt.Sprintf("type=%s", "Data"))
	// filterParams.AddQ(fmt.Sprintf("status=%s", "Ready"))
	filterParams.AddQ(fmt.Sprintf("vmInstanceUuid=%s", "3c430ff7da8440479920032a72f88ca0"))
	lastParams := param.NewQueryParam()
	lastParams.AddQ(fmt.Sprintf("lastVmInstanceUuid=%s", "3c430ff7da8440479920032a72f88ca0"))
	lastParams.AddQ("vmInstanceUuid is null")
	lastParams.AddQ("status!=Deleted")
	type args struct {
		params param.QueryParam
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		// {accountLogin, accountLoginCli, args{params}, false},
		// {accountLogin + "Filter", accountLoginCli, args{filterParams}, false},
		// // {accessKeyAuth, accessKeyAuthCli, args{params}, false},
		// // {accessKeyAuth + "Filter", accessKeyAuthCli, args{filterParams}, false},
		{accountLogin + "Last", accountLoginCli, args{lastParams}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.QueryVolume(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.QueryVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.QueryVolume() = %v", got)
			t.Logf("======================================")
			for _, r := range got {
				t.Logf("%s\t%s\t%s\t%s\t%s\t%d\t%d\t%s", r.UUID, r.Name, r.Type, r.State, r.Status, r.ActualSize, r.Size, r.Description)
			}
			t.Logf("======================================")
		})
	}
}

func TestQueryVolume(t *testing.T) {
	p := param.NewQueryParam()
	p.AddQ("hostUuid=4f8a562c919041e6980dc9f523ee0e8e")
	volume, err := accountLoginCli.QueryVolume(p)
	if err != nil {
		t.Errorf("ZSClient.QueryVolume() error = %v", err)
		return
	}
	t.Logf("ZSClient.QueryVolume() = %v", volume)
}

func TestZSClient_PageVolume(t *testing.T) {
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
			got, got1, err := tt.cli.PageVolume(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.PageVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.PageVolume() got = %v, total %v", got, got1)
		})
	}
}

func TestZSClient_GetVolume(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{"3f5c5267f2b04e6895238d402c7be42a"}, false},
		{accessKeyAuth, accessKeyAuthCli, args{volumeID[accessKeyAuth]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := tt.cli.GetVolume(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.GetVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.GetVolume() = %v", r)
			t.Logf("======================================")
			t.Logf("%s\t%s\t%s\t%s\t%s\t%d\t%d\t%s", r.UUID, r.Name, r.Type, r.State, r.Status, r.ActualSize, r.Size, r.Description)
			t.Logf("======================================")
		})
	}
}

func TestZSClient_GetVolumeFormat(t *testing.T) {
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
			got, err := tt.cli.GetVolumeFormat()
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.GetVolumeFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.GetVolumeFormat() = %v", got)
		})
	}
}

func TestZSClient_GetVolumeCapabilities(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{volumeID[accountLogin]}, false},
		{accessKeyAuth, accessKeyAuthCli, args{volumeID[accessKeyAuth]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := tt.cli.GetVolumeCapabilities(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.GetVolumeCapabilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.GetVolumeCapabilities() = %v", r)
			t.Logf("======================================")
			t.Logf("%t\t%t", r.MigrationInCurrentPrimaryStorage, r.MigrationToOtherPrimaryStorage)
			t.Logf("======================================")
		})
	}
}

func TestZSClient_SyncVolumeSize(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{volumeID[accountLogin]}, false},
		{accessKeyAuth, accessKeyAuthCli, args{volumeID[accessKeyAuth]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.SyncVolumeSize(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.SyncVolumeSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.SyncVolumeSize() = %v", got)
		})
	}
}

func TestZSClient_ResizeRootVolume(t *testing.T) {
	type args struct {
		uuid string
		size int64
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{rootVolumeID[accountLogin], 0}, true},
		{accessKeyAuth, accessKeyAuthCli, args{rootVolumeID[accessKeyAuth], 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.ResizeRootVolume(tt.args.uuid, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.ResizeRootVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.ResizeRootVolume() = %v", got)
		})
	}
}

func TestZSClient_ResizeDataVolume(t *testing.T) {
	type args struct {
		uuid string
		size int64
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{volumeID[accountLogin], 2}, true},
		{accessKeyAuth, accessKeyAuthCli, args{volumeID[accessKeyAuth], 2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.ResizeDataVolume(tt.args.uuid, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.ResizeDataVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.ResizeDataVolume() = %v", got)
		})
	}
}

func TestZSClient_UpdateVolume(t *testing.T) {
	desc := "Just a Data Volume for go sdk[sdk:go-1.0.0]"
	params := param.UpdateVolumeParam{
		UpdateVolume: param.UpdateVolumeDetailParam{
			Name:        "Data-for-sdk",
			Description: &desc,
		},
	}
	type args struct {
		uuid   string
		params param.UpdateVolumeParam
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
			r, err := tt.cli.UpdateVolume(tt.args.uuid, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.UpdateVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.UpdateVolume() = %v", r)
			t.Logf("======================================")
			t.Logf("%s\t%s\t%s\t%s\t%s\t%d\t%d\t%s", r.UUID, r.Name, r.Type, r.State, r.Status, r.ActualSize, r.Size, r.Description)
			t.Logf("======================================")
		})
	}
}

func TestZSClient_SetVolumeQoS(t *testing.T) {
	params := param.SetVolumeQoSParam{
		SetVolumeQoS: param.SetVolumeQoSDetailParam{
			Mode:            "total",
			VolumeBandwidth: 1024,
			TotalBandwidth:  10000.0,
			ReadBandwidth:   1024,
			WriteBandwidth:  1024,
			TotalIOPS:       1000.0,
			ReadIOPS:        1024,
			WriteIOPS:       1024,
		},
	}
	type args struct {
		uuid   string
		params param.SetVolumeQoSParam
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{volumeID[accountLogin], params}, true},
		{accessKeyAuth, accessKeyAuthCli, args{volumeID[accessKeyAuth], params}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cli.SetVolumeQoS(tt.args.uuid, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.SetVolumeQoS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestZSClient_GetVolumeQoS(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{volumeID[accountLogin]}, true},
		{accessKeyAuth, accessKeyAuthCli, args{volumeID[accessKeyAuth]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := tt.cli.GetVolumeQoS(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.GetVolumeQoS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.GetVolumeQoS() = %v", r)
			t.Logf("======================================")
			t.Logf("%s\t%d\t%d\t%d\t%d\t%d\t%d", r.VolumeUuid,
				r.VolumeBandwidth, r.VolumeBandwidthRead, r.VolumeBandwidthWrite,
				r.IopsTotal, r.IopsRead, r.IopsWrite)
			t.Logf("======================================")
		})
	}
}

func TestZSClient_DeleteVolumeQoS(t *testing.T) {
	type args struct {
		uuid string
		mode string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{volumeID[accountLogin], "total"}, true},
		{accessKeyAuth, accessKeyAuthCli, args{volumeID[accessKeyAuth], "total"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cli.DeleteVolumeQoS(tt.args.uuid, tt.args.mode); (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.DeleteVolumeQoS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZSClient_GetDataVolumeAttachableVm(t *testing.T) {
	type args struct {
		volumeUuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{volumeID[accountLogin]}, false},
		{accessKeyAuth, accessKeyAuthCli, args{volumeID[accessKeyAuth]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.GetDataVolumeAttachableVm(tt.args.volumeUuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.GetDataVolumeAttachableVm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.GetDataVolumeAttachableVm() = %v", got)
		})
	}
}

func TestZSClient_AttachDataVolumeToVm(t *testing.T) {
	type args struct {
		volumeUuid     string
		vmInstanceUuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{volumeID[accountLogin], vmID[accountLogin]}, false},
		{accessKeyAuth, accessKeyAuthCli, args{volumeID[accessKeyAuth], vmID[accessKeyAuth]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.AttachDataVolumeToVm(tt.args.volumeUuid, tt.args.vmInstanceUuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.AttachDataVolumeToVm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.AttachDataVolumeToVm() = %v", got)
		})
	}
}

func TestZSClient_DetachDataVolumeFromVm(t *testing.T) {
	type args struct {
		uuid   string
		vmUuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{volumeID[accountLogin], vmID[accountLogin]}, false},
		{accessKeyAuth, accessKeyAuthCli, args{volumeID[accessKeyAuth], vmID[accessKeyAuth]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.DetachDataVolumeFromVm(tt.args.uuid, tt.args.vmUuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.DetachDataVolumeFromVm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.DetachDataVolumeFromVm() = %v", got)
		})
	}
}

func TestZSClient_DeleteDataVolume(t *testing.T) {
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
		{accountLogin, accountLoginCli, args{volumeID[accountLogin], param.DeleteModePermissive}, false},
		{accessKeyAuth, accessKeyAuthCli, args{volumeID[accessKeyAuth], param.DeleteModePermissive}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cli.DeleteDataVolume(tt.args.uuid, tt.args.deleteMode); (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.DeleteDataVolume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZSClient_RecoverDataVolume(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{volumeID[accountLogin]}, false},
		{accessKeyAuth, accessKeyAuthCli, args{volumeID[accessKeyAuth]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := tt.cli.RecoverDataVolume(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.RecoverDataVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.RecoverDataVolume() = %v", r)
			t.Logf("======================================")
			t.Logf("%s\t%s\t%s\t%s\t%s\t%d\t%d\t%s", r.UUID, r.Name, r.Type, r.State, r.Status, r.ActualSize, r.Size, r.Description)
			t.Logf("======================================")
		})
	}
}

func TestZSClient_ExpungeDataVolume(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{volumeID[accountLogin]}, false},
		{accessKeyAuth, accessKeyAuthCli, args{volumeID[accessKeyAuth]}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cli.ExpungeDataVolume(tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.ExpungeDataVolume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
