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
