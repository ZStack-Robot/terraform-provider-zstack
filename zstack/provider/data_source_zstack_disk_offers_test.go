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

func TestAccZStackDiskOffersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_disk_offers" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.#", "3"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.name", "mediumDiskOffering"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.uuid", "12eafeadb422451f944e15f78658f629"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.allocator_strategy", "DefaultPrimaryStorageAllocationStrategy"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.description", "Medium Disk Offering"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.disk_size", "104857600"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.sort_key", "0"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.type", "DefaultDiskOfferingType"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.state", "Enabled"),
				),
			},
		},
	})
}

func TestAccZStackDiskOffersDataSourceFilterByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_disk_offers" "test" { name="mediumDiskOffering" }`,
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.name", "mediumDiskOffering"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.uuid", "12eafeadb422451f944e15f78658f629"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.allocator_strategy", "DefaultPrimaryStorageAllocationStrategy"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.description", "Medium Disk Offering"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.disk_size", "104857600"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.sort_key", "0"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.type", "DefaultDiskOfferingType"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.state", "Enabled"),
				),
			},
		},
	})
}

func TestAccZStackDiskOffersDataSourceFilterByNamePattern(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_disk_offers" "test" {name_pattern="m%"}`,
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.name", "mediumDiskOffering"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.uuid", "12eafeadb422451f944e15f78658f629"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.allocator_strategy", "DefaultPrimaryStorageAllocationStrategy"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.description", "Medium Disk Offering"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.disk_size", "104857600"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.sort_key", "0"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.type", "DefaultDiskOfferingType"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.state", "Enabled"),
				),
			},
		},
	})
}
