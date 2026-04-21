// Copyright (c) ZStack.io, Inc.

package provider

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBugTrackerProviderSourceDropsLegacyFactoryNames(t *testing.T) {
	providerPath := filepath.Join("provider.go")
	data, err := os.ReadFile(providerPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", providerPath, err)
	}

	legacyNames := []string{
		"ZStackl3NetworkDataSource",
		"ZStackvmsDataSource",
		"ZStackmnNodeDataSource",
		"ZStackvrouterDataSource",
	}

	for _, legacyName := range legacyNames {
		if strings.Contains(string(data), legacyName) {
			t.Fatalf("provider.go still references legacy factory name %q", legacyName)
		}
	}
}

func TestBugTrackerInstanceSourceDropsLegacyTypeName(t *testing.T) {
	instancePath := filepath.Join("resource_zstack_instance.go")
	data, err := os.ReadFile(instancePath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", instancePath, err)
	}

	if strings.Contains(string(data), "type vmResource struct") {
		t.Fatal("resource_zstack_instance.go still uses legacy vmResource type name")
	}
}

func TestBugTrackerRandomizedAcceptanceNamesUseHelper(t *testing.T) {
	files := []string{
		"resource_zstack_account_test.go",
		"resource_zstack_affinity_group_test.go",
		"resource_zstack_auto_scaling_group_test.go",
		"resource_zstack_cluster_test.go",
		"resource_zstack_iam2_project_test.go",
		"resource_zstack_l2vlan_network_test.go",
		"resource_zstack_load_balancer_listener_test.go",
		"resource_zstack_load_balancer_test.go",
		"resource_zstack_port_forwarding_rule_test.go",
		"resource_zstack_ssh_key_pair_test.go",
		"resource_zstack_virtual_router_image_test.go",
		"resource_zstack_virtual_router_offering_test.go",
		"resource_zstack_vip_test.go",
		"resource_zstack_volume_test.go",
		"resource_zstack_zone_test.go",
	}

	for _, file := range files {
		data, err := os.ReadFile(filepath.Join(file))
		if err != nil {
			t.Fatalf("failed to read %s: %v", file, err)
		}

		text := string(data)
		if strings.Contains(text, `= "acc-test-`) {
			t.Fatalf("%s still hardcodes an acc-test-* resource name", file)
		}
		if !strings.Contains(text, "testAccName(") {
			t.Fatalf("%s should use testAccName helper for randomized resource names", file)
		}
	}
}

func TestBugTrackerEmptyIDPatternIgnoresComments(t *testing.T) {
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "comment_only.go")
	content := `package provider

// func (r *resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
// 	state := model{Uuid: types.StringValue("")}
// }
`
	if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}

	if violations := scanForEmptyIDPattern(t, file); len(violations) != 0 {
		t.Fatalf("expected commented code to be ignored, got %d violations", len(violations))
	}
}

func TestBugTrackerEmptyIDPatternDetectsRealViolation(t *testing.T) {
	tempDir := t.TempDir()
	file := filepath.Join(tempDir, "real_violation.go")
	content := `package provider

func (r *resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := model{Uuid: types.StringValue("")}
	_ = state
}
`
	if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}

	if violations := scanForEmptyIDPattern(t, file); len(violations) == 0 {
		t.Fatal("expected a real empty UUID violation to be detected")
	}
}

func TestBugTrackerZoneAcceptanceCoversUpdateStep(t *testing.T) {
	zoneTestPath := filepath.Join("resource_zstack_zone_test.go")
	data, err := os.ReadFile(zoneTestPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", zoneTestPath, err)
	}

	text := string(data)
	if !strings.Contains(text, `description = "Updated acceptance test zone"`) {
		t.Fatal("resource_zstack_zone_test.go should include an explicit update step for BUG-024 coverage")
	}
	if !strings.Contains(text, `knownvalue.StringExact("Disabled")`) {
		t.Fatal("resource_zstack_zone_test.go should verify updated zone state in the update step")
	}
}
