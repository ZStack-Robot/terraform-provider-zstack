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

func TestGuestToolsResource_Schema(t *testing.T) {
	var r guestToolsResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"instance_uuid"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"id", "guest_tools_version", "guest_tools_status"}
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

func TestGuestToolsResource_Metadata(t *testing.T) {
	var r guestToolsResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_guest_tools_attachment" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccGuestToolsAttachmentResource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Images) == 0 {
		t.Skip("no images in env data")
	}
	if len(env.InstanceOfferings) == 0 {
		t.Skip("no instance_offerings in env data")
	}
	if len(env.L3Networks) == 0 {
		t.Skip("no l3_networks in env data")
	}

	// Pick ttylinux image (smallest, fastest boot); fall back to any Ready image.
	const ttylinuxUUID = "dfc919110d734009bea2f04a5e8ac9ef"
	var imageUUID string
	for _, img := range env.Images {
		if envStr(img, "uuid") == ttylinuxUUID && envStr(img, "status") == "Ready" {
			imageUUID = ttylinuxUUID
			break
		}
	}
	if imageUUID == "" {
		for _, img := range env.Images {
			if envStr(img, "status") == "Ready" {
				imageUUID = envStr(img, "uuid")
				break
			}
		}
	}
	if imageUUID == "" {
		t.Skip("no Ready images in env data")
	}

	offeringUUID := envStr(env.InstanceOfferings[0], "uuid")
	l3UUID := envStr(env.L3Networks[0], "uuid")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGuestToolsAttachmentDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create VM then attach guest tools ISO.
			// The VM created by Terraform starts in Running state.
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_instance" "gt_test_vm" {
  name                   = "acc-test-gt-vm"
  image_uuid             = %q
  instance_offering_uuid = %q
  expunge                = true
  network_interfaces = [
    {
      l3_network_uuid = %q
      default_l3      = true
    }
  ]
}

resource "zstack_guest_tools_attachment" "test" {
  instance_uuid = zstack_instance.gt_test_vm.uuid
}
`, imageUUID, offeringUUID, l3UUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_guest_tools_attachment.test", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
		},
	})
}
