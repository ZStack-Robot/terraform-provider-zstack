package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccReservedIpResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test creating a single image resource
			{
				Config: providerConfig + `
				resource "zstack_reserved_ip" "test" {
					l3_network_uuid         = "a5e77b2972e64316878993af7a695400"
					start_ip = "172.26.111.250"
					end_ip = "172.26.111.253"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the resource attributes
					resource.TestCheckResourceAttr("zstack_reserved_ip.test", "l3_network_uuid", "a5e77b2972e64316878993af7a695400"),
					resource.TestCheckResourceAttr("zstack_reserved_ip.test", "start_ip", "172.26.111.250"),
					resource.TestCheckResourceAttr("zstack_reserved_ip.test", "end_ip", "172.26.111.253"),
				),
			},
		},
	})
}
