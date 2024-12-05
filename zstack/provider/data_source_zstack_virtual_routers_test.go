// Copyright (c) ZStack.io, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Run go testing with TF_ACC environment variable set. Edit vscode settings.json and insert
//   "go.testEnvVars": {
//        "TF_ACC": "1"
//   },

func TestAccZStackVirtualRoutersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_virtual_routers" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of image returned
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.#", "1"),

					// Verify the first virtual router instances to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.name", "vr"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.uuid", "6d4dd9fe12ca4fb6a6626858b415f1a1"),

					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.hypervisor_type", "KVM"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.appliance_vm_type", "vpcvrouter"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.state", "Running"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.status", "Connected"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.agent_port", "7272"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.type", "ApplianceVm"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.ha_status", "NoHa"),

					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.zone_uuid", "d29f4847a99f4dea83bc446c8fe6e64c"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.cluster_uuid", "37c25209578c495ca176f60ad0cd97fa"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.management_network_uuid", "50e8c0d69681447fbe347c8dae2b1bef"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.image_uuid", "93005c8a2a314a489635eca8c30794d4"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.host_uuid", "66c622bef5b14b76a1a9992b87ddbe1c"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.instance_offering_uuid", "9b56381daf204657adba1a4f0ec7ef39"),

					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.platform", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.architecture", "x86_64"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.cpu_num", "1"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.memory_size", "1073741824"),

					// Verify the first nic of vm instance to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.vm_nics.0.ip", "172.30.9.60"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.vm_nics.0.mac", "fa:21:df:98:58:00"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.vm_nics.0.netmask", "255.255.0.0"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.vm_nics.0.gateway", "172.30.0.1"),

					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.vm_nics.1.ip", "192.168.110.1"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.vm_nics.1.mac", "fa:59:ac:82:b3:01"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.vm_nics.1.netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("data.zstack_virtual_routers.test", "virtual_router.0.vm_nics.1.gateway", "192.168.110.1"),
				),
			},
		},
	})
}
