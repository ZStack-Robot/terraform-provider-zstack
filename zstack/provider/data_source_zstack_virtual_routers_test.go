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

var _ = fmt.Sprintf

func TestAccZStackVirtualRoutersDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.VirtualRouters) == 0 {
		t.Skip("no virtual routers in env data")
	}
	vr := env.VirtualRouters[0]

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_virtual_routers" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router"), knownvalue.ListSizeExact(len(env.VirtualRouters))),

					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(envStr(vr, "name"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(vr, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("hypervisor_type"), knownvalue.StringExact(envStr(vr, "hypervisor_type"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("appliance_vm_type"), knownvalue.StringExact(envStr(vr, "appliance_vm_type"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(vr, "state"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("status"), knownvalue.StringExact(envStr(vr, "status"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("agent_port"), knownvalue.StringExact(envStr(vr, "agent_port"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(vr, "type"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("ha_status"), knownvalue.StringExact(envStr(vr, "ha_status"))),

					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("zone_uuid"), knownvalue.StringExact(envStr(vr, "zone_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("cluster_uuid"), knownvalue.StringExact(envStr(vr, "cluster_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("management_network_uuid"), knownvalue.StringExact(envStr(vr, "management_network_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("image_uuid"), knownvalue.StringExact(envStr(vr, "image_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("host_uuid"), knownvalue.StringExact(envStr(vr, "host_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("instance_offering_uuid"), knownvalue.StringExact(envStr(vr, "instance_offering_uuid"))),

					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("platform"), knownvalue.StringExact(envStr(vr, "platform"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("architecture"), knownvalue.StringExact(envStr(vr, "architecture"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("cpu_num"), knownvalue.StringExact(envStr(vr, "cpu_num"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("memory_size"), knownvalue.StringExact(envStr(vr, "memory_size"))),

					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("vm_nics").AtSliceIndex(0).AtMapKey("ip"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("vm_nics").AtSliceIndex(0).AtMapKey("mac"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("vm_nics").AtSliceIndex(0).AtMapKey("netmask"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.zstack_virtual_routers.test", tfjsonpath.New("virtual_router").AtSliceIndex(0).AtMapKey("vm_nics").AtSliceIndex(0).AtMapKey("gateway"), knownvalue.NotNull()),
				},
			},
		},
	})
}
