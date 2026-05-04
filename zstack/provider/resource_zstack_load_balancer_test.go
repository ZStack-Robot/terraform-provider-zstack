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

func TestAccLoadBalancerResource_disappears(t *testing.T) {
	env := loadEnvData(t)
	name := testAccName("lb-disappears")

	if len(env.L3Networks) == 0 {
		t.Skip("no l3_networks in env.json, skipping load balancer acceptance test")
	}

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_vips" "test" {
}

resource "zstack_load_balancer" "test" {
  name     = %q
  vip_uuid = data.zstack_vips.test.vips.0.uuid
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					stateCheckLoadBalancerDisappears("zstack_load_balancer.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccLoadBalancerResource(t *testing.T) {
	env := loadEnvData(t)
	name := testAccName("lb")
	updatedName := name + "-updated"

	if len(env.L3Networks) == 0 {
		t.Skip("no l3_networks in env.json, skipping load balancer acceptance test")
	}

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_vips" "test" {
}

resource "zstack_load_balancer" "test" {
  name        = %q
  description = "acceptance load balancer"
  vip_uuid    = data.zstack_vips.test.vips.0.uuid
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_load_balancer.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_load_balancer.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_load_balancer.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance load balancer")),
					statecheck.ExpectKnownValue("zstack_load_balancer.test", tfjsonpath.New("vip_uuid"), knownvalue.NotNull()),
				},
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_vips" "test" {
}

resource "zstack_load_balancer" "test" {
  name        = %q
  description = "acceptance load balancer updated"
  vip_uuid    = data.zstack_vips.test.vips.0.uuid
}
`, updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_load_balancer.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("zstack_load_balancer.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance load balancer updated")),
					statecheck.ExpectKnownValue("zstack_load_balancer.test", tfjsonpath.New("vip_uuid"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:                         "zstack_load_balancer.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_load_balancer.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
