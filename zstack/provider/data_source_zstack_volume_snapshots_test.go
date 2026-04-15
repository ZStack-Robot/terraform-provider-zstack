// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackVolumeSnapshotDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.VolumeSnapshots) == 0 {
		t.Skip("no volume_snapshots in env data")
	}
	item := env.VolumeSnapshots[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_volume_snapshots" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_volume_snapshots.test", "snapshots.#", fmt.Sprintf("%d", len(env.VolumeSnapshots))),
					resource.TestCheckResourceAttr("data.zstack_volume_snapshots.test", "snapshots.0.name", envStr(item, "name")),
					resource.TestCheckResourceAttr("data.zstack_volume_snapshots.test", "snapshots.0.uuid", envStr(item, "uuid")),
				),
			},
		},
	})
}

func TestAccZStackVolumeSnapshotDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	if len(env.VolumeSnapshots) == 0 {
		t.Skip("no volume_snapshots in env data")
	}
	item := env.VolumeSnapshots[0]
	name := envStr(item, "name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_volume_snapshots" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_volume_snapshots.test", "snapshots.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_volume_snapshots.test", "snapshots.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_volume_snapshots.test", "snapshots.0.uuid", envStr(item, "uuid")),
				),
			},
		},
	})
}
