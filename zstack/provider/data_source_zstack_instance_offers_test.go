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

func TestAccZStackInstanceOffersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_instance_offers" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.name", "InstanceOffering-1"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.uuid", "4fb8a154b03d418ea771ec74d3273da3"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.allocator_strategy", "LeastVmPreferredHostAllocatorStrategy"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.cpu_num", "1"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.cpu_speed", "0"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.memory_size", "1073741824"),

					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.sort_key", "0"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.type", "UserVm"),
				),
			},
		},
	})
}

func TestAccZStackInstanceOffersDataSourceFilterByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_instance_offers" "test" { name = "InstanceOffering-1"}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.name", "InstanceOffering-1"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.uuid", "4fb8a154b03d418ea771ec74d3273da3"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.allocator_strategy", "LeastVmPreferredHostAllocatorStrategy"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.cpu_num", "1"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.cpu_speed", "0"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.memory_size", "1073741824"),

					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.sort_key", "0"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.type", "UserVm"),
				),
			},
		},
	})
}

func TestAccZStackInstanceOffersDataSourceFilterByNamePattern(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_instance_offers" "test" { name_pattern = "InstanceO%"}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.name", "InstanceOffering-1"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.uuid", "4fb8a154b03d418ea771ec74d3273da3"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.allocator_strategy", "LeastVmPreferredHostAllocatorStrategy"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.cpu_num", "1"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.cpu_speed", "0"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.memory_size", "1073741824"),

					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.sort_key", "0"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.type", "UserVm"),
				),
			},
		},
	})
}
