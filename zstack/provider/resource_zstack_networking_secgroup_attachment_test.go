// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackSecurityGroupAttachment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSecGroupAttachmentConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zstack_networking_secgroup_attachment.test_attach", "secgroup_uuid", "e5f2f001d9834a87bd7a8ff65b2c0b22"),
					resource.TestCheckResourceAttr("zstack_networking_secgroup_attachment.test_attach", "nic_uuid", "8432b2608a8a4f639cc56a8c65ee42ef"),
					resource.TestCheckResourceAttrSet("zstack_networking_secgroup_attachment.test_attach", "id"),
				),
			},
		},
	})
}

func testAccSecGroupAttachmentConfig() string {
	return `
resource "zstack_networking_secgroup_attachment" "test_attach" {
  secgroup_uuid = "e5f2f001d9834a87bd7a8ff65b2c0b22"
  nic_uuid      = "8432b2608a8a4f639cc56a8c65ee42ef"
}
`
}
