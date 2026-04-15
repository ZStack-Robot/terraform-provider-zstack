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

func TestClusterResource_Schema(t *testing.T) {
	var r clusterResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "zone_uuid", "hypervisor_type"}
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
}

func TestClusterResource_Metadata(t *testing.T) {
	var r clusterResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_cluster" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccClusterResource(t *testing.T) {
	env := loadEnvData(t)

	if len(env.Zones) == 0 {
		t.Skip("no zones in env.json, skipping cluster acceptance test")
	}

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckClusterDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_cluster" "test" {
  name            = "acc-test-cluster"
  zone_uuid       = "%s"
  hypervisor_type = "KVM"
}
`, envStr(env.Zones[0], "uuid")),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_cluster.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_cluster.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-cluster")),
					statecheck.ExpectKnownValue("zstack_cluster.test", tfjsonpath.New("hypervisor_type"), knownvalue.StringExact("KVM")),
				},
			},
			{
				ResourceName:      "zstack_cluster.test",
				ImportState:       true,
				ImportStateIdFunc:       importStateUUID("zstack_cluster.test"),
				ImportStateVerify: true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}

func TestAccClusterResource_disappears(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Zones) == 0 {
		t.Skip("no zones in env data")
	}
	zoneUUID := envStr(env.Zones[0], "uuid")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckClusterDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_cluster" "test_disappears" {
  name            = "acc-test-cluster-disappears"
  description     = "Disappears test cluster"
  hypervisor_type = "KVM"
  zone_uuid       = %q
}
`, zoneUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_cluster.test_disappears", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					stateCheckClusterDisappears("zstack_cluster.test_disappears"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
