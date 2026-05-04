// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

func TestBUG019DataDisksStatePersistsVolumeUUIDs(t *testing.T) {
	ctx := context.Background()
	state := vmInstanceDataSourceModel{}

	dataDisks, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: instanceDiskModelAttrTypes()}, []diskModel{{
		OfferingUuid:       types.StringValue("offering-1"),
		PrimaryStorageUuid: types.StringNull(),
	}})
	if diags.HasError() {
		t.Fatalf("failed to seed data_disks state: %v", diags)
	}
	state.DataDisks = dataDisks

	vm := &view.VmInstanceInventoryView{
		BaseInfoView: view.BaseInfoView{UUID: "vm-1"},
		AllVolumes: []view.VolumeInventoryView{
			{BaseInfoView: view.BaseInfoView{UUID: "root-volume"}, Type: "Root", PrimaryStorageUuid: "ps-root"},
			{BaseInfoView: view.BaseInfoView{UUID: "data-volume-1"}, Type: "Data", PrimaryStorageUuid: "ps-data"},
		},
	}

	updatedState, err := syncInstanceDataDisksFromVM(ctx, state, vm)
	if err != nil {
		t.Fatalf("syncInstanceDataDisksFromVM returned error: %v", err)
	}

	var disks []diskModel
	if diags := updatedState.DataDisks.ElementsAs(ctx, &disks, false); diags.HasError() {
		t.Fatalf("failed to decode updated data_disks state: %v", diags)
	}
	if len(disks) != 1 {
		t.Fatalf("expected 1 data disk in state, got %d", len(disks))
	}
	if got := disks[0].VolumeUuid.ValueString(); got != "data-volume-1" {
		t.Fatalf("expected data disk uuid to be persisted from VM inventory, got %q", got)
	}
	if got := disks[0].PrimaryStorageUuid.ValueString(); got != "ps-data" {
		t.Fatalf("expected data disk primary storage uuid to be refreshed, got %q", got)
	}
}

func TestBUG019DeleteUsesStateDataDiskUUIDs(t *testing.T) {
	ctx := context.Background()
	state := vmInstanceDataSourceModel{}

	dataDisks, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: instanceDiskModelAttrTypes()}, []diskModel{{
		VolumeUuid:         types.StringValue("data-volume-1"),
		OfferingUuid:       types.StringValue("offering-1"),
		PrimaryStorageUuid: types.StringValue("ps-data"),
	}})
	if diags.HasError() {
		t.Fatalf("failed to seed data_disks state: %v", diags)
	}
	state.DataDisks = dataDisks

	volumeUUIDs, err := instanceDataVolumeUUIDsFromState(ctx, state)
	if err != nil {
		t.Fatalf("instanceDataVolumeUUIDsFromState returned error: %v", err)
	}
	if len(volumeUUIDs) != 1 || volumeUUIDs[0] != "data-volume-1" {
		t.Fatalf("expected delete path to use state volume uuid, got %#v", volumeUUIDs)
	}
}

func TestBUG019NormalizeNetworkInterfacesSetsStaticIPFromObservedNIC(t *testing.T) {
	vm := &view.VmInstanceInventoryView{
		VmNics: []view.VmNicInventoryView{
			{
				L3NetworkUuid: "l3-1",
				Ip:            "172.24.202.21",
			},
		},
		DefaultL3NetworkUuid: "l3-1",
	}

	nics := normalizeNetworkInterfacesFromVM(vm)
	if len(nics) != 1 {
		t.Fatalf("expected 1 network interface, got %d", len(nics))
	}

	if got := nics[0].L3NetworkUuid.ValueString(); got != "l3-1" {
		t.Fatalf("expected l3_network_uuid to be l3-1, got %q", got)
	}

	if got := nics[0].StaticIp.ValueString(); got != "172.24.202.21" {
		t.Fatalf("expected static_ip to be observed nic ip, got %q", got)
	}

	if got := nics[0].DefaultL3.ValueBool(); !got {
		t.Fatalf("expected default_l3 to be true")
	}
}

func TestBUG019BuildUpdatedStateFromVMKeepsIdentityAndCoreFields(t *testing.T) {
	ctx := context.Background()
	current := vmInstanceDataSourceModel{}

	vm := &view.VmInstanceInventoryView{
		BaseInfoView:         view.BaseInfoView{UUID: "vm-123", Name: "vm-name"},
		Description:          "vm-desc",
		ImageUuid:            "img-123",
		MemorySize:           2 * 1024 * 1024 * 1024,
		CpuNum:               2,
		DefaultL3NetworkUuid: "l3-1",
		VmNics:               []view.VmNicInventoryView{{L3NetworkUuid: "l3-1", Ip: "172.24.202.21", Gateway: "172.24.0.1", Netmask: "255.255.0.0"}},
	}

	updated, err := buildUpdatedStateFromVM(ctx, current, vm)
	if err != nil {
		t.Fatalf("buildUpdatedStateFromVM returned error: %v", err)
	}

	if got := updated.Uuid.ValueString(); got != "vm-123" {
		t.Fatalf("expected uuid vm-123, got %q", got)
	}
	if got := updated.Name.ValueString(); got != "vm-name" {
		t.Fatalf("expected name vm-name, got %q", got)
	}
	if got := updated.Description.ValueString(); got != "vm-desc" {
		t.Fatalf("expected description vm-desc, got %q", got)
	}
	if got := updated.MemorySize.ValueInt64(); got != 2048 {
		t.Fatalf("expected memory_size 2048MB, got %d", got)
	}
	if got := updated.CPUNum.ValueInt64(); got != 2 {
		t.Fatalf("expected cpu_num 2, got %d", got)
	}
}
