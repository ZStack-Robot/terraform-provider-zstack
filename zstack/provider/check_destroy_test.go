// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

var testAccCheckZoneDestroy = testAccCheckResourceDestroyByGet("zstack_zone", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetZone(id)
	return err
})

var testAccCheckAccountDestroy = testAccCheckResourceDestroyByGet("zstack_account", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetAccount(id)
	return err
})

var testAccCheckAffinityGroupDestroy = testAccCheckResourceDestroyByGet("zstack_affinity_group", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetAffinityGroup(id)
	return err
})

var testAccCheckClusterDestroy = testAccCheckResourceDestroyByGet("zstack_cluster", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetCluster(id)
	return err
})

var testAccCheckSshKeyPairDestroy = testAccCheckResourceDestroyByGet("zstack_ssh_key_pair", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetSshKeyPair(id)
	return err
})

var testAccCheckImageDestroy = testAccCheckResourceDestroyByGet("zstack_image", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetImage(id)
	return err
})

var testAccCheckVolumeDestroy = testAccCheckResourceDestroyByGet("zstack_volume", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetVolume(id)
	return err
})

var testAccCheckIAM2ProjectDestroy = testAccCheckResourceDestroyByGet("zstack_iam2_project", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetIAM2Project(id)
	return err
})

var testAccCheckL2VlanNetworkDestroy = testAccCheckResourceDestroyByGet("zstack_l2vlan_network", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetL2VlanNetwork(id)
	return err
})

var testAccCheckLoadBalancerDestroy = testAccCheckResourceDestroyByGet("zstack_load_balancer", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetLoadBalancer(id)
	return err
})

var testAccCheckLoadBalancerListenerDestroy = testAccCheckResourceDestroyByGet("zstack_load_balancer_listener", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetLoadBalancerListener(id)
	return err
})

var testAccCheckVirtualRouterImageDestroy = testAccCheckResourceDestroyByGet("zstack_virtual_router_image", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetImage(id)
	return err
})

var testAccCheckVirtualRouterOfferingDestroy = testAccCheckResourceDestroyByGet("zstack_virtual_router_offering", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetVirtualRouterOffering(id)
	return err
})

var testAccCheckAutoScalingGroupDestroy = testAccCheckResourceDestroyByGet("zstack_auto_scaling_group", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetAutoScalingGroup(id)
	return err
})

var testAccCheckPortForwardingRuleDestroy = testAccCheckResourceDestroyByGet("zstack_port_forwarding_rule", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetPortForwardingRule(id)
	return err
})

func testAccCheckVipDestroy(s *terraform.State) error {
	cli := testAccClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zstack_vip" {
			continue
		}
		id := rs.Primary.Attributes["uuid"]
		if id == "" {
			id = rs.Primary.ID
		}
		if id == "" {
			continue
		}
		q := param.NewQueryParam()
		q.AddQ(fmt.Sprintf("uuid=%s", id))
		vips, err := cli.QueryVip(&q)
		if err != nil {
			if isZStackNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking zstack_vip %s destroyed: %w", id, err)
		}
		if len(vips) > 0 {
			return fmt.Errorf("zstack_vip %s still exists", id)
		}
	}
	return nil
}

func testAccCheckReservedIpDestroy(s *terraform.State) error {
	cli := testAccClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zstack_reserved_ip" {
			continue
		}
		id := rs.Primary.Attributes["uuid"]
		if id == "" {
			id = rs.Primary.ID
		}
		if id == "" {
			continue
		}
		var reservedIpRanges []view.ReservedIpRangeInventoryView
		_, err := cli.Zql(context.Background(), fmt.Sprintf("query reservedIpRange where uuid='%s'", id), &reservedIpRanges, "inventories")
		if err != nil {
			if isZStackNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking zstack_reserved_ip %s destroyed: %w", id, err)
		}
		if len(reservedIpRanges) > 0 {
			return fmt.Errorf("zstack_reserved_ip %s still exists", id)
		}
	}
	return nil
}

func testAccCheckSecGroupAttachmentDestroy(s *terraform.State) error {
	cli := testAccClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zstack_networking_secgroup_attachment" {
			continue
		}
		sgUUID := rs.Primary.Attributes["secgroup_uuid"]
		nicUUID := rs.Primary.Attributes["nic_uuid"]
		if sgUUID == "" || nicUUID == "" {
			continue
		}
		candidate, err := cli.GetCandidateVmNicForSecurityGroup(sgUUID)
		if err != nil {
			if isZStackNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking secgroup attachment destroyed: %w", err)
		}
		if candidate != nil && candidate.UUID == nicUUID {
			continue
		}
		return fmt.Errorf("zstack_networking_secgroup_attachment %s/%s still attached", sgUUID, nicUUID)
	}
	return nil
}
