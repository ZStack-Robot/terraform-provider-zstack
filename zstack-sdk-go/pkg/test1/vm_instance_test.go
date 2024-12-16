// # Copyright (c) ZStack.io, Inc.

package test

import (
	"fmt"
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/util/jsonutils"
)

func TestCreateVmInstance(t *testing.T) {

	network, err := accountLoginCli.QueryL3Network(param.NewQueryParam())
	if err != nil {
		golog.Errorf("TestZSClient_QueryL3Network %v: %v", network, err)
	}
	rootDiskSize := int64(1073741824)
	instance, err := accountLoginCli.CreateVmInstance(param.CreateVmInstanceParam{
		BaseParam: param.BaseParam{
			SystemTags: []string{"resourceConfig::vm::vm.clock.track::guest", "cdroms::Empty::None::None"},
			UserTags:   nil,
			RequestIp:  "",
		},
		Params: param.CreateVmInstanceDetailParam{
			Name: "test-uuid",
			//InstanceOfferingUUID:            "",
			ImageUUID:            "968e87334a12422fbe78c8b72bcfab68",
			L3NetworkUuids:       []string{network[0].UUID},
			Type:                 "",
			RootDiskOfferingUuid: "",
			RootDiskSize:         &rootDiskSize,
			//	DataDiskOfferingUuids:           []string{"04229f19712d41cb990ab4b9252d9f93"},
			DataDiskSizes:                   []int64{10240},
			ZoneUuid:                        "",
			ClusterUUID:                     "",
			HostUuid:                        "",
			PrimaryStorageUuidForRootVolume: nil,
			Description:                     "Description",
			DefaultL3NetworkUuid:            network[0].UUID,
			ResourceUuid:                    "56644230e0384ef6b84764530ef306cd",
			TagUuids:                        nil,
			Strategy:                        "",
			MemorySize:                      1073741824,
			CpuNum:                          3,
		},
	})
	if err != nil {
		t.Errorf("TestCreateVmInstance %v", err)
	}
	golog.Println(instance)
}

func TestStartVmInstance(t *testing.T) {
	resp, err := accountLoginCli.StartVmInstance("139115955fb8498ebddb271e4911b80a",
		&param.StartVmInstanceParam{
			StartVmInstance: param.StartVmInstanceDetailParam{
				HostUuid: "6a86f8ac4655479bb065540ce1345708",
			},
		})
	if err != nil {
		t.Errorf("TestStartVmInstance : %v", err)
	}
	golog.Println(resp)
}

func TestStopVmInstance(t *testing.T) {
	resp, err := accountLoginCli.StopVmInstance("09b65cc19e3b47b8b4c7347a5baa4c2a", param.StopVmInstanceParam{
		StopVmInstance: param.StopVmInstanceDetailParam{
			Type:   param.Grace,
			StopHA: true,
		},
	})
	if err != nil {
		t.Errorf("TestStopVmInstance : %v", err)
	}
	golog.Println(resp)
}

func TestRebootVmInstance(t *testing.T) {
	resp, err := accountLoginCli.RebootVmInstance("09b65cc19e3b47b8b4c7347a5baa4c2a")
	if err != nil {
		t.Errorf("TestRebootVmInstance : %v", err)
	}
	golog.Println(resp)
}

func TestPauseVmInstance(t *testing.T) {
	resp, err := accountLoginCli.PauseVmInstance("09b65cc19e3b47b8b4c7347a5baa4c2a")
	if err != nil {
		t.Errorf("TestPauseVmInstance : %v", err)
	}
	golog.Println(resp)
}

func TestResumeVmInstance(t *testing.T) {
	resp, err := accountLoginCli.ResumeVmInstance("09b65cc19e3b47b8b4c7347a5baa4c2a")
	if err != nil {
		t.Errorf("TestResumeVmInstance : %v", err)
	}
	golog.Println(resp)
}

func TestGetVmAttachableDataVolume(t *testing.T) {
	data, err := accountLoginCli.GetVmAttachableDataVolume("09b65cc19e3b47b8b4c7347a5baa4c2a")
	if err != nil {
		t.Errorf("TestGetVmAttachableDataVolume : %v", err)
	}
	golog.Infof("%v", data)
}

func TestGetVmAttachableL3Network(t *testing.T) {
	network, err := accountLoginCli.GetVmAttachableL3Network("2a2f85dc3ef14e27907fbc24440d6b0b")
	if err != nil {
		t.Errorf("TestGetVmAttachableL3Network : %v", err)
	}
	golog.Infof("%v", network)
}

func TestUpdateVmInstance(t *testing.T) {
	data, err := accountLoginCli.UpdateVmInstance("e29e11da127e42fb8844d4a0e421c84a", param.UpdateVmInstanceParam{
		UpdateVmInstance: param.UpdateVmInstanceDetailParam{
			Name:                 "",
			Description:          nil,
			State:                "",
			DefaultL3NetworkUuid: "",
			Platform:             "Windows",
			CpuNum:               nil,
			MemorySize:           nil,
			GuestOsType:          "Windows 8",
		},
	})
	if err != nil {
		t.Errorf("TestUpdateVmInstance : %v", err)
	}
	golog.Infof("%v", data)
}

