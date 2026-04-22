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

func TestDiskOfferingResource_Schema(t *testing.T) {
	var r diskOfferingResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "disk_size"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "description"}
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

func TestDiskOfferingResource_Metadata(t *testing.T) {
	var r diskOfferingResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_disk_offer" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccDiskOfferingResource_disappears(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDiskOfferingDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_disk_offer" "test" {
  name      = "acc-test-disk-offer"
  disk_size = 10
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					stateCheckDiskOfferingDisappears("zstack_disk_offer.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Note: Update Step not applicable — all user-settable attributes have RequiresReplace.
func TestAccDiskOfferingResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDiskOfferingDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_disk_offer" "test" {
  name      = "acc-test-disk-offer"
  disk_size = 10
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_disk_offer.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_disk_offer.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-disk-offer")),
					statecheck.ExpectKnownValue("zstack_disk_offer.test", tfjsonpath.New("disk_size"), knownvalue.Int64Exact(10)),
				},
			},
			{
				ResourceName:                         "zstack_disk_offer.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_disk_offer.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
