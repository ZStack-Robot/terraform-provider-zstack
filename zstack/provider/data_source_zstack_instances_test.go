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

func TestAccZStackvmInstancesDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.VmInstances) == 0 {
		t.Skip("no vm instances in env data")
	}
	vm := env.VmInstances[0]

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_instances" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances"), knownvalue.ListSizeExact(len(env.VmInstances))),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(envStr(vm, "name"))),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(vm, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("architecture"), knownvalue.StringExact(envStr(vm, "architecture"))),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("cluster_uuid"), knownvalue.StringExact(envStr(vm, "cluster_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("cpu_num"), knownvalue.StringExact(envStr(vm, "cpu_num"))),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("hypervisor_type"), knownvalue.StringExact(envStr(vm, "hypervisor_type"))),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("image_uuid"), knownvalue.StringExact(envStr(vm, "image_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("memory_size"), knownvalue.StringExact(envStr(vm, "memory_size"))),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("platform"), knownvalue.StringExact(envStr(vm, "platform"))),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(vm, "state"))),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(vm, "type"))),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("zone_uuid"), knownvalue.StringExact(envStr(vm, "zone_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("all_volumes").AtSliceIndex(0).AtMapKey("volume_uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("vm_nics").AtSliceIndex(0).AtMapKey("ip"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccZStackvmInstancesDataSourceFilterByNameRegex(t *testing.T) {
	env := loadEnvData(t)
	if len(env.VmInstances) == 0 {
		t.Skip("no vm instances in env data")
	}
	vm := env.VmInstances[0]
	name := envStr(vm, "name")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_instances" "test" { name = %q }`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(vm, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(vm, "state"))),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(vm, "type"))),
				},
			},
		},
	})
}

func TestAccZStackvmInstancesDataSourceFilterByNamePattern(t *testing.T) {
	env := loadEnvData(t)
	if len(env.VmInstances) == 0 {
		t.Skip("no vm instances in env data")
	}
	vm := env.VmInstances[0]
	name := envStr(vm, "name")
	// Use first 3 runes for pattern (handles CJK characters)
	runes := []rune(name)
	pattern := string(runes[:3]) + "%"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_instances" "test" { name = %q }`, pattern),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(vm, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_instances.test", tfjsonpath.New("vminstances").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(vm, "state"))),
				},
			},
		},
	})
}
