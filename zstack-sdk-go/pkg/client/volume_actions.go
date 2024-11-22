// Copyright (c) ZStack.io, Inc.

package client

import (
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

func (cli *ZSClient) CreateDataVolume(params param.CreateDataVolumeParam) (*view.VolumeView, error) {
	volume := view.VolumeView{}
	return &volume, cli.Post("v1/volumes/data", &params, &volume)
}

func (cli *ZSClient) DeleteDataVolume(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/volumes", uuid, string(deleteMode))
}

func (cli *ZSClient) ExpungeDataVolume(uuid string) error {
	return cli.Put("v1/volumes", uuid, map[string]struct{}{"expungeDataVolume": {}}, nil)
}

func (cli *ZSClient) RecoverDataVolume(uuid string) (*view.VolumeView, error) {
	volume := view.VolumeView{}
	return &volume, cli.Put("v1/volumes", uuid, map[string]struct{}{"recoverDataVolume": {}}, &volume)
}

func (cli *ZSClient) ChangeVolumeState(uuid string, params param.ChangeVolumeStateParam) (*view.VolumeView, error) {
	volume := view.VolumeView{}
	return &volume, cli.Put("v1/volumes", uuid, &params, &volume)
}

func (cli *ZSClient) CreateDataVolumeFromVolumeTemplate(imageUuid string, params param.CreateDataVolumeFromVolumeTemplateParam) (*view.VolumeView, error) {
	volume := view.VolumeView{}
	return &volume, cli.Post(fmt.Sprintf("v1/volumes/data/from/data-volume-templates/%s", imageUuid), &params, &volume)
}

func (cli *ZSClient) CreateDataVolumeFromVolumeSnapshot(volumeSnapshotUuid string, params param.CreateDataVolumeFromVolumeSnapshotParam) (*view.VolumeView, error) {
	volume := view.VolumeView{}
	return &volume, cli.Post(fmt.Sprintf("v1/volumes/data/from/volume-snapshots/%s", volumeSnapshotUuid), &params, &volume)
}

func (cli *ZSClient) QueryVolume(params param.QueryParam) ([]view.VolumeView, error) {
	volumes := []view.VolumeView{}
	return volumes, cli.List("v1/volumes", &params, &volumes)
}

func (cli *ZSClient) PageVolume(params param.QueryParam) ([]view.VolumeView, int, error) {
	volumes := []view.VolumeView{}
	total, err := cli.Page("v1/volumes", &params, &volumes)
	return volumes, total, err
}

func (cli *ZSClient) GetVolume(uuid string) (*view.VolumeView, error) {
	volume := view.VolumeView{}
	return &volume, cli.Get("v1/volumes", uuid, nil, &volume)
}

func (cli *ZSClient) GetVolumeFormat() ([]view.VolumeFormatView, error) {
	params := param.NewQueryParam()
	formats := []view.VolumeFormatView{}
	return formats, cli.ListWithRespKey("v1/volumes/formats", "formats", &params, &formats)
}

func (cli *ZSClient) GetVolumeCapabilities(uuid string) (*view.VolumeCapabilitiesView, error) {
	vc := view.VolumeCapabilitiesView{}
	return &vc, cli.GetWithSpec("v1/volumes", uuid, "capabilities", "capabilities", nil, &vc)
}

func (cli *ZSClient) SyncVolumeSize(uuid string) (*view.VolumeView, error) {
	volume := view.VolumeView{}
	return &volume, cli.Put("v1/volumes", uuid, map[string]struct{}{"syncVolumeSize": {}}, &volume)
}

func (cli *ZSClient) ResizeRootVolume(uuid string, size int64) (*view.VolumeView, error) {
	volume := view.VolumeView{}
	return &volume, cli.Put("v1/volumes/resize", uuid, map[string]map[string]int64{
		"resizeRootVolume": {"size": size},
	}, &volume)
}

func (cli *ZSClient) ResizeDataVolume(uuid string, size int64) (*view.VolumeView, error) {
	volume := view.VolumeView{}
	return &volume, cli.Put("v1/volumes/data/resize", uuid, map[string]map[string]int64{
		"resizeDataVolume": {"size": size},
	}, &volume)
}

func (cli *ZSClient) UpdateVolume(uuid string, params param.UpdateVolumeParam) (*view.VolumeView, error) {
	volume := view.VolumeView{}
	return &volume, cli.Put("v1/volumes", uuid, &params, &volume)
}

func (cli *ZSClient) SetVolumeQoS(uuid string, params param.SetVolumeQoSParam) error {
	return cli.Put("v1/volumes", uuid, &params, nil)
}

func (cli *ZSClient) GetVolumeQoS(uuid string) (*view.VolumeQoSView, error) {
	vq := view.VolumeQoSView{}
	return &vq, cli.GetWithSpec("v1/volumes", uuid, "qos", "", nil, &vq)
}

func (cli *ZSClient) DeleteVolumeQoS(uuid string, mode string) error {
	return cli.DeleteWithSpec("v1/volumes", uuid, "qos", fmt.Sprintf("mode=%s", mode), nil)
}

func (cli *ZSClient) GetDataVolumeAttachableVm(volumeUuid string) ([]view.VmInstanceInventoryView, error) {
	params := param.NewQueryParam()
	vms := []view.VmInstanceInventoryView{}
	return vms, cli.List(fmt.Sprintf("v1/volumes/%s/candidate-vm-instances", volumeUuid), &params, &vms)
}

func (cli *ZSClient) AttachDataVolumeToVm(volumeUuid, vmInstanceUuid string) (*view.VolumeView, error) {
	volume := view.VolumeView{}
	return &volume, cli.Post(fmt.Sprintf("v1/volumes/%s/vm-instances/%s", volumeUuid, vmInstanceUuid), nil, &volume)
}

func (cli *ZSClient) DetachDataVolumeFromVm(uuid, vmUuid string) (*view.VolumeView, error) {
	volume := view.VolumeView{}
	return &volume, cli.DeleteWithSpec("v1/volumes", uuid, "vm-instances", fmt.Sprintf("vmUuid=%s", vmUuid), &volume)
}
