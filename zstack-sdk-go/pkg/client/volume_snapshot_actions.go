// Copyright (c) ZStack.io, Inc.

package client

import (
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

func (cli *ZSClient) CreateVolumeSnapshot(volumeUuid string, params param.VolumeSnapshotParam) (*view.VolumeSnapshotView, error) {
	snapshot := view.VolumeSnapshotView{}
	return &snapshot, cli.Post(fmt.Sprintf("v1/volumes/%s/volume-snapshots", volumeUuid), &params, &snapshot)
}

func (cli *ZSClient) QueryVolumeSnapshot(params param.QueryParam) ([]view.VolumeSnapshotView, error) {
	snapshots := []view.VolumeSnapshotView{}
	return snapshots, cli.List("v1/volume-snapshots", &params, &snapshots)
}

func (cli *ZSClient) PageVolumeSnapshot(params param.QueryParam) ([]view.VolumeSnapshotView, int, error) {
	snapshots := []view.VolumeSnapshotView{}
	total, err := cli.Page("v1/volume-snapshots", &params, &snapshots)
	return snapshots, total, err
}

func (cli *ZSClient) GetVolumeSnapshot(uuid string) (*view.VolumeSnapshotView, error) {
	snapshot := view.VolumeSnapshotView{}
	return &snapshot, cli.Get("v1/volume-snapshots", uuid, nil, &snapshot)
}

func (cli *ZSClient) QueryVolumeSnapshotTree(params param.QueryParam) ([]view.VolumeSnapshotTreeView, error) {
	trees := []view.VolumeSnapshotTreeView{}
	return trees, cli.List("v1/volume-snapshots/trees", &params, &trees)
}

func (cli *ZSClient) GetVolumeSnapshotTree(uuid string) (*view.VolumeSnapshotTreeView, error) {
	tree := view.VolumeSnapshotTreeView{}
	return &tree, cli.Get("v1/volume-snapshots/trees", uuid, nil, &tree)
}

func (cli *ZSClient) UpdateVolumeSnapshot(uuid string, params param.UpdateVolumeSnapshotParam) (*view.VolumeSnapshotView, error) {
	snapshot := view.VolumeSnapshotView{}
	return &snapshot, cli.Put("v1/volume-snapshots", uuid, &params, &snapshot)
}

func (cli *ZSClient) DeleteVolumeSnapshot(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/volume-snapshots", uuid, string(deleteMode))
}

func (cli *ZSClient) RevertVolumeFromSnapshot(uuid string) error {
	return cli.Put("v1/volume-snapshots", uuid, map[string]struct{}{"revertVolumeFromSnapshot": {}}, nil)
}

func (cli *ZSClient) GetVolumeSnapshotSize(uuid string) (*view.VolumeSnapshotSizeView, error) {
	size := view.VolumeSnapshotSizeView{}
	return &size, cli.PutWithRespKey("v1/volume-snapshots", uuid, "", map[string]struct{}{"getVolumeSnapshotSize": {}}, &size)
}

func (cli *ZSClient) ShrinkVolumeSnapshot(uuid string) (*view.VolumeSnapshotShrinkResultView, error) {
	shrinkResult := view.VolumeSnapshotShrinkResultView{}
	return &shrinkResult, cli.Put("v1/volume-snapshots/shrink", uuid, map[string]struct{}{"shrinkVolumeSnapshot": {}}, &shrinkResult)
}
