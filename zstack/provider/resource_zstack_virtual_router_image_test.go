// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestVirtualRouterImageResource_Schema(t *testing.T) {
	var r virtualRouterImageResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "url"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid"}
	for _, attr := range computed {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing computed attribute %q", attr)
		}
		if !a.IsComputed() {
			t.Errorf("attribute %q should be computed", attr)
		}
	}

	optional := []string{"description", "backup_storage_uuids", "virtio", "boot_mode", "guest_os_type", "platform", "architecture"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}
}

func TestVirtualRouterImageResource_Metadata(t *testing.T) {
	var r virtualRouterImageResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_virtual_router_image" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccVirtualRouterImageResource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.BackupStorages) == 0 {
		t.Skip("no backup storages in env data")
	}
	bsUUID := envStr(env.BackupStorages[0], "uuid")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVirtualRouterImageDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_virtual_router_image" "test" {
  name                 = "acc-test-vr-image"
  url                  = "http://192.168.200.100/mirror/diskimages/CentOS-7-x86_64-Cloudinit-8G-official.qcow2"
  platform             = "Linux"
  architecture         = "x86_64"
  backup_storage_uuids = [%q]
}
`, bsUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_virtual_router_image.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_virtual_router_image.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-vr-image")),
				},
			},
			{
				ResourceName:            "zstack_virtual_router_image.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"backup_storage_uuids", "architecture", "guest_os_type", "virtio", "boot_mode", "media_type"},
			},
		},
	})
}
