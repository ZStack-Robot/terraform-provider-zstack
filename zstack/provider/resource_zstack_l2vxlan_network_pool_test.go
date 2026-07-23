// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

func TestL2VxlanNetworkPoolResourceSchema(t *testing.T) {
	var r l2VxlanNetworkPoolResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	for _, name := range []string{"name", "zone_uuid", "physical_interface"} {
		attribute, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Fatalf("schema missing required attribute %q", name)
		}
		if !attribute.IsRequired() {
			t.Errorf("attribute %q should be required", name)
		}
	}

	for _, name := range []string{"uuid", "virtual_network_id", "attached_cluster_uuids"} {
		attribute, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Fatalf("schema missing computed attribute %q", name)
		}
		if !attribute.IsComputed() {
			t.Errorf("attribute %q should be computed", name)
		}
	}
}

func TestL2VxlanNetworkPoolResourceMetadata(t *testing.T) {
	var r l2VxlanNetworkPoolResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_l2vxlan_network_pool" {
		t.Fatalf("unexpected type name: %s", resp.TypeName)
	}
}

func TestL2VxlanNetworkPoolModelFromViewPreservesCreateOnlyFields(t *testing.T) {
	prior := l2VxlanNetworkPoolResourceModel{
		ResourceUuid: types.StringValue("custom-uuid"),
		TagUuids:     stringSliceToList([]string{"tag-1"}),
		SystemTags:   stringSliceToList([]string{"system-tag"}),
	}

	state := l2VxlanNetworkPoolModelFromView(&view.L2VxlanNetworkPoolInventoryView{
		BaseInfoView:         view.BaseInfoView{UUID: "pool-uuid", Name: "pool"},
		ZoneUuid:             "zone-uuid",
		PhysicalInterface:    "bond0",
		Type:                 "L2VxlanNetworkPool",
		VSwitchType:          "LinuxBridge",
		VirtualNetworkId:     100,
		AttachedClusterUuids: []string{"cluster-uuid"},
	}, prior)

	if state.ResourceUuid.ValueString() != "custom-uuid" {
		t.Fatalf("resource_uuid was not preserved: %s", state.ResourceUuid.ValueString())
	}
	if got := listToStringSlice(state.TagUuids); len(got) != 1 || got[0] != "tag-1" {
		t.Fatalf("tag_uuids were not preserved: %#v", got)
	}
	if got := listToStringSlice(state.SystemTags); len(got) != 1 || got[0] != "system-tag" {
		t.Fatalf("system_tags were not preserved: %#v", got)
	}
	if got := listToStringSlice(state.AttachedClusterUuids); len(got) != 1 || got[0] != "cluster-uuid" {
		t.Fatalf("attached_cluster_uuids were not mapped: %#v", got)
	}
}

func TestL2VxlanNetworkPoolCreateParam(t *testing.T) {
	createParam := l2VxlanNetworkPoolCreateParam(l2VxlanNetworkPoolResourceModel{
		Name:              types.StringValue("pool"),
		ZoneUuid:          types.StringValue("zone-uuid"),
		PhysicalInterface: types.StringValue("bond0"),
		TagUuids:          stringSliceToList([]string{"tag-1"}),
		SystemTags:        stringSliceToList([]string{"system-tag"}),
	})

	if createParam.Params.PhysicalInterface == nil || *createParam.Params.PhysicalInterface != "bond0" {
		t.Fatalf("physicalInterface was not passed: %#v", createParam.Params.PhysicalInterface)
	}
	if len(createParam.Params.TagUuids) != 1 || createParam.Params.TagUuids[0] != "tag-1" {
		t.Fatalf("tag UUIDs were not passed: %#v", createParam.Params.TagUuids)
	}
	if len(createParam.SystemTags) != 1 || createParam.SystemTags[0] != "system-tag" {
		t.Fatalf("system tags were not passed: %#v", createParam.SystemTags)
	}
}

