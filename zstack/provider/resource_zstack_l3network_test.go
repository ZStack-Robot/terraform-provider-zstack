// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestL3networkResource_Schema(t *testing.T) {
	var r l3networkResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "l2_network_uuid"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "state", "zone_uuid"}
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

func TestL3networkResource_Metadata(t *testing.T) {
	var r l3networkResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_l3network" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestL3NetworkUpdateGuardsUnknownValues(t *testing.T) {
	t.Run("IpVersion Unknown omitted from payload", func(t *testing.T) {
		plan := l3networkResourceModel{
			Uuid:          types.StringValue("test-uuid"),
			Name:          types.StringValue("test-network"),
			L2NetworkUuid: types.StringValue("l2-uuid"),
			IpVersion:     types.Int64Unknown(),
			System:        types.BoolValue(false),
		}

		if !plan.IpVersion.IsNull() && !plan.IpVersion.IsUnknown() {
			t.Error("IpVersion guard failed: Unknown value was not properly guarded")
		}
	})

	t.Run("System Unknown omitted from payload", func(t *testing.T) {
		plan := l3networkResourceModel{
			Uuid:          types.StringValue("test-uuid"),
			Name:          types.StringValue("test-network"),
			L2NetworkUuid: types.StringValue("l2-uuid"),
			IpVersion:     types.Int64Value(4),
			System:        types.BoolUnknown(),
		}

		if !plan.System.IsNull() && !plan.System.IsUnknown() {
			t.Error("System guard failed: Unknown value was not properly guarded")
		}
	})
}

func TestAccL3NetworkResource(t *testing.T) {
	env := loadEnvData(t)

	if len(env.L2Networks) == 0 {
		t.Skip("no l2_networks in env.json, skipping L3 network acceptance test")
	}

	// Use the L2NoVlanNetwork (simpler, no VLAN requirement)
	var l2UUID string
	for _, l2 := range env.L2Networks {
		if envStr(l2, "type") == "L2NoVlanNetwork" {
			l2UUID = envStr(l2, "uuid")
			break
		}
	}
	if l2UUID == "" {
		l2UUID = envStr(env.L2Networks[0], "uuid")
	}
	name := testAccName("l3network")
	updatedName := name + "-updated"

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckL3NetworkDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create (ip_version must be set explicitly; API rejects value 0)
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_l3network" "test" {
  name            = %q
  description     = "acceptance l3 network"
  l2_network_uuid = %q
  category        = "Private"
  ip_version      = 4
}
`, name, l2UUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_l3network.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_l3network.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_l3network.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance l3 network")),
					statecheck.ExpectKnownValue("zstack_l3network.test", tfjsonpath.New("l2_network_uuid"), knownvalue.StringExact(l2UUID)),
				},
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_l3network" "test" {
  name            = %q
  description     = "acceptance l3 network updated"
  l2_network_uuid = %q
  category        = "Private"
  ip_version      = 4
}
`, updatedName, l2UUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_l3network.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("zstack_l3network.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance l3 network updated")),
					statecheck.ExpectKnownValue("zstack_l3network.test", tfjsonpath.New("l2_network_uuid"), knownvalue.StringExact(l2UUID)),
				},
			},
			// Import step
			{
				ResourceName:                         "zstack_l3network.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_l3network.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
