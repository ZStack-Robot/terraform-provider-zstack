// Copyright (c) ZStack.io, Inc.

package client

import (
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/util/jsonutils"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// PageImage Pagination
func (cli *ZSClient) PageImage(params param.QueryParam) ([]view.ImageView, int, error) {
	var images []view.ImageView
	total, err := cli.Page("v1/images", &params, &images)
	return images, total, err
}

// QueryImage Query Image
func (cli *ZSClient) QueryImage(params param.QueryParam) ([]view.ImageView, error) {
	var images []view.ImageView
	return images, cli.List("v1/images", &params, &images)
}

// QueryImage Query All Images
func (cli *ZSClient) ListAllImages() ([]view.ImageView, error) {
	params := param.NewQueryParam()
	var images []view.ImageView
	return images, cli.ListAll("v1/images", &params, &images)
}

// GetImage Query Image
func (cli *ZSClient) GetImage(uuid string) (*view.ImageView, error) {
	image := view.ImageView{}
	return &image, cli.Get("v1/images", uuid, nil, &image)
}

// AddImage Add Image
func (cli *ZSClient) AddImage(imageParam param.AddImageParam) (*view.ImageView, error) {
	image := view.ImageView{}
	return &image, cli.Post("v1/images", imageParam, &image)
}

func (cli *ZSClient) AddImageAsync(imageParam param.AddImageParam) (string, error) {
	// 调用 PostWithAsync，传入必要的参数
	resource := "v1/images"
	responseKey := ""      // 不需要解析具体的返回数据，只获取 ApiId
	var retVal interface{} // retVal 为 nil，因为异步情况下不需要解析具体数据

	// 调用异步函数
	apiId, err := cli.PostWithAsync(resource, responseKey, imageParam, retVal, true)
	if err != nil {
		return "", err
	}

	return apiId, nil
}

// UpdateImage Edit Image
func (cli *ZSClient) UpdateImage(uuid string, params param.UpdateImageParam) (view.ImageView, error) {
	image := view.ImageView{}
	return image, cli.Put("v1/images", uuid, params, &image)
}

// DeleteImage Delete Image
func (cli *ZSClient) DeleteImage(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/images", uuid, string(deleteMode))
}

// UpdateImageVirtio Update Image Information
func (cli *ZSClient) UpdateImageVirtio(params param.UpdateImageVirtioParam) (view.ImageView, error) {
	image := view.ImageView{}
	return image, cli.Put("v1/images", params.UpdateImage.UUID, params, &image)
}

// UpdateArchitectureParam Modify Virtual Machine Image CPU Architecture
func (cli *ZSClient) UpdateArchitectureParam(params param.UpdateImageArchitectureParam) (view.ImageView, error) {
	image := view.ImageView{}
	return image, cli.Put("v1/images", params.UpdateImage.UUID, params, &image)
}

// UpdateImagePlatform Update Image Information
func (cli *ZSClient) UpdateImagePlatform(params param.UpdateImagePlatformParam) (view.ImageView, error) {
	image := view.ImageView{}
	return image, cli.Put("v1/images", params.UpdateImage.UUID, params, &image)
}

// ExpungeImage Permanently Delete Image
func (cli *ZSClient) ExpungeImage(imageId string) error {
	params := map[string]interface{}{
		"expungeImage": jsonutils.NewDict(),
	}
	return cli.Put("v1/images", imageId, jsonutils.Marshal(params), nil)
}

// RecoverImage Recover Image
func (cli *ZSClient) RecoverImage(params param.RecoverImageParam) (view.ImageView, error) {
	image := view.ImageView{}
	return image, cli.Put("v1/images", params.ImageUuid, params, &image)
}

// ChangeImageState Modify Image State
func (cli *ZSClient) ChangeImageState(params param.ChangeImageStateParam) (view.ImageView, error) {
	image := view.ImageView{}
	return image, cli.Put("v1/images", params.ImageUuid, params, &image)
}

// SyncImageSize Refresh Image Size Information
func (cli *ZSClient) SyncImageSize(params param.SyncImageSizeParam) (view.ImageView, error) {
	image := view.ImageView{}
	return image, cli.Put("v1/images", params.ImageUuid, params, &image)
}

// GetCandidateBackupStorageForCreatingImage Get Backup Storage Candidates for Creating Image
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

// CreateRootVolumeTemplateFromRootVolume Create Root Volume Template from Root Volume
func (cli *ZSClient) CreateRootVolumeTemplateFromRootVolume(params param.CreateRootVolumeTemplateFromRootVolumeParam) (view.ImageView, error) {
	image := view.ImageView{}
	resource := fmt.Sprintf("v1/images/root-volume-templates/from/volumes/%s", params.RootVolumeUuid)
	return image, cli.Post(resource, params, &image)
}

// CreateRootVolumeTemplateFromVolumeSnapshot Create Root Volume Template from Volume Snapshot
func (cli *ZSClient) CreateRootVolumeTemplateFromVolumeSnapshot(params param.CreateRootVolumeTemplateFromVolumeSnapshotParam) (view.ImageView, error) {
	image := view.ImageView{}
	resource := fmt.Sprintf("v1/images/root-volume-templates/from/volume-snapshots/%s", params.SnapshotUuid)
	return image, cli.Post(resource, params, &image)
}

// CreateDataVolumeTemplateFromVolume Create Data Volume Template from Volume
func (cli *ZSClient) CreateDataVolumeTemplateFromVolume(params param.CreateDataVolumeTemplateFromVolumeParam) (view.ImageView, error) {
	image := view.ImageView{}
	resource := fmt.Sprintf("v1/images/data-volume-templates/from/volumes/%s", params.VolumeUuid)
	return image, cli.Post(resource, params, &image)
}

// CreateDataVolumeTemplateFromVolumeSnapshot Create Data Volume Template from Volume Snapshot
func (cli *ZSClient) CreateDataVolumeTemplateFromVolumeSnapshot(params param.CreateDataVolumeTemplateFromVolumeSnapshotParam) (view.ImageView, error) {
	image := view.ImageView{}
	resource := fmt.Sprintf("v1/images/data-volume-templates/from/volume-snapshots/%s", params.SnapshotUuid)
	return image, cli.Post(resource, params, &image)
}

// GetImageQga Get Image Qga
func (cli *ZSClient) GetImageQga(uuid string) (view.GetImageQgaView, error) {
	resp := view.GetImageQgaView{}
	return resp, cli.GetWithSpec("v1/images", uuid, "qga", "", nil, &resp)
}

// SetImageQga Set Image Qga
func (cli *ZSClient) SetImageQga(params param.SetImageQgaParam) (error, error) {
	return cli.Put("v1/images", params.Uuid, params, nil), nil
}

// SetImageBootMode Set Image Boot Mode
func (cli *ZSClient) SetImageBootMode(params param.SetImageBootModeRequest) error {
	return cli.Put("v1/images", params.Uuid, params, nil)
}

// GetUploadImageJobDetails Get Upload Image Job Details
func (cli *ZSClient) GetUploadImageJobDetails(params param.GetUploadImageJobDetailsParam) (*view.GetUploadImageJobDetailsResponse, error) {
	//resp := make([]view.ExistingJobDetails, 0)
	//success := false
	var response view.GetUploadImageJobDetailsResponse
	err := cli.GetWithSpec("v1/images/upload-job/details", params.ImageId, "", "", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload image job details: %v", err)
	}
	return &response, nil
	//return success, cli.GetWithSpec("v1/images/upload-job/details", params.ImageId, "", "success", nil, &success)
}

// GetCandidateVmForAttachingIso Get Candidate VMs for Attaching ISO
func (cli *ZSClient) GetCandidateVmForAttachingIso(uuid string, p *param.QueryParam) ([]view.VmInstanceInventoryView, error) {
	resp := make([]view.VmInstanceInventoryView, 0)
	return resp, cli.List("v1/images/iso/"+uuid+"/vm-candidates", p, &resp)
}
