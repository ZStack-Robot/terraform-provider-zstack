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

func TestAccZStackClusterDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_clusters" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of clusters returned
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.#", "1"),

					// Verify the first Cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.hypervisortype", "KVM"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.name", "cluster1"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.type", "zstack"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.zone_uuid", "1df948fedd3b45dd89e9549348280e17"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.uuid", "8087f700b6474a6fb916e0ba139f767c"),
				),
			},
		},
	})
}

func TestAccZStackClusterDataSourceFilterByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_clusters" "test" { name ="cluster1" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of clusters returned
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.#", "1"),

					// Verify the first Cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.hypervisortype", "KVM"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.name", "cluster1"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.type", "zstack"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.zone_uuid", "1df948fedd3b45dd89e9549348280e17"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.uuid", "8087f700b6474a6fb916e0ba139f767c"),
				),
			},
		},
	})
}

func TestAccZStackClusterDataSourceFilterByNamePattern(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_clusters" "test" { name_pattern = "clu%" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of clusters returned
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.#", "1"),

					// Verify the first Cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.hypervisortype", "KVM"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.name", "cluster1"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.type", "zstack"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.zone_uuid", "1df948fedd3b45dd89e9549348280e17"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.uuid", "8087f700b6474a6fb916e0ba139f767c"),
				),
			},
		},
	})
}
