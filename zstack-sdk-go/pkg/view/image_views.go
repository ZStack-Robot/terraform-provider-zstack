package view

import "zstack.io/zstack-sdk-go/pkg/param"

type ImageView struct {
	BaseInfoView
	BaseTimeView

	State             string                   `json:"state" `       //"state": "Enabled", 镜像的启动状态
	Status            string                   `json:"status" `      // "status": "Ready", 镜像的就绪状态
	Size              int64                    `json:"size" `        //镜像大小
	ActualSize        int64                    `json:"actualSize" `  //镜像真实容量
	Md5Sum            string                   `json:"md5Sum" `      //镜像的md5值
	Url               string                   `json:"url" `         //被添加镜像的URL地址
	MediaType         string                   `json:"mediaType" `   //镜像的类型
	GuestOsType       string                   `json:"guestOsType" ` //镜像对应客户机操作系统的类型
	Type              string                   `json:"type" `
	Platform          string                   `json:"platform" `    //镜像的系统平台,Linux,Windows,WindowsVirtio,Other,Paravirtualization
	Architecture      param.Architecture       `json:"architecture"` //x86_64,aarch64,mips64el
	Format            string                   `json:"format"`       //镜像格式 qcow2
	System            string                   `json:"system" `      //是否系统镜像（如，云路由镜像）
	Virtio            bool                     `json:"virtio"`
	BackupStorageRefs []ImageBackupStorageRefs `json:"backupStorageRefs" `
	SystemTags        []string                 `json:"systemTags"`
}

type ImageBackupStorageRefs struct {
	ImageUuid         string `json:"imageUuid" `         //镜像UUID
	BackupStorageUuid string `json:"backupStorageUuid" ` //镜像存储UUID
	InstallPath       string `json:"installPath" `       //安装路径
	ExportUrl         string `json:"exportUrl" `
	ExportMd5Sum      string `json:"exportMd5Sum" `
	State             string `json:"state" `      // "status": "Ready"
	CreateDate        string `json:"createDate" ` //创建时间
	LastOpDate        string `json:"lastOpDate" ` //最后一次修改时间
}

type GuestOsTypeView struct {
	Platform string         `json:"platform"`
	Children []PlatformView `json:"children"`
}
type PlatformView struct {
	GuestName string        `json:"guestName"`
	Children  []ReleaseView `json:"children"`
}

type ReleaseView struct {
	Uuid      string `json:"uuid"`
	Platform  string `json:"platform"`
	Name      string `json:"name"`
	OsRelease string `json:"osRelease"`
	Version   string `json:"version"`
}

type GetImageQgaView struct {
	Enable bool `json:"enable"`
}

type GetUploadImageJobDetailsResponse struct {
	Success            bool               `json:"success" `
	ExistingJobDetails ExistingJobDetails `json:"existingJobDetails" `
}
type ExistingJobDetails struct {
	LongJobUuid    string `json:"longJobUuid"`
	LongJobState   string `json:"longJobState"`
	ImageUuid      string `json:"imageUuid" `
	ImageUploadUrl string `json:"imageUploadUrl" `
	Offset         int64  `json:"offset"`
}

type GuestOsNameView struct {
	Name string `json:"name"` //操作系统名称
}
