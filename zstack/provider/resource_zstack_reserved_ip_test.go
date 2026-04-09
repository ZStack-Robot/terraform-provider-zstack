// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"net"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// decrementIP returns the IP address one less than the given IP.
func decrementIP(ip net.IP) net.IP {
	result := make(net.IP, len(ip))
	copy(result, ip)
	for i := len(result) - 1; i >= 0; i-- {
		if result[i] > 0 {
			result[i]--
			break
		}
		result[i] = 255
	}
	return result
}

func TestAccReservedIpResource(t *testing.T) {
	env := loadEnvData(t)

	// Find a Public L3 network and its IP range
	var l3UUID, startIP, endIP string
	for _, l3 := range env.L3Networks {
		if envStr(l3, "category") != "Public" {
			continue
		}
		l3uuid := envStr(l3, "uuid")
		for _, ipr := range env.IpRanges {
			if envStr(ipr, "l3_network_uuid") == l3uuid {
				l3UUID = l3uuid
				startIP = envStr(ipr, "start_ip")
				endIP = envStr(ipr, "end_ip")
				break
			}
		}
		if l3UUID != "" {
			break
		}
	}

	if l3UUID == "" || startIP == "" || endIP == "" {
		t.Skip("no Public L3 network with IP range found in env data")
	}

	// Use the last 2 IPs from the range: endIP-1 and endIP
	end := net.ParseIP(endIP)
	if end == nil {
		t.Fatalf("cannot parse end_ip %q", endIP)
	}
	reserveEnd := end.String()
	reserveStart := decrementIP(end.To4()).String()

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckReservedIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
				resource "zstack_reserved_ip" "test" {
					l3_network_uuid = %q
					start_ip = %q
					end_ip   = %q
				}`, l3UUID, reserveStart, reserveEnd),

				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_reserved_ip.test", tfjsonpath.New("l3_network_uuid"), knownvalue.StringExact(l3UUID)),
					statecheck.ExpectKnownValue("zstack_reserved_ip.test", tfjsonpath.New("start_ip"), knownvalue.StringExact(reserveStart)),
					statecheck.ExpectKnownValue("zstack_reserved_ip.test", tfjsonpath.New("end_ip"), knownvalue.StringExact(reserveEnd)),
				},
			},
			{
				ResourceName:      "zstack_reserved_ip.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
