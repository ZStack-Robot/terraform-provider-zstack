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

func TestPortForwardingRuleResource_Schema(t *testing.T) {
	var r portForwardingRuleResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "vip_uuid", "vip_port_start", "protocol_type"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "vip_ip", "guest_ip", "state"}
	for _, attr := range computed {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing computed attribute %q", attr)
		}
		if !a.IsComputed() {
			t.Errorf("attribute %q should be computed", attr)
		}
	}

	optional := []string{"description", "vip_port_end", "private_port_start", "private_port_end", "vm_nic_uuid", "allowed_cidr"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}
}

func TestPortForwardingRuleResource_Metadata(t *testing.T) {
	var r portForwardingRuleResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_port_forwarding_rule" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccPortForwardingRuleResource(t *testing.T) {
	env := loadEnvData(t)

	// Port forwarding requires a VIP. Try to find one from env.json.
	// If VIPs are not available, we create a minimal config that references a VIP data source.
	if len(env.L3Networks) == 0 {
		t.Skip("no l3_networks in env.json, skipping port forwarding rule acceptance test")
	}

	// Use a pre-existing VIP if available via data source
	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPortForwardingRuleDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
data "zstack_vips" "test" {
}

resource "zstack_port_forwarding_rule" "test" {
  name           = "acc-test-pf-rule"
  vip_uuid       = data.zstack_vips.test.vips.0.uuid
  vip_port_start = 8080
  protocol_type  = "TCP"
  allowed_cidr   = "0.0.0.0/0"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_port_forwarding_rule.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_port_forwarding_rule.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-pf-rule")),
					statecheck.ExpectKnownValue("zstack_port_forwarding_rule.test", tfjsonpath.New("protocol_type"), knownvalue.StringExact("TCP")),
					statecheck.ExpectKnownValue("zstack_port_forwarding_rule.test", tfjsonpath.New("vip_port_start"), knownvalue.StringExact("8080")),
				},
			},
			{
				ResourceName:      "zstack_port_forwarding_rule.test",
				ImportState:       true,
				ImportStateIdFunc:       importStateUUID("zstack_port_forwarding_rule.test"),
				ImportStateVerify: true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