func TestDestroyVmInstance(t *testing.T) {
	err := accountLoginCli.DestroyVmInstance("09b65cc19e3b47b8b4c7347a5baa4c2a1", param.DeleteModePermissive)
	if err != nil {
		t.Errorf("TestDestroyVmInstance : %v", err)
	}
}

func TestQueryVmInstance(t *testing.T) {
	data, err := accountLoginCli.QueryVmInstance(param.NewQueryParam())
	if err != nil {
		t.Errorf("TestQueryVmInstance : %v", err)
	}
	golog.Info(jsonutils.Marshal(data))
}

func TestGetVmInstance(t *testing.T) {
	data, err := accountLoginCli.GetVmInstance("10c9a1aa52f5404899586ff6bd484685")
	if err != nil {
		t.Errorf("TestGetVmInstance : %v", err)
	}
	golog.Info(jsonutils.Marshal(data))
}

func TestGetVmConsoleAddress(t *testing.T) {
	address, err := accountLoginCli.GetVmConsoleAddress("10c9a1aa52f5404899586ff6bd484685")
	if err != nil {
		t.Errorf("TestGetVmConsoleAddress GetVmConsoleAddress error: %v", err)
	}
	golog.Println(address)
}

func TestGetInstanceConsolePassword(t *testing.T) {
	password, err := accountLoginCli.GetInstanceConsolePassword("10c9a1aa52f5404899586ff6bd484685")
	if err != nil {
		t.Errorf("TestGetVmConsoleAddress GetVmConsoleAddress error: %v", err)
	}
	golog.Println("pwd: " + password)
}

func TestLiveMigrateVM(t *testing.T) {
	data, err := accountLoginCli.LiveMigrateVM("11e6aac3f81c49e98e924b6f57b44da0", "6a86f8ac4655479bb065540ce1345708", false)
	if err != nil {
		t.Errorf("TestLiveMigrateVM LiveMigrateVM error: %v", err)
	}
	golog.Println(jsonutils.Marshal(data))
}

func TestGetVmMigrationCandidateHosts(t *testing.T) {
	hosts, err := accountLoginCli.GetVmMigrationCandidateHosts("b4db493c7dd14e729f3d4719afd836b8")
	if err != nil {
		t.Errorf("TestGetVmMigrationCandidateHosts GetVmMigrationCandidateHosts error: %v", err)
	}
	golog.Println(jsonutils.Marshal(hosts))
}

func TestGetVmQga(t *testing.T) {
	qga, err := accountLoginCli.GetVmQga("6d2c6006a9b14c9eb6caa36b097bfbae")
	if err != nil {
		t.Errorf("TestGetVmMigrationCandidateHosts GetVmMigrationCandidateHosts error: %v", err)
	}
	golog.Println(jsonutils.Marshal(qga))
}

func TestSetVmQga(t *testing.T) {
	err := accountLoginCli.SetVmQga(param.UpdateVmInstanceQgaParam{
		SetVmQga: param.SetVmQgaParam{
			Enable: true,
		},
		UUID: "2d6c502f0177451db6f9320e83dc6a6a",
	})
	fmt.Println(err)
}

func TestGetVmSshKey(t *testing.T) {
	key, err := accountLoginCli.GetVmSshKey("6d2c6006a9b14c9eb6caa36b097bfbae")
	if err != nil {
		t.Errorf("TestGetVmMigrationCandidateHosts GetVmMigrationCandidateHosts error: %v", err)
	}
	golog.Println(jsonutils.Marshal(key))
}

func TestSetVmSshKey(t *testing.T) {
	accountLoginCli.SetVmSshKey(param.UpdateVmInstanceSshKeyParam{
		UUID: "6d2c6006a9b14c9eb6caa36b097bfbae",
		SetVmSshKey: param.SetSshKeyParam{
			SshKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCaaV5YUvz9nx54+pvxI\\\ne5L5uQFHFQsvpwRdVRfMObIgWgcliB9vl4hMCPHXfaKqJD79jBpwJWpUBPebKF7vgevWqFJeU\\\ngR/LBHTfOnRrEjVsSzanaGGzfjbrwMHdZ5YJVhDTE376+OuXz1Wu5M1mwcarJpcanmqNgyz8Yh\\\nYjc50xKDusDVvtpLKxdC6WvhR0+7gaDJKkukip1Up8doOUeNUe2cObJfMoOgi2lNrtKorGp1O7\\\nNv+mdTflboYizgQOCFReiW/1ipPjX06OMZZ3Tsx3ZwBib5ocDpLV9CjONvnDBygWb30wydVoUS\\\np1hKIzlWPkfyWHjxCf9pvLcHGUXZ root@10-0-98-199",
		},
	})
}

