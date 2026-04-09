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

func TestAccZStackNetworkingSecGroupRulesDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.SecurityGroupRules) == 0 {
		t.Skip("no security group rules in env data")
	}

	// Find the first rule to use for filter values
	rule := env.SecurityGroupRules[0]
	action := envStr(rule, "action")
	protocol := envStr(rule, "protocol")

	resource.ParallelTest(t, resource.TestCase{
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_networking_secgroup_rules.test", tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("action"), knownvalue.StringExact(action)),
					statecheck.ExpectKnownValue("data.zstack_networking_secgroup_rules.test", tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("protocol"), knownvalue.StringExact(protocol)),
					statecheck.ExpectKnownValue("data.zstack_networking_secgroup_rules.test", tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(rule, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_networking_secgroup_rules.test", tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(rule, "type"))),
					statecheck.ExpectKnownValue("data.zstack_networking_secgroup_rules.test", tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(rule, "state"))),
					statecheck.ExpectKnownValue("data.zstack_networking_secgroup_rules.test", tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("dst_port_range"), knownvalue.StringExact(envStr(rule, "dst_port_range"))),
					statecheck.ExpectKnownValue("data.zstack_networking_secgroup_rules.test", tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("src_ip_range"), knownvalue.StringExact(envStr(rule, "src_ip_range"))),
				},
			},
		},
	})
}
