// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestVirtualRouterOfferingResource_Schema(t *testing.T) {
	var r virtualRouterOfferingResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "cpu_num", "memory_size", "management_network_uuid", "zone_uuid", "image_uuid"}
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

	optional := []string{"description", "public_network_uuid", "is_default"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}
}

func TestVirtualRouterOfferingResource_Metadata(t *testing.T) {
	var r virtualRouterOfferingResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_virtual_router_offer" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccVirtualRouterOfferingResource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Zones) == 0 {
		t.Skip("no zones in env data")
	}

	// Find the management network UUID and a VR image UUID from existing VR offerings
	var mgmtNetUUID, imageUUID string
	if len(env.VirtualRouterOfferings) > 0 {
		mgmtNetUUID = envStr(env.VirtualRouterOfferings[0], "management_network_uuid")
		imageUUID = envStr(env.VirtualRouterOfferings[0], "image_uuid")
	}
	if mgmtNetUUID == "" || imageUUID == "" {
		t.Skip("no virtual router offerings in env data to derive management_network_uuid and image_uuid")
	}

	zoneUUID := envStr(env.Zones[0], "uuid")

	tfresource.Test(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_virtual_router_offer" "test" {
  name                    = "acc-test-vr-offering"
  cpu_num                 = 1
  memory_size             = 512
  zone_uuid               = %q
  management_network_uuid = %q
  image_uuid              = %q
}
`, zoneUUID, mgmtNetUUID, imageUUID),
				Check: tfresource.ComposeAggregateTestCheckFunc(
					tfresource.TestCheckResourceAttrSet("zstack_virtual_router_offer.test", "uuid"),
					tfresource.TestCheckResourceAttr("zstack_virtual_router_offer.test", "name", "acc-test-vr-offering"),
				),
			},
		},
	})
}
