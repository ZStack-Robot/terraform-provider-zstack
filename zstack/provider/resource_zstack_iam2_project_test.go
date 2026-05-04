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

func TestIAM2ProjectResource_Schema(t *testing.T) {
	var r iam2ProjectResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "state"}
	for _, attr := range computed {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing computed attribute %q", attr)
		}
		if !a.IsComputed() {
			t.Errorf("attribute %q should be computed", attr)
		}
	}

	optional := []string{"description"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}
}

func TestIAM2ProjectResource_Metadata(t *testing.T) {
	var r iam2ProjectResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_iam2_project" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccIAM2ProjectResource(t *testing.T) {
	_ = loadEnvData(t)
	name := testAccName("iam2-project")
	updatedName := name + "-updated"

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIAM2ProjectDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_iam2_project" "test" {
  name        = %q
  description = "acceptance test IAM2 project"
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_iam2_project.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_iam2_project.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_iam2_project.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance test IAM2 project")),
				},
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_iam2_project" "test" {
  name        = %q
  description = "acceptance test IAM2 project updated"
}
`, updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_iam2_project.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("zstack_iam2_project.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance test IAM2 project updated")),
				},
			},
			{
				ResourceName:                         "zstack_iam2_project.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_iam2_project.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}

func TestAccIAM2ProjectResource_disappears(t *testing.T) {
	_ = loadEnvData(t)
	name := testAccName("project-disappears")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckIAM2ProjectDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_iam2_project" "test_disappears" {
  name        = %q
  description = "Disappears test project"
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_iam2_project.test_disappears", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					stateCheckIAM2ProjectDisappears("zstack_iam2_project.test_disappears"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
