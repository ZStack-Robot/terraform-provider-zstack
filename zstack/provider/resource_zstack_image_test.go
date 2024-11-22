// Copyright (c) ZStack.io, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCreateImageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test creating a single image resource
			{
				Config: providerConfig + `
				resource "zstack_image" "test" {
					name        = "example-image"
					description = "A test image for creation"
					url         = "http://192.168.200.100/mirror/diskimages/CentOS-7-x86_64-Cloudinit-8G-official.qcow2"
					guest_os_type = "Centos 7"
					platform    = "Linux"
					format      = "qcow2"
					architecture = "x86_64"
					virtio      = true
					backup_storage_uuids = ["5565902ddccc4aada737e0ca6b844ff4"]
					boot_mode   = "legacy"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the resource attributes
					resource.TestCheckResourceAttr("zstack_image.test", "name", "example-image"),
					resource.TestCheckResourceAttr("zstack_image.test", "description", "A test image for creation"),
					resource.TestCheckResourceAttr("zstack_image.test", "url", "http://192.168.200.100/mirror/diskimages/CentOS-7-x86_64-Cloudinit-8G-official.qcow2"),
					resource.TestCheckResourceAttr("zstack_image.test", "guest_os_type", "Centos 7"),
					resource.TestCheckResourceAttr("zstack_image.test", "platform", "Linux"),
					resource.TestCheckResourceAttr("zstack_image.test", "format", "qcow2"),
					resource.TestCheckResourceAttr("zstack_image.test", "architecture", "x86_64"),
					resource.TestCheckResourceAttr("zstack_image.test", "virtio", "true"),
					resource.TestCheckResourceAttr("zstack_image.test", "boot_mode", "legacy"),
					resource.TestCheckResourceAttr("zstack_image.test", "backup_storage_uuids.#", "1"),
					resource.TestCheckResourceAttr("zstack_image.test", "backup_storage_uuids.0", "5565902ddccc4aada737e0ca6b844ff4"),
				),
			},
		},
	})
}
