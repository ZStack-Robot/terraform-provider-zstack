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

func TestVipResource_Schema(t *testing.T) {
	var r vipResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "l3_network_uuid"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	if _, ok := resp.Schema.Attributes["uuid"]; !ok {
		t.Fatal("schema missing computed attribute uuid")
	}

	optional := []string{"description", "ip_range_uuid", "vip"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}
}

func TestVipResource_Metadata(t *testing.T) {
	var r vipResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_vip" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccVipResource(t *testing.T) {
	env := loadEnvData(t)

	// Find a Public L3 network
	var l3UUID string
	for _, l3 := range env.L3Networks {
		if envStr(l3, "category") == "Public" {
			l3UUID = envStr(l3, "uuid")
			break
		}
	}
	if l3UUID == "" {
		t.Skip("no Public L3 network found in env data")
	}

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVipDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_vip" "test" {
  name            = "acc-test-vip"
  l3_network_uuid = %q
}
`, l3UUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_vip.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_vip.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-vip")),
					statecheck.ExpectKnownValue("zstack_vip.test", tfjsonpath.New("l3_network_uuid"), knownvalue.StringExact(l3UUID)),
				},
			},
			{
				ResourceName:      "zstack_vip.test",
				ImportState:       true,
				ImportStateIdFunc:       importStateUUID("zstack_vip.test"),
				ImportStateVerify: true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
