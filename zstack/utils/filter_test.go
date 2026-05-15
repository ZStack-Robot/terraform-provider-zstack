// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type filterTestBase struct {
	UUID string `json:"uuid"`
}

type filterTestAffinityGroup struct {
	filterTestBase
	Name     string `json:"name,omitempty"`
	ZoneUuid string `json:"zoneUuid,omitempty"`
}

type filterTestInstance struct {
	CpuNum int `json:"cpuNum,omitempty"`
}

type filterTestVolume struct {
	DiskOfferingUuid   string `json:"diskOfferingUuid,omitempty"`
	PrimaryStorageUuid string `json:"primaryStorageUuid,omitempty"`
	VmInstanceUuid     string `json:"vmInstanceUuid,omitempty"`
	LastVmInstanceUuid string `json:"lastVmInstanceUuid,omitempty"`
}

type filterTestVolumeSnapshot struct {
	ParentUuid         string `json:"parentUuid,omitempty"`
	PrimaryStorageUuid string `json:"primaryStorageUuid,omitempty"`
	TreeUuid           string `json:"treeUuid,omitempty"`
	VolumeUuid         string `json:"volumeUuid,omitempty"`
}

type filterTestVIP struct {
	UUID               string   `json:"uuid,omitempty"`
	PeerL3NetworkUuids []string `json:"peerL3NetworkUuids,omitempty"`
}

type filterTestSecurityGroup struct {
	UUID                   string                        `json:"uuid,omitempty"`
	AttachedL3NetworkUuids []string                      `json:"attachedL3NetworkUuids,omitempty"`
	Rules                  []filterTestSecurityGroupRule `json:"rules,omitempty"`
}

type filterTestSecurityGroupRule struct {
	SrcIpRange string `json:"srcIpRange,omitempty"`
	DstIpRange string `json:"dstIpRange,omitempty"`
}

type filterTestImage struct {
	UUID              string                            `json:"uuid,omitempty"`
	Format            string                            `json:"format,omitempty"`
	Type              string                            `json:"type,omitempty"`
	BackupStorageRefs []filterTestImageBackupStorageRef `json:"backupStorageRefs,omitempty"`
}

type filterTestImageBackupStorageRef struct {
	BackupStorageUuid string `json:"backupStorageUuid,omitempty"`
}

type filterTestEIP struct {
	VmNicUuid string `json:"vmNicUuid,omitempty"`
}

type filterTestTerraformModel struct {
	ZoneUuid types.String `tfsdk:"zone_uuid"`
}

func TestFilterResourceResolvesMappedCamelCaseFieldByJSONTag(t *testing.T) {
	resources := []filterTestAffinityGroup{
		{filterTestBase: filterTestBase{UUID: "ag-1"}, Name: "ag-a", ZoneUuid: "zone-a"},
		{filterTestBase: filterTestBase{UUID: "ag-2"}, Name: "ag-b", ZoneUuid: "zone-b"},
	}

	filtered, diags := FilterResource(context.Background(), resources, map[string][]string{
		"zone_uuid": {"zone-a"},
	}, "affinity_group")
	if diags.HasError() {
		t.Fatalf("FilterResource returned diagnostics: %v", diags)
	}

	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered resource, got %d", len(filtered))
	}
	if filtered[0].UUID != "ag-1" {
		t.Fatalf("expected ag-1, got %s", filtered[0].UUID)
	}
}

func TestFilterResourceResolvesEmbeddedFieldByJSONTag(t *testing.T) {
	resources := []filterTestAffinityGroup{
		{filterTestBase: filterTestBase{UUID: "ag-1"}, Name: "ag-a", ZoneUuid: "zone-a"},
		{filterTestBase: filterTestBase{UUID: "ag-2"}, Name: "ag-b", ZoneUuid: "zone-b"},
	}

	filtered, diags := FilterResource(context.Background(), resources, map[string][]string{
		"uuid": {"ag-2"},
	}, "affinity_group")
	if diags.HasError() {
		t.Fatalf("FilterResource returned diagnostics: %v", diags)
	}

	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered resource, got %d", len(filtered))
	}
	if filtered[0].UUID != "ag-2" {
		t.Fatalf("expected ag-2, got %s", filtered[0].UUID)
	}
}

func TestFilterResourceResolvesInstanceCpuNumByJSONTag(t *testing.T) {
	resources := []filterTestInstance{
		{CpuNum: 2},
		{CpuNum: 4},
	}

	filtered, diags := FilterResource(context.Background(), resources, map[string][]string{
		"cpu_num": {"4"},
	}, "instance")
	if diags.HasError() {
		t.Fatalf("FilterResource returned diagnostics: %v", diags)
	}

	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered resource, got %d", len(filtered))
	}
	if filtered[0].CpuNum != 4 {
		t.Fatalf("expected cpu_num 4, got %d", filtered[0].CpuNum)
	}
}

