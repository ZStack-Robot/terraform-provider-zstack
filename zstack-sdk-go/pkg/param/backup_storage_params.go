// Copyright (c) ZStack.io, Inc.

package param

type BackupStorageType string

const (
	ImageStoreBackupStorage BackupStorageType = "ImageStoreBackupStorage"
)

type ExportImageFromBackupStorageParam struct {
	BaseParam
	BackupStorageUuid            string                                  `json:"backupStorageUuid"`
	ExportImageFromBackupStorage ExportImageFromBackupStorageDetailParam `json:"exportImageFromBackupStorage"`
}

type ExportImageFromBackupStorageDetailParam struct {
	ImageUuid string `json:"imageUuid"`
}

type DeleteExportedImageFromBackupStorageParam struct {
	BackupStorageUuid string `json:"backupStorageUuid"`
	ImageUuid         string `json:"imageUuid"`
}

type AddImageStoreBackupStorageParam struct {
	BaseParam
	Params AddImageStoreBackupStorageDetailParam `json:"params"`
}

type AddImageStoreBackupStorageDetailParam struct {
	Hostname     string            `json:"hostname"`
	Username     string            `json:"username"`
	Password     string            `json:"password"`
	SshPort      int               `json:"sshPort"`
	Url          string            `json:"url"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Type         BackupStorageType `json:"type"`
	ImportImages bool              `json:"importImages"`
	ResourceUuid string            `json:"resourceUuid"`
}

type UpdateImageStoreBackupStorageParam struct {
	BaseParam
	UpdateImageStoreBackupStorage UpdateImageStoreBackupStorageDetailParam `json:"updateImageStoreBackupStorage"`
}

type UpdateImageStoreBackupStorageDetailParam struct {
	UUID        string  `json:"uuid"`
	Username    *string `json:"username"`
	Password    *string `json:"password"`
	Hostname    *string `json:"hostname"`
	SshPort     *int    `json:"sshPort"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type ChangeBackupStorageStateParam struct {
	BaseParam
	ChangeBackupStorageState ChangeBackupStorageStateDetailParam `json:"changeBackupStorageState"`
}

type ChangeBackupStorageStateDetailParam struct {
	StateEvent StateEvent `json:"stateEvent"`
}
