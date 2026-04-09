// Copyright (c) ZStack.io, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Note: Virtual router images are queried via ZQL with systemTag='applianceType::vrouter'.
// Not all environments have VR images. These tests only verify that the data source
// can be called without error; if no images exist, the check is lenient.

func TestAccZStackVirtualRouterImagesDataSource(t *testing.T) {
	_ = loadEnvData(t) // ensure env.json exists

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_virtual_router_images" "test" {}`,
			},
		},
	})
}

func TestAccZStackVirtualRouterImagesDataSourceFilterByName(t *testing.T) {
	_ = loadEnvData(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_virtual_router_images" "test" { name_pattern = "%" }`,
			},
		},
	})
}

func TestAccZStackVirtualRouterImagesDataSourceFilterByNamePattern(t *testing.T) {
	_ = loadEnvData(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_virtual_router_images" "test" { name_pattern = "%" }`,
			},
		},
	})
}
