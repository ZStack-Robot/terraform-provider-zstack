// Copyright (c) ZStack.io, Inc.

package client

import (
	"strings"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

func (cli *ZSClient) CreateVolumeSnapshotGroup(params param.VolumeSnapshotGroupParam) (*view.VolumeSnapshotGroupView, error) {
	group := view.VolumeSnapshotGroupView{}
	return &group, cli.Post("v1/volume-snapshots/group", &params, &group)
}

func (cli *ZSClient) DeleteVolumeSnapshotGroup(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/volume-snapshots/group", uuid, string(deleteMode))
}

func (cli *ZSClient) UpdateVolumeSnapshotGroup(uuid string, params param.UpdateVolumeSnapshotGroupParam) (*view.VolumeSnapshotGroupView, error) {
	group := view.VolumeSnapshotGroupView{}
	return &group, cli.Put("v1/volume-snapshots/group", uuid, &params, &group)
}

func (cli *ZSClient) QueryVolumeSnapshotGroup(params param.QueryParam) ([]view.VolumeSnapshotGroupView, error) {
	groups := []view.VolumeSnapshotGroupView{}
	return groups, cli.List("v1/volume-snapshots/group", &params, &groups)
}

func (cli *ZSClient) PageVolumeSnapshotGroup(params param.QueryParam) ([]view.VolumeSnapshotGroupView, int, error) {
	groups := []view.VolumeSnapshotGroupView{}
	total, err := cli.Page("v1/volume-snapshots/group", &params, &groups)
	return groups, total, err
}

func (cli *ZSClient) GetVolumeSnapshotGroup(uuid string) (*view.VolumeSnapshotGroupView, error) {
	group := view.VolumeSnapshotGroupView{}
	return &group, cli.Get("v1/volume-snapshots/group", uuid, nil, &group)
}

func (cli *ZSClient) CheckVolumeSnapshotGroupAvailability(uuids []string) ([]view.VolumeSnapshotGroupAvailabilityView, error) {
	params := param.NewQueryParam()
	params.Add("uuids", strings.Join(uuids, ","))
	availabilitys := []view.VolumeSnapshotGroupAvailabilityView{}
	return availabilitys, cli.ListWithRespKey("v1/volume-snapshots/groups/availabilities", "results", &params, &availabilitys)
}

func (cli *ZSClient) RevertVmFromSnapshotGroup(uuid string) error {
	return cli.Put("v1/volume-snapshots/group", uuid, map[string]struct{}{"revertVmFromSnapshotGroup": {}}, nil)
}

func (cli *ZSClient) UngroupVolumeSnapshotGroup(uuid string) error {
	return cli.Delete("v1/volume-snapshots/ungroup", uuid, "")
}
