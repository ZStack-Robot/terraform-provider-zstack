package client

import (
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/util/jsonutils"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// PageImage 分页
func (cli *ZSClient) PageImage(params param.QueryParam) ([]view.ImageView, int, error) {
	var images []view.ImageView
	total, err := cli.Page("v1/images", &params, &images)
	return images, total, err
}

// QueryImage 查询镜像
func (cli *ZSClient) QueryImage(params param.QueryParam) ([]view.ImageView, error) {
	var images []view.ImageView
	return images, cli.List("v1/images", &params, &images)
}

// QueryImage 查询所有镜像
func (cli *ZSClient) ListAllImages() ([]view.ImageView, error) {
	params := param.NewQueryParam()
	var images []view.ImageView
	return images, cli.ListAll("v1/images", &params, &images)
}

// GetImage 查询镜像
func (cli *ZSClient) GetImage(uuid string) (*view.ImageView, error) {
	image := view.ImageView{}
	return &image, cli.Get("v1/images", uuid, nil, &image)
}

// AddImage 添加镜像
func (cli *ZSClient) AddImage(imageParam param.AddImageParam) (*view.ImageView, error) {
	image := view.ImageView{}
	return &image, cli.Post("v1/images", imageParam, &image)
}

// UpdateImage 编辑镜像
func (cli *ZSClient) UpdateImage(uuid string, params param.UpdateImageParam) (view.ImageView, error) {
	image := view.ImageView{}
	return image, cli.Put("v1/images", uuid, params, &image)
}

// DeleteImage 删除镜像
func (cli *ZSClient) DeleteImage(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/images", uuid, string(deleteMode))
}

// UpdateImageVirtio 更新镜像信息
func (cli *ZSClient) UpdateImageVirtio(params param.UpdateImageVirtioParam) (view.ImageView, error) {
	image := view.ImageView{}
	return image, cli.Put("v1/images", params.UpdateImage.UUID, params, &image)
}

// UpdateArchitectureParam 修改虚拟机镜像CPU架构
func (cli *ZSClient) UpdateArchitectureParam(params param.UpdateImageArchitectureParam) (view.ImageView, error) {
	image := view.ImageView{}
	return image, cli.Put("v1/images", params.UpdateImage.UUID, params, &image)
}

// UpdateImagePlatform 更新镜像信息
func (cli *ZSClient) UpdateImagePlatform(params param.UpdateImagePlatformParam) (view.ImageView, error) {
	image := view.ImageView{}
	return image, cli.Put("v1/images", params.UpdateImage.UUID, params, &image)
}

// ExpungeImage 彻底删除镜像
func (cli *ZSClient) ExpungeImage(imageId string) error {
	params := map[string]interface{}{
		"expungeImage": jsonutils.NewDict(),
	}
	return cli.Put("v1/images", imageId, jsonutils.Marshal(params), nil)
}

// RecoverImage 恢复镜像
func (cli *ZSClient) RecoverImage(params param.RecoverImageParam) (view.ImageView, error) {
	image := view.ImageView{}
	return image, cli.Put("v1/images", params.ImageUuid, params, &image)
}

// ChangeImageState 修改镜像状态
func (cli *ZSClient) ChangeImageState(params param.ChangeImageStateParam) (view.ImageView, error) {
	image := view.ImageView{}
	return image, cli.Put("v1/images", params.ImageUuid, params, &image)
}

// SyncImageSize 刷新镜像大小信息
func (cli *ZSClient) SyncImageSize(params param.SyncImageSizeParam) (view.ImageView, error) {
	image := view.ImageView{}
	return image, cli.Put("v1/images", params.ImageUuid, params, &image)
}

// GetCandidateBackupStorageForCreatingImage 获取镜像服务器候选
func (cli *ZSClient) GetCandidateBackupStorageForCreatingImage(params param.GetCandidateBackupStorageForCreatingImageParam) ([]view.ImageView, error) {
	resource := "v1/images"
	switch params.CandidateBackupStorageType {
	case param.CandidateBackupStorageTypeVolumes:
		resource = fmt.Sprintf("v1/images/volumes/%s", params.VolumeUuid)
	case param.CandidateBackupStorageTypeVolumeSnapshots:
		resource = fmt.Sprintf("v1/images/volume-snapshot/%s", params.VolumeSnapshotUuid)
	}

	resp := make([]view.ImageView, 0)
	return resp, cli.GetWithSpec(resource, "", fmt.Sprintf("candidate-backup-storage?volumeUuid=%s&volumeSnapshotUuid=%s", params.VolumeUuid, params.VolumeSnapshotUuid), responseKeyInventories, nil, &resp)
}

// CreateRootVolumeTemplateFromRootVolume 从根云盘创建根云盘镜像
func (cli *ZSClient) CreateRootVolumeTemplateFromRootVolume(params param.CreateRootVolumeTemplateFromRootVolumeParam) (view.ImageView, error) {
	image := view.ImageView{}
	resource := fmt.Sprintf("v1/images/root-volume-templates/from/volumes/%s", params.RootVolumeUuid)
	return image, cli.Post(resource, params, &image)
}

// CreateRootVolumeTemplateFromVolumeSnapshot 从云盘快照创建根云盘镜像
func (cli *ZSClient) CreateRootVolumeTemplateFromVolumeSnapshot(params param.CreateRootVolumeTemplateFromVolumeSnapshotParam) (view.ImageView, error) {
	image := view.ImageView{}
	resource := fmt.Sprintf("v1/images/root-volume-templates/from/volume-snapshots/%s", params.SnapshotUuid)
	return image, cli.Post(resource, params, &image)
}

// CreateDataVolumeTemplateFromVolume 从云盘创建数据云盘镜像
func (cli *ZSClient) CreateDataVolumeTemplateFromVolume(params param.CreateDataVolumeTemplateFromVolumeParam) (view.ImageView, error) {
	image := view.ImageView{}
	resource := fmt.Sprintf("v1/images/data-volume-templates/from/volumes/%s", params.VolumeUuid)
	return image, cli.Post(resource, params, &image)
}

// CreateDataVolumeTemplateFromVolumeSnapshot 从云盘快照创建数据云盘镜像
func (cli *ZSClient) CreateDataVolumeTemplateFromVolumeSnapshot(params param.CreateDataVolumeTemplateFromVolumeSnapshotParam) (view.ImageView, error) {
	image := view.ImageView{}
	resource := fmt.Sprintf("v1/images/data-volume-templates/from/volume-snapshots/%s", params.SnapshotUuid)
	return image, cli.Post(resource, params, &image)
}

// GetImageQga 获取镜像Qga
func (cli *ZSClient) GetImageQga(uuid string) (view.GetImageQgaView, error) {
	resp := view.GetImageQgaView{}
	return resp, cli.GetWithSpec("v1/images", uuid, "qga", "", nil, &resp)
}

// SetImageQga 设置镜像Qga
func (cli *ZSClient) SetImageQga(params param.SetImageQgaParam) (error, error) {
	return cli.Put("v1/images", params.Uuid, params, nil), nil
}

// SetImageBootMode 设置镜像启动模式
func (cli *ZSClient) SetImageBootMode(params param.SetImageBootModeRequest) error {
	return cli.Put("v1/images", params.Uuid, params, nil)
}

// GetUploadImageJobDetails 获取上传镜像任务详情
func (cli *ZSClient) GetUploadImageJobDetails(params param.GetUploadImageJobDetailsParam) (bool, error) {
	//resp := make([]view.ExistingJobDetails, 0)
	success := false
	return success, cli.GetWithSpec("v1/images/upload-job/details", params.ImageId, "", "success", nil, &success)
}

// GetCandidateVmForAttachingIso 获取ISO可加载云主机列表
func (cli *ZSClient) GetCandidateVmForAttachingIso(uuid string, p *param.QueryParam) ([]view.VmInstanceInventoryView, error) {
	resp := make([]view.VmInstanceInventoryView, 0)
	return resp, cli.List("v1/images/iso/"+uuid+"/vm-candidates", p, &resp)
}
