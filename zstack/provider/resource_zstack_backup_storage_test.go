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

func TestBackupStorageResource_Schema(t *testing.T) {
	var r backupStorageResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "type"}
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

func TestBackupStorageResource_Metadata(t *testing.T) {
	var r backupStorageResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_backup_storage" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccBackupStorageResource(t *testing.T) {
	env := loadEnvData(t)

	// Use a dedicated test host to avoid "duplicate ImageStore IP" error
	hostname := "172.24.189.211"

	// Zone UUID for attach test
	if len(env.Zones) == 0 {
		t.Fatal("env.json has no zones")
	}
	zoneUUID := env.Zones[0]["uuid"].(string)

	createConfig := providerConfig() + fmt.Sprintf(`
resource "zstack_backup_storage" "test" {
  name     = "acc-test-backup-storage"
  type     = "ImageStoreBackupStorage"
  hostname = "%s"
  username = "root"
  password = "password"
  url      = "/tmp/acc-test-bs"
  attached_zone_uuids = ["%s"]
}
`, hostname, zoneUUID)

	updateConfig := providerConfig() + fmt.Sprintf(`
resource "zstack_backup_storage" "test" {
  name     = "acc-test-backup-storage-updated"
  type     = "ImageStoreBackupStorage"
  hostname = "%s"
  username = "root"
  password = "password"
  url      = "/tmp/acc-test-bs"
  attached_zone_uuids = ["%s"]
}
`, hostname, zoneUUID)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckBackupStorageDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: createConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_backup_storage.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_backup_storage.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-backup-storage")),
					statecheck.ExpectKnownValue("zstack_backup_storage.test", tfjsonpath.New("type"), knownvalue.StringExact("ImageStoreBackupStorage")),
				},
			},
			// Step 2: Update name (in-place)
			{
				Config: updateConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_backup_storage.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-backup-storage-updated")),
				},
			},
			// Step 3: Import
			{
				ResourceName:                         "zstack_backup_storage.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_backup_storage.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"hostname", "username", "password", "ssh_port", "mon_urls", "pool_name"},
			},
		},
	})
}
