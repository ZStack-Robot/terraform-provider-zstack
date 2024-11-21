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

func TestAccZStackBackupStorageDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_backupstorages" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of backupstorage returned
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.#", "1"),

					// Verify the first backupstorage to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.name", "imagestore"),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.status", "Connected"),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.available_capacity", "395321520128"),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.total_capacity", "464423182336"),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.uuid", "530c16460d974b8da73edae3d7b7ac41"),
				),
			},
		},
	})
}

func TestAccZStackBackupStorageDataSourceFilterByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_backupstorages" "test" { name ="imagestore" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of clusters returned
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.#", "1"),

					// Verify the first Cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.name", "imagestore"),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.status", "Connected"),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.available_capacity", "395321421824"),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.total_capacity", "464423182336"),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.uuid", "530c16460d974b8da73edae3d7b7ac41"),
				),
			},
		},
	})
}

func TestAccZStackBackupStorageDataSourceFilterByNamePattern(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_backupstorages" "test" { name_pattern ="imag%" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of clusters returned
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.#", "1"),

					// Verify the first Cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.name", "imagestore"),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.status", "Connected"),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.available_capacity", "395321405440"),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.total_capacity", "464423182336"),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.uuid", "530c16460d974b8da73edae3d7b7ac41"),
				),
			},
		},
	})
}
