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

	config := func(poolResourceName, rangeResourceName string, includeAttachment bool) string {
		attachmentConfig := ""
		if includeAttachment {
			attachmentConfig = fmt.Sprintf(`
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
`, clusterUuid, clusterUuid, vtepCidr)
		}

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

`, poolResourceName, zoneUuid, physicalInterface, rangeResourceName, startVni, endVni) + attachmentConfig
	}

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVxlanPoolAndRangeDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: config(poolName, rangeName, true),
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
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
				Config: config(poolName+"-updated", rangeName+"-updated", true),
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("zstack_l2vxlan_network_pool.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("zstack_vni_range.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("zstack_l2_network_cluster_attachment.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
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
			{
				ResourceName:                         "zstack_l2_network_cluster_attachment.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdL2NetworkClusterAttachment("zstack_l2_network_cluster_attachment.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
			},
			{
				Config: config(poolName+"-updated", rangeName+"-updated", false),
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("zstack_l2vxlan_network_pool.test", plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction("zstack_vni_range.test", plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction("zstack_l2_network_cluster_attachment.test", plancheck.ResourceActionDestroy),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: testAccCheckVxlanAttachmentDeletedSafely(
					"zstack_l2vxlan_network_pool.test",
					"zstack_vni_range.test",
					clusterUuid,
				),
				ConfigStateChecks: []statecheck.StateCheck{
					poolUuidUnchanged.AddStateValue("zstack_l2vxlan_network_pool.test", tfjsonpath.New("uuid")),
					rangeUuidUnchanged.AddStateValue("zstack_vni_range.test", tfjsonpath.New("uuid")),
				},
			},
		},
	})
}

func TestAccL2VxlanNetworkPoolImportPlanNoop(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("acceptance test skipped unless TF_ACC is set")
	}

	zoneUuid := requireAcceptanceEnv(t, "ZSTACK_TEST_VXLAN_ZONE_UUID")
	clusterUuid := requireAcceptanceEnv(t, "ZSTACK_TEST_VXLAN_CLUSTER_UUID")
	physicalInterface := requireAcceptanceEnv(t, "ZSTACK_TEST_VXLAN_PHYSICAL_INTERFACE")
	vtepCidr := requireAcceptanceEnv(t, "ZSTACK_TEST_VXLAN_VTEP_CIDR")
	startVni, endVni := testAccFreeVniRange(t, 100)
	poolName := testAccName("vxlan-import")
	rangeName := testAccName("vni-import")
	description := "Terraform VXLAN import acceptance test"

	cli := testAccClientLoggedIn()
	pool, err := cli.CreateL2VxlanNetworkPool(param.CreateL2VxlanNetworkPoolParam{
		Params: param.CreateL2VxlanNetworkPoolParamDetail{
			Name:              poolName,
			Description:       stringPtr(description),
			ZoneUuid:          zoneUuid,
			PhysicalInterface: stringPtr(physicalInterface),
			VSwitchType:       stringPtr("LinuxBridge"),
		},
	})
	if err != nil {
		t.Fatalf("create VXLAN pool import fixture: %v", err)
	}
	t.Cleanup(func() {
		_ = cli.DeleteL2Network(pool.UUID, param.DeleteModePermissive)
	})

	vniRange, err := cli.CreateVniRange(pool.UUID, param.CreateVniRangeParam{
		Params: param.CreateVniRangeParamDetail{
			Name:        rangeName,
			Description: stringPtr(description),
			StartVni:    startVni,
			EndVni:      endVni,
		},
	})
	if err != nil {
		t.Fatalf("create VNI range import fixture: %v", err)
	}
	t.Cleanup(func() {
		_ = cli.DeleteVniRange(vniRange.UUID, param.DeleteModePermissive)
	})

	attachmentSystemTag := fmt.Sprintf(
		"l2NetworkUuid::%s::clusterUuid::%s::cidr::{%s}",
		pool.UUID,
		clusterUuid,
		vtepCidr,
	)
	if _, err := cli.AttachL2NetworkToCluster(
		pool.UUID,
		clusterUuid,
		attachL2NetworkToClusterParam(stringPtr("LinuxBridge"), []string{attachmentSystemTag}),
	); err != nil {
		t.Fatalf("create VXLAN attachment import fixture: %v", err)
	}
	t.Cleanup(func() {
		_ = cli.DetachL2NetworkFromCluster(pool.UUID, clusterUuid, param.DeleteModePermissive)
	})

	config := providerConfig() + fmt.Sprintf(`
import {
  to = zstack_l2vxlan_network_pool.test
  id = %q
}

import {
  to = zstack_vni_range.test
  id = %q
}

import {
  to = zstack_l2_network_cluster_attachment.test
  id = %q
}

resource "zstack_l2vxlan_network_pool" "test" {
  name               = %q
  description        = %q
  zone_uuid          = %q
  physical_interface = %q
  vswitch_type       = "LinuxBridge"
}

resource "zstack_vni_range" "test" {
  name        = %q
  description = %q
  start_vni   = %d
  end_vni     = %d
  pool_uuid   = zstack_l2vxlan_network_pool.test.uuid
}

resource "zstack_l2_network_cluster_attachment" "test" {
  l2_network_uuid  = zstack_l2vxlan_network_pool.test.uuid
  cluster_uuid     = %q
  l2_provider_type = "LinuxBridge"
  system_tags      = [%q]
}
`,
		pool.UUID,
		vniRange.UUID,
		l2NetworkClusterAttachmentID(pool.UUID, clusterUuid),
		poolName,
		description,
		zoneUuid,
		physicalInterface,
		rangeName,
		description,
		startVni,
		endVni,
		clusterUuid,
		attachmentSystemTag,
	)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVxlanPoolAndRangeDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: config,
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_l2vxlan_network_pool.test", tfjsonpath.New("uuid"), knownvalue.StringExact(pool.UUID)),
					statecheck.ExpectKnownValue("zstack_vni_range.test", tfjsonpath.New("uuid"), knownvalue.StringExact(vniRange.UUID)),
					statecheck.ExpectKnownValue(
						"zstack_l2_network_cluster_attachment.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact(l2NetworkClusterAttachmentID(pool.UUID, clusterUuid)),
					),
				},
			},
			{
				Config:   config,
				PlanOnly: true,
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
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

func testAccCheckVxlanAttachmentDeletedSafely(
	poolResourceName, rangeResourceName, clusterUuid string,
) tfresource.TestCheckFunc {
	return func(state *terraform.State) error {
		poolState, ok := state.RootModule().Resources[poolResourceName]
		if !ok {
			return fmt.Errorf("VXLAN pool %s is missing from Terraform state", poolResourceName)
		}
		rangeState, ok := state.RootModule().Resources[rangeResourceName]
		if !ok {
			return fmt.Errorf("VNI range %s is missing from Terraform state", rangeResourceName)
		}

		poolUuid := poolState.Primary.Attributes["uuid"]
		rangeUuid := rangeState.Primary.Attributes["uuid"]
		cli := testAccClientLoggedIn()

		pool, err := cli.GetL2VxlanNetworkPool(poolUuid)
		if err != nil {
			return fmt.Errorf("VXLAN pool %s was affected while deleting its attachment: %w", poolUuid, err)
		}
		for _, attachedClusterUuid := range pool.AttachedClusterUuids {
			if attachedClusterUuid == clusterUuid {
				return fmt.Errorf("VXLAN pool %s is still attached to cluster %s", poolUuid, clusterUuid)
			}
		}
		if _, err := cli.GetVniRange(rangeUuid); err != nil {
			return fmt.Errorf("VNI range %s was affected while deleting its attachment: %w", rangeUuid, err)
		}
		if _, err := cli.GetCluster(clusterUuid); err != nil {
			return fmt.Errorf("cluster %s was affected while deleting the attachment: %w", clusterUuid, err)
		}
		return nil
	}
}
