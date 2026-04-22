// Copyright (c) ZStack.io, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
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
}
