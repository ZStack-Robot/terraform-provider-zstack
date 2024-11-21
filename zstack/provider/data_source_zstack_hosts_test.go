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

func TestAccZStackHostsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_hosts" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of hosts returned
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.#", "3"),

					// Verify the first Cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.architecture", "aarch64"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.cluster_uuid", "8087f700b6474a6fb916e0ba139f767c"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.name", "172.26.103.172"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.status", "Connected"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.type", "KVM"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.zone_uuid", "1df948fedd3b45dd89e9549348280e17"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.uuid", "73613363ae224154831f0c85f845856c"),
				),
			},
		},
	})
}

func TestAccZStackHostDataSourceFilterByname(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_hosts" "test" { name ="172.26.107.217" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of clusters returned
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.#", "1"),

					// Verify the first Cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.architecture", "aarch64"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.cluster_uuid", "8087f700b6474a6fb916e0ba139f767c"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.name", "172.26.107.217"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.status", "Connected"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.type", "KVM"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.zone_uuid", "1df948fedd3b45dd89e9549348280e17"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.uuid", "b5d0a89e52374621b190cc1978d7d586"),
				),
			},
		},
	})
}

func TestAccZStackHostDataSourceFilterBynamePattern(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_hosts" "test" { name_pattern="172.26.%" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of clusters returned
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.#", "2"),

					// Verify the first Cluster to ensure all attributes are set
					// Verify the first Cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.1.architecture", "aarch64"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.1.cluster_uuid", "8087f700b6474a6fb916e0ba139f767c"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.1.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.1.name", "172.26.107.217"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.1.status", "Connected"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.1.type", "KVM"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.1.zone_uuid", "1df948fedd3b45dd89e9549348280e17"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.1.uuid", "b5d0a89e52374621b190cc1978d7d586"),
				),
			},
		},
	})
}
