// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackvmInstancesDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.VmInstances) == 0 {
		t.Skip("no vm instances in env data")
	}
	vm := env.VmInstances[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_instances" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.#", fmt.Sprintf("%d", len(env.VmInstances))),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.name", envStr(vm, "name")),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.uuid", envStr(vm, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.architecture", envStr(vm, "architecture")),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.cluster_uuid", envStr(vm, "cluster_uuid")),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.cpu_num", envStr(vm, "cpu_num")),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.hypervisor_type", envStr(vm, "hypervisor_type")),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.image_uuid", envStr(vm, "image_uuid")),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.memory_size", envStr(vm, "memory_size")),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.platform", envStr(vm, "platform")),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.state", envStr(vm, "state")),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.type", envStr(vm, "type")),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.zone_uuid", envStr(vm, "zone_uuid")),
					resource.TestCheckResourceAttrSet("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_uuid"),
					resource.TestCheckResourceAttrSet("data.zstack_instances.test", "vminstances.0.vm_nics.0.ip"),
				),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_instances" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.uuid", envStr(vm, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.state", envStr(vm, "state")),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.type", envStr(vm, "type")),
				),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_instances" "test" { name = %q }`, pattern),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.uuid", envStr(vm, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.state", envStr(vm, "state")),
				),
			},
		},
	})
}
