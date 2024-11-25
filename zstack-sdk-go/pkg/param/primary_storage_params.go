// Copyright (c) ZStack.io, Inc.

package param

type PrimaryStorageType string

const (
	PrimaryStorageLocalStorage PrimaryStorageType = "LocalStorage"
	PrimaryStorageCeph         PrimaryStorageType = "Ceph"
)
