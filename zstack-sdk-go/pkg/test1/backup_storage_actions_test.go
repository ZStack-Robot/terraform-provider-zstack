// # Copyright (c) ZStack.io, Inc.

package test

import (
	"fmt"
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/util/jsonutils"
)

func TestQueryBackupStorage(t *testing.T) {
	accountLoginCli.Login()
	defer accountLoginCli.Logout()
	queryParam := param.NewQueryParam()
	fmt.Println("Query All===============================")
	storage, err := accountLoginCli.QueryBackupStorage(queryParam)
	if err != nil {
		golog.Errorf("TestQueryBackupStorage %v", err)
	}

	fmt.Printf("data : %v", storage)

	//条件查询
	fmt.Println("Query condition===============================")
	queryParam.Set("q", "uuid=ce255d96a94c47009163eb8622dea20a")
	storage, err = accountLoginCli.QueryBackupStorage(queryParam)
	if err != nil {
		golog.Errorf("TestQueryBackupStorage %v", err)
	}
	fmt.Printf("data : %v", storage)
}

func TestExportImageFromBackupStorage(t *testing.T) {
	storage, err := accountLoginCli.ExportImageFromBackupStorage(param.ExportImageFromBackupStorageParam{
		BackupStorageUuid: "7a912545634b4ddc86c40af82c14b452",
		ExportImageFromBackupStorage: param.ExportImageFromBackupStorageDetailParam{
			ImageUuid: "40172a4601b14fcfb822c40a70da3955",
		},
	})
	if err != nil {
		golog.Errorf("TestExportImageFromBackupStorage %v", err)
		return
	}
	fmt.Printf("data : %+v", storage)
}

func TestDeleteExportedImageFromBackupStorage(t *testing.T) {
	err := accountLoginCli.DeleteExportedImageFromBackupStorage(param.DeleteExportedImageFromBackupStorageParam{
		BackupStorageUuid: "6fe7cac7389649ecbc0fc1a554323251",
		ImageUuid:         "f75bef77811c42d4810914a3f0771b76",
	})
	if err != nil {
		golog.Errorf("TestDeleteExportedImageFromBackupStorage %v", err)
		return
	}
	fmt.Printf("data : %+v", "success")
}

func TestPageBackupStorage(t *testing.T) {
	queryParam := param.NewQueryParam()
	queryParam.AddQ("__systemTag__?=" + "onlybackup,allowbackup")
	queryParam.AddQ("type?=" + "ImageStoreBackupStorage")
	queryParam.AddQ("totalCapacity>" + "107374182400")
	queryParam.AddQ("state=Enabled")
	queryParam.AddQ("status=Connected")
	fmt.Println("Query All===============================")
	storage, total, err := accountLoginCli.PageBackupStorage(queryParam)
	if err != nil {
		golog.Errorf("TestQueryBackupStorage %v", err)
	}

	fmt.Printf("data : %v", storage)
	fmt.Printf("total : %v", total)

}

func TestAddImageStoreBackupStorage(t *testing.T) {
	storage, err := accountLoginCli.AddImageStoreBackupStorage(param.AddImageStoreBackupStorageParam{
		BaseParam: param.BaseParam{
			SystemTags: []string{"onlybackup"},
			UserTags:   nil,
			RequestIp:  "",
		},
		Params: param.AddImageStoreBackupStorageDetailParam{
			Hostname:     "172.32.1.183",
			Username:     "root",
			Password:     "zstack@2022",
			SshPort:      22,
			Url:          "/cdpdata",
			Name:         "172.32.1.183-test",
			Description:  "ddd",
			ImportImages: false,
		},
	})
	if err != nil {
		golog.Errorf("TestAddImageStoreBackupStorage %v", err)
		return
	}
	fmt.Printf("data : %+v", storage)
	t2, err := accountLoginCli.AttachBackupStorageToZone("6e8191bfd57745f282f78cb013b732b6", storage.UUID)
	if err != nil {
		golog.Errorf("TestAttachBackupStorageToZone %v", err)
		return
	}
	fmt.Printf("data : %+v", t2)
}

