// Copyright (c) ZStack.io, Inc.

package param

type CreateVirtualRouterInstanceParam struct {
	BaseParam
	Params CreateVirtualRouterInstanceDetailParam `json:"params"`
}

type CreateVirtualRouterInstanceDetailParam struct {
	Name                            string   `json:"name" `                            // VM instance name
	VirtualRouterOfferingUuid       string   `json:"virtualRouterOfferingUuid" `       // Instance offering UUID, specifies CPU, memory, etc.
	Description                     string   `json:"description" `                     // Detailed description of the VM instance
	ResourceUuid                    string   `json:"resourceUuid" `                    // Resource UUID, if specified, the VM will use this value as its UUID.
	ZoneUuid                        string   `json:"zoneUuid" `                        // Zone UUID, if specified, the VM will be created in the specified zone.
	ClusterUUID                     string   `json:"clusterUuid" `                     // Cluster UUID, if specified, the VM will be created in the specified cluster, higher priority than zoneUuid.
	HostUuid                        string   `json:"hostUuid" `                        // Host UUID, if specified, the VM will be created on the specified host, higher priority than zoneUuid and clusterUuid.
	PrimaryStorageUuidForRootVolume *string  `json:"primaryStorageUuidForRootVolume" ` // Primary storage UUID, if specified, the root volume will be created on the specified primary storage.
	RootVolumeSystemTags            []string `json:"rootVolumeSystemTags"`
	TagUuids                        []string `json:"tagUuids" ` // List of tag UUIDs
}
