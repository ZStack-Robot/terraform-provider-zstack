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
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.#", "8"),

					// Verify the first Cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.hypervisortype", "baremetal"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.name", "信创"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.type", "baremetal"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.zoneuuid", "8aa8ddb83f2c47088791478dfbbe5f65"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.uuid", "70b8ae9c48d4494593af4efe3b23f7f7"),
				),
			},
		},
	})
}
