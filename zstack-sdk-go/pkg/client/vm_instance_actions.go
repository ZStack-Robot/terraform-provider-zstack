package client

import (
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// CreateVmInstance 创建云主机
func (cli *ZSClient) CreateVmInstance(params param.CreateVmInstanceParam) (*view.VmInstanceInventoryView, error) {
	resp := view.VmInstanceInventoryView{}
	if err := cli.Post("v1/vm-instances", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateVmInstanceFromVolume 从云盘创建云主机
func (cli *ZSClient) CreateVmInstanceFromVolume(params param.CreateVmFromVolumeParam) (*view.VmInstanceInventoryView, error) {
	resp := view.VmInstanceInventoryView{}
	if err := cli.Post("v1/vm-instances/from/volume", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DestroyVmInstance 删除云主机
func (cli *ZSClient) DestroyVmInstance(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/vm-instances", uuid, string(deleteMode))
}

// ExpungeVmInstance 彻底删除云主机
func (cli *ZSClient) ExpungeVmInstance(uuid string) error {
	params := map[string]struct{}{
		"expungeVmInstance": {},
	}
	return cli.Put("v1/vm-instances", uuid, params, nil)
}

// QueryVmInstance 查询云主机
func (cli *ZSClient) QueryVmInstance(params param.QueryParam) ([]view.VmInstanceInventoryView, error) {
	var resp []view.VmInstanceInventoryView
	return resp, cli.List("v1/vm-instances", &params, &resp)
}

// PageVmInstance 分页查询云主机
func (cli *ZSClient) PageVmInstance(params param.QueryParam) ([]view.VmInstanceInventoryView, int, error) {
	var resp []view.VmInstanceInventoryView
	page, err := cli.Page("v1/vm-instances", &params, &resp)
	return resp, page, err
}

func (cli *ZSClient) GetVmInstance(uuid string) (*view.VmInstanceInventoryView, error) {
	var resp view.VmInstanceInventoryView
	if err := cli.Get("v1/vm-instances", uuid, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// StartVmInstance 启动云主机
func (cli *ZSClient) StartVmInstance(uuid string, params *param.StartVmInstanceParam) (*view.VmInstanceInventoryView, error) {
	resp := view.VmInstanceInventoryView{}
	if params == nil {
		return &resp, cli.Put("v1/vm-instances", uuid, map[string]struct{}{
			"startVmInstance": {},
		}, &resp)
	}
	if err := cli.Put("v1/vm-instances", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// StopVmInstance 停止云主机
func (cli *ZSClient) StopVmInstance(uuid string, params param.StopVmInstanceParam) (*view.VmInstanceInventoryView, error) {
	resp := view.VmInstanceInventoryView{}
	if err := cli.Put("v1/vm-instances", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RebootVmInstance 重启云主机
func (cli *ZSClient) RebootVmInstance(uuid string) (*view.VmInstanceInventoryView, error) {
	params := map[string]struct{}{
		"rebootVmInstance": {},
	}
	resp := view.VmInstanceInventoryView{}
	if err := cli.Put("v1/vm-instances", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// PauseVmInstance 暂停云主机
func (cli *ZSClient) PauseVmInstance(uuid string) (*view.VmInstanceInventoryView, error) {
	params := map[string]struct{}{
		"pauseVmInstance": {},
	}
	resp := view.VmInstanceInventoryView{}
	if err := cli.Put("v1/vm-instances", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ResumeVmInstance 恢复暂停的云主机
func (cli *ZSClient) ResumeVmInstance(uuid string) (*view.VmInstanceInventoryView, error) {
	params := map[string]struct{}{
		"resumeVmInstance": {},
	}
	resp := view.VmInstanceInventoryView{}
	if err := cli.Put("v1/vm-instances", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetVmGuestToolsInfo 获取云主机内部增强工具的信息
func (cli *ZSClient) GetVmGuestToolsInfo(uuid string) (*view.VmGuestToolsInfoView, error) {
	var resp view.VmGuestToolsInfoView
	if err := cli.GetWithSpec("v1/vm-instances", uuid, "guest-tools-infos", "", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetLatestGuestToolsForVm 获取云主机可用的最新增强工具
func (cli *ZSClient) GetLatestGuestToolsForVm(uuid string) (*view.LatestGuestToolsView, error) {
	var resp view.LatestGuestToolsView
	if err := cli.GetWithSpec("v1/vm-instances", uuid, "latest-guest-tools", responseKeyInventory, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetVmBootOrder 获取云主机启动设备列表
func (cli *ZSClient) GetVmBootOrder(uuid string) ([]string, error) {
	var resp []string
	return resp, cli.GetWithSpec("v1/vm-instances", uuid, "boot-orders", "orders", nil, &resp)
}

// AttachGuestToolsIsoToVm 为云主机加载增强工具镜像
func (cli *ZSClient) AttachGuestToolsIsoToVm(uuid string) error {
	params := map[string]struct{}{
		"attachGuestToolsIsoToVm": {},
	}
	return cli.Put("v1/vm-instances", uuid, params, nil)
}

// GetVmAttachableDataVolume 获取云主机可加载云盘列表
func (cli *ZSClient) GetVmAttachableDataVolume(uuid string) ([]view.VolumeView, error) {
	resource := fmt.Sprintf("v1/vm-instances/%s/data-volume-candidates", uuid)
	var resp []view.VolumeView
	params := param.NewQueryParam()
	return resp, cli.List(resource, &params, &resp)
}

// GetVmAttachableL3Network 获取云主机可加载L3网络列表
func (cli *ZSClient) GetVmAttachableL3Network(uuid string) ([]view.L3NetworkInventoryView, error) {
	resource := fmt.Sprintf("v1/vm-instances/%s/l3-networks-candidates", uuid)
	var resp []view.L3NetworkInventoryView
	params := param.NewQueryParam()
	return resp, cli.List(resource, &params, &resp)
}

// CloneVmInstance 克隆云主机到指定物理机
func (cli *ZSClient) CloneVmInstance(uuid string, params param.CloneVmInstanceParam) (*view.CloneVmInstanceResult, error) {
	var resp view.CloneVmInstanceResult
	if err := cli.PutWithRespKey("v1/vm-instances", uuid, "result", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateVmInstance 更新云主机信息
func (cli *ZSClient) UpdateVmInstance(uuid string, params param.UpdateVmInstanceParam) (*view.VmInstanceInventoryView, error) {
	var resp view.VmInstanceInventoryView
	if err := cli.Put("v1/vm-instances", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetVmConsoleAddress 获取云主机控制台地址和访问协议
func (cli *ZSClient) GetVmConsoleAddress(instanceUUID string) (*view.VMConsoleAddressView, error) {
	var resp view.VMConsoleAddressView
	if err := cli.GetWithSpec("v1/vm-instances", instanceUUID, "console-addresses", "", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetInstanceConsolePassword 获取云主机控制台密码
func (cli *ZSClient) GetInstanceConsolePassword(instanceUUID string) (string, error) {
	resp := view.GetVmConsolePasswordView{}
	return resp.ConsolePassword, cli.GetWithSpec("v1/vm-instances", instanceUUID, "console-passwords", "", nil, &resp)
}

// LiveMigrateVM 热迁移云主机
// hostUUID is the target host's uuid, if empty will choose a host automatically by cloud.
func (cli *ZSClient) LiveMigrateVM(instanceUUID, hostUUID string, autoConverge bool) (*view.VmInstanceInventoryView, error) {
	type migrateVM struct {
		HostUUID string `json:"hostUuid"`
		Strategy string `json:"strategy"`
	}
	migratePara := migrateVM{
		HostUUID: hostUUID,
	}
	if autoConverge {
		migratePara.Strategy = "auto-converge"
	}
	params := map[string]migrateVM{
		"migrateVm": migratePara,
	}
	resp := view.VmInstanceInventoryView{}
	if err := cli.Put("v1/vm-instances", instanceUUID, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetVmMigrationCandidateHosts 获取可热迁移的物理机列表
func (cli *ZSClient) GetVmMigrationCandidateHosts(instanceUUID string) ([]view.HostInventoryView, error) {
	var resp []view.HostInventoryView
	queryParam := param.NewQueryParam()
	return resp, cli.List(fmt.Sprintf("v1/vm-instances/%s/migration-target-hosts", instanceUUID), &queryParam, &resp)
}

// GetVmStartingCandidateClustersHosts 获取云主机可启动目的地列表
func (cli *ZSClient) GetVmStartingCandidateClustersHosts(instanceUUID string) ([]view.HostInventoryView, error) {
	var resp []view.HostInventoryView
	queryParam := param.NewQueryParam()
	return resp, cli.ListWithRespKey(fmt.Sprintf("v1/vm-instances/%s/starting-target-hosts", instanceUUID), "hosts", &queryParam, &resp)
}

// GetVmQga 获取云主机Qga
func (cli *ZSClient) GetVmQga(uuid string) (view.VMQgaView, error) {
	resp := view.VMQgaView{}
	return resp, cli.GetWithSpec("v1/vm-instances", uuid, "qga", "", nil, &resp)
}

// SetVmQga 设置云主机Qga
func (cli *ZSClient) SetVmQga(params param.UpdateVmInstanceQgaParam) error {
	return cli.Put("v1/vm-instances", params.UUID, params, nil)
}

// SetVmBootMode 设置云主机启动模式
func (cli *ZSClient) SetVmBootMode(uuid string, params param.SetVmBootModeParam) error {
	return cli.Put("v1/vm-instances", uuid, params, nil)
}

// GetVmSshKey 获取云主机SSH Key
func (cli *ZSClient) GetVmSshKey(uuid string) (view.VMSshKeyView, error) {
	resp := view.VMSshKeyView{}
	return resp, cli.GetWithSpec("v1/vm-instances", uuid, "ssh-keys", "", nil, &resp)
}

// SetVmSshKey 设置云主机SSH Key
func (cli *ZSClient) SetVmSshKey(params param.UpdateVmInstanceSshKeyParam) error {
	return cli.Put("v1/vm-instances", params.UUID, params, nil)
}

// DeleteVmSshKey 删除云主机SSH Key
func (cli *ZSClient) DeleteVmSshKey(uuid string, mode param.DeleteMode) error {
	return cli.DeleteWithSpec("v1/vm-instances", uuid, "ssh-keys", fmt.Sprintf("mode=%s", mode), nil)
}

// ChangeVmPassword 变更云主机密码
func (cli *ZSClient) ChangeVmPassword(params param.UpdateVmInstanceChangePwdParam) error {
	return cli.Put("v1/vm-instances", params.UUID, params, nil)
}

// GetCandidateIsoForAttachingVm 获取云主机可加载ISO列表
func (cli *ZSClient) GetCandidateIsoForAttachingVm(uuid string, p *param.QueryParam) ([]view.ImageView, error) {
	resp := make([]view.ImageView, 0)
	return resp, cli.List("v1/vm-instances/"+uuid+"/iso-candidates", p, &resp)
}

// AttachIsoToVmInstance 加载ISO到云主机
func (cli *ZSClient) AttachIsoToVmInstance(isoUUID, instanceUUID, cdRomUUID string) (*view.VmInstanceInventoryView, error) {
	var resp view.VmInstanceInventoryView
	p := param.BaseParam{
		SystemTags: []string{fmt.Sprintf("cdromUuid::%s", cdRomUUID)},
	}
	if err := cli.Post("v1/vm-instances/"+instanceUUID+"/iso/"+isoUUID, p, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DetachIsoFromVmInstance 卸载云主机上的ISO
func (cli *ZSClient) DetachIsoFromVmInstance(instanceUUID, isoUUID string) (*view.VmInstanceInventoryView, error) {
	var resp view.VmInstanceInventoryView
	if err := cli.DeleteWithSpec("v1/vm-instances", instanceUUID, "iso", fmt.Sprintf("isoUuid=%s&deleteMode=%s", isoUUID, param.DeleteModePermissive), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetVmClockTrack 设置云主机时钟同步
func (cli *ZSClient) SetVmClockTrack(uuid string, params param.UpdateVmInstanceClockTrackParam) (*view.VmInstanceInventoryView, error) {
	var resp view.VmInstanceInventoryView
	if err := cli.Put("v1/vm-instances", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// QueryVmCdRom 查询CDROM清单
func (cli *ZSClient) QueryVmCdRom(p param.QueryParam) ([]view.VMCDRomView, error) {
	resp := make([]view.VMCDRomView, 0)
	return resp, cli.List("v1/vm-instances/cdroms", &p, &resp)
}

// PageVmCdRom 分页查询CDROM清单
func (cli *ZSClient) PageVmCdRom(p param.QueryParam) ([]view.VMCDRomView, int, error) {
	resp := make([]view.VMCDRomView, 0)
	num, err := cli.Page("v1/vm-instances/cdroms", &p, &resp)
	return resp, num, err
}

// SetVmInstanceDefaultCdRom 设置云主机默认CDROM
func (cli *ZSClient) SetVmInstanceDefaultCdRom(vmInstanceUUID, cdRomUUID string) (*view.VMCDRomView, error) {
	var resp view.VMCDRomView
	if err := cli.Put("v1/vm-instances/", vmInstanceUUID+"/cdroms/"+cdRomUUID, map[string]interface{}{"setVmInstanceDefaultCdRom": nil}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateVmCdRom 修改CDROM
func (cli *ZSClient) UpdateVmCdRom(cdRomUUID string, params param.UpdateVmCdRomParam) (*view.VMCDRomView, error) {
	var resp view.VMCDRomView
	if err := cli.Put("v1/vm-instances/cdroms", cdRomUUID, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteVmCdRom 删除CDROM
func (cli *ZSClient) DeleteVmCdRom(cdRomUUID string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/vm-instances/cdroms", cdRomUUID, string(deleteMode))
}

// CreateVmCdRom 为云主机创建CDROM
func (cli *ZSClient) CreateVmCdRom(params param.CreateVmCdRomParam) (*view.VMCDRomView, error) {
	var resp view.VMCDRomView
	if err := cli.Post("v1/vm-instances/cdroms", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCandidatePrimaryStoragesForCreatingVm 获取可选择的主存储
func (cli *ZSClient) GetCandidatePrimaryStoragesForCreatingVm(params param.QueryParam) (*view.GetCandidatePrimaryStoragesForCreatingVmView, error) {
	resp := new(view.GetCandidatePrimaryStoragesForCreatingVmView)
	return resp, cli.ListWithRespKey("v1/vm-instances/candidate-storages", "", &params, resp)
}
