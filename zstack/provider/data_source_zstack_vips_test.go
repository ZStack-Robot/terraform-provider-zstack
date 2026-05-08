// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccZStackVipDataSource(t *testing.T) {
	env := loadEnvData(t)
	userVips := userVIPs(env.Vips)
	if len(userVips) == 0 {
		t.Skip("no user vips in env data")
	}
	item := userVips[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_vips" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_vips.test", "vips.#", fmt.Sprintf("%d", len(userVips))),
					testCheckVIPPresent("data.zstack_vips.test", envStr(item, "uuid"), envStr(item, "name")),
				),
			},
		},
	})
}

func TestAccZStackVipDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	item, ok := preferredUserVIP(env.Vips)
	if !ok {
		t.Skip("no user vips in env data")
	}
	name := envStr(item, "name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_vips" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_vips.test", "vips.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_vips.test", "vips.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_vips.test", "vips.0.uuid", envStr(item, "uuid")),
				),
			},
		},
	})
}

func userVIPs(vips []map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(vips))
	for _, vip := range vips {
		if !envBool(vip, "system") {
			result = append(result, vip)
		}
	}
	return result
}

func preferredUserVIP(vips []map[string]interface{}) (map[string]interface{}, bool) {
	if vip, ok := userVIPByUseFor(vips, "LoadBalancer"); ok {
		return vip, true
	}
	userVips := userVIPs(vips)
	if len(userVips) == 0 {
		return nil, false
	}
	return userVips[0], true
}

func userVIPByUseFor(vips []map[string]interface{}, useFor string) (map[string]interface{}, bool) {
	for _, vip := range vips {
		if !envBool(vip, "system") && envStr(vip, "use_for") == useFor {
			return vip, true
		}
	}
	return nil, false
}

func testCheckVIPPresent(resourceName, uuid, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		count, err := strconv.Atoi(rs.Primary.Attributes["vips.#"])
		if err != nil {
			return fmt.Errorf("invalid %s vips count %q: %w", resourceName, rs.Primary.Attributes["vips.#"], err)
		}
		for i := 0; i < count; i++ {
			prefix := fmt.Sprintf("vips.%d", i)
			if rs.Primary.Attributes[prefix+".uuid"] == uuid && rs.Primary.Attributes[prefix+".name"] == name {
				return nil
			}
		}

		return fmt.Errorf("VIP %s (%s) not found in %s", name, uuid, resourceName)
	}
}
