// Copyright (c) ZStack.io, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackTagDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.SystemTags) == 0 {
		t.Skip("no system_tags in env data")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_tags" "test" { tag_type = "system" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.zstack_tags.test", "system_tags.#"),
				),
			},
		},
	})
}

func TestAccZStackTagDataSourceUserTags(t *testing.T) {
	env := loadEnvData(t)
	if len(env.UserTags) == 0 {
		t.Skip("no user_tags in env data")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_tags" "test" { tag_type = "user" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.zstack_tags.test", "user_tags.#"),
				),
			},
		},
	})
}
