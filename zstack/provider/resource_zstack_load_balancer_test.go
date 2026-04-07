// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"

	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLoadBalancerResource_Schema(t *testing.T) {
	var r loadBalancerResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "vip_uuid"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "state", "type", "server_group_uuid"}
	for _, attr := range computed {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing computed attribute %q", attr)
		}
		if !a.IsComputed() {
			t.Errorf("attribute %q should be computed", attr)
		}
	}

	optional := []string{"description"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}
}

func TestLoadBalancerResource_Metadata(t *testing.T) {
	var r loadBalancerResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_load_balancer" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccLoadBalancerResource(t *testing.T) {
	env := loadEnvData(t)

	if len(env.L3Networks) == 0 {
		t.Skip("no l3_networks in env.json, skipping load balancer acceptance test")
	}

	tfresource.Test(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
data "zstack_vips" "test" {
}

resource "zstack_load_balancer" "test" {
  name     = "acc-test-lb"
  vip_uuid = data.zstack_vips.test.vips.0.uuid
}
`,
				Check: tfresource.ComposeAggregateTestCheckFunc(
					tfresource.TestCheckResourceAttrSet("zstack_load_balancer.test", "uuid"),
					tfresource.TestCheckResourceAttr("zstack_load_balancer.test", "name", "acc-test-lb"),
					tfresource.TestCheckResourceAttrSet("zstack_load_balancer.test", "vip_uuid"),
				),
			},
		},
	})
}
