// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

// FieldMapping
var FieldMapping = map[string]map[string]string{
	"backup_storage": {
		"total_capacity":     "totalCapacity",
		"available_capacity": "availableCapacity",
		"state":              "state",
		"status":             "status",
	},

	"instance": {
		"cluster_uuid": "clusterUuid",
		//"cpu_num":         "cpuNum",
		"cpu_num":         "CPUNum",
		"host_uuid":       "hostUuid",
		"hypervisor_type": "hypervisorType",
		"image_uuid":      "imageUuid",
		"memory_size":     "memorySize",
		"zone_uuid":       "zoneUuid",
	},
	"cluster": {
		"hypervisor_type": "hypervisorType",
		"zone_uuid":       "zoneUuid",
	},
	"host": {
		"cluster_uuid": "clusterUuid",
		"managementip": "managementIp",
		"zone_uuid":    "zoneUuid",
	},
	"disk_offer": {
		"disk_size":          "diskSize",
		"allocator_strategy": "allocatorStrategy",
	},
	"image": {
		"guest_os_type":        "guestOsType",
		"image_format":         "imageFormat",
		"image_type":           "imageType",
		"media_type":           "mediaType",
		"boot_mode":            "bootMode",
		"backup_storage_uuids": "backupStorageUuids",
	},
	"instance_offer": {
		"cpu_num":            "cpuNum",
		"cpu_speed":          "cpuSpeed",
		"memory_size":        "memorySize",
		"allocator_strategy": "allocatorStrategy",
	},
	"l2network": {
		"physical_interface": "physicalInterface",
		"zone_uuid":          "zoneUuid",
	},
	"l3network": {},
	"vip": {
		"l3_network_uuid":       "l3NetworkUuid",
		"peer_l3_network_uuids": "peerL3NetworkUuids",
		"use_for":               "useFor",
	},
	"virtual_router_offer": {
		"allocator_strategy":      "allocatorStrategy",
		"cpu_num":                 "cpuNum",
		"image_uuid":              "imageUuid",
		"is_default":              "isDefault",
		"management_network_uuid": "managementNetworkUuid",
		"memory_size":             "memorySize",
		"public_network_uuid":     "publicNetworkUuid",
		"zone_uuid":               "zoneUuid",
	},
	"virtual_router_instance": {
		"agent_port":        "agentPort",
		"appliance_vm_type": "applianceVmType",
		"cluster_uuid":      "clusterUuid",
		//"cpu_num":                 "cpuNum",
		"cpu_num":                 "CPUNum",
		"ha_status":               "haStatus",
		"host_uuid":               "hostUuid",
		"hypervisor_type":         "hypervisorType",
		"image_uuid":              "imageUuid",
		"instance_offering_uuid":  "instanceOfferingUuid",
		"management_network_uuid": "managementNetworkUuid",
		"memory_size":             "memorySize",
	},
	"virtual_router_image": {
		"guest_os_type": "guestOsType",
	},
	"zone": {},
	"primary_storage": {
		"total_capacity":              "totalCapacity",
		"available_capacity":          "availableCapacity",
		"total_physical_capacity":     "totalPhysicalCapacity",
		"available_physical_capacity": "availablePhysicalCapacity",
		"system_used_capacity":        "systemUsedCapacity",
	},
	"disks": {
		"disk_offering_uuid":   "diskOfferingUuid",
		"is_shareable":         "isShareable",
		"primary_storage_uuid": "primaryStorageUuid",
		"vm_instance_uuid":     "vmInstanceUuid",
	},
	"qga": {
		"instance_uuid":       "instanceUuid",
		"guest_tools_version": "guestToolsVersion",
		"guest_tools_status":  "guestToolsStatus",
	},
	"script": {
		"script_type":    "scriptType",
		"script_content": "scriptContent",
		"render_params":  "renderParams",
		"script_timeout": "scriptTimeout",
	},
	"user_tags": {
		"resource_type": "resourceType",
		"resource_uuid": "resourceUuid",
	},
	"system_tags": {
		"resource_type": "resourceType",
		"resource_uuid": "resourceUuid",
	},
	"tag": {},
	"security_group": {
		"ip_version":                "ipVersion",
		"src_ip_range":              "srcIpRange",
		"dst_ip_range":              "dstIpRange",
		"attached_l3_network_uuids": "attachedL3NetworkUuids",
	},
	"security_group_rule": {
		"ip_version":                 "ipVersion",
		"src_ip_range":               "srcIpRange",
		"dst_ip_range":               "dstIpRange",
		"remote_security_group_uuid": "remoteSecurityGroupUuid",
		"security_group_uuid":        "securityGroupUuid",
		"dst_port_range":             "dstPortRange",
	},
	"sdn_controller": {
		"vendor_type": "vendorType",
	},
}

// GetFieldMapping
func GetFieldMapping(dataSourceName string) map[string]string {
	return FieldMapping[dataSourceName]
}
