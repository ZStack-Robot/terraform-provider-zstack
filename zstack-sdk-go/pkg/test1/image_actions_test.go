// # Copyright (c) ZStack.io, Inc.

package test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/util/httputils"
	"zstack.io/zstack-sdk-go/pkg/util/multipart"
	"zstack.io/zstack-sdk-go/pkg/util/progress"
)

func TestPageImage(t *testing.T) {
	params := param.NewQueryParam()
	params.Start(0).Limit(2)
	result, total, err := accountLoginCli.PageImage(params)
	if err != nil {
		golog.Errorf("ZSClient.QueryImage error:%v", err)
		return
	}
	golog.Infof("======================================")
	for _, r := range result {
		golog.Infof("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s %d", r.UUID, r.Name, r.Platform, r.GuestOsType, r.Format, r.Status, r.Size, r.Description, total)
	}
}

func TestZSClient_QueryImage3(t *testing.T) {
	params := param.NewQueryParam()
	//params.AddQ("backupStorage.zone.uuid=6e8191bfd57745f282f78cb013b732b6")
	params.AddQ("format!=vmtx")
	params.AddQ("system=false")
	//	params.AddQ("backupStorageRefs.exportUrl!=null")
	//params.AddQ("backupStorage.zone.uuid=6e8191bfd57745f282f78cb013b732b6")
	result, err := accountLoginCli.QueryImage(params)
	if err != nil {
		golog.Errorf("ZSClient.QueryImage error:%v", err)
		return
	}
	golog.Infof("======================================")
	for _, r := range result {
		golog.Infof("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s", r.UUID, r.Name, r.Platform, r.GuestOsType, r.Format, r.Status, r.Size, r.Description)
	}
	golog.Infof("======================================")
}

func TestZSClient_QueryImage2(t *testing.T) {
	params := param.NewQueryParam()
	//params.AddQ("state=Enabled")
	//params.AddQ("type=zstack")
	//params.AddQ("format!=vmtx")
	//params.AddQ("status=Ready")
	//params.AddQ("system=false")
	params.AddQ("name=fd")
	//params.AddQ("mediaType=DataVolumeTemplate")
	//params.AddQ("backupStorage.zone.uuid=6e8191bfd57745f282f78cb013b732b6")
	result, err := accountLoginCli.QueryImage(params)
	if err != nil {
		golog.Errorf("ZSClient.QueryImage error:%v", err)
		return
	}
	golog.Infof("======================================")
	for _, r := range result {
		golog.Infof("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s", r.UUID, r.Name, r.Platform, r.GuestOsType, r.Format, r.Status, r.Size, r.Description)
	}
	golog.Infof("======================================")
}

func TestZSClient_QueryImage(t *testing.T) {
	//查询所有
	params := param.NewQueryParam()
	result, err := accountLoginCli.QueryImage(params)
	if err != nil {
		golog.Errorf("ZSClient.QueryImage error:%v", err)
		return
	}
	golog.Infof("======================================")
	for _, r := range result {
		golog.Infof("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s", r.UUID, r.Name, r.Platform, r.GuestOsType, r.Format, r.Status, r.Size, r.Description)
	}
	golog.Infof("======================================")

	//带过滤条件查询
	params.AddQ(fmt.Sprintf("platform=%s", "Windows"))
	params.Start(0).Limit(2).Count(true)
	result, err = accountLoginCli.QueryImage(params)
	if err != nil {
		golog.Errorf("ZSClient.QueryImage error:%v", err)
		return
	}
	golog.Infof("======================================")
	for _, r := range result {
		golog.Infof("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s", r.UUID, r.Name, r.Platform, r.GuestOsType, r.Format, r.Status, r.Size, r.Description)
	}
}

