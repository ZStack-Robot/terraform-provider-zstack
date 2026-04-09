// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSecGroupAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_networking_secgroup_attachment" "test_attach" {
  secgroup_uuid = %q
  nic_uuid      = %q
}
`, sgUUID, nicUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_networking_secgroup_attachment.test_attach", tfjsonpath.New("secgroup_uuid"), knownvalue.StringExact(sgUUID)),
					statecheck.ExpectKnownValue("zstack_networking_secgroup_attachment.test_attach", tfjsonpath.New("nic_uuid"), knownvalue.StringExact(nicUUID)),
					statecheck.ExpectKnownValue("zstack_networking_secgroup_attachment.test_attach", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
		},
	})
}