func TestAccL2VxlanNetworkPoolAndVniRange(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("acceptance test skipped unless TF_ACC is set")
	}

	zoneUuid := requireAcceptanceEnv(t, "ZSTACK_TEST_VXLAN_ZONE_UUID")
	clusterUuid := requireAcceptanceEnv(t, "ZSTACK_TEST_VXLAN_CLUSTER_UUID")
	physicalInterface := requireAcceptanceEnv(t, "ZSTACK_TEST_VXLAN_PHYSICAL_INTERFACE")
	vtepCidr := requireAcceptanceEnv(t, "ZSTACK_TEST_VXLAN_VTEP_CIDR")
	startVni, endVni := testAccFreeVniRange(t, 100)
	poolName := testAccName("vxlan-pool")
	rangeName := testAccName("vni-range")
	poolUuidUnchanged := statecheck.CompareValue(compare.ValuesSame())
	rangeUuidUnchanged := statecheck.CompareValue(compare.ValuesSame())

	config := func(poolResourceName, rangeResourceName string) string {
		return providerConfig() + fmt.Sprintf(`
resource "zstack_l2vxlan_network_pool" "test" {
  name               = %q
  description        = "Terraform VXLAN acceptance test"
  zone_uuid          = %q
  physical_interface = %q
  vswitch_type       = "LinuxBridge"
}

resource "zstack_vni_range" "test" {
  name        = %q
  description = "Terraform VNI range acceptance test"
  start_vni   = %d
  end_vni     = %d
  pool_uuid   = zstack_l2vxlan_network_pool.test.uuid
}

resource "zstack_l2_network_cluster_attachment" "test" {
  l2_network_uuid  = zstack_l2vxlan_network_pool.test.uuid
  cluster_uuid     = %q
  l2_provider_type = "LinuxBridge"
  system_tags = [format(
    "l2NetworkUuid::%%s::clusterUuid::%%s::cidr::{%%s}",
    zstack_l2vxlan_network_pool.test.uuid,
    %q,
    %q,
  )]
}
`, poolResourceName, zoneUuid, physicalInterface, rangeResourceName, startVni, endVni, clusterUuid, clusterUuid, vtepCidr)
	}

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVxlanPoolAndRangeDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: config(poolName, rangeName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_l2vxlan_network_pool.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					poolUuidUnchanged.AddStateValue("zstack_l2vxlan_network_pool.test", tfjsonpath.New("uuid")),
					statecheck.ExpectKnownValue("zstack_l2vxlan_network_pool.test", tfjsonpath.New("physical_interface"), knownvalue.StringExact(physicalInterface)),
					statecheck.ExpectKnownValue("zstack_vni_range.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					rangeUuidUnchanged.AddStateValue("zstack_vni_range.test", tfjsonpath.New("uuid")),
					statecheck.ExpectKnownValue("zstack_vni_range.test", tfjsonpath.New("start_vni"), knownvalue.Int64Exact(int64(startVni))),
					statecheck.ExpectKnownValue("zstack_vni_range.test", tfjsonpath.New("end_vni"), knownvalue.Int64Exact(int64(endVni))),
					statecheck.ExpectKnownValue("zstack_l2_network_cluster_attachment.test", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				Config: config(poolName+"-updated", rangeName+"-updated"),
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("zstack_l2vxlan_network_pool.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("zstack_vni_range.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("zstack_l2_network_cluster_attachment.test", plancheck.ResourceActionNoop),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					poolUuidUnchanged.AddStateValue("zstack_l2vxlan_network_pool.test", tfjsonpath.New("uuid")),
					statecheck.ExpectKnownValue("zstack_l2vxlan_network_pool.test", tfjsonpath.New("name"), knownvalue.StringExact(poolName+"-updated")),
					rangeUuidUnchanged.AddStateValue("zstack_vni_range.test", tfjsonpath.New("uuid")),
					statecheck.ExpectKnownValue("zstack_vni_range.test", tfjsonpath.New("name"), knownvalue.StringExact(rangeName+"-updated")),
				},
			},
			{
				ResourceName:                         "zstack_l2vxlan_network_pool.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_l2vxlan_network_pool.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"resource_uuid", "tag_uuids", "system_tags", "attached_cluster_uuids"},
			},
			{
				ResourceName:                         "zstack_vni_range.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_vni_range.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"resource_uuid", "tag_uuids", "system_tags"},
			},
		},
	})
}

func requireAcceptanceEnv(t *testing.T, name string) string {
	t.Helper()
	value := os.Getenv(name)
	if value == "" {
		t.Skipf("acceptance test requires %s", name)
	}
	return value
}

func testAccFreeVniRange(t *testing.T, size int) (int, int) {
	t.Helper()
	query := param.NewQueryParam()
	ranges, err := testAccClientLoggedIn().QueryVniRange(&query)
	if err != nil {
		t.Fatalf("query VNI ranges: %v", err)
	}

	for start := 10000; start+size-1 <= maxVni; start += size {
		end := start + size - 1
		available := true
		for _, existing := range ranges {
			if start <= existing.EndVni && end >= existing.StartVni {
				available = false
				break
			}
		}
		if available {
			return start, end
		}
	}

	t.Fatal("no free VNI range available for acceptance test")
	return 0, 0
}

func testAccCheckVxlanPoolAndRangeDestroy(state *terraform.State) error {
	cli := testAccClientLoggedIn()
	for _, resourceState := range state.RootModule().Resources {
		uuid := resourceState.Primary.Attributes["uuid"]
		if uuid == "" {
			uuid = resourceState.Primary.ID
		}
		if uuid == "" {
			continue
		}

		var err error
		switch resourceState.Type {
		case "zstack_l2vxlan_network_pool":
			_, err = cli.GetL2VxlanNetworkPool(uuid)
		case "zstack_vni_range":
			_, err = cli.GetVniRange(uuid)
		default:
			continue
		}
		if err == nil {
			return fmt.Errorf("%s %s still exists", resourceState.Type, uuid)
		}
		if !isZStackNotFoundError(err) {
			return fmt.Errorf("check %s %s destroyed: %w", resourceState.Type, uuid, err)
		}
	}
	return nil
}
