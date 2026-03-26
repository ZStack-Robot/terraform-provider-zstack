// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackImageDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Images) == 0 {
		t.Skip("no images in env data")
	}
	img := env.Images[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_images" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.#", fmt.Sprintf("%d", len(env.Images))),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.name", envStr(img, "name")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.architecture", envStr(img, "architecture")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.format", envStr(img, "format")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.guest_os_type", envStr(img, "guest_os_type")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.platform", envStr(img, "platform")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.state", envStr(img, "state")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.status", envStr(img, "status")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.uuid", envStr(img, "uuid")),
				),
			},
		},
	})
}

func TestAccZStackImageDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Images) == 0 {
		t.Skip("no images in env data")
	}
	img := env.Images[0]
	name := envStr(img, "name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_images" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.uuid", envStr(img, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.architecture", envStr(img, "architecture")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.format", envStr(img, "format")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.guest_os_type", envStr(img, "guest_os_type")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.platform", envStr(img, "platform")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.state", envStr(img, "state")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.status", envStr(img, "status")),
				),
			},
		},
	})
}

func TestAccZStackImageDataSourceFilterByNamePattern(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Images) == 0 {
		t.Skip("no images in env data")
	}
	img := env.Images[0]
	name := envStr(img, "name")
	pattern := name[:3] + "%"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_images" "test" { name_pattern = %q }`, pattern),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.uuid", envStr(img, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.architecture", envStr(img, "architecture")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.format", envStr(img, "format")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.guest_os_type", envStr(img, "guest_os_type")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.platform", envStr(img, "platform")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.state", envStr(img, "state")),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.status", envStr(img, "status")),
				),
			},
		},
	})
}
