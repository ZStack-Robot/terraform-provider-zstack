// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccZStackImageDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Images) == 0 {
		t.Skip("no images in env data")
	}
	img := env.Images[0]

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_images" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images"), knownvalue.ListSizeExact(len(env.Images))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(envStr(img, "name"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("architecture"), knownvalue.StringExact(envStr(img, "architecture"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("format"), knownvalue.StringExact(envStr(img, "format"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("guest_os_type"), knownvalue.StringExact(envStr(img, "guest_os_type"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("platform"), knownvalue.StringExact(envStr(img, "platform"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(img, "state"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("status"), knownvalue.StringExact(envStr(img, "status"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(img, "uuid"))),
				},
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

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_images" "test" { name = %q }`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(img, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("architecture"), knownvalue.StringExact(envStr(img, "architecture"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("format"), knownvalue.StringExact(envStr(img, "format"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("guest_os_type"), knownvalue.StringExact(envStr(img, "guest_os_type"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("platform"), knownvalue.StringExact(envStr(img, "platform"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(img, "state"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("status"), knownvalue.StringExact(envStr(img, "status"))),
				},
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

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_images" "test" { name_pattern = %q }`, pattern),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(img, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("architecture"), knownvalue.StringExact(envStr(img, "architecture"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("format"), knownvalue.StringExact(envStr(img, "format"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("guest_os_type"), knownvalue.StringExact(envStr(img, "guest_os_type"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("platform"), knownvalue.StringExact(envStr(img, "platform"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(img, "state"))),
					statecheck.ExpectKnownValue("data.zstack_images.test", tfjsonpath.New("images").AtSliceIndex(0).AtMapKey("status"), knownvalue.StringExact(envStr(img, "status"))),
				},
			},
		},
	})
}
