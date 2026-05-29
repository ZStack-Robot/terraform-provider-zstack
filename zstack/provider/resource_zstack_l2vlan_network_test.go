// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"os"
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

func TestAccL2VlanNetworkResource_disappears(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("acceptance test skipped unless TF_ACC is set")
	}

	_, zoneUuid := testAccLiveCluster(t)
	name := testAccName("l2vlan-disappears")
	vlan := testAccFreeL2Vlan(t, zoneUuid, "eth0")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckL2VlanNetworkDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_l2vlan_network" "test" {
  name              = %q
  vlan              = %d
  zone_uuid         = %q
  physical_interface = "eth0"
}
`, name, vlan, zoneUuid),
				ConfigStateChecks: []statecheck.StateCheck{
					stateCheckL2VlanNetworkDisappears("zstack_l2vlan_network.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccL2VlanNetworkResource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("acceptance test skipped unless TF_ACC is set")
	}

	_, zoneUuid := testAccLiveCluster(t)
	name := testAccName("l2vlan")
	updatedName := name + "-updated"
	vlan := testAccFreeL2Vlan(t, zoneUuid, "eth0")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckL2VlanNetworkDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_l2vlan_network" "test" {
  name              = %q
  description       = "acceptance l2 vlan network"
  vlan              = %d
  zone_uuid         = %q
  physical_interface = "eth0"
}
`, name, vlan, zoneUuid),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_l2vlan_network.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_l2vlan_network.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_l2vlan_network.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance l2 vlan network")),
					statecheck.ExpectKnownValue("zstack_l2vlan_network.test", tfjsonpath.New("vlan"), knownvalue.Int64Exact(int64(vlan))),
				},
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_l2vlan_network" "test" {
  name              = %q
  description       = "acceptance l2 vlan network updated"
  vlan              = %d
  zone_uuid         = %q
  physical_interface = "eth0"
}
`, updatedName, vlan, zoneUuid),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_l2vlan_network.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("zstack_l2vlan_network.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance l2 vlan network updated")),
					statecheck.ExpectKnownValue("zstack_l2vlan_network.test", tfjsonpath.New("vlan"), knownvalue.Int64Exact(int64(vlan))),
				},
			},
			{
				ResourceName:                         "zstack_l2vlan_network.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_l2vlan_network.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
