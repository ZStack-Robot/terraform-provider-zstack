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

func TestAccZStackvmInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_instances" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of vm instances returned
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.#", "1"),

					// Verify the first vm instances to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.architecture", "aarch64"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.cluster_uuid", "8087f700b6474a6fb916e0ba139f767c"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.cpu_num", "2"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.host_uuid", "fd4d0c9f335c47e6a804345dc4a74601"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.hypervisor_type", "KVM"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.image_uuid", "4c2b64fe5d1b43ccb15410be32227076"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.memory_size", "2147483648"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.name", "ZStack组件服务监控套件"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.platform", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.state", "Running"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.type", "UserVm"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.uuid", "7fea17727cff4aae9aed09cd6e5f3f67"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.zone_uuid", "1df948fedd3b45dd89e9549348280e17"),

					// Verify the first volume of vm instance to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_actual_size", "304099328"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_description", "Root volume for VM[uuid:7fea17727cff4aae9aed09cd6e5f3f67]"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_format", "qcow2"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_size", "107374182400"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_status", "Ready"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_type", "Root"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_uuid", "b5850a9e9bdd409e909ea206857565df"),

					// Verify the first nic of vm instance to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.vm_nics.0.ip", "172.26.111.252"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.vm_nics.0.mac", "fa:9a:6d:16:84:00"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.vm_nics.0.netmask", "255.255.0.0"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.vm_nics.0.gateway", "172.26.0.1"),
				),
			},
		},
	})
}

func TestAccZStackvmInstancesDataSourceFilterByNameRegex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_instances" "test" { name = "ZStack组件服务监控套件"}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of vm instances returned
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.#", "1"),

					// Verify the first vm instances to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.architecture", "aarch64"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.cluster_uuid", "8087f700b6474a6fb916e0ba139f767c"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.cpu_num", "2"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.host_uuid", "fd4d0c9f335c47e6a804345dc4a74601"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.hypervisor_type", "KVM"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.image_uuid", "4c2b64fe5d1b43ccb15410be32227076"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.memory_size", "2147483648"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.name", "ZStack组件服务监控套件"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.platform", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.state", "Running"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.type", "UserVm"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.uuid", "7fea17727cff4aae9aed09cd6e5f3f67"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.zone_uuid", "1df948fedd3b45dd89e9549348280e17"),

					// Verify the first volume of vm instance to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_actual_size", "304099328"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_description", "Root volume for VM[uuid:7fea17727cff4aae9aed09cd6e5f3f67]"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_format", "qcow2"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_size", "107374182400"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_status", "Ready"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_type", "Root"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_uuid", "b5850a9e9bdd409e909ea206857565df"),

					// Verify the first nic of vm instance to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.vm_nics.0.ip", "172.26.111.252"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.vm_nics.0.mac", "fa:9a:6d:16:84:00"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.vm_nics.0.netmask", "255.255.0.0"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.vm_nics.0.gateway", "172.26.0.1"),
				),
			},
		},
	})
}

func TestAccZStackvmInstancesDataSourceFilterByNamePattern(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_instances" "test" { name = "ZStack组件%"}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of vm instances returned
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.#", "1"),

					// Verify the first vm instances to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.architecture", "aarch64"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.cluster_uuid", "8087f700b6474a6fb916e0ba139f767c"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.cpu_num", "2"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.host_uuid", "fd4d0c9f335c47e6a804345dc4a74601"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.hypervisor_type", "KVM"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.image_uuid", "4c2b64fe5d1b43ccb15410be32227076"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.memory_size", "2147483648"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.name", "ZStack组件服务监控套件"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.platform", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.state", "Running"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.type", "UserVm"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.uuid", "7fea17727cff4aae9aed09cd6e5f3f67"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.zone_uuid", "1df948fedd3b45dd89e9549348280e17"),

					// Verify the first volume of vm instance to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_actual_size", "304099328"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_description", "Root volume for VM[uuid:7fea17727cff4aae9aed09cd6e5f3f67]"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_format", "qcow2"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_size", "107374182400"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_status", "Ready"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_type", "Root"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.all_volumes.0.volume_uuid", "b5850a9e9bdd409e909ea206857565df"),

					// Verify the first nic of vm instance to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.vm_nics.0.ip", "172.26.111.252"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.vm_nics.0.mac", "fa:9a:6d:16:84:00"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.vm_nics.0.netmask", "255.255.0.0"),
					resource.TestCheckResourceAttr("data.zstack_instances.test", "vminstances.0.vm_nics.0.gateway", "172.26.0.1"),
				),
			},
		},
	})
}
