// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

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

func TestVolumeResource_Schema(t *testing.T) {
	var r volumeResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Error("schema should not be empty")
	}

	if _, ok := resp.Schema.Attributes["uuid"]; !ok {
		t.Fatal("schema should include uuid")
	}

	if _, ok := resp.Schema.Attributes["name"]; !ok {
		t.Fatal("schema should include name")
	}
}

func TestAccVolumeResource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.DiskOfferings) == 0 {
		t.Skip("no disk offerings in env data")
	}
	doUUID := envStr(env.DiskOfferings[0], "uuid")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVolumeDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_volume" "test" {
  name               = "acc-test-volume"
  disk_offering_uuid = %q
}
`, doUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_volume.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_volume.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-volume")),
				},
			},
			{
				ResourceName:      "zstack_volume.test",
				ImportState:       true,
				ImportStateIdFunc:       importStateUUID("zstack_volume.test"),
				ImportStateVerify: true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
