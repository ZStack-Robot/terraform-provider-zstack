// Copyright (c) ZStack.io, Inc.

package provider

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

func TestAccSlbOfferingLifecycle(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("set TF_ACC to run against ZStack")
	}
	cli := testAccClientLoggedIn()
	q := param.NewQueryParam()
	offers, err := cli.QuerySlbOffering(&q)
	if err != nil {
		t.Fatal(err)
	}
	if len(offers) == 0 {
		t.Skip("an existing SLB offering is needed as a reference for image and network")
	}
	ref := offers[0]
	if ref.ZoneUuid == "" || ref.ManagementNetworkUuid == "" || ref.ImageUuid == "" {
		t.Fatal("reference SLB offering is missing zone, network, or image")
	}
	name := testAccName("slb-offering")
	names := []string{name, name + "-replaced"}
	// Also clean up a successful API create if Terraform fails before recording state.
	t.Cleanup(func() {
		for _, n := range names {
			query := param.NewQueryParam()
			query.AddQ("name=" + n)
			remaining, err := cli.QuerySlbOffering(&query)
			if err != nil {
				t.Errorf("cleanup query: %v", err)
				continue
			}
			for _, offer := range remaining {
				if offer.Name != n {
					continue
				}
				if err := cli.DeleteInstanceOffering(offer.UUID, param.DeleteModePermissive); err != nil {
					t.Errorf("cleanup %s: %v", offer.UUID, err)
				}
			}
			remaining, err = cli.QuerySlbOffering(&query)
			if err != nil || len(remaining) != 0 {
				t.Errorf("cleanup verification for %s: %d remaining, error %v", n, len(remaining), err)
			}
		}
	})
	config := func(n string) string {
		// Credentials are taken from the environment, never embedded in HCL.
		return fmt.Sprintf(`provider "zstack" {}
resource "zstack_slb_offering" "test" {
 name = %q
 description = "Temporary Terraform SLB acceptance test"
 cpu_num = %d
 memory_size = %d
 zone_uuid = %q
 management_network_uuid = %q
 image_uuid = %q
}
data "zstack_slb_offerings" "by_uuid" { uuid = zstack_slb_offering.test.uuid }
data "zstack_slb_offerings" "by_name" { name = zstack_slb_offering.test.name }
data "zstack_slb_offerings" "by_pattern" {
 name_pattern = "${zstack_slb_offering.test.name}%%"
 filter {
  name = "memory_size"
  values = ["%d"]
 }
}
data "zstack_slb_offerings" "empty" { uuid = "00000000000000000000000000000000" }
data "zstack_slb_offerings" "filtered_empty" {
 uuid = zstack_slb_offering.test.uuid
 filter {
  name = "state"
  values = ["Disabled"]
 }
}
output "empty_count" { value = length(data.zstack_slb_offerings.empty.slb_offers) }
output "filtered_empty_count" { value = length(data.zstack_slb_offerings.filtered_empty.slb_offers) }
`, n, ref.CpuNum, ref.MemorySize/(1024*1024), ref.ZoneUuid, ref.ManagementNetworkUuid, ref.ImageUuid, ref.MemorySize/(1024*1024))
	}
	var currentID string
	var createdIDs []string
	check := func(n string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("zstack_slb_offering.test", "name", n),
			resource.TestCheckResourceAttr("zstack_slb_offering.test", "type", "SLB"),
			resource.TestCheckResourceAttr("zstack_slb_offering.test", "memory_size", fmt.Sprint(ref.MemorySize/(1024*1024))),
			resource.TestCheckResourceAttr("data.zstack_slb_offerings.by_uuid", "slb_offers.#", "1"),
			resource.TestCheckResourceAttr("data.zstack_slb_offerings.by_name", "slb_offers.#", "1"),
			resource.TestCheckResourceAttr("data.zstack_slb_offerings.by_pattern", "slb_offers.#", "1"),
			resource.TestCheckResourceAttrPair("data.zstack_slb_offerings.by_uuid", "slb_offers.0.uuid", "zstack_slb_offering.test", "uuid"),
			resource.TestCheckResourceAttrPair("data.zstack_slb_offerings.by_name", "slb_offers.0.uuid", "zstack_slb_offering.test", "uuid"),
			resource.TestCheckOutput("empty_count", "0"),
			resource.TestCheckOutput("filtered_empty_count", "0"),
			func(s *terraform.State) error {
				id := s.RootModule().Resources["zstack_slb_offering.test"].Primary.Attributes["uuid"]
				if id == "" || id == currentID {
					return fmt.Errorf("expected newly created UUID, got %q", id)
				}
				currentID = id
				createdIDs = append(createdIDs, id)
				t.Logf("Created temporary SLB offering %s (%s)", n, id)
				return nil
			},
		)
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: func(_ *terraform.State) error {
			for _, id := range createdIDs {
				if _, err := findResourceByGet(cli.GetSlbOffering, id); !errors.Is(err, ErrResourceNotFound) {
					return fmt.Errorf("offering %s should be absent after destroy: %v", id, err)
				}
			}
			return nil
		},
		Steps: []resource.TestStep{
			{Config: config(names[0]), Check: check(names[0])},
			{ResourceName: "zstack_slb_offering.test", ImportState: true, ImportStateIdFunc: importStateIdFromUUID("zstack_slb_offering.test"), ImportStateVerify: true, ImportStateVerifyIdentifierAttribute: "uuid"},
			{Config: config(names[0]), PlanOnly: true},
			{Config: config(names[1]), Check: check(names[1])},
			{Config: config(names[1]), PreConfig: func() {
				if currentID == "" {
					t.Fatal("no test offering UUID to delete")
				}
				if err := cli.DeleteInstanceOffering(currentID, param.DeleteModePermissive); err != nil {
					t.Fatal(err)
				}
			}, Check: check(names[1])},
		},
	})
}