func TestFilterResourceResolvesVolumeFieldsByJSONTag(t *testing.T) {
	resources := []filterTestVolume{
		{DiskOfferingUuid: "disk-1", PrimaryStorageUuid: "ps-1", VmInstanceUuid: "vm-1", LastVmInstanceUuid: "last-vm-1"},
		{DiskOfferingUuid: "disk-2", PrimaryStorageUuid: "ps-2", VmInstanceUuid: "vm-2", LastVmInstanceUuid: "last-vm-2"},
	}

	tests := []struct {
		name      string
		filterKey string
		value     string
		wantVM    string
	}{
		{name: "disk offering uuid", filterKey: "disk_offering_uuid", value: "disk-2", wantVM: "vm-2"},
		{name: "primary storage uuid", filterKey: "primary_storage_uuid", value: "ps-2", wantVM: "vm-2"},
		{name: "vm instance uuid", filterKey: "vm_instance_uuid", value: "vm-2", wantVM: "vm-2"},
		{name: "last vm instance uuid", filterKey: "last_vm_instance_uuid", value: "last-vm-2", wantVM: "vm-2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered, diags := FilterResource(context.Background(), resources, map[string][]string{
				tt.filterKey: {tt.value},
			}, "volume")
			if diags.HasError() {
				t.Fatalf("FilterResource returned diagnostics: %v", diags)
			}

			if len(filtered) != 1 {
				t.Fatalf("expected 1 filtered resource, got %d", len(filtered))
			}
			if filtered[0].VmInstanceUuid != tt.wantVM {
				t.Fatalf("expected %s, got %s", tt.wantVM, filtered[0].VmInstanceUuid)
			}
		})
	}
}

func TestFilterResourceResolvesVolumeSnapshotFieldsByJSONTag(t *testing.T) {
	resources := []filterTestVolumeSnapshot{
		{ParentUuid: "parent-1", PrimaryStorageUuid: "ps-1", TreeUuid: "tree-1", VolumeUuid: "vol-1"},
		{ParentUuid: "parent-2", PrimaryStorageUuid: "ps-2", TreeUuid: "tree-2", VolumeUuid: "vol-2"},
	}

	tests := []struct {
		name      string
		filterKey string
		value     string
		wantVol   string
	}{
		{name: "parent uuid", filterKey: "parent_uuid", value: "parent-2", wantVol: "vol-2"},
		{name: "primary storage uuid", filterKey: "primary_storage_uuid", value: "ps-2", wantVol: "vol-2"},
		{name: "tree uuid", filterKey: "tree_uuid", value: "tree-2", wantVol: "vol-2"},
		{name: "volume uuid", filterKey: "volume_uuid", value: "vol-2", wantVol: "vol-2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered, diags := FilterResource(context.Background(), resources, map[string][]string{
				tt.filterKey: {tt.value},
			}, "volume_snapshot")
			if diags.HasError() {
				t.Fatalf("FilterResource returned diagnostics: %v", diags)
			}

			if len(filtered) != 1 {
				t.Fatalf("expected 1 filtered resource, got %d", len(filtered))
			}
			if filtered[0].VolumeUuid != tt.wantVol {
				t.Fatalf("expected %s, got %s", tt.wantVol, filtered[0].VolumeUuid)
			}
		})
	}
}

func TestFilterResourceResolvesStringSliceFieldByJSONTag(t *testing.T) {
	resources := []filterTestVIP{
		{UUID: "vip-1", PeerL3NetworkUuids: []string{"l3-a", "l3-b"}},
		{UUID: "vip-2", PeerL3NetworkUuids: []string{"l3-c"}},
	}

	filtered, diags := FilterResource(context.Background(), resources, map[string][]string{
		"peer_l3_network_uuids": {"l3-b"},
	}, "vip")
	if diags.HasError() {
		t.Fatalf("FilterResource returned diagnostics: %v", diags)
	}

	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered resource, got %d", len(filtered))
	}
	if filtered[0].UUID != "vip-1" {
		t.Fatalf("expected vip-1, got %s", filtered[0].UUID)
	}
}

func TestFilterResourceResolvesSecurityGroupAttachedL3NetworkAliases(t *testing.T) {
	resources := []filterTestSecurityGroup{
		{UUID: "sg-1", AttachedL3NetworkUuids: []string{"l3-a"}},
		{UUID: "sg-2", AttachedL3NetworkUuids: []string{"l3-b", "l3-c"}},
	}

	tests := []struct {
		name      string
		filterKey string
	}{
		{name: "schema spelling", filterKey: "attached_l3network_uuids"},
		{name: "canonical spelling", filterKey: "attached_l3_network_uuids"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered, diags := FilterResource(context.Background(), resources, map[string][]string{
				tt.filterKey: {"l3-c"},
			}, "security_group")
			if diags.HasError() {
				t.Fatalf("FilterResource returned diagnostics: %v", diags)
			}

			if len(filtered) != 1 {
				t.Fatalf("expected 1 filtered resource, got %d", len(filtered))
			}
			if filtered[0].UUID != "sg-2" {
				t.Fatalf("expected sg-2, got %s", filtered[0].UUID)
			}
		})
	}
}

