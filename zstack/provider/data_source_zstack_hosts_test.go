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
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.#", "20"),

					// Verify the first Cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.architecture", "aarch64"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.cluster_uuid", "e6e983b6dacc4a12b2f24047f7d52db7"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.name", "172.26.100.31"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.status", "Connected"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.type", "KVM"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.zone_uuid", "8aa8ddb83f2c47088791478dfbbe5f65"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.uuid", "01b243dc3caa4dada0bc20bb545887b6"),
				),
			},
		},
	})
}

func TestAccZStackHostDataSourceFilterByname_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_hosts" "test" { name_regex="172.26.100.31" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of clusters returned
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.#", "1"),

					// Verify the first Cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.architecture", "aarch64"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.cluster_uuid", "e6e983b6dacc4a12b2f24047f7d52db7"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.name", "172.26.100.31"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.status", "Connected"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.type", "KVM"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.zone_uuid", "8aa8ddb83f2c47088791478dfbbe5f65"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.uuid", "01b243dc3caa4dada0bc20bb545887b6"),
				),
			},
		},
	})
}
