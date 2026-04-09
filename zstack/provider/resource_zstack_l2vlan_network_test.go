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

func TestL2VlanNetworkResource_Schema(t *testing.T) {
	var r l2VlanNetworkResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "vlan", "zone_uuid"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "type"}
	for _, attr := range computed {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing computed attribute %q", attr)
		}
		if !a.IsComputed() {
			t.Errorf("attribute %q should be computed", attr)
		}
	}

	optional := []string{"description", "physical_interface", "vswitch_type", "attached_cluster_uuids"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}
}

func TestL2VlanNetworkResource_Metadata(t *testing.T) {
	var r l2VlanNetworkResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_l2vlan_network" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccL2VlanNetworkResource(t *testing.T) {
	env := loadEnvData(t)

	if len(env.Zones) == 0 {
		t.Skip("no zones in env.json, skipping L2 VLAN network acceptance test")
	}

	zoneUuid := envStr(env.Zones[0], "uuid")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckL2VlanNetworkDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_l2vlan_network" "test" {
  name              = "acc-test-l2vlan"
  vlan              = 3999
  zone_uuid         = "` + zoneUuid + `"
  physical_interface = "eth0"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_l2vlan_network.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_l2vlan_network.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-l2vlan")),
					statecheck.ExpectKnownValue("zstack_l2vlan_network.test", tfjsonpath.New("vlan"), knownvalue.StringExact("3999")),
				},
			},
			{
				ResourceName:      "zstack_l2vlan_network.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
