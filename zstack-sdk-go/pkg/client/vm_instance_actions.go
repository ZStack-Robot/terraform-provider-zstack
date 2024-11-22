// Copyright (c) ZStack.io, Inc.

package client

import (
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// CreateVmInstance Create a VM instance
func (cli *ZSClient) CreateVmInstance(params param.CreateVmInstanceParam) (*view.VmInstanceInventoryView, error) {
	resp := view.VmInstanceInventoryView{}
	if err := cli.Post("v1/vm-instances", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateVmInstanceFromVolume Create a VM instance from a volume
func (cli *ZSClient) CreateVmInstanceFromVolume(params param.CreateVmFromVolumeParam) (*view.VmInstanceInventoryView, error) {
	resp := view.VmInstanceInventoryView{}
	if err := cli.Post("v1/vm-instances/from/volume", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DestroyVmInstance Delete a VM instance
func (cli *ZSClient) DestroyVmInstance(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/vm-instances", uuid, string(deleteMode))
}

// ExpungeVmInstance Permanently delete a VM instance
func (cli *ZSClient) ExpungeVmInstance(uuid string) error {
	params := map[string]struct{}{
		"expungeVmInstance": {},
	}
	return cli.Put("v1/vm-instances", uuid, params, nil)
}

// QueryVmInstance Query VM instances
func (cli *ZSClient) QueryVmInstance(params param.QueryParam) ([]view.VmInstanceInventoryView, error) {
	var resp []view.VmInstanceInventoryView
	return resp, cli.List("v1/vm-instances", &params, &resp)
}

// PageVmInstance Paginate VM instances
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

// StartVmInstance Start a VM instance
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

// StopVmInstance Stop a VM instance
func (cli *ZSClient) StopVmInstance(uuid string, params param.StopVmInstanceParam) (*view.VmInstanceInventoryView, error) {
	resp := view.VmInstanceInventoryView{}
	if err := cli.Put("v1/vm-instances", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RebootVmInstance Reboot a VM instance
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

// PauseVmInstance Pause a VM instance
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

// ResumeVmInstance Resume a paused VM instance
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

// GetVmGuestToolsInfo Get information about VM guest tools
func (cli *ZSClient) GetVmGuestToolsInfo(uuid string) (*view.VmGuestToolsInfoView, error) {
	var resp view.VmGuestToolsInfoView
	if err := cli.GetWithSpec("v1/vm-instances", uuid, "guest-tools-infos", "", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetLatestGuestToolsForVm Get the latest guest tools available for the VM
func (cli *ZSClient) GetLatestGuestToolsForVm(uuid string) (*view.LatestGuestToolsView, error) {
	var resp view.LatestGuestToolsView
	if err := cli.GetWithSpec("v1/vm-instances", uuid, "latest-guest-tools", responseKeyInventory, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetVmBootOrder retrieves the boot device order of a VM.
func (cli *ZSClient) GetVmBootOrder(uuid string) ([]string, error) {
	var resp []string
	return resp, cli.GetWithSpec("v1/vm-instances", uuid, "boot-orders", "orders", nil, &resp)
}

// AttachGuestToolsIsoToVm attaches the guest tools ISO to a VM.
func (cli *ZSClient) AttachGuestToolsIsoToVm(uuid string) error {
	params := map[string]struct{}{
		"attachGuestToolsIsoToVm": {},
	}
	return cli.Put("v1/vm-instances", uuid, params, nil)
}

// GetVmAttachableDataVolume retrieves a list of data volumes that can be attached to a VM.
func (cli *ZSClient) GetVmAttachableDataVolume(uuid string) ([]view.VolumeView, error) {
	resource := fmt.Sprintf("v1/vm-instances/%s/data-volume-candidates", uuid)
	var resp []view.VolumeView
	params := param.NewQueryParam()
	return resp, cli.List(resource, &params, &resp)
}

// GetVmAttachableL3Network retrieves a list of L3 networks that can be attached to a VM.
func (cli *ZSClient) GetVmAttachableL3Network(uuid string) ([]view.L3NetworkInventoryView, error) {
	resource := fmt.Sprintf("v1/vm-instances/%s/l3-networks-candidates", uuid)
	var resp []view.L3NetworkInventoryView
	params := param.NewQueryParam()
	return resp, cli.List(resource, &params, &resp)
}

// CloneVmInstance clones a VM to a specified host.
func (cli *ZSClient) CloneVmInstance(uuid string, params param.CloneVmInstanceParam) (*view.CloneVmInstanceResult, error) {
	var resp view.CloneVmInstanceResult
	if err := cli.PutWithRespKey("v1/vm-instances", uuid, "result", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateVmInstance updates VM information.
func (cli *ZSClient) UpdateVmInstance(uuid string, params param.UpdateVmInstanceParam) (*view.VmInstanceInventoryView, error) {
	var resp view.VmInstanceInventoryView
	if err := cli.Put("v1/vm-instances", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetVmConsoleAddress retrieves the console address and protocol of a VM.
func (cli *ZSClient) GetVmConsoleAddress(instanceUUID string) (*view.VMConsoleAddressView, error) {
	var resp view.VMConsoleAddressView
	if err := cli.GetWithSpec("v1/vm-instances", instanceUUID, "console-addresses", "", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetInstanceConsolePassword retrieves the console password of a VM.
func (cli *ZSClient) GetInstanceConsolePassword(instanceUUID string) (string, error) {
	resp := view.GetVmConsolePasswordView{}
	return resp.ConsolePassword, cli.GetWithSpec("v1/vm-instances", instanceUUID, "console-passwords", "", nil, &resp)
}

// LiveMigrateVM performs a live migration of a VM.
// hostUUID specifies the target host. If empty, the cloud chooses a host automatically.
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

// GetVmMigrationCandidateHosts retrieves a list of hosts eligible for live migration.
func (cli *ZSClient) GetVmMigrationCandidateHosts(instanceUUID string) ([]view.HostInventoryView, error) {
	var resp []view.HostInventoryView
	queryParam := param.NewQueryParam()
	return resp, cli.List(fmt.Sprintf("v1/vm-instances/%s/migration-target-hosts", instanceUUID), &queryParam, &resp)
}

// GetVmStartingCandidateClustersHosts retrieves a list of potential destinations for starting a VM.
func (cli *ZSClient) GetVmStartingCandidateClustersHosts(instanceUUID string) ([]view.HostInventoryView, error) {
	var resp []view.HostInventoryView
	queryParam := param.NewQueryParam()
	return resp, cli.ListWithRespKey(fmt.Sprintf("v1/vm-instances/%s/starting-target-hosts", instanceUUID), "hosts", &queryParam, &resp)
}

// GetVmQga retrieves the Qga (Guest Agent) information of a VM.
func (cli *ZSClient) GetVmQga(uuid string) (view.VMQgaView, error) {
	resp := view.VMQgaView{}
	return resp, cli.GetWithSpec("v1/vm-instances", uuid, "qga", "", nil, &resp)
}

// SetVmQga sets the Qga (Guest Agent) configuration for a VM.
func (cli *ZSClient) SetVmQga(params param.UpdateVmInstanceQgaParam) error {
	return cli.Put("v1/vm-instances", params.UUID, params, nil)
}

// SetVmBootMode sets the boot mode for a VM.
func (cli *ZSClient) SetVmBootMode(uuid string, params param.SetVmBootModeParam) error {
	return cli.Put("v1/vm-instances", uuid, params, nil)
}

// GetVmSshKey retrieves the SSH key associated with a virtual machine.
func (cli *ZSClient) GetVmSshKey(uuid string) (view.VMSshKeyView, error) {
	resp := view.VMSshKeyView{}
	return resp, cli.GetWithSpec("v1/vm-instances", uuid, "ssh-keys", "", nil, &resp)
}

// SetVmSshKey sets the SSH key for a virtual machine.
func (cli *ZSClient) SetVmSshKey(params param.UpdateVmInstanceSshKeyParam) error {
	return cli.Put("v1/vm-instances", params.UUID, params, nil)
}

// DeleteVmSshKey deletes the SSH key of a virtual machine.
func (cli *ZSClient) DeleteVmSshKey(uuid string, mode param.DeleteMode) error {
	return cli.DeleteWithSpec("v1/vm-instances", uuid, "ssh-keys", fmt.Sprintf("mode=%s", mode), nil)
}

// ChangeVmPassword changes the password of a virtual machine.
func (cli *ZSClient) ChangeVmPassword(params param.UpdateVmInstanceChangePwdParam) error {
	return cli.Put("v1/vm-instances", params.UUID, params, nil)
}

// GetCandidateIsoForAttachingVm retrieves the list of ISO images that can be attached to a virtual machine.
func (cli *ZSClient) GetCandidateIsoForAttachingVm(uuid string, p *param.QueryParam) ([]view.ImageView, error) {
	resp := make([]view.ImageView, 0)
	return resp, cli.List("v1/vm-instances/"+uuid+"/iso-candidates", p, &resp)
}

// AttachIsoToVmInstance attaches an ISO image to a virtual machine.
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

// DetachIsoFromVmInstance detaches an ISO image from a virtual machine.
func (cli *ZSClient) DetachIsoFromVmInstance(instanceUUID, isoUUID string) (*view.VmInstanceInventoryView, error) {
	var resp view.VmInstanceInventoryView
	if err := cli.DeleteWithSpec("v1/vm-instances", instanceUUID, "iso", fmt.Sprintf("isoUuid=%s&deleteMode=%s", isoUUID, param.DeleteModePermissive), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetVmClockTrack configures clock synchronization for a virtual machine.
func (cli *ZSClient) SetVmClockTrack(uuid string, params param.UpdateVmInstanceClockTrackParam) (*view.VmInstanceInventoryView, error) {
	var resp view.VmInstanceInventoryView
	if err := cli.Put("v1/vm-instances", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// QueryVmCdRom queries the list of CD-ROMs for virtual machines.
func (cli *ZSClient) QueryVmCdRom(p param.QueryParam) ([]view.VMCDRomView, error) {
	resp := make([]view.VMCDRomView, 0)
	return resp, cli.List("v1/vm-instances/cdroms", &p, &resp)
}

// PageVmCdRom queries CD-ROMs for virtual machines with pagination.
func (cli *ZSClient) PageVmCdRom(p param.QueryParam) ([]view.VMCDRomView, int, error) {
	resp := make([]view.VMCDRomView, 0)
	num, err := cli.Page("v1/vm-instances/cdroms", &p, &resp)
	return resp, num, err
}

// SetVmInstanceDefaultCdRom sets the default CD-ROM for a virtual machine.
func (cli *ZSClient) SetVmInstanceDefaultCdRom(vmInstanceUUID, cdRomUUID string) (*view.VMCDRomView, error) {
	var resp view.VMCDRomView
	if err := cli.Put("v1/vm-instances/", vmInstanceUUID+"/cdroms/"+cdRomUUID, map[string]interface{}{"setVmInstanceDefaultCdRom": nil}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateVmCdRom updates a CD-ROM for a virtual machine.
func (cli *ZSClient) UpdateVmCdRom(cdRomUUID string, params param.UpdateVmCdRomParam) (*view.VMCDRomView, error) {
	var resp view.VMCDRomView
	if err := cli.Put("v1/vm-instances/cdroms", cdRomUUID, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteVmCdRom deletes a CD-ROM from a virtual machine.
func (cli *ZSClient) DeleteVmCdRom(cdRomUUID string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/vm-instances/cdroms", cdRomUUID, string(deleteMode))
}

// CreateVmCdRom creates a new CD-ROM for a virtual machine.
func (cli *ZSClient) CreateVmCdRom(params param.CreateVmCdRomParam) (*view.VMCDRomView, error) {
	var resp view.VMCDRomView
	if err := cli.Post("v1/vm-instances/cdroms", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCandidatePrimaryStoragesForCreatingVm retrieves the list of candidate primary storages for creating a virtual machine.
func (cli *ZSClient) GetCandidatePrimaryStoragesForCreatingVm(params param.QueryParam) (*view.GetCandidatePrimaryStoragesForCreatingVmView, error) {
	resp := new(view.GetCandidatePrimaryStoragesForCreatingVmView)
	return resp, cli.ListWithRespKey("v1/vm-instances/candidate-storages", "", &params, resp)
}