func TestDeleteVmSshKey(t *testing.T) {
	err := accountLoginCli.DeleteVmSshKey("6d2c6006a9b14c9eb6caa36b097bfbae", param.DeleteModeEnforcing)
	if err != nil {
		t.Errorf("TestGetVmMigrationCandidateHosts GetVmMigrationCandidateHosts error: %v", err)
	}
}

func TestChangeVmPassword(t *testing.T) {
	err := accountLoginCli.ChangeVmPassword(param.UpdateVmInstanceChangePwdParam{
		UUID: "6d2c6006a9b14c9eb6caa36b097bfbae",
		ChangeVmPassword: param.ChangeVmPasswordParam{
			Password: "root",
			Account:  "root",
		},
	})
	if err != nil {
		t.Errorf("TestChangeVmPassword ChangeVmPassword error: %v", err)
	}
}

func TestGetCandidateIsoForAttachingVm(t *testing.T) {
	p := param.NewQueryParam()
	vm, err := accountLoginCli.GetCandidateIsoForAttachingVm("ef94a22342cb43a2b571f208f8dbfe93", &p)
	if err != nil {
		t.Errorf("TestGetCandidateIsoForAttachingVm GetCandidateIsoForAttachingVm error: %v", err)
		return
	}
	golog.Println(jsonutils.Marshal(vm))
}

func TestGetCandidateVmForAttachingIso(t *testing.T) {
	p := param.NewQueryParam()
	iso, err := accountLoginCli.GetCandidateVmForAttachingIso("300c24f5d4bc4758aa27bb5d303c0930", &p)
	if err != nil {
		t.Errorf("TestGetCandidateVmForAttachingIso GetCandidateVmForAttachingIso error: %v", err)
		return
	}
	golog.Println(jsonutils.Marshal(iso))
}

func TestAttachIsoToVmInstance(t *testing.T) {
	instance, err := accountLoginCli.AttachIsoToVmInstance("54bb6848ba7b4719b8ef93083d6e3896", "8f6924be22b54a3db872b17bc7fa08c6", "")
	if err != nil {
		t.Errorf("TestAttachIsoToVmInstance AttachIsoToVmInstance error: %v", err)
		return
	}
	golog.Println(jsonutils.Marshal(instance))
}

func TestDetachIsoFromVmInstance(t *testing.T) {
	resp, err := accountLoginCli.DetachIsoFromVmInstance("5d560219635e4283a2be9238e72fd84a", "3c89ad443f73444bb24e0b3524e37b86")
	if err != nil {
		t.Errorf("TestDetachIsoFromVmInstance DetachIsoFromVmInstance error: %v", err)
		return
	}
	golog.Println(jsonutils.Marshal(resp))
}

func TestSetVmClockTrack(t *testing.T) {
	resp, err := accountLoginCli.SetVmClockTrack("fe2fe11506224d0faa0ef697a826919b", param.UpdateVmInstanceClockTrackParam{
		BaseParam: param.BaseParam{},
		SetVmClockTrack: param.UpdateVmInstanceClockTrackDetailParam{
			Track:             "guest",
			SyncAfterVMResume: true,
			IntervalInSeconds: 60,
		},
	})
	if err != nil {
		t.Errorf("TestSetVmClockTrack SetVmClockTrack error: %v", err)
		return
	}
	golog.Println(jsonutils.Marshal(resp))
}

func TestQueryVmCdRom(t *testing.T) {
	p := param.NewQueryParam()
	p.AddQ("vmInstanceUuid=" + "6772515251ba402cad1cbf430fffab50")
	rom, err := accountLoginCli.QueryVmCdRom(p)
	if err != nil {
		t.Errorf("TestQueryVmCdRom QueryVmCdRom error: %v", err)
		return
	}
	golog.Println(jsonutils.Marshal(rom))
}

func TestPageVmCdRom(t *testing.T) {
	p := param.NewQueryParam()
	p.AddQ("vmInstanceUuid=" + "6772515251ba402cad1cbf430fffab50")
	rom, num, err := accountLoginCli.PageVmCdRom(p)
	if err != nil {
		t.Errorf("TestPageVmCdRom PageVmCdRom error: %v", err)
		return
	}
	golog.Println(num, jsonutils.Marshal(rom))
}

func TestCreateVmCdRom(t *testing.T) {
	rom, err := accountLoginCli.CreateVmCdRom(param.CreateVmCdRomParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVmCdRomDetailParam{
			Name:           "CDROM-3",
			Description:    "dd",
			VmInstanceUuid: "6772515251ba402cad1cbf430fffab50",
			IsoUuid:        "cf930105a98c44c3856b20e13d6e42db",
		},
	})
	if err != nil {
		t.Errorf("TestCreateVmCdRom CreateVmCdRom error: %v", err)
		return
	}
	golog.Println(jsonutils.Marshal(rom))
}

func TestDeleteVmCdRom(t *testing.T) {
	err := accountLoginCli.DeleteVmCdRom("667f7c34019347f6a59c12696d6b0687", param.DeleteModePermissive)
	if err != nil {
		t.Errorf("TestDeleteVmCdRom DeleteVmCdRom error: %v", err)
		return
	}
}
