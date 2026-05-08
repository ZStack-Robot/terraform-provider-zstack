// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
	matchingRules := securityGroupRulesByActionProtocol(env.SecurityGroupRules, action, protocol)

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
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.#", fmt.Sprintf("%d", len(matchingRules))),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.0.action", action),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.0.protocol", protocol),
					testCheckSecurityGroupRulePresent("data.zstack_networking_secgroup_rules.test", envStr(rule, "uuid")),
				),
			},
		},
	})
}

func TestAccZStackNetworkingSecGroupRulesDataSourceStableOrdering(t *testing.T) {
	env := loadEnvData(t)
	if len(env.SecurityGroupRules) == 0 {
		t.Skip("no security group rules in env data")
	}

	rule := env.SecurityGroupRules[0]
	action := envStr(rule, "action")
	protocol := envStr(rule, "protocol")
	matchingRules := securityGroupRulesByActionProtocol(env.SecurityGroupRules, action, protocol)
	if len(matchingRules) < 2 {
		t.Skipf("need at least 2 security group rules matching action=%s protocol=%s in env data", action, protocol)
	}

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
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.#", fmt.Sprintf("%d", len(matchingRules))),
					testCheckSecurityGroupRulesOrdered("data.zstack_networking_secgroup_rules.test"),
				),
			},
		},
	})
}

func securityGroupRulesByActionProtocol(rules []map[string]interface{}, action, protocol string) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		if envStr(rule, "action") == action && envStr(rule, "protocol") == protocol {
			result = append(result, rule)
		}
	}
	return result
}

func testCheckSecurityGroupRulePresent(resourceName, uuid string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		count, err := strconv.Atoi(rs.Primary.Attributes["rules.#"])
		if err != nil {
			return fmt.Errorf("invalid %s rules count %q: %w", resourceName, rs.Primary.Attributes["rules.#"], err)
		}
		for i := 0; i < count; i++ {
			if rs.Primary.Attributes[fmt.Sprintf("rules.%d.uuid", i)] == uuid {
				return nil
			}
		}

		return fmt.Errorf("security group rule %s not found in %s", uuid, resourceName)
	}
}

func testCheckSecurityGroupRulesOrdered(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		count, err := strconv.Atoi(rs.Primary.Attributes["rules.#"])
		if err != nil {
			return fmt.Errorf("invalid %s rules count %q: %w", resourceName, rs.Primary.Attributes["rules.#"], err)
		}
		for i := 1; i < count; i++ {
			prevPrefix := fmt.Sprintf("rules.%d", i-1)
			currPrefix := fmt.Sprintf("rules.%d", i)
			prevSecurityGroupUUID := rs.Primary.Attributes[prevPrefix+".security_group_uuid"]
			currSecurityGroupUUID := rs.Primary.Attributes[currPrefix+".security_group_uuid"]
			prevUUID := rs.Primary.Attributes[prevPrefix+".uuid"]
			currUUID := rs.Primary.Attributes[currPrefix+".uuid"]

			if prevSecurityGroupUUID > currSecurityGroupUUID ||
				(prevSecurityGroupUUID == currSecurityGroupUUID && prevUUID > currUUID) {
				return fmt.Errorf(
					"%s rules are not ordered by security_group_uuid ASC, uuid ASC at index %d: previous=(%s,%s), current=(%s,%s)",
					resourceName,
					i,
					prevSecurityGroupUUID,
					prevUUID,
					currSecurityGroupUUID,
					currUUID,
				)
			}
		}
		return nil
	}
}
