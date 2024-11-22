// Copyright (c) ZStack.io, Inc.

package client

import (
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// QueryBackupStorage queries backup storage
func (cli *ZSClient) QueryBackupStorage(params param.QueryParam) ([]view.BackupStorageInventoryView, error) {
	var views []view.BackupStorageInventoryView
	return views, cli.List("v1/backup-storage", &params, &views)
}

// PageBackupStorage queries backup storage with pagination
func (cli *ZSClient) PageBackupStorage(params param.QueryParam) ([]view.BackupStorageInventoryView, int, error) {
	var views []view.BackupStorageInventoryView
	total, err := cli.Page("v1/backup-storage", &params, &views)
	return views, total, err
}

// GetBackupStorage retrieves a specific backup storage by UUID
func (cli *ZSClient) GetBackupStorage(uuid string) (*view.BackupStorageInventoryView, error) {
	var resp view.BackupStorageInventoryView
	if err := cli.Get("v1/backup-storage", uuid, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ExportImageFromBackupStorage exports an image from backup storage
func (cli *ZSClient) ExportImageFromBackupStorage(params param.ExportImageFromBackupStorageParam) (view.ExportImageFromBackupStorageResultView, error) {
	var resultView view.ExportImageFromBackupStorageResultView
	return resultView, cli.PutWithSpec("v1/backup-storage", params.BackupStorageUuid, "actions", "", params, &resultView)
}

// DeleteExportedImageFromBackupStorage deletes an exported image from backup storage
func (cli *ZSClient) DeleteExportedImageFromBackupStorage(params param.DeleteExportedImageFromBackupStorageParam) error {
	return cli.DeleteWithSpec("v1/backup-storage", params.BackupStorageUuid, "exported-images/"+params.ImageUuid, "", nil)
}

// AddImageStoreBackupStorage adds an image store backup storage
func (cli *ZSClient) AddImageStoreBackupStorage(params param.AddImageStoreBackupStorageParam) (*view.BackupStorageInventoryView, error) {
	var resp view.BackupStorageInventoryView
	if err := cli.Post("v1/backup-storage/image-store", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AttachBackupStorageToZone attaches backup storage to a zone
func (cli *ZSClient) AttachBackupStorageToZone(zoneUuid, backupStorageUuid string) (*view.BackupStorageInventoryView, error) {
	var resp view.BackupStorageInventoryView
	if err := cli.Post(fmt.Sprintf("v1/zones/%s/backup-storage/%s", zoneUuid, backupStorageUuid), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteBackupStorage deletes a backup storage
func (cli *ZSClient) DeleteBackupStorage(uuid string) error {
	return cli.Delete("v1/backup-storage", uuid, string(param.DeleteModePermissive))
}

// QueryImageStoreBackupStorage queries image store backup storage
func (cli *ZSClient) QueryImageStoreBackupStorage(params param.QueryParam) ([]view.BackupStorageInventoryView, error) {
	var views []view.BackupStorageInventoryView
	return views, cli.List("v1/backup-storage/image-store", &params, &views)
}

// PageImageStoreBackupStorage queries image store backup storage with pagination
func (cli *ZSClient) PageImageStoreBackupStorage(params param.QueryParam) ([]view.BackupStorageInventoryView, int, error) {
	var views []view.BackupStorageInventoryView
	total, err := cli.Page("v1/backup-storage/image-store", &params, &views)
	return views, total, err
}

// GetImageStoreBackupStorage retrieves a specific image store backup storage by UUID
func (cli *ZSClient) GetImageStoreBackupStorage(uuid string) (*view.BackupStorageInventoryView, error) {
	var resp view.BackupStorageInventoryView
	if err := cli.Get("v1/backup-storage/image-store", uuid, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ReconnectImageStoreBackupStorage reconnects an image store backup storage
func (cli *ZSClient) ReconnectImageStoreBackupStorage(uuid string) (*view.BackupStorageInventoryView, error) {
	var resp view.BackupStorageInventoryView
	if err := cli.Put("v1/backup-storage/image-store", uuid, map[string]interface{}{"reconnectImageStoreBackupStorage": map[string]string{}}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateImageStoreBackupStorage updates information for an image store backup storage
func (cli *ZSClient) UpdateImageStoreBackupStorage(uuid string, params param.UpdateImageStoreBackupStorageParam) (*view.BackupStorageInventoryView, error) {
	var resp view.BackupStorageInventoryView
	if err := cli.Put("v1/backup-storage/image-store", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ReclaimSpaceFromImageStore reclaims space from an image store backup storage
func (cli *ZSClient) ReclaimSpaceFromImageStore(uuid string) (*view.GcResultView, error) {
	var resp view.GcResultView
	if err := cli.PutWithRespKey("v1/backup-storage/image-store", uuid, "gcResult", map[string]interface{}{"reclaimSpaceFromImageStore": map[string]string{}}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ChangeBackupStorageState changes the availability state of a backup storage
func (cli *ZSClient) ChangeBackupStorageState(uuid string, state param.StateEvent) (*view.BackupStorageInventoryView, error) {
	var resp view.BackupStorageInventoryView
	if err := cli.Put("v1/backup-storage", uuid, param.ChangeBackupStorageStateParam{
		ChangeBackupStorageState: param.ChangeBackupStorageStateDetailParam{
			StateEvent: state,
		},
	}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
