// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"os"
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

	optional := []string{"cpu_mode", "hostname"}
	for _, attr := range optional {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
		if !a.IsOptional() {
			t.Errorf("attribute %q should be optional", attr)
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
	name := testAccName("instance")
	updatedName := name + "-updated"

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

	createConfig := func(name, description, platform, guestOsType string) string {
		return providerConfig() + fmt.Sprintf(`
resource "zstack_instance" "test" {
  name                   = %q
  description            = %q
  image_uuid             = %q
  instance_offering_uuid = %q
  platform               = %q
  guest_os_type          = %q
  expunge                = true
  network_interfaces = [
    {
      l3_network_uuid = %q
      default_l3      = true
    }
  ]
}
`, name, description, imageUUID, offeringUUID, platform, guestOsType, l3UUID)
	}

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckInstanceDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: createConfig(name, "acceptance instance", "Linux", "CentOS 7.6"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance instance")),
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("platform"), knownvalue.StringExact("Linux")),
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("guest_os_type"), knownvalue.StringExact("CentOS 7.6")),
				},
			},
			// Step 2: Update mutable VM metadata
			{
				Config: createConfig(updatedName, "acceptance instance updated", "Linux", "CentOS 7.9"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance instance updated")),
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("platform"), knownvalue.StringExact("Linux")),
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("guest_os_type"), knownvalue.StringExact("CentOS 7.9")),
				},
			},
			// Step 3: Import
			{
				Config:                               createConfig(updatedName, "acceptance instance updated", "Linux", "CentOS 7.9"),
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

func TestAccInstanceResourceRootDiskSizeState(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("acceptance test skipped unless TF_ACC is set")
	}

	imageName := os.Getenv("ZSTACK_TEST_IMAGE_NAME")
	l3Name := os.Getenv("ZSTACK_TEST_L3_NAME")
	if (imageName == "") != (l3Name == "") {
		t.Fatal("ZSTACK_TEST_IMAGE_NAME and ZSTACK_TEST_L3_NAME must be set together")
	}

	var lookupConfig, imageRef, l3Ref string
	if imageName != "" && l3Name != "" {
		lookupConfig = fmt.Sprintf(`
data "zstack_images" "test" {
  name = %q
}

data "zstack_l3networks" "test" {
  name = %q
}
`, imageName, l3Name)
		imageRef = "data.zstack_images.test.images[0].uuid"
		l3Ref = "data.zstack_l3networks.test.l3networks[0].uuid"
	} else {
		env := loadEnvData(t)
		if len(env.Images) == 0 {
			t.Skip("no images in env data")
		}
		if len(env.L3Networks) == 0 {
			t.Skip("no l3_networks in env data")
		}

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

		imageRef = fmt.Sprintf("%q", imageUUID)
		l3Ref = fmt.Sprintf("%q", l3UUID)
	}

	name := testAccName("instance-root-disk")
	config := func(description string) string {
		return providerConfig() + lookupConfig + fmt.Sprintf(`
resource "zstack_instance" "test" {
  name          = %q
  description   = %q
  image_uuid    = %s
  cpu_num       = 1
  memory_size   = 2048
  strategy      = "CreateStopped"
  expunge       = true

  network_interfaces = [
    {
      l3_network_uuid = %s
      default_l3      = true
    }
  ]

  root_disk = {
    size = 50
  }
}
`, name, description, imageRef, l3Ref)
	}

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckInstanceDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: config("zstac-86122 initial"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("description"), knownvalue.StringExact("zstac-86122 initial")),
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("root_disk").AtMapKey("size"), knownvalue.Int64Exact(50)),
				},
			},
			{
				Config: config("zstac-86122 updated"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("description"), knownvalue.StringExact("zstac-86122 updated")),
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("root_disk").AtMapKey("size"), knownvalue.Int64Exact(50)),
				},
			},
		},
	})
}

func TestAccInstanceResourceHostnameState(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("acceptance test skipped unless TF_ACC is set")
	}

	imageName := os.Getenv("ZSTACK_TEST_IMAGE_NAME")
	l3Name := os.Getenv("ZSTACK_TEST_L3_NAME")
	if imageName == "" || l3Name == "" {
		t.Skip("ZSTACK_TEST_IMAGE_NAME and ZSTACK_TEST_L3_NAME must be set for hostname acceptance test")
	}

	hostname := testAccName("guest-host")
	name := testAccName("instance-hostname")
	config := providerConfig() + fmt.Sprintf(`
data "zstack_images" "test" {
  name = %q
}

data "zstack_l3networks" "test" {
  name = %q
}

resource "zstack_instance" "test" {
  name        = %q
  hostname    = %q
  description = "zstac-86125 hostname test"
  image_uuid  = data.zstack_images.test.images[0].uuid
  cpu_num     = 1
  memory_size = 2048
  strategy    = "CreateStopped"
  expunge     = true

  network_interfaces = [
    {
      l3_network_uuid = data.zstack_l3networks.test.l3networks[0].uuid
      default_l3      = true
    }
  ]
}
`, imageName, l3Name, name, hostname)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckInstanceDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: config,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_instance.test", tfjsonpath.New("hostname"), knownvalue.StringExact(hostname)),
				},
			},
		},
	})
}
