// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

func TestAccZStackVolumeDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Volumes) == 0 {
		t.Skip("no volumes in env data")
	}
	_ = env

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_volumes" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.zstack_volumes.test", "volumes.#"),
					resource.TestCheckResourceAttrSet("data.zstack_volumes.test", "volumes.0.name"),
					resource.TestCheckResourceAttrSet("data.zstack_volumes.test", "volumes.0.uuid"),
				),
			},
		},
	})
}

func TestAccZStackVolumeDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Volumes) == 0 {
		t.Skip("no volumes in env data")
	}
	item := env.Volumes[0]
	name := envStr(item, "name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_volumes" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_volumes.test", "volumes.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_volumes.test", "volumes.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_volumes.test", "volumes.0.uuid", envStr(item, "uuid")),
				),
			},
		},
	})
}

func TestAccZStackVolumeDataSourceFilterByStorageAndVM(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("acceptance test skipped unless TF_ACC is set")
	}

	cli := testAccClientLoggedIn()
	params := param.NewQueryParam()
	volumes, err := cli.QueryVolume(&params)
	if err != nil {
		t.Fatalf("QueryVolume error: %v", err)
	}

	var selectedUUID string
	var primaryStorageUuid string
	var vmInstanceUuid string
	for _, volume := range volumes {
		if volume.PrimaryStorageUuid != "" && volume.VmInstanceUuid != "" {
			selectedUUID = volume.UUID
			primaryStorageUuid = volume.PrimaryStorageUuid
			vmInstanceUuid = volume.VmInstanceUuid
			break
		}
	}
	if selectedUUID == "" {
		t.Skip("no volume with primary_storage_uuid and vm_instance_uuid in live environment")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_volumes" "test" {
  filter {
    name   = "primary_storage_uuid"
    values = [%q]
  }

  filter {
    name   = "vm_instance_uuid"
    values = [%q]
  }

  filter {
    name   = "uuid"
    values = [%q]
  }
}
`, primaryStorageUuid, vmInstanceUuid, selectedUUID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_volumes.test", "volumes.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_volumes.test", "volumes.0.uuid", selectedUUID),
					resource.TestCheckResourceAttr("data.zstack_volumes.test", "volumes.0.primary_storage_uuid", primaryStorageUuid),
					resource.TestCheckResourceAttr("data.zstack_volumes.test", "volumes.0.vm_instance_uuid", vmInstanceUuid),
				),
			},
		},
	})
}
