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

func TestVolumeSnapshotResource_Schema(t *testing.T) {
	var r volumeSnapshotResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "volume_uuid"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "tree_uuid", "parent_uuid", "primary_storage_uuid", "volume_type", "format", "latest", "size", "state", "status", "distance", "group_uuid"}
	for _, attr := range computed {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing computed attribute %q", attr)
		}
		if !a.IsComputed() {
			t.Errorf("attribute %q should be computed", attr)
		}
	}

	optional := []string{"description", "revert"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}
}

func TestVolumeSnapshotResource_Metadata(t *testing.T) {
	var r volumeSnapshotResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_volume_snapshot" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccVolumeSnapshotResource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.DiskOfferings) == 0 {
		t.Skip("no disk_offerings in env data")
	}
	if len(env.Images) == 0 || len(env.InstanceOfferings) == 0 || len(env.L3Networks) == 0 {
		t.Skip("volume_snapshot test requires images, instance_offerings, and l3_networks for VM creation")
	}

	// Pick ttylinux image (fastest boot); fall back to any Ready image.
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
	diskOfferingUUID := envStr(env.DiskOfferings[0], "uuid")

	// A standalone data volume has status NotInstantiated and cannot be snapshotted.
	// Attach the data volume to a Running VM to instantiate it, then snapshot.
	vmAndVolBlock := fmt.Sprintf(`
resource "zstack_instance" "snap_test_vm" {
  name                   = "acc-test-snap-vm"
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

resource "zstack_volume" "snap_data_vol" {
  name               = "acc-test-snap-vol"
  disk_offering_uuid = %q
  vm_instance_uuid   = zstack_instance.snap_test_vm.uuid
}
`, imageUUID, offeringUUID, l3UUID, diskOfferingUUID)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVolumeSnapshotDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create snapshot of the attached (instantiated) data volume
			{
				Config: providerConfig() + vmAndVolBlock + `
resource "zstack_volume_snapshot" "test" {
  name        = "acc-test-snapshot"
  volume_uuid = zstack_volume.snap_data_vol.uuid
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_volume_snapshot.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_volume_snapshot.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-snapshot")),
				},
			},
			// Step 2: Update snapshot name
			{
				Config: providerConfig() + vmAndVolBlock + `
resource "zstack_volume_snapshot" "test" {
  name        = "acc-test-snapshot-updated"
  volume_uuid = zstack_volume.snap_data_vol.uuid
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_volume_snapshot.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-snapshot-updated")),
				},
			},
			// Step 3: Import
			{
				Config: providerConfig() + vmAndVolBlock + `
resource "zstack_volume_snapshot" "test" {
  name        = "acc-test-snapshot-updated"
  volume_uuid = zstack_volume.snap_data_vol.uuid
}
`,
				ResourceName:                         "zstack_volume_snapshot.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_volume_snapshot.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"revert"},
			},
		},
	})
}
