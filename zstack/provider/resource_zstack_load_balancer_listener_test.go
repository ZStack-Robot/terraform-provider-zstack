// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestLoadBalancerListenerResource_Schema(t *testing.T) {
	var r loadBalancerListenerResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "load_balancer_uuid", "protocol", "load_balancer_port", "instance_port"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "server_group_uuid"}
	for _, attr := range computed {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing computed attribute %q", attr)
		}
		if !a.IsComputed() {
			t.Errorf("attribute %q should be computed", attr)
		}
	}

	optional := []string{"description", "security_policy_type"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}
}

func TestLoadBalancerListenerResource_Metadata(t *testing.T) {
	var r loadBalancerListenerResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_load_balancer_listener" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccLoadBalancerListenerResource(t *testing.T) {
	env := loadEnvData(t)

	if len(env.L3Networks) == 0 {
		t.Skip("no l3_networks in env.json, skipping load balancer listener acceptance test")
	}

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLoadBalancerListenerDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
data "zstack_vips" "test" {
}

resource "zstack_load_balancer" "test" {
  name     = "acc-test-lb-for-listener"
  vip_uuid = data.zstack_vips.test.vips.0.uuid
}

resource "zstack_load_balancer_listener" "test" {
  name               = "acc-test-lb-listener"
  load_balancer_uuid = zstack_load_balancer.test.uuid
  protocol           = "tcp"
  load_balancer_port = 80
  instance_port      = 8080
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_load_balancer_listener.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_load_balancer_listener.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-lb-listener")),
					statecheck.ExpectKnownValue("zstack_load_balancer_listener.test", tfjsonpath.New("protocol"), knownvalue.StringExact("tcp")),
					statecheck.ExpectKnownValue("zstack_load_balancer_listener.test", tfjsonpath.New("load_balancer_port"), knownvalue.StringExact("80")),
					statecheck.ExpectKnownValue("zstack_load_balancer_listener.test", tfjsonpath.New("instance_port"), knownvalue.StringExact("8080")),
					statecheck.ExpectKnownValue("zstack_load_balancer_listener.test", tfjsonpath.New("load_balancer_uuid"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:      "zstack_load_balancer_listener.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
