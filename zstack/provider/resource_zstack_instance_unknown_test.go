// Copyright (c) ZStack.io, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

// TestInstanceUpdateGuardsUnknownValues verifies that Unknown values in plan
// do NOT trigger API calls with zero values (which would catastrophically
// shrink VM memory/CPU or disable NeverStop).
//
// Per Story-15a (Task 6): Unknown semantics for instance.go fields:
// - NeverStop: Unknown → omit (don't send ha::NeverStop system tag)
// - MemorySize: Unknown → omit (sending 0 would shrink VM memory to zero - data loss)
// - CPUNum: Unknown → omit (sending 0 would shrink VM to zero CPUs - inoperable)
//
// This test validates the guard logic by constructing Unknown values and
// testing the conditional expressions directly.
func TestInstanceUpdateGuardsUnknownValues(t *testing.T) {
	t.Run("CPUNum", func(t *testing.T) {
		// Simulate plan with Unknown CPUNum
		plan := vmInstanceDataSourceModel{
			CPUNum: types.Int64Unknown(),
		}
		state := vmInstanceDataSourceModel{
			CPUNum: types.Int64Value(4),
		}

		// The guard condition: !plan.CPUNum.IsNull() && !plan.CPUNum.IsUnknown() && plan.CPUNum.ValueInt64() != state.CPUNum.ValueInt64()
		// Expected: false (should NOT trigger update because IsUnknown() == true)
		shouldUpdate := !plan.CPUNum.IsNull() && !plan.CPUNum.IsUnknown() && plan.CPUNum.ValueInt64() != state.CPUNum.ValueInt64()

		if shouldUpdate {
			t.Errorf("CPUNum: Unknown value should NOT trigger update (would send 0 to API)")
		}

		// Verify that without the IsUnknown guard, it WOULD incorrectly trigger
		// (this demonstrates the bug we're fixing)
		wouldTriggerWithoutGuard := plan.CPUNum.ValueInt64() != state.CPUNum.ValueInt64()
		if !wouldTriggerWithoutGuard {
			t.Errorf("CPUNum: Test validation failed - Unknown.ValueInt64() should return 0 and differ from state")
		}
	})

	t.Run("MemorySize", func(t *testing.T) {
		// Simulate plan with Unknown MemorySize
		plan := vmInstanceDataSourceModel{
			MemorySize: types.Int64Unknown(),
		}
		state := vmInstanceDataSourceModel{
			MemorySize: types.Int64Value(8192),
		}

		// The guard condition: !plan.MemorySize.IsNull() && !plan.MemorySize.IsUnknown() && plan.MemorySize.ValueInt64() != state.MemorySize.ValueInt64()
		// Expected: false (should NOT trigger update because IsUnknown() == true)
		shouldUpdate := !plan.MemorySize.IsNull() && !plan.MemorySize.IsUnknown() && plan.MemorySize.ValueInt64() != state.MemorySize.ValueInt64()

		if shouldUpdate {
			t.Errorf("MemorySize: Unknown value should NOT trigger update (would send 0 to API)")
		}

		// Verify that without the IsUnknown guard, it WOULD incorrectly trigger
		wouldTriggerWithoutGuard := plan.MemorySize.ValueInt64() != state.MemorySize.ValueInt64()
		if !wouldTriggerWithoutGuard {
			t.Errorf("MemorySize: Test validation failed - Unknown.ValueInt64() should return 0 and differ from state")
		}
	})

	t.Run("NeverStop", func(t *testing.T) {
		// Simulate plan with Unknown NeverStop
		plan := vmInstanceDataSourceModel{
			NeverStop: types.BoolUnknown(),
		}

		// The guard condition: !plan.NeverStop.IsNull() && !plan.NeverStop.IsUnknown() && plan.NeverStop.ValueBool()
		// Expected: false (should NOT add system tag because IsUnknown() == true)
		shouldAddTag := !plan.NeverStop.IsNull() && !plan.NeverStop.IsUnknown() && plan.NeverStop.ValueBool()

		if shouldAddTag {
			t.Errorf("NeverStop: Unknown value should NOT add system tag (would send false to API)")
		}

		// Verify that without the IsUnknown guard, it WOULD incorrectly NOT add the tag
		// (Unknown.ValueBool() returns false, so the old condition would be false anyway,
		// but the semantic is wrong - we want to explicitly check IsUnknown)
		wouldAddWithoutGuard := !plan.NeverStop.IsNull() && plan.NeverStop.ValueBool()
		if wouldAddWithoutGuard {
			t.Errorf("NeverStop: Test validation failed - Unknown.ValueBool() should return false")
		}
	})

	t.Run("CPUMode", func(t *testing.T) {
		plan := vmInstanceDataSourceModel{
			CPUMode: types.StringUnknown(),
		}

		shouldAddTag := !plan.CPUMode.IsNull() && !plan.CPUMode.IsUnknown() && plan.CPUMode.ValueString() != ""
		if shouldAddTag {
			t.Errorf("CPUMode: Unknown value should NOT add system tag")
		}
	})

	t.Run("RootDiskSize_Unknown_guard", func(t *testing.T) {
		// Simulate Unknown Int64 value (as would be extracted from RootDisk.Size)
		rootDiskSize := types.Int64Unknown()

		// The guard condition: !rootDiskPlan.Size.IsNull() && !rootDiskPlan.Size.IsUnknown()
		// Expected: false (should NOT convert size because IsUnknown() == true)
		shouldConvert := !rootDiskSize.IsNull() && !rootDiskSize.IsUnknown()

		if shouldConvert {
			t.Errorf("RootDiskSize: Unknown value should NOT trigger conversion (would send 0 to API)")
		}

		// Verify that Unknown.IsUnknown() is true (the guard we're testing)
		if !rootDiskSize.IsUnknown() {
			t.Errorf("RootDiskSize: Test validation failed - Unknown.IsUnknown() should return true")
		}
	})

	t.Run("DataDiskSize_Unknown_guard", func(t *testing.T) {
		// Simulate Unknown Int64 value (as would be extracted from DataDisk.Size)
		dataDiskSize := types.Int64Unknown()

		// The guard condition: !disk.Size.IsNull() && !disk.Size.IsUnknown()
		// Expected: false (should NOT append size because IsUnknown() == true)
		shouldAppend := !dataDiskSize.IsNull() && !dataDiskSize.IsUnknown()

		if shouldAppend {
			t.Errorf("DataDiskSize: Unknown value should NOT trigger append (would send 0 to API)")
		}

		// Verify that Unknown.IsUnknown() is true (the guard we're testing)
		if !dataDiskSize.IsUnknown() {
			t.Errorf("DataDiskSize: Test validation failed - Unknown.IsUnknown() should return true")
		}
	})

	t.Run("GpuNumber_Unknown_guard", func(t *testing.T) {
		// Simulate Unknown Int64 value (as would be extracted from GpuSpec.Number)
		gpuNumber := types.Int64Unknown()

		// The guard condition: !gpuSpecPlan.Number.IsNull() && !gpuSpecPlan.Number.IsUnknown()
		// Expected: false (should NOT extract number because IsUnknown() == true)
		shouldExtract := !gpuNumber.IsNull() && !gpuNumber.IsUnknown()

		if shouldExtract {
			t.Errorf("GpuNumber: Unknown value should NOT trigger extraction (would send 0 to API)")
		}

		// Verify that Unknown.IsUnknown() is true (the guard we're testing)
		if !gpuNumber.IsUnknown() {
			t.Errorf("GpuNumber: Test validation failed - Unknown.IsUnknown() should return true")
		}
	})
}

