// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackVirtualRoutersDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.VirtualRouters) == 0 {
		t.Skip("no virtual routers in env data")
	}
	vr := env.VirtualRouters[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_virtual_routers" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.#", fmt.Sprintf("%d", len(env.VirtualRouters))),

					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.name", envStr(vr, "name")),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.uuid", envStr(vr, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.hypervisor_type", envStr(vr, "hypervisor_type")),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.appliance_vm_type", envStr(vr, "appliance_vm_type")),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.state", envStr(vr, "state")),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.status", envStr(vr, "status")),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.agent_port", envStr(vr, "agent_port")),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.type", envStr(vr, "type")),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.ha_status", envStr(vr, "ha_status")),

					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.zone_uuid", envStr(vr, "zone_uuid")),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.cluster_uuid", envStr(vr, "cluster_uuid")),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.management_network_uuid", envStr(vr, "management_network_uuid")),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.image_uuid", envStr(vr, "image_uuid")),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.host_uuid", envStr(vr, "host_uuid")),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.instance_offering_uuid", envStr(vr, "instance_offering_uuid")),

					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.platform", envStr(vr, "platform")),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.architecture", envStr(vr, "architecture")),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.cpu_num", envStr(vr, "cpu_num")),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.memory_size", envStr(vr, "memory_size")),

					resource.TestCheckResourceAttrSet("data.zstack_virtual_routers.test", "virtual_router.0.vm_nics.0.ip"),
					resource.TestCheckResourceAttrSet("data.zstack_virtual_routers.test", "virtual_router.0.vm_nics.0.mac"),
					resource.TestCheckResourceAttrSet("data.zstack_virtual_routers.test", "virtual_router.0.vm_nics.0.netmask"),
					resource.TestCheckResourceAttrSet("data.zstack_virtual_routers.test", "virtual_router.0.vm_nics.0.gateway"),
				),
			},
		},
	})
}
