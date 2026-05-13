// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

func TestL2NetworkClusterAttachmentResource_Schema(t *testing.T) {
	var r l2NetworkClusterAttachmentResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	required := []string{"l2_network_uuid", "cluster_uuid"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	id, ok := resp.Schema.Attributes["id"]
	if !ok {
		t.Fatal("schema missing computed attribute \"id\"")
	}
	if !id.IsComputed() {
		t.Error("attribute \"id\" should be computed")
	}
}

func TestL2NetworkClusterAttachmentResource_Metadata(t *testing.T) {
	var r l2NetworkClusterAttachmentResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_l2_network_cluster_attachment" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestParseL2NetworkClusterAttachmentID(t *testing.T) {
	l2NetworkUuid, clusterUuid, err := parseL2NetworkClusterAttachmentID("l2-uuid:cluster-uuid")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l2NetworkUuid != "l2-uuid" || clusterUuid != "cluster-uuid" {
		t.Fatalf("unexpected parsed id: %q %q", l2NetworkUuid, clusterUuid)
	}

	invalidIDs := []string{"", "l2-uuid", ":cluster-uuid", "l2-uuid:"}
	for _, id := range invalidIDs {
		if _, _, err := parseL2NetworkClusterAttachmentID(id); err == nil {
			t.Fatalf("expected error for invalid id %q", id)
		}
	}
}

func TestAccL2NetworkClusterAttachmentResource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("acceptance test skipped unless TF_ACC is set")
	}

	clusterUuid, zoneUuid := testAccLiveCluster(t)

	name := testAccName("l2net-attach")
	vlan := testAccFreeL2Vlan(t, zoneUuid, "eth0")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckL2VlanNetworkDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_l2vlan_network" "test" {
  name               = %q
  vlan               = %d
  zone_uuid          = %q
  physical_interface = "eth0"
}

resource "zstack_l2_network_cluster_attachment" "test" {
  l2_network_uuid = zstack_l2vlan_network.test.uuid
  cluster_uuid    = %q
}
`, name, vlan, zoneUuid, clusterUuid),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_l2_network_cluster_attachment.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_l2_network_cluster_attachment.test", tfjsonpath.New("cluster_uuid"), knownvalue.StringExact(clusterUuid)),
				},
			},
			{
				ResourceName:                         "zstack_l2_network_cluster_attachment.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdL2NetworkClusterAttachment("zstack_l2_network_cluster_attachment.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_l2vlan_network" "test" {
  name               = %q
  vlan               = %d
  zone_uuid          = %q
  physical_interface = "eth0"
}
`, name, vlan, zoneUuid),
				Check: testAccCheckL2NetworkDetached("zstack_l2vlan_network.test", clusterUuid),
			},
		},
	})
}

func importStateIdL2NetworkClusterAttachment(resourceName string) tfresource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", resourceName)
		}
		l2NetworkUuid := rs.Primary.Attributes["l2_network_uuid"]
		clusterUuid := rs.Primary.Attributes["cluster_uuid"]
		if l2NetworkUuid == "" || clusterUuid == "" {
			return "", fmt.Errorf("l2_network_uuid or cluster_uuid attribute is empty for %s", resourceName)
		}
		return l2NetworkClusterAttachmentID(l2NetworkUuid, clusterUuid), nil
	}
}

func testAccCheckL2NetworkDetached(resourceName, clusterUuid string) tfresource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}
		l2NetworkUuid := rs.Primary.Attributes["uuid"]
		if l2NetworkUuid == "" {
			return fmt.Errorf("uuid attribute is empty for %s", resourceName)
		}

		l2Network, err := testAccClientLoggedIn().GetL2Network(l2NetworkUuid)
		if err != nil {
			return fmt.Errorf("get L2 network %s: %w", l2NetworkUuid, err)
		}
		for _, attachedClusterUuid := range l2Network.AttachedClusterUuids {
			if attachedClusterUuid == clusterUuid {
				return fmt.Errorf("L2 network %s is still attached to cluster %s", l2NetworkUuid, clusterUuid)
			}
		}
		return nil
	}
}

func testAccFreeL2Vlan(t *testing.T, zoneUuid, physicalInterface string) int {
	t.Helper()

	used := map[int]struct{}{}
	queryParam := param.NewQueryParam()
	l2VlanNetworks, err := testAccClientLoggedIn().QueryL2VlanNetwork(&queryParam)
	if err != nil {
		t.Logf("query existing L2 VLAN networks failed, falling back to time-based VLAN: %v", err)
		return 3900 + int(time.Now().UnixNano()%100)
	}
	for _, l2 := range l2VlanNetworks {
		if l2.ZoneUuid == zoneUuid && l2.PhysicalInterface == physicalInterface {
			used[l2.Vlan] = struct{}{}
		}
	}

	for vlan := 3900; vlan <= 4094; vlan++ {
		if _, ok := used[vlan]; !ok {
			return vlan
		}
	}
	t.Skipf("no free VLAN ID in range 3900-4094 for zone %s interface %s", zoneUuid, physicalInterface)
	return 0
}

func testAccLiveCluster(t *testing.T) (string, string) {
	t.Helper()

	queryParam := param.NewQueryParam()
	queryParam.AddQ("state=Enabled")
	clusters, err := testAccClientLoggedIn().QueryCluster(&queryParam)
	if err != nil {
		t.Fatalf("query enabled clusters: %v", err)
	}

	for _, cluster := range clusters {
		if cluster.UUID != "" && cluster.ZoneUuid != "" {
			return cluster.UUID, cluster.ZoneUuid
		}
	}
	t.Skip("no enabled cluster with zone_uuid found in live ZStack environment")
	return "", ""
}
