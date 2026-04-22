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

func TestVolumeBackupResource_Schema(t *testing.T) {
	var r volumeBackupResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}
	// Check required attributes
	required := []string{"name", "volume_uuid", "backup_storage_uuid"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}
	// Check computed attributes
	computed := []string{"uuid", "type", "state", "status", "size"}
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

func TestVolumeBackupResource_Metadata(t *testing.T) {
	var r volumeBackupResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_volume_backup" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccVolumeBackupResource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.BackupStorages) == 0 && len(env.ImageStoreBackupStorages) == 0 {
		t.Skip("no backup_storages in env data")
	}
	if len(env.DiskOfferings) == 0 {
		t.Skip("no disk_offerings in env data")
	}

	// Pick backup storage UUID — prefer ImageStoreBackupStorage.
	var backupStorageUUID string
	for _, bs := range env.ImageStoreBackupStorages {
		backupStorageUUID = envStr(bs, "uuid")
		break
	}
	if backupStorageUUID == "" {
		for _, bs := range env.BackupStorages {
			backupStorageUUID = envStr(bs, "uuid")
			break
		}
	}

	diskOfferingUUID := envStr(env.DiskOfferings[0], "uuid")

	// ZStack volume backup requires the volume to be attached to a running VM.
	// We create a VM and then back up its root volume is not possible without
	// exposing root_volume_uuid from zstack_instance. Instead we create a data
	// volume, attach it to a VM (via vm_instance_uuid), then back it up.
	// Approach: create a data volume attached to a VM, then backup it.
	if len(env.Images) == 0 || len(env.InstanceOfferings) == 0 || len(env.L3Networks) == 0 {
		t.Skip("volume_backup test requires images, instance_offerings, and l3_networks")
	}

	const ttylinuxUUID = "dfc919110d734009bea2f04a5e8ac9ef"
	var imageUUID string
	for _, img := range env.Images {
		if envStr(img, "uuid") == ttylinuxUUID && envStr(img, "status") == "Ready" {
			imageUUID = ttylinuxUUID
			break
		}
	}
	if imageUUID == "" {
		for _, img := range env.Images {
			if envStr(img, "status") == "Ready" {
				imageUUID = envStr(img, "uuid")
				break
			}
		}
	}
	if imageUUID == "" {
		t.Skip("no Ready images in env data")
	}

	offeringUUID := envStr(env.InstanceOfferings[0], "uuid")
	l3UUID := envStr(env.L3Networks[0], "uuid")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVolumeBackupDestroy,
		Steps: []tfresource.TestStep{
			// Create a VM (Running), attach a data volume to it, then back up the data volume.
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_instance" "backup_test_vm" {
  name                   = "acc-test-backup-vm"
  image_uuid             = %q
  instance_offering_uuid = %q
  expunge                = true
  network_interfaces = [
    {
      l3_network_uuid = %q
      default_l3      = true
    }
  ]
}

resource "zstack_volume" "backup_data_vol" {
  name              = "acc-test-backup-vol"
  disk_offering_uuid = %q
  vm_instance_uuid  = zstack_instance.backup_test_vm.uuid
}

resource "zstack_volume_backup" "test" {
  name                = "acc-test-volume-backup"
  volume_uuid         = zstack_volume.backup_data_vol.uuid
  backup_storage_uuid = %q
}
`, imageUUID, offeringUUID, l3UUID, diskOfferingUUID, backupStorageUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_volume_backup.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_volume_backup.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-volume-backup")),
				},
			},
		},
	})
}
