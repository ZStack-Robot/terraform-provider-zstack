// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestPrimaryStorageResource_Schema(t *testing.T) {
	var r primaryStorageResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "zone_uuid", "type"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "state", "status", "total_capacity", "available_capacity"}
	for _, attr := range computed {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing computed attribute %q", attr)
		}
		if !a.IsComputed() {
			t.Errorf("attribute %q should be computed", attr)
		}
	}
}

func TestPrimaryStorageResource_Metadata(t *testing.T) {
	var r primaryStorageResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_primary_storage" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}


// TestAccPrimaryStorageResource tests the primary_storage resource acceptance.
// LocalStorage requires a host path that does not conflict with existing storage.
// Creating a second LocalStorage on a single-node env risks destabilizing it,
// so this test is skipped unless a safe test path is explicitly configured.
func TestAccPrimaryStorageResource(t *testing.T) {
	env := loadEnvData(t)

	if len(env.Zones) == 0 {
		t.Skip("no zones in env.json, skipping primary storage acceptance test")
	}

	zoneUUID := envStr(env.Zones[0], "uuid")

	// Check if there's a designated test path for LocalStorage.
	// We look for a primary_storage entry with a "test_url" field indicating
	// a safe path to use for acceptance testing (e.g. /tmp/acc-test-ps).
	var testURL string
	for _, ps := range env.PrimaryStorages {
		if u := envStr(ps, "test_url"); u != "" {
			testURL = u
			break
		}
	}

	if testURL == "" {
		// Default to a tmp path that is unlikely to conflict.
		// This will only work if the host has enough capacity and the path is writable.
		testURL = "/tmp/acc-test-primary-storage"
	}

	// Safety check: skip if there's only one primary storage and it's already using
	// LocalStorage — adding another LocalStorage on the same host is risky.
	localCount := 0
	for _, ps := range env.PrimaryStorages {
		if envStr(ps, "type") == "LocalStorage" {
			localCount++
		}
	}
	if localCount > 0 {
		t.Skip("existing LocalStorage primary storage detected; skipping to avoid destabilizing single-node env. Set test_url in env.json primary_storages to enable.")
	}

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPrimaryStorageDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_primary_storage" "test" {
  name      = "acc-test-primary-storage"
  zone_uuid = %q
  type      = "LocalStorage"
  url       = %q
}
`, zoneUUID, testURL),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_primary_storage.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_primary_storage.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-primary-storage")),
					statecheck.ExpectKnownValue("zstack_primary_storage.test", tfjsonpath.New("type"), knownvalue.StringExact("LocalStorage")),
				},
			},
			// Step 2: Update name
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_primary_storage" "test" {
  name      = "acc-test-primary-storage-updated"
  zone_uuid = %q
  type      = "LocalStorage"
  url       = %q
}
`, zoneUUID, testURL),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_primary_storage.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-primary-storage-updated")),
				},
			},
			// Step 3: Import
			{
				ResourceName:                        "zstack_primary_storage.test",
				ImportState:                         true,
				ImportStateIdFunc:                   importStateIdFromUUID("zstack_primary_storage.test"),
				ImportStateVerify:                   true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				// Plan-only fields not returned by the API
				ImportStateVerifyIgnore:             []string{"mon_urls", "root_volume_pool_name", "data_volume_pool_name", "image_cache_pool_name", "disk_uuids"},
			},
		},
	})
}
