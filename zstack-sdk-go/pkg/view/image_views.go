// Copyright (c) ZStack.io, Inc.

package view

import "zstack.io/zstack-sdk-go/pkg/param"

type ImageView struct {
	BaseInfoView
	BaseTimeView

	State             string                   `json:"state"`       // "state": "Enabled", Image boot state
	Status            string                   `json:"status"`      // "status": "Ready", Image ready state
	Size              int64                    `json:"size"`        // Image size
	ActualSize        int64                    `json:"actualSize"`  // Image actual capacity
	Md5Sum            string                   `json:"md5Sum"`      // Image md5 value
	Url               string                   `json:"url"`         // URL of the added image
	MediaType         string                   `json:"mediaType"`   // Image type
	GuestOsType       string                   `json:"guestOsType"` // Guest OS type corresponding to the image
	Type              string                   `json:"type"`
	Platform          string                   `json:"platform"`     // Image system platform, Linux, Windows, WindowsVirtio, Other, Paravirtualization
	Architecture      param.Architecture       `json:"architecture"` // x86_64, aarch64, mips64el
	Format            string                   `json:"format"`       // Image format qcow2
	System            string                   `json:"system"`       // Whether it is a system image (e.g., cloud router image)
	Virtio            bool                     `json:"virtio"`
	BackupStorageRefs []ImageBackupStorageRefs `json:"backupStorageRefs"`
	SystemTags        []string                 `json:"systemTags"`
}

type ImageBackupStorageRefs struct {
	ImageUuid         string `json:"imageUuid"`         // Image UUID
	BackupStorageUuid string `json:"backupStorageUuid"` // Image storage UUID
	InstallPath       string `json:"installPath"`       // Installation path
	ExportUrl         string `json:"exportUrl"`
	ExportMd5Sum      string `json:"exportMd5Sum"`
	State             string `json:"state"`      // "status": "Ready"
	CreateDate        string `json:"createDate"` // Creation time
	LastOpDate        string `json:"lastOpDate"` // Last modification time
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
	Success            bool               `json:"success"`
	ExistingJobDetails ExistingJobDetails `json:"existingJobDetails"`
}

type ExistingJobDetails struct {
	LongJobUuid    string `json:"longJobUuid"`
	LongJobState   string `json:"longJobState"`
	ImageUuid      string `json:"imageUuid"`
	ImageUploadUrl string `json:"imageUploadUrl"`
	Offset         int64  `json:"offset"`
}

type GuestOsNameView struct {
	Name string `json:"name"` // Operating system name
}
