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

func TestTagAttachmentResource_Schema(t *testing.T) {
	var r tagAttachmentResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"tag_uuid", "resource_uuids"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	if _, ok := resp.Schema.Attributes["tokens"]; !ok {
		t.Fatal("schema missing optional attribute tokens")
	}
}

func TestTagAttachmentResource_Metadata(t *testing.T) {
	var r tagAttachmentResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_tag_attachment" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccTagAttachmentResource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Zones) == 0 {
		t.Skip("no zones in env data, required for tag attachment target")
	}
	zoneUUID := envStr(env.Zones[0], "uuid")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTagAttachmentDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_tag" "test_for_attach" {
  name  = "acc-test-tag-for-attach"
  value = "attach-value"
  type  = "simple"
  color = "#FF0000"
}

resource "zstack_tag_attachment" "test" {
  tag_uuid       = zstack_tag.test_for_attach.uuid
  resource_uuids = [%q]
}
`, zoneUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_tag_attachment.test", tfjsonpath.New("tag_uuid"), knownvalue.NotNull()),
				},
			},
		},
	})
}
