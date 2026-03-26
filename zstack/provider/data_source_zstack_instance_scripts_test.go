// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackInstanceScriptsDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.InstanceScripts) == 0 {
		t.Skip("no instance scripts in env data")
	}
	s := env.InstanceScripts[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_scripts" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.#", fmt.Sprintf("%d", len(env.InstanceScripts))),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.name", envStr(s, "name")),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.uuid", envStr(s, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.platform", envStr(s, "platform")),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.script_timeout", envStr(s, "script_timeout")),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.script_type", envStr(s, "script_type")),
				),
			},
		},
	})
}

func TestAccZStackInstanceScriptsDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	if len(env.InstanceScripts) == 0 {
		t.Skip("no instance scripts in env data")
	}
	s := env.InstanceScripts[0]
	name := envStr(s, "name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_scripts" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.uuid", envStr(s, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.platform", envStr(s, "platform")),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.script_timeout", envStr(s, "script_timeout")),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.script_type", envStr(s, "script_type")),
				),
			},
		},
	})
}

func TestAccZStackInstanceScriptsDataSourceFilterByNamePattern(t *testing.T) {
	env := loadEnvData(t)
	if len(env.InstanceScripts) == 0 {
		t.Skip("no instance scripts in env data")
	}
	s := env.InstanceScripts[0]
	name := envStr(s, "name")
	pattern := string([]rune(name)[:3]) + "%"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_scripts" "test" { name_pattern = %q }`, pattern),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.uuid", envStr(s, "uuid")),
				),
			},
		},
	})
}
