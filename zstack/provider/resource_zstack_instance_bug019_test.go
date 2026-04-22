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