type Image struct {
	Name        string `json:"name"`
	UUID        string `json:"uuid"`
	Format      string `json:"format"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Status      string `json:"status"`
	State       string `json:"state"`
}

type QueryResult struct {
	Results []struct {
		Inventories []Image `json:"inventories"`
	} `json:"results"`
}

func TestZSCliebt_QueryByZql(t *testing.T) {

	//var reservedIpRanges []view.ReservedIpRangeInventoryView
	var queryResult QueryResult
	//	_, err := accountLoginCli.Zql(fmt.Sprintf("query Image "), &virtualRouterImages, "inventories")

	_, err := accountLoginCli.Zql(
		"query Image where __systemTag__='applianceType::vrouter'",
		&queryResult,
		"results", "inventories",
	)
	if err != nil {
		golog.Errorf("failed to execute ZQL query: %w", err)
	}

	//_, err := accountLoginCli.Zql(fmt.Sprintf("query Image"), &virtualRouterImages, "inventories")
	// 提取结果
	if len(queryResult.Results) > 0 {
		inventories := queryResult.Results[0].Inventories
		fmt.Printf("Query Response: %+v\n", inventories)
	} else {
		fmt.Println("No inventories found.")
	}

}

func TestZSClient_GetImage(t *testing.T) {
	r, err := accountLoginCli.GetImage("c7ad6441c19c4a4c8084ddb0c13200f2")
	if err != nil {
		golog.Errorf("ZSClient.GetImage error:%v", err)
		return
	}
	golog.Infof("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s", r.UUID, r.Name, r.Platform, r.GuestOsType, r.Format, r.Status, r.Size, r.Description)
}

func TestZSClient_CreateImage(t *testing.T) {
	storage, err := accountLoginCli.QueryBackupStorage(param.QueryParam{})
	if err != nil {
		golog.Errorf("QueryBackupStorage error:%v", err)
	}

	imageParam := param.AddImageParam{
		BaseParam: param.BaseParam{
			SystemTags: []string{"bootMode::Legacy"},
		},
		Params: param.AddImageDetailParam{
			Name:               "CentOS-6.8-i386-LiveCD",
			Description:        "接口image",
			Url:                "http://172.20.15.213/rds/V3.14.1-p2/zstack-rds-3.14.1-p2_x86.qcow2",
			MediaType:          param.RootVolumeTemplate,
			GuestOsType:        "Linux",
			System:             false,
			Format:             param.Qcow2,
			Platform:           "Linux",
			BackupStorageUuids: []string{storage[0].UUID},
			Type:               "",
			ResourceUuid:       "",
			Architecture:       "x86_64",
			Virtio:             false,
		},
	}

	var apiId string
	//_, err = accountLoginCli.AddImage(imageParam)
	_, err = accountLoginCli.AddImage(imageParam)
	if err != nil {
		golog.Errorf("ZSClient.CreateImage error:%v", err)
		return
	}

	//golog.Infof("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s", r.UUID, r.Name, r.Platform, r.GuestOsType, r.Format, r.Status, r.Size, r.Description)

	golog.Infof("======================================")
	golog.Infof(apiId)
	//创建失败情况
	/*
		imageParam = param.AddImageParam{
			BaseParam: param.BaseParam{
				SystemTags: []string{"bootMode::Legacy"},
			},
			Params: param.AddImageDetailParam{
				Name:               "image-4",
				Description:        "接口image",
				Url:                "https://image.baidu.com/search/detail?tn=baiduimagedetail&word=%E6%B8%90%E5%8F%98%E9%A3%8E%E6%A0%BC%E6%8F%92%E7%94%BB&album_tab=%E8%AE%BE%E8%AE%A1%E7%B4%A0%E6%9D%90&album_id=409&ie=utf-8&fr=albumsdetail&cs=4036010509,3445021118&pi=144521&pn=1&ic=0&objurl=https%3A%2F%2Ft7.baidu.com%2Fit%2Fu%3D4036010509%2C3445021118%26fm%3D193%26f%3DGIF",
				MediaType:          param.RootVolumeTemplate,
				GuestOsType:        "Windows 10",
				System:             false,
				Format:             param.Qcow2,
				Platform:           "Windows",
				BackupStorageUuids: []string{"26684790e4734a0bbb506f40907f57da"},
				Type:               "",
				ResourceUuid:       "",
				Architecture:       "x86_64",
				Virtio:             false,
			},
		}

		_, err = accountLoginCli.AddImage(imageParam)
		if err != nil {
			golog.Errorf("ZSClient.CreateImage error:%v", err)
		}

		golog.Infof("======================================")
		//删除
		err = accountLoginCli.DeleteImage(r.UUID, param.DeleteModeEnforcing)
		if err != nil {
			golog.Errorf("ZSClient.DeleteImage error:%v", err)
		}

		//彻底删除
		err = accountLoginCli.ExpungeImage(r.UUID)
		if err != nil {
			golog.Errorf("ZSClient.ExpungeImage error:%v", err)
		}*/
}

func TestZSClient_CreateImage2(t *testing.T) {
	storage, err := accountLoginCli.QueryBackupStorage(param.QueryParam{})
	if err != nil {
		golog.Errorf("QueryBackupStorage error:%v", err)
	}

	imageName := "image-chenjh001"
	imageParam := param.AddImageParam{
		BaseParam: param.BaseParam{
			//SystemTags: []string{"bootMode::Legacy"},
		},
		Params: param.AddImageDetailParam{
			Name:               imageName,
			Description:        "接口image",
			Url:                fmt.Sprintf("upload://%s", imageName),
			MediaType:          param.RootVolumeTemplate,
			GuestOsType:        "Windows 10",
			System:             false,
			Format:             param.Qcow2,
			Platform:           "Windows",
			BackupStorageUuids: []string{storage[0].UUID},
			//Type:               "",
			//ResourceUuid:       "",
			SystemTags:   []string{"qemuga", "bootMode::Legacy"},
			Architecture: "x86_64",
			Virtio:       false,
		},
	}

	r, err := accountLoginCli.AddImage(imageParam)
	if err != nil {
		golog.Errorf("ZSClient.CreateImage error:%v", err)
		return
	}
	golog.Infof("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s", r.UUID, r.Name, r.Platform, r.GuestOsType, r.Format, r.Status, r.Size, r.Description)
	file, err := ioutil.ReadFile("D:\\a6cc8eeeb9cc49bcba635900be7264e3")
	if err != nil {
		return
	}
	// fmt.Println(file)
	newReader := bytes.NewReader(file)
	// newReader := strings.NewReader("11111111111111111111D:\\image-d1f501b3887a6a084feb66d0a995215731f665555554e4.qcow2")
	body := multipart.NewReader(newReader, "", imageName)
	header := http.Header{}
	header.Add("X-IMAGE-UUID", r.UUID)
	header.Add("X-IMAGE-SIZE", fmt.Sprintf("%d", newReader.Size()))
	header.Add("Content-Type", body.FormDataContentType())
	//r := multicloud.NewProgress(size, 99, body, callback)
	/*r := progress.NewProgress(reader.Size(), 99, body, func(progress float32) {

	})*/
	re := progress.NewProgress(newReader.Size(), 99, body, func(progress float32) {
		golog.Infof("progress:%f", progress)
	})

	url := r.BackupStorageRefs[0].InstallPath //"http://172.20.20.132:8001/imagestore/upload"
	resp, err := httputils.Request(httputils.GetTimeoutClient(0), context.Background(), "POST", url, header, re, true)
	if err != nil {
		golog.Error(err)
	}
	defer resp.Body.Close()
}

func TestZSClient_UpdateImage(t *testing.T) {
	desc := "test test"
	params := param.UpdateImageParam{
		UpdateImage: param.UpdateImageDetailParam{
			Name:        "image4chenjtest",
			Description: &desc,
		},
	}
	type args struct {
		uuid   string
		params param.UpdateImageParam
	}
	tests := []struct {
		name    string
		cli     *client.ZSClient
		args    args
		wantErr bool
	}{
		{accountLogin, accountLoginCli, args{imageID[accountLogin], params}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.UpdateImage("c1c74f05841d47fd8438f0641c99c682", tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.UpdateImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("ZSClient.UpdateImage() = %v", got)
		})
	}
}

/*func TestDeleteImageTest(t *testing.T) { http://172.20.20.132:8001/imagestore/upload
	image, err := accountLoginCli.ExpungeImage("85d6eff243f14df69029b0b21945f2b0", param.ExpungeImageParam{
		BackupStorageUuids: []string{"85d6eff243f14df69029b0b21945f2b0"},
	})
	fmt.Printf("%v \n %v", image, err)
}*/

func TestZSClient_GuestOsTypeList(t *testing.T) {
	data, err := accountLoginCli.GuestOsTypeList()
	if err != nil {
		golog.Errorf("ZSClient.GuestOsTypeList error:%v", err)
		return
	}
	fmt.Printf("%v\n", data)
}

func TestGuestNameList(t *testing.T) {
	data, err := accountLoginCli.GuestNameList()
	if err != nil {
		golog.Errorf("ZSClient.GuestOsTypeList error:%v", err)
		return
	}
	fmt.Printf("%v\n", data)
}

func TestZSClient_UpdatePlatform(t *testing.T) {
	platform, err := accountLoginCli.UpdateImagePlatform(param.UpdateImagePlatformParam{
		BaseParam: param.BaseParam{},
		UpdateImage: param.UpdateImagePlatformDetailParam{
			Platform:    "Linux",
			GuestOsType: "CentOS 6",
			UUID:        "e77a9d5431124157bed5935d6165c78d",
		},
	})
	if err != nil {
		golog.Errorf("TestZSClient_UpdatePlatform error:%v", err)
	}
	golog.Println(platform)
}

func TestZSClient_UpdateVirtio(t *testing.T) {
	virtio, err := accountLoginCli.UpdateImageVirtio(param.UpdateImageVirtioParam{
		BaseParam: param.BaseParam{},
		UpdateImage: param.UpdateImageVirtioDetailParam{
			Virtio: false,
			UUID:   "e77a9d5431124157bed5935d6165c78d",
		},
	})
	if err != nil {
		golog.Errorf("TestZSClient_UpdateVirtio error:%v", err)
	}
	fmt.Println(virtio)
}

func TestChangeImageState(t *testing.T) {
	state, err := accountLoginCli.ChangeImageState(param.ChangeImageStateParam{
		BaseParam: param.BaseParam{},
		ImageUuid: "e77a9d5431124157bed5935d6165c78d",
		ChangeImageState: param.ChangeImageStateDetailParam{
			StateEvent: param.StateEventEnable,
		},
	})
	if err != nil {
		golog.Errorf("TestZSClient_UpdateVirtio error:%v", err)
	}
	fmt.Println(state)
}

func TestSyncImageSize(t *testing.T) {
	size, err := accountLoginCli.SyncImageSize(param.SyncImageSizeParam{
		BaseParam:     param.BaseParam{},
		ImageUuid:     "4c9d3864089f4611be17faff64b7f965",
		SyncImageSize: param.SyncImageSizeDetailParam{},
	})

	if err != nil {
		golog.Errorf("TestSyncImageSize error:%v", err)
		return
	}
	fmt.Println(size)
}

func TestGetCandidateBackupStorageForCreatingImage(t *testing.T) {
	data, err := accountLoginCli.GetCandidateBackupStorageForCreatingImage(param.GetCandidateBackupStorageForCreatingImageParam{
		BaseParam:                  param.BaseParam{},
		CandidateBackupStorageType: param.CandidateBackupStorageTypeDefault,
		VolumeUuid:                 "",
		VolumeSnapshotUuid:         "",
	})
	if err != nil {
		golog.Errorf("TestGetCandidateBackupStorageForCreatingImage error:%v", err)
	}
	fmt.Println(data)
}

func TestCreateRootVolumeTemplateFromRootVolume(t *testing.T) {
	data, err := accountLoginCli.CreateRootVolumeTemplateFromRootVolume(param.CreateRootVolumeTemplateFromRootVolumeParam{
		BaseParam:      param.BaseParam{},
		RootVolumeUuid: "153796cecf784546879411e9a35b6f3d",
		Params: param.CreateRootVolumeTemplateFromRootVolumeDetailParam{
			Name:               "root-data-image",
			Description:        "CreateRootVolumeTemplateFromRootVolume",
			BackupStorageUuids: []string{"6fe7cac7389649ecbc0fc1a554323251"},
		},
	})
	if err != nil {
		golog.Errorf("TestCreateRootVolumeTemplateFromRootVolume error:%v", err)
	}
	fmt.Println(data)
}

func TestCreateDataVolumeTemplateFromVolume(t *testing.T) {
	data, err := accountLoginCli.CreateDataVolumeTemplateFromVolume(param.CreateDataVolumeTemplateFromVolumeParam{
		BaseParam:  param.BaseParam{},
		VolumeUuid: "153796cecf784546879411e9a35b6f3d",
		Params: param.CreateDataVolumeTemplateFromVolumeDetailParam{
			Name:               "data-data-image",
			Description:        "CreateDataVolumeTemplateFromVolume",
			BackupStorageUuids: []string{"6fe7cac7389649ecbc0fc1a554323251"},
		},
	})
	if err != nil {
		golog.Errorf("TestCreateRootVolumeTemplateFromRootVolume error:%v", err)
	}
	fmt.Println(data)
}

func TestCreateRootVolumeTemplateFromVolumeSnapshot(t *testing.T) {
	data, err := accountLoginCli.CreateRootVolumeTemplateFromVolumeSnapshot(param.CreateRootVolumeTemplateFromVolumeSnapshotParam{
		BaseParam:    param.BaseParam{},
		SnapshotUuid: "78d67e786dbb4863a2b532db563f47ca",
		Params: param.CreateRootVolumeTemplateFromVolumeSnapshotDetailParams{
			Name:               "CreateRootVolumeTemplateFromVolumeSnapshotDetailParams",
			Description:        "CreateRootVolumeTemplateFromVolumeSnapshotDetailParams",
			GuestOsType:        "Windows 10",
			BackupStorageUuids: []string{"26684790e4734a0bbb506f40907f57da"},
			Platform:           "Windows",
			System:             false,
			ResourceUuid:       "",
			Architecture:       "x86_64",
			TagUuids:           nil,
		},
	})
	if err != nil {
		golog.Errorf("TestCreateRootVolumeTemplateFromVolumeSnapshot error:%v", err)
	}
	fmt.Println(data)
}

func TestGetImageQga(t *testing.T) {
	data, err := accountLoginCli.GetImageQga("c7ad6441c19c4a4c8084ddb0c13200f2")
	if err != nil {
		golog.Errorf("TestGetImageQga error:%v", err)
	}
	fmt.Println(data)
}

func TestSetImageQga(t *testing.T) {
	var data, err = accountLoginCli.SetImageQga(param.SetImageQgaParam{
		BaseParam: param.BaseParam{},
		Uuid:      "c7ad6441c19c4a4c8084ddb0c13200f2",
		SetImageQga: param.SetImageQgaDetailParam{
			Enable: false,
		},
	})
	if err != nil {
		golog.Errorf("TestSetImageQga error:%v", err)
	}
	fmt.Println(data)
}

func TestSetImageBootMode(t *testing.T) {
	err := accountLoginCli.SetImageBootMode(param.SetImageBootModeRequest{
		BaseParam: param.BaseParam{},
		Uuid:      "968e87334a12422fbe78c8b72bcfab68",
		SetImageBootMode: param.SetImageBootModeParams{
			BootMode: param.Legacy,
		},
	})
	if err != nil {
		golog.Errorf("TestSetImageBootMode error:%v", err)
	}
}

func TestGetUploadImageJobDetails(t *testing.T) {
	response, err := accountLoginCli.GetUploadImageJobDetails(param.GetUploadImageJobDetailsParam{
		BaseParam: param.BaseParam{},
		ImageId:   "6718db9f407f4c58b5a7bbae689080bb",
	})

	if err != nil {
		golog.Errorf("TestGetUploadImageJobDetails error:%v", err)
	}

	if response == nil {
		t.Fatal("TestGetUploadImageJobDetails failed: response is nil") // 确保 response 不为 nil
	}

	if response.Success {
		jobDetail := response.ExistingJobDetails
		fmt.Printf("Job UUID: %s, State: %s, Image UUID: %s, Upload URL: %s, Offset: %d\n",
			jobDetail.LongJobUuid, jobDetail.LongJobState, jobDetail.ImageUuid, jobDetail.ImageUploadUrl, jobDetail.Offset)
	} else {
		fmt.Println("Failed to retrieve job details.")
	}
}