func TestFilterResourceResolvesNestedSecurityGroupRuleFields(t *testing.T) {
	resources := []filterTestSecurityGroup{
		{UUID: "sg-1", Rules: []filterTestSecurityGroupRule{{SrcIpRange: "10.0.0.0/24", DstIpRange: "192.168.1.0/24"}}},
		{UUID: "sg-2", Rules: []filterTestSecurityGroupRule{{SrcIpRange: "172.16.0.0/16", DstIpRange: "192.168.2.0/24"}}},
	}

	tests := []struct {
		name      string
		filterKey string
		value     string
		wantUUID  string
	}{
		{name: "source ip range", filterKey: "src_ip_range", value: "172.16.0.0/16", wantUUID: "sg-2"},
		{name: "destination ip range", filterKey: "dst_ip_range", value: "192.168.1.0/24", wantUUID: "sg-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered, diags := FilterResource(context.Background(), resources, map[string][]string{
				tt.filterKey: {tt.value},
			}, "security_group")
			if diags.HasError() {
				t.Fatalf("FilterResource returned diagnostics: %v", diags)
			}

			if len(filtered) != 1 {
				t.Fatalf("expected 1 filtered resource, got %d", len(filtered))
			}
			if filtered[0].UUID != tt.wantUUID {
				t.Fatalf("expected %s, got %s", tt.wantUUID, filtered[0].UUID)
			}
		})
	}
}

func TestFilterResourceResolvesNestedSliceFieldPath(t *testing.T) {
	resources := []filterTestImage{
		{UUID: "image-1", Format: "qcow2", Type: "RootVolumeTemplate", BackupStorageRefs: []filterTestImageBackupStorageRef{{BackupStorageUuid: "bs-1"}}},
		{UUID: "image-2", Format: "raw", Type: "DataVolumeTemplate", BackupStorageRefs: []filterTestImageBackupStorageRef{{BackupStorageUuid: "bs-2"}}},
	}

	tests := []struct {
		name      string
		filterKey string
		value     string
		wantUUID  string
	}{
		{name: "backup storage uuids", filterKey: "backup_storage_uuids", value: "bs-2", wantUUID: "image-2"},
		{name: "image format", filterKey: "image_format", value: "qcow2", wantUUID: "image-1"},
		{name: "image type", filterKey: "image_type", value: "DataVolumeTemplate", wantUUID: "image-2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered, diags := FilterResource(context.Background(), resources, map[string][]string{
				tt.filterKey: {tt.value},
			}, "image")
			if diags.HasError() {
				t.Fatalf("FilterResource returned diagnostics: %v", diags)
			}

			if len(filtered) != 1 {
				t.Fatalf("expected 1 filtered resource, got %d", len(filtered))
			}
			if filtered[0].UUID != tt.wantUUID {
				t.Fatalf("expected %s, got %s", tt.wantUUID, filtered[0].UUID)
			}
		})
	}
}

func TestFilterResourceResolvesUnmappedSnakeCaseByJSONTag(t *testing.T) {
	resources := []filterTestEIP{
		{VmNicUuid: "nic-a"},
		{VmNicUuid: "nic-b"},
	}

	filtered, diags := FilterResource(context.Background(), resources, map[string][]string{
		"vm_nic_uuid": {"nic-b"},
	}, "unknown")
	if diags.HasError() {
		t.Fatalf("FilterResource returned diagnostics: %v", diags)
	}

	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered resource, got %d", len(filtered))
	}
	if filtered[0].VmNicUuid != "nic-b" {
		t.Fatalf("expected nic-b, got %s", filtered[0].VmNicUuid)
	}
}

func TestFilterResourceResolvesTerraformSDKTag(t *testing.T) {
	resources := []filterTestTerraformModel{
		{ZoneUuid: types.StringValue("zone-a")},
		{ZoneUuid: types.StringValue("zone-b")},
	}

	filtered, diags := FilterResource(context.Background(), resources, map[string][]string{
		"zone_uuid": {"zone-b"},
	}, "unknown")
	if diags.HasError() {
		t.Fatalf("FilterResource returned diagnostics: %v", diags)
	}

	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered resource, got %d", len(filtered))
	}
	if filtered[0].ZoneUuid.ValueString() != "zone-b" {
		t.Fatalf("expected zone-b, got %s", filtered[0].ZoneUuid.ValueString())
	}
}

func TestFilterResourceInvalidFilterKey(t *testing.T) {
	resources := []filterTestAffinityGroup{
		{filterTestBase: filterTestBase{UUID: "ag-1"}, Name: "ag-a", ZoneUuid: "zone-a"},
	}

	_, diags := FilterResource(context.Background(), resources, map[string][]string{
		"missing": {"value"},
	}, "affinity_group")

	if !diags.HasError() {
		t.Fatal("expected invalid filter key diagnostic")
	}
}
