// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackNetworkingSecGroupRulesDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.SecurityGroupRules) == 0 {
		t.Skip("no security group rules in env data")
	}

	// Find the first rule to use for filter values
	rule := env.SecurityGroupRules[0]
	action := envStr(rule, "action")
	protocol := envStr(rule, "protocol")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_networking_secgroup_rules" "test" {
	filter {
		name   = "action"
		values = [%q]
	}
	filter {
		name   = "protocol"
		values = [%q]
	}
}`, action, protocol),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.0.action", action),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.0.protocol", protocol),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.0.uuid", envStr(rule, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.0.type", envStr(rule, "type")),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.0.state", envStr(rule, "state")),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.0.dst_port_range", envStr(rule, "dst_port_range")),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.0.src_ip_range", envStr(rule, "src_ip_range")),
				),
			},
		},
	})
}
