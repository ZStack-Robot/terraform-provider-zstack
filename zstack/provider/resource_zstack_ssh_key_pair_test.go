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

func TestSshKeyPairResource_Schema(t *testing.T) {
	var r sshKeyPairResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "public_key"}
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

	optional := []string{"description"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}
}

func TestSshKeyPairResource_Metadata(t *testing.T) {
	var r sshKeyPairResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_ssh_key_pair" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccSshKeyPairResource(t *testing.T) {
	_ = loadEnvData(t)
	name := testAccName("ssh-key-pair")
	updatedName := name + "-updated"

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSshKeyPairDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_ssh_key_pair" "test" {
  name        = %q
  description = "acceptance SSH key pair"
  public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC7 test@example.com"
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_ssh_key_pair.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_ssh_key_pair.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_ssh_key_pair.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance SSH key pair")),
				},
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_ssh_key_pair" "test" {
  name        = %q
  description = "acceptance SSH key pair updated"
  public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC7 test@example.com"
}
`, updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_ssh_key_pair.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("zstack_ssh_key_pair.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance SSH key pair updated")),
				},
			},
			{
				ResourceName:                         "zstack_ssh_key_pair.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_ssh_key_pair.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}

func TestAccSshKeyPairResource_disappears(t *testing.T) {
	_ = loadEnvData(t)
	name := testAccName("ssh-disappears")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSshKeyPairDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_ssh_key_pair" "test_disappears" {
  name = %q
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_ssh_key_pair.test_disappears", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					stateCheckSshKeyPairDisappears("zstack_ssh_key_pair.test_disappears"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