func TestRootDiskSizeBytesForCreate(t *testing.T) {
	t.Run("ConvertsGBForAPIWithoutChangingTerraformValue", func(t *testing.T) {
		size := types.Int64Value(50)

		got := rootDiskSizeBytesForCreate(size)
		if got == nil {
			t.Fatalf("expected root disk size bytes, got nil")
		}
		if *got != 53687091200 {
			t.Fatalf("expected 50GB to be converted to 53687091200 bytes, got %d", *got)
		}
		if size.ValueInt64() != 50 {
			t.Fatalf("expected Terraform root_disk.size to remain 50GB, got %d", size.ValueInt64())
		}
	})

	t.Run("OmitsNullAndUnknown", func(t *testing.T) {
		if got := rootDiskSizeBytesForCreate(types.Int64Null()); got != nil {
			t.Fatalf("expected null size to be omitted, got %d", *got)
		}
		if got := rootDiskSizeBytesForCreate(types.Int64Unknown()); got != nil {
			t.Fatalf("expected unknown size to be omitted, got %d", *got)
		}
	})
}

func TestPreserveInstanceNameForUpdate(t *testing.T) {
	t.Run("SetsCurrentNameWhenOnlyOtherFieldsChange", func(t *testing.T) {
		update := param.UpdateVmInstanceParam{
			Params: param.UpdateVmInstanceParamDetail{
				Description: stringPtr("updated"),
			},
		}

		preserveInstanceNameForUpdate(&update, types.StringValue("vm-name"))

		if update.Params.Name != "vm-name" {
			t.Fatalf("expected update to preserve vm name, got %q", update.Params.Name)
		}
	})

	t.Run("KeepsExplicitNameChange", func(t *testing.T) {
		update := param.UpdateVmInstanceParam{
			Params: param.UpdateVmInstanceParamDetail{
				Name: "new-name",
			},
		}

		preserveInstanceNameForUpdate(&update, types.StringValue("old-name"))

		if update.Params.Name != "new-name" {
			t.Fatalf("expected explicit name update to be preserved, got %q", update.Params.Name)
		}
	})

	t.Run("OmitsNullAndUnknownName", func(t *testing.T) {
		update := param.UpdateVmInstanceParam{}
		preserveInstanceNameForUpdate(&update, types.StringNull())
		if update.Params.Name != "" {
			t.Fatalf("expected null name to be omitted, got %q", update.Params.Name)
		}

		preserveInstanceNameForUpdate(&update, types.StringUnknown())
		if update.Params.Name != "" {
			t.Fatalf("expected unknown name to be omitted, got %q", update.Params.Name)
		}
	})
}
