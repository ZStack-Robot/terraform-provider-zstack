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
				Config: providerConfig + `data "zstack_zsclusters" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of clusters returned
					resource.TestCheckResourceAttr("data.zstack_zsclusters.test", "clusters.#", "1"),

					// Verify the first Cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_zsclusters.test", "clusters.0.hypervisortype", "KVM"),
					resource.TestCheckResourceAttr("data.zstack_zsclusters.test", "clusters.0.name", "Cluster-1"),
					resource.TestCheckResourceAttr("data.zstack_zsclusters.test", "clusters.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_zsclusters.test", "clusters.0.type", "zstack"),
					resource.TestCheckResourceAttr("data.zstack_zsclusters.test", "clusters.0.zoneuuid", "4981061bd27c42c7bc063c4c4529518a"),
					resource.TestCheckResourceAttr("data.zstack_zsclusters.test", "clusters.0.uuid", "b286789480254a208e6327136bb3dcb3"),
				),
			},
		},
	})
}
