package client

import (
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// QueryBackupStorage 查询镜像服务器
func (cli *ZSClient) QueryBackupStorage(params param.QueryParam) ([]view.BackupStorageInventoryView, error) {
	var views []view.BackupStorageInventoryView
	return views, cli.List("v1/backup-storage", &params, &views)
}

// PageBackupStorage 分页查询镜像服务器
func (cli *ZSClient) PageBackupStorage(params param.QueryParam) ([]view.BackupStorageInventoryView, int, error) {
	var views []view.BackupStorageInventoryView
	total, err := cli.Page("v1/backup-storage", &params, &views)
	return views, total, err
}

// GetBackupStorage 查询镜像服务器
func (cli *ZSClient) GetBackupStorage(uuid string) (*view.BackupStorageInventoryView, error) {
	var resp view.BackupStorageInventoryView
	if err := cli.Get("v1/backup-storage", uuid, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ExportImageFromBackupStorage 从镜像服务器导出镜像
func (cli *ZSClient) ExportImageFromBackupStorage(params param.ExportImageFromBackupStorageParam) (view.ExportImageFromBackupStorageResultView, error) {
	var resultView view.ExportImageFromBackupStorageResultView
	return resultView, cli.PutWithSpec("v1/backup-storage", params.BackupStorageUuid, "actions", "", params, &resultView)
}

// DeleteExportedImageFromBackupStorage 从镜像服务器删除导出的镜像
func (cli *ZSClient) DeleteExportedImageFromBackupStorage(params param.DeleteExportedImageFromBackupStorageParam) error {
	return cli.DeleteWithSpec("v1/backup-storage", params.BackupStorageUuid, "exported-images/"+params.ImageUuid, "", nil)
}

// AddImageStoreBackupStorage 添加镜像仓库服务器-本地服务器
func (cli *ZSClient) AddImageStoreBackupStorage(params param.AddImageStoreBackupStorageParam) (*view.BackupStorageInventoryView, error) {
	var resp view.BackupStorageInventoryView
	if err := cli.Post("v1/backup-storage/image-store", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AttachBackupStorageToZone 挂载镜像服务器至区域
func (cli *ZSClient) AttachBackupStorageToZone(zoneUuid, backupStorageUuid string) (*view.BackupStorageInventoryView, error) {
	var resp view.BackupStorageInventoryView
	if err := cli.Post(fmt.Sprintf("v1/zones/%s/backup-storage/%s", zoneUuid, backupStorageUuid), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteBackupStorage 删除镜像服务器
func (cli *ZSClient) DeleteBackupStorage(uuid string) error {
	return cli.Delete("v1/backup-storage", uuid, string(param.DeleteModePermissive))
}

// QueryImageStoreBackupStorage 查询镜像仓库服务器
func (cli *ZSClient) QueryImageStoreBackupStorage(params param.QueryParam) ([]view.BackupStorageInventoryView, error) {
	var views []view.BackupStorageInventoryView
	return views, cli.List("v1/backup-storage/image-store", &params, &views)
}

// PageImageStoreBackupStorage 分页查询镜像仓库服务器
func (cli *ZSClient) PageImageStoreBackupStorage(params param.QueryParam) ([]view.BackupStorageInventoryView, int, error) {
	var views []view.BackupStorageInventoryView
	total, err := cli.Page("v1/backup-storage/image-store", &params, &views)
	return views, total, err
}

// GetImageStoreBackupStorage 查询镜像仓库服务器
func (cli *ZSClient) GetImageStoreBackupStorage(uuid string) (*view.BackupStorageInventoryView, error) {
	var resp view.BackupStorageInventoryView
	if err := cli.Get("v1/backup-storage/image-store", uuid, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ReconnectImageStoreBackupStorage 重连镜像仓库服务器
func (cli *ZSClient) ReconnectImageStoreBackupStorage(uuid string) (*view.BackupStorageInventoryView, error) {
	var resp view.BackupStorageInventoryView
	if err := cli.Put("v1/backup-storage/image-store", uuid, map[string]interface{}{"reconnectImageStoreBackupStorage": map[string]string{}}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateImageStoreBackupStorage 更新镜像仓库服务器信息
func (cli *ZSClient) UpdateImageStoreBackupStorage(uuid string, params param.UpdateImageStoreBackupStorageParam) (*view.BackupStorageInventoryView, error) {
	var resp view.BackupStorageInventoryView
	if err := cli.Put("v1/backup-storage/image-store", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ReclaimSpaceFromImageStore 从镜像仓库服务器回收空间
func (cli *ZSClient) ReclaimSpaceFromImageStore(uuid string) (*view.GcResultView, error) {
	var resp view.GcResultView
	if err := cli.PutWithRespKey("v1/backup-storage/image-store", uuid, "gcResult", map[string]interface{}{"reclaimSpaceFromImageStore": map[string]string{}}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ChangeBackupStorageState 更改镜像服务器可用状态(
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
