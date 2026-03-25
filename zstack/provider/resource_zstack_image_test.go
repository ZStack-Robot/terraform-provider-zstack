// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCreateImageResource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.BackupStorages) == 0 {
		t.Skip("no backup storages in env data")
	}
	bsUUID := envStr(env.BackupStorages[0], "uuid")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zstack_image.test", "name", "example-image"),
					resource.TestCheckResourceAttr("zstack_image.test", "description", "A test image for creation"),
					resource.TestCheckResourceAttr("zstack_image.test", "guest_os_type", "Centos 7"),
					resource.TestCheckResourceAttr("zstack_image.test", "platform", "Linux"),
					resource.TestCheckResourceAttr("zstack_image.test", "format", "qcow2"),
					resource.TestCheckResourceAttr("zstack_image.test", "architecture", "x86_64"),
					resource.TestCheckResourceAttr("zstack_image.test", "boot_mode", "Legacy"),
					resource.TestCheckResourceAttr("zstack_image.test", "backup_storage_uuids.#", "1"),
					resource.TestCheckResourceAttr("zstack_image.test", "backup_storage_uuids.0", bsUUID),
				),
			},
		},
	})
}
