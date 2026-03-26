// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackSecurityGroupAttachment(t *testing.T) {
	env := loadEnvData(t)
	if len(env.SecurityGroups) == 0 {
		t.Skip("no security groups in env data")
	}
	if len(env.VmNics) == 0 {
		t.Skip("no vm nics in env data - need a NIC UUID")
	}
	sgUUID := envStr(env.SecurityGroups[0], "uuid")
	nicUUID := envStr(env.VmNics[0], "uuid")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_networking_secgroup_attachment" "test_attach" {
  secgroup_uuid = %q
  nic_uuid      = %q
}
`, sgUUID, nicUUID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("zstack_networking_secgroup_attachment.test_attach", "secgroup_uuid", sgUUID),
					resource.TestCheckResourceAttr("zstack_networking_secgroup_attachment.test_attach", "nic_uuid", nicUUID),
					resource.TestCheckResourceAttrSet("zstack_networking_secgroup_attachment.test_attach", "id"),
				),
			},
		},
	})
}
