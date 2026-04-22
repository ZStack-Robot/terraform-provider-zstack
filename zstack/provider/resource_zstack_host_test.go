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

func TestHostResource_Schema(t *testing.T) {
	var r hostResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "management_ip", "cluster_uuid", "username", "password"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "zone_uuid", "hypervisor_type", "status", "architecture"}
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

func TestHostResource_Metadata(t *testing.T) {
	var r hostResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_host" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

// TestAccHostResource tests the host resource acceptance.
// Creating a host requires a spare physical host IP + SSH credentials not in the test env.
// This test is skipped unless the env provides a "spare_host_ip" in hosts entries.
func TestAccHostResource(t *testing.T) {
	env := loadEnvData(t)

	if len(env.Clusters) == 0 {
		t.Skip("no clusters in env.json, skipping host acceptance test")
	}

	clusterUUID := envStr(env.Clusters[0], "uuid")

	// Look for a spare host entry (one with spare_ip field set)
	var spareHostIP, spareHostUser, spareHostPass string
	for _, h := range env.Hosts {
		if ip := envStr(h, "spare_ip"); ip != "" {
			spareHostIP = ip
			spareHostUser = envStr(h, "spare_username")
			spareHostPass = envStr(h, "spare_password")
			break
		}
	}

	if spareHostIP == "" {
		t.Skip("no spare host available for Create test (set spare_ip/spare_username/spare_password in env.json hosts entry)")
	}

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckHostDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_host" "test" {
  name          = "acc-test-host"
  management_ip = %q
  cluster_uuid  = %q
  username      = %q
  password      = %q
}
`, spareHostIP, clusterUUID, spareHostUser, spareHostPass),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_host.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_host.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-host")),
				},
			},
			// Step 2: Update name
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_host" "test" {
  name          = "acc-test-host-updated"
  management_ip = %q
  cluster_uuid  = %q
  username      = %q
  password      = %q
}
`, spareHostIP, clusterUUID, spareHostUser, spareHostPass),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_host.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-host-updated")),
				},
			},
			// Step 3: Import
			{
				ResourceName:                        "zstack_host.test",
				ImportState:                         true,
				ImportStateIdFunc:                   importStateIdFromUUID("zstack_host.test"),
				ImportStateVerify:                   true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				// username and password are sensitive write-only and not returned by the API
				ImportStateVerifyIgnore:             []string{"username", "password"},
			},
		},
	})
}
