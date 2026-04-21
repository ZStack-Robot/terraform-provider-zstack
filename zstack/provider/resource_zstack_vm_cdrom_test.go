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

func TestVmCdRomResource_Schema(t *testing.T) {
	var r vmCdRomResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}
	// Check required attributes
	required := []string{"name", "vm_instance_uuid"}
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
	computed := []string{"uuid", "device_id"}
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

func TestVmCdRomResource_Metadata(t *testing.T) {
	var r vmCdRomResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_vm_cdrom" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccVmCdRomResource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Images) == 0 {
		t.Skip("no images in env data")
	}
	if len(env.InstanceOfferings) == 0 {
		t.Skip("no instance_offerings in env data")
	}
	if len(env.L3Networks) == 0 {
		t.Skip("no l3_networks in env data")
	}

	// Pick ttylinux image (smallest, fastest boot); fall back to any Ready image.
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

	// Create the VM in Stopped state: CD-ROM creation requires the VM to be Stopped.
	vmBlock := fmt.Sprintf(`
resource "zstack_instance" "cdrom_test_vm" {
  name                   = "acc-test-cdrom-vm"
  image_uuid             = %q
  instance_offering_uuid = %q
  strategy               = "CreateStopped"
  expunge                = true
  network_interfaces = [
    {
      l3_network_uuid = %q
      default_l3      = true
    }
  ]
}
`, imageUUID, offeringUUID, l3UUID)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVmCdRomDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: providerConfig() + vmBlock + `
resource "zstack_vm_cdrom" "test" {
  name             = "acc-test-cdrom"
  vm_instance_uuid = zstack_instance.cdrom_test_vm.uuid
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_vm_cdrom.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_vm_cdrom.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-cdrom")),
				},
			},
			// Step 2: Update name — SKIPPED: BUG-11 (UpdateVmCdRom SDK returns empty struct)
			// Step 3: Import
			{
				Config:                               providerConfig() + vmBlock + `
resource "zstack_vm_cdrom" "test" {
  name             = "acc-test-cdrom"
  vm_instance_uuid = zstack_instance.cdrom_test_vm.uuid
}
`,
				ResourceName:                         "zstack_vm_cdrom.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_vm_cdrom.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
