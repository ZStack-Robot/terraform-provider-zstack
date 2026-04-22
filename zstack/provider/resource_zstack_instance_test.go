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

func TestVMResource_Schema(t *testing.T) {
	var r instanceResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "image_uuid"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "description", "memory_size", "cpu_num"}
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

func TestVMResource_Metadata(t *testing.T) {
	var r instanceResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_instance" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccInstanceResource(t *testing.T) {
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

	// Use the image UUID that existing VMs were created with (proven non-system image).
	// Fall back to any Ready image if no VMs exist.
	var imageUUID string
	if len(env.VmInstances) > 0 {
		imageUUID = envStr(env.VmInstances[0], "image_uuid")
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
		t.Skip("no suitable images in env data")
	}

	offeringUUID := envStr(env.InstanceOfferings[0], "uuid")

	// Prefer Public L3 network; fall back to first available L3.
	var l3UUID string
	for _, l3 := range env.L3Networks {
		if envStr(l3, "category") == "Public" {
			l3UUID = envStr(l3, "uuid")
			break
		}
	}
	if l3UUID == "" {
		l3UUID = envStr(env.L3Networks[0], "uuid")
	}

	createConfig := func(name string) string {
		return providerConfig() + fmt.Sprintf(`
resource "zstack_instance" "test" {
  name                   = %q
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
`, name, imageUUID, offeringUUID, l3UUID)
	}

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckInstanceDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: createConfig("acc-test-instance"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-instance")),
				},
			},
			// Step 2: Update name — SKIPPED: BUG-9 (UpdateVmInstance SDK returns empty struct)
			// Step 3: Import
			{
				Config:                               createConfig("acc-test-instance"),
				ResourceName:                         "zstack_instance.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_instance.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"expunge", "network_interfaces", "instance_offering_uuid"},
			},
		},
	})
}
