// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackNetworkingSecGroupsDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.SecurityGroups) == 0 {
		t.Skip("no security groups in env data")
	}
	sg := env.SecurityGroups[0]
	name := envStr(sg, "name")
	pattern := string([]rune(name)[:1]) + "%"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_networking_secgroups" "test" {
	name_pattern = %q
	filter {
		name = "state"
		values = [%q]
	}
}`, pattern, envStr(sg, "state")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_networking_secgroups.test", "networking_secgroups.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroups.test", "networking_secgroups.0.state", envStr(sg, "state")),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroups.test", "networking_secgroups.0.uuid", envStr(sg, "uuid")),
				),
			},
		},
	})
}
