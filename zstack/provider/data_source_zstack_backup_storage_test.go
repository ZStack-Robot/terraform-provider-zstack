// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Run go testing with TF_ACC environment variable set. Edit vscode settings.json and insert
//   "go.testEnvVars": {
//        "TF_ACC": "1"
//   },

func TestAccZStackBackupStorageDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_backupstorage" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of backupstorage returned
					resource.TestCheckResourceAttr("data.zstack_backupstorage.test", "backupstorages.#", "2"),

					// Verify the first backupstorage to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_backupstorage.test", "backupstorages.0.name", "BS-1"),
					resource.TestCheckResourceAttr("data.zstack_backupstorage.test", "backupstorages.0.status", "Connected"),
					resource.TestCheckResourceAttr("data.zstack_backupstorage.test", "backupstorages.0.state", "Enabled"),
					//resource.TestCheckResourceAttr("data.zstack_backupstorage.test", "backupstorages.0.availablecapacity", "65421975552"),
					resource.TestCheckResourceAttr("data.zstack_backupstorage.test", "backupstorages.0.totalcapacity", "229057531904"),
					resource.TestCheckResourceAttr("data.zstack_backupstorage.test", "backupstorages.0.uuid", "7a912545634b4ddc86c40af82c14b452"),
				),
			},
		},
	})
}

func TestAccZStackBackupStorageDataSourceFilterByname_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_backupstorage" "test" { name_regex="image" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of clusters returned
					resource.TestCheckResourceAttr("data.zstack_backupstorage.test", "backupstorages.#", "1"),

					// Verify the first Cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_backupstorage.test", "backupstorages.0.name", "image"),
					resource.TestCheckResourceAttr("data.zstack_backupstorage.test", "backupstorages.0.status", "Connected"),
					resource.TestCheckResourceAttr("data.zstack_backupstorage.test", "backupstorages.0.state", "Enabled"),
					//resource.TestCheckResourceAttr("data.zstack_backupstorage.test", "backupstorages.0.availablecapacity", "65519099904"),
					resource.TestCheckResourceAttr("data.zstack_backupstorage.test", "backupstorages.0.totalcapacity", "229057531904"),
					resource.TestCheckResourceAttr("data.zstack_backupstorage.test", "backupstorages.0.uuid", "936d51b80fdb4a0ea9c742bcecab56e0"),
				),
			},
		},
	})
}