func TestAttachBackupStorageToZone(t *testing.T) {
	storage, err := accountLoginCli.AttachBackupStorageToZone("6e8191bfd57745f282f78cb013b732b6", "baa6bfa820694c2f92af932589cd26f5")
	if err != nil {
		golog.Errorf("TestAttachBackupStorageToZone %v", err)
		return
	}
	fmt.Printf("data : %+v", storage)
}

func TestDeleteBackupStorage(t *testing.T) {
	err := accountLoginCli.DeleteBackupStorage("3a99049f0de543c4bead3c4b25c05456")
	if err != nil {
		golog.Errorf("TestDeleteBackupStorage %v", err)
		return
	}
	fmt.Printf("data : %+v", "success")
}

func TestQueryImageStoreBackupStorage(t *testing.T) {
	queryParam := param.NewQueryParam()
	fmt.Println("Query All===============================")
	storage, err := accountLoginCli.QueryImageStoreBackupStorage(queryParam)
	if err != nil {
		golog.Errorf("TestQueryBackupStorage %v", err)
	}

	golog.Infof("%v", jsonutils.Marshal(storage))
}

func TestPageImageStoreBackupStorage(t *testing.T) {
	queryParam := param.NewQueryParam()
	fmt.Println("Query All===============================")
	storage, total, err := accountLoginCli.PageImageStoreBackupStorage(queryParam)
	if err != nil {
		golog.Errorf("TestQueryBackupStorage %v", err)
	}

	fmt.Printf("data : %v", storage)
	fmt.Printf("total : %v", total)
}

func TestGetImageStoreBackupStorage(t *testing.T) {
	fmt.Println("Query All===============================")
	storage, err := accountLoginCli.GetImageStoreBackupStorage("6e028b03cf1745d58bcc3acf6138bd7e")
	if err != nil {
		golog.Errorf("TestQueryBackupStorage %v", err)
	}

	fmt.Printf("data : %v", storage)
}

func TestReconnectImageStoreBackupStorage(t *testing.T) {
	storage, err := accountLoginCli.ReconnectImageStoreBackupStorage("eafc789c12ef405aaf7168787f9e6db7")
	if err != nil {
		golog.Errorf("TestQueryBackupStorage %v", err)
	}

	fmt.Printf("data : %v", storage)
}

func TestUpdateImageStoreBackupStorage(t *testing.T) {
	name := "f发顺丰"
	storage, err := accountLoginCli.UpdateImageStoreBackupStorage("eafc789c12ef405aaf7168787f9e6db7", param.UpdateImageStoreBackupStorageParam{
		BaseParam: param.BaseParam{},
		UpdateImageStoreBackupStorage: param.UpdateImageStoreBackupStorageDetailParam{
			UUID:        "eafc789c12ef405aaf7168787f9e6db7",
			Username:    nil,
			Password:    nil,
			Hostname:    nil,
			SshPort:     nil,
			Name:        &name,
			Description: &name,
		},
	})
	if err != nil {
		golog.Errorf("TestQueryBackupStorage %v", err)
		return
	}
	golog.Infof("%v", jsonutils.Marshal(storage))
}

func TestReclaimSpaceFromImageStore(t *testing.T) {
	storage, err := accountLoginCli.ReclaimSpaceFromImageStore("eafc789c12ef405aaf7168787f9e6db7")
	if err != nil {
		golog.Errorf("TestQueryBackupStorage %v", err)
		return
	}
	golog.Infof("%v", jsonutils.Marshal(storage))
}

func TestChangeBackupStorageState(t *testing.T) {
	storage, err := accountLoginCli.ChangeBackupStorageState("eafc789c12ef405aaf7168787f9e6db7", param.StateEventEnable)
	if err != nil {
		golog.Errorf("TestQueryBackupStorage %v", err)
		return
	}
	golog.Infof("%v", jsonutils.Marshal(storage))
}
