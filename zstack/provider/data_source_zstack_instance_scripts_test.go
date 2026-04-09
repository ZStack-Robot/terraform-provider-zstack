// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccZStackInstanceScriptsDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.InstanceScripts) == 0 {
		t.Skip("no instance scripts in env data")
	}
	s := env.InstanceScripts[0]

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_scripts" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_scripts.test", tfjsonpath.New("scripts"), knownvalue.ListSizeExact(len(env.InstanceScripts))),
					statecheck.ExpectKnownValue("data.zstack_scripts.test", tfjsonpath.New("scripts").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(envStr(s, "name"))),
					statecheck.ExpectKnownValue("data.zstack_scripts.test", tfjsonpath.New("scripts").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(s, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_scripts.test", tfjsonpath.New("scripts").AtSliceIndex(0).AtMapKey("platform"), knownvalue.StringExact(envStr(s, "platform"))),
					statecheck.ExpectKnownValue("data.zstack_scripts.test", tfjsonpath.New("scripts").AtSliceIndex(0).AtMapKey("script_timeout"), knownvalue.StringExact(envStr(s, "script_timeout"))),
					statecheck.ExpectKnownValue("data.zstack_scripts.test", tfjsonpath.New("scripts").AtSliceIndex(0).AtMapKey("script_type"), knownvalue.StringExact(envStr(s, "script_type"))),
				},
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

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_scripts" "test" { name = %q }`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_scripts.test", tfjsonpath.New("scripts"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.zstack_scripts.test", tfjsonpath.New("scripts").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_scripts.test", tfjsonpath.New("scripts").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(s, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_scripts.test", tfjsonpath.New("scripts").AtSliceIndex(0).AtMapKey("platform"), knownvalue.StringExact(envStr(s, "platform"))),
					statecheck.ExpectKnownValue("data.zstack_scripts.test", tfjsonpath.New("scripts").AtSliceIndex(0).AtMapKey("script_timeout"), knownvalue.StringExact(envStr(s, "script_timeout"))),
					statecheck.ExpectKnownValue("data.zstack_scripts.test", tfjsonpath.New("scripts").AtSliceIndex(0).AtMapKey("script_type"), knownvalue.StringExact(envStr(s, "script_type"))),
				},
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

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_scripts" "test" { name_pattern = %q }`, pattern),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_scripts.test", tfjsonpath.New("scripts").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_scripts.test", tfjsonpath.New("scripts").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(s, "uuid"))),
				},
			},
		},
	})
}
