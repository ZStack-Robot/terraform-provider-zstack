// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccCreateImageResource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.BackupStorages) == 0 {
		t.Skip("no backup storages in env data")
	}
	bsUUID := envStr(env.BackupStorages[0], "uuid")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "zstack_image" "test" {
					name        = "example-image"
					description = "A test image for creation"
					url         = "http://192.168.200.100/mirror/diskimages/CentOS-7-x86_64-Cloudinit-8G-official.qcow2"
					guest_os_type = "Centos 7"
					platform    = "Linux"
					format      = "qcow2"
					media_type  = "RootVolumeTemplate"
					architecture = "x86_64"
					backup_storage_uuids = [%q]
					boot_mode   = "Legacy"
				}`, bsUUID),

				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_image.test", tfjsonpath.New("name"), knownvalue.StringExact("example-image")),
					statecheck.ExpectKnownValue("zstack_image.test", tfjsonpath.New("description"), knownvalue.StringExact("A test image for creation")),
					statecheck.ExpectKnownValue("zstack_image.test", tfjsonpath.New("guest_os_type"), knownvalue.StringExact("Centos 7")),
					statecheck.ExpectKnownValue("zstack_image.test", tfjsonpath.New("platform"), knownvalue.StringExact("Linux")),
					statecheck.ExpectKnownValue("zstack_image.test", tfjsonpath.New("format"), knownvalue.StringExact("qcow2")),
					statecheck.ExpectKnownValue("zstack_image.test", tfjsonpath.New("architecture"), knownvalue.StringExact("x86_64")),
					statecheck.ExpectKnownValue("zstack_image.test", tfjsonpath.New("boot_mode"), knownvalue.StringExact("Legacy")),
					statecheck.ExpectKnownValue("zstack_image.test", tfjsonpath.New("backup_storage_uuids"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("zstack_image.test", tfjsonpath.New("backup_storage_uuids").AtSliceIndex(0), knownvalue.StringExact(bsUUID)),
				},
			},
			{
				ResourceName:            "zstack_image.test",
				ImportState:             true,
				ImportStateIdFunc:       importStateUUID("zstack_image.test"),
				ImportStateVerify:       true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore: []string{"description", "guest_os_type", "platform", "format", "media_type", "backup_storage_uuids", "architecture", "virtio", "boot_mode", "expunge", "system", "last_updated"},
			},
		},
	})
}
