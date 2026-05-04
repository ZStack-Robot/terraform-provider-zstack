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

func TestInstanceOfferingResource_Schema(t *testing.T) {
	var r instanceOfferingResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "cpu_num", "memory_size"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "description", "type"}
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

func TestInstanceOfferingResource_Metadata(t *testing.T) {
	var r instanceOfferingResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_instance_offering" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccInstanceOfferingResource_disappears(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckInstanceOfferingDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_instance_offering" "test" {
  name        = "acc-test-instance-offer"
  cpu_num     = 1
  memory_size = 1073741824
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					stateCheckInstanceOfferingDisappears("zstack_instance_offering.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Note: Update Step not applicable — all user-settable attributes have RequiresReplace.
func TestAccInstanceOfferingResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckInstanceOfferingDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_instance_offering" "test" {
  name        = "acc-test-instance-offer"
  cpu_num     = 1
  memory_size = 1073741824
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_instance_offering.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_instance_offering.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-instance-offer")),
					statecheck.ExpectKnownValue("zstack_instance_offering.test", tfjsonpath.New("cpu_num"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("zstack_instance_offering.test", tfjsonpath.New("memory_size"), knownvalue.Int64Exact(1073741824)),
				},
			},
			{
				ResourceName:                         "zstack_instance_offering.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_instance_offering.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
