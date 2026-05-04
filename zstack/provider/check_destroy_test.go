// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"strings"

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

var testAccCheckL3NetworkDestroy = testAccCheckResourceDestroyByQuery("zstack_l3network", func(cli *client.ZSClient, q *param.QueryParam) ([]view.L3NetworkInventoryView, error) {
	return cli.QueryL3Network(q)
})

var testAccCheckSubnetIpRangeDestroy = testAccCheckResourceDestroyByQuery("zstack_subnet_ip_range", func(cli *client.ZSClient, q *param.QueryParam) ([]view.IpRangeInventoryView, error) {
	return cli.QueryIpRange(q)
})

var testAccCheckL2VxlanNetworkDestroy = testAccCheckResourceDestroyByQuery("zstack_l2vxlan_network", func(cli *client.ZSClient, q *param.QueryParam) ([]view.L2VxlanNetworkInventoryView, error) {
	return cli.QueryL2VxlanNetwork(q)
})

var testAccCheckEipDestroy = testAccCheckResourceDestroyByQuery("zstack_eip", func(cli *client.ZSClient, q *param.QueryParam) ([]view.EipInventoryView, error) {
	return cli.QueryEip(q)
})

var testAccCheckHostDestroy = testAccCheckResourceDestroyByGet("zstack_host", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetHost(id)
	return err
})

var testAccCheckPrimaryStorageDestroy = testAccCheckResourceDestroyByGet("zstack_primary_storage", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetPrimaryStorage(id)
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

var testAccCheckInstanceOfferingDestroy = testAccCheckResourceDestroyByGet("zstack_instance_offering", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetInstanceOffering(id)
	return err
})

var testAccCheckDiskOfferingDestroy = testAccCheckResourceDestroyByGet("zstack_disk_offering", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetDiskOffering(id)
	return err
})

var testAccCheckTagDestroy = testAccCheckResourceDestroyByGet("zstack_tag", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetTag(id)
	return err
})

func testAccCheckCertificateDestroy(s *terraform.State) error {
	cli := testAccClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zstack_certificate" {
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
		items, err := cli.QueryCertificate(&q)
		if err != nil {
			if isZStackNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking zstack_certificate %s destroyed: %w", id, err)
		}
		if len(items) > 0 {
			return fmt.Errorf("zstack_certificate %s still exists", id)
		}
	}
	return nil
}

func testAccCheckWebhookDestroy(s *terraform.State) error {
	cli := testAccClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zstack_webhook" {
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
		items, err := cli.QueryWebhook(&q)
		if err != nil {
			if isZStackNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking zstack_webhook %s destroyed: %w", id, err)
		}
		if len(items) > 0 {
			return fmt.Errorf("zstack_webhook %s still exists", id)
		}
	}
	return nil
}

func testAccCheckVipDestroy(s *terraform.State) error {
	cli := testAccClientLoggedIn()
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
	cli := testAccClientLoggedIn()
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
			// Zql returns "key not found" when query returns empty inventories (SDK bug).
			// Treat this as not-found since the resource was successfully deleted.
			if isZStackNotFoundError(err) || strings.Contains(err.Error(), "key not found") {
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

func testAccCheckUserDestroy(s *terraform.State) error {
	cli := testAccClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zstack_user" {
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
		items, err := cli.QueryUser(&q)
		if err != nil {
			if isZStackNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking zstack_user %s destroyed: %w", id, err)
		}
		if len(items) > 0 {
			return fmt.Errorf("zstack_user %s still exists", id)
		}
	}
	return nil
}

func testAccCheckRoleDestroy(s *terraform.State) error {
	cli := testAccClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zstack_role" {
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
		items, err := cli.QueryRole(&q)
		if err != nil {
			if isZStackNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking zstack_role %s destroyed: %w", id, err)
		}
		if len(items) > 0 {
			return fmt.Errorf("zstack_role %s still exists", id)
		}
	}
	return nil
}

func testAccCheckPolicyDestroy(s *terraform.State) error {
	cli := testAccClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zstack_policy" {
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
		items, err := cli.QueryPolicy(&q)
		if err != nil {
			if isZStackNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking zstack_policy %s destroyed: %w", id, err)
		}
		if len(items) > 0 {
			return fmt.Errorf("zstack_policy %s still exists", id)
		}
	}
	return nil
}

var testAccCheckSecurityGroupDestroy = testAccCheckResourceDestroyByGet("zstack_networking_secgroup", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetSecurityGroup(id)
	return err
})

var testAccCheckSecurityGroupRuleDestroy = testAccCheckResourceDestroyByGet("zstack_networking_secgroup_rule", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetSecurityGroupRule(id)
	return err
})

func testAccCheckSecGroupAttachmentDestroy(s *terraform.State) error {
	cli := testAccClientLoggedIn()
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

func testAccCheckIAM2VirtualIDDestroy(s *terraform.State) error {
	cli := testAccClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zstack_iam2_virtual_id" {
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
		items, err := cli.QueryIAM2VirtualID(&q)
		if err != nil {
			if isZStackNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking zstack_iam2_virtual_id %s destroyed: %w", id, err)
		}
		if len(items) > 0 {
			return fmt.Errorf("zstack_iam2_virtual_id %s still exists", id)
		}
	}
	return nil
}

func testAccCheckIAM2OrganizationDestroy(s *terraform.State) error {
	cli := testAccClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zstack_iam2_organization" {
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
		items, err := cli.QueryIAM2Organization(&q)
		if err != nil {
			if isZStackNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking zstack_iam2_organization %s destroyed: %w", id, err)
		}
		if len(items) > 0 {
			return fmt.Errorf("zstack_iam2_organization %s still exists", id)
		}
	}
	return nil
}

func testAccCheckSNSTopicDestroy(s *terraform.State) error {
	cli := testAccClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zstack_sns_topic" {
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
		items, err := cli.QuerySNSTopic(&q)
		if err != nil {
			if isZStackNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking zstack_sns_topic %s destroyed: %w", id, err)
		}
		if len(items) > 0 {
			return fmt.Errorf("zstack_sns_topic %s still exists", id)
		}
	}
	return nil
}

func testAccCheckSNSEmailEndpointDestroy(s *terraform.State) error {
	// SNS email endpoints do not support deletion via API — Delete only removes from Terraform state.
	return nil
}

func testAccCheckSNSHttpEndpointDestroy(s *terraform.State) error {
	// SNS HTTP endpoints do not support deletion via API — Delete only removes from Terraform state.
	return nil
}

func testAccCheckGlobalConfigDestroy(s *terraform.State) error {
	// GlobalConfig resources are not destroyed — they are reset to default.
	return nil
}

func testAccCheckAlarmDestroy(s *terraform.State) error {
	cli := testAccClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zstack_alarm" {
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
		items, err := cli.QueryAlarm(&q)
		if err != nil {
			if isZStackNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking zstack_alarm %s destroyed: %w", id, err)
		}
		if len(items) > 0 {
			return fmt.Errorf("zstack_alarm %s still exists", id)
		}
	}
	return nil
}

func testAccCheckSchedulerTriggerDestroy(s *terraform.State) error {
	cli := testAccClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zstack_scheduler_trigger" {
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
		items, err := cli.QuerySchedulerTrigger(&q)
		if err != nil {
			if isZStackNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking zstack_scheduler_trigger %s destroyed: %w", id, err)
		}
		if len(items) > 0 {
			return fmt.Errorf("zstack_scheduler_trigger %s still exists", id)
		}
	}
	return nil
}

func testAccCheckSchedulerJobDestroy(s *terraform.State) error {
	cli := testAccClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zstack_scheduler_job" {
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
		items, err := cli.QuerySchedulerJob(&q)
		if err != nil {
			if isZStackNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking zstack_scheduler_job %s destroyed: %w", id, err)
		}
		if len(items) > 0 {
			return fmt.Errorf("zstack_scheduler_job %s still exists", id)
		}
	}
	return nil
}

func testAccCheckTagAttachmentDestroy(s *terraform.State) error {
	// Tag attachments are removed when the tag or resources are deleted.
	// No direct destroy check needed.
	return nil
}

func testAccCheckVirtualRouterInstanceDestroy(s *terraform.State) error {
	cli := testAccClientLoggedIn()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zstack_virtual_router_instance" {
			continue
		}
		id := rs.Primary.Attributes["uuid"]
		if id == "" {
			id = rs.Primary.ID
		}
		if id == "" {
			continue
		}
		_, err := cli.GetVirtualRouterVm(id)
		if err != nil {
			if isZStackNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking zstack_virtual_router_instance %s destroyed: %w", id, err)
		}
		return fmt.Errorf("zstack_virtual_router_instance %s still exists", id)
	}
	return nil
}

// Story-10: Self-contained resources destroy check functions

var testAccCheckAccessKeyDestroy = testAccCheckResourceDestroyByQuery("zstack_access_key", func(cli *client.ZSClient, q *param.QueryParam) ([]view.AccessKeyInventoryView, error) {
	return cli.QueryAccessKey(q)
})

var testAccCheckMonitorTemplateDestroy = testAccCheckResourceDestroyByQuery("zstack_monitor_template", func(cli *client.ZSClient, q *param.QueryParam) ([]view.MonitorTemplateInventoryView, error) {
	return cli.QueryMonitorTemplate(q)
})

var testAccCheckLogServerDestroy = testAccCheckResourceDestroyByQuery("zstack_log_server", func(cli *client.ZSClient, q *param.QueryParam) ([]view.LogServerInventoryView, error) {
	return cli.QueryLogServer(q)
})

var testAccCheckCdpPolicyDestroy = testAccCheckResourceDestroyByQuery("zstack_cdp_policy", func(cli *client.ZSClient, q *param.QueryParam) ([]view.CdpPolicyInventoryView, error) {
	return cli.QueryCdpPolicy(q)
})

var testAccCheckPriceTableDestroy = testAccCheckResourceDestroyByQuery("zstack_price_table", func(cli *client.ZSClient, q *param.QueryParam) ([]view.PriceTableInventoryView, error) {
	return cli.QueryPriceTable(q)
})

var testAccCheckStackTemplateDestroy = testAccCheckResourceDestroyByQuery("zstack_stack_template", func(cli *client.ZSClient, q *param.QueryParam) ([]view.StackTemplateInventoryView, error) {
	return cli.QueryStackTemplate(q)
})

var testAccCheckPreconfigurationTemplateDestroy = testAccCheckResourceDestroyByQuery("zstack_preconfiguration_template", func(cli *client.ZSClient, q *param.QueryParam) ([]view.PreconfigurationTemplateInventoryView, error) {
	return cli.QueryPreconfigurationTemplate(q)
})

// Story-11: Compute & Storage destroy check functions

var testAccCheckScriptDestroy = testAccCheckResourceDestroyByGet("zstack_instance_scripts", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetGuestVmScript(id)
	return err
})

func testAccCheckInstanceDestroy(s *terraform.State) error {
	cli := testAccClientLoggedIn()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zstack_instance" {
			continue
		}
		id := rs.Primary.Attributes["uuid"]
		if id == "" {
			id = rs.Primary.ID
		}
		if id == "" {
			continue
		}
		vm, err := cli.GetVmInstance(id)
		if err != nil {
			if isZStackNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error checking zstack_instance %s destroyed: %w", id, err)
		}
		// VM may be in "Destroyed" state in recycle bin — that counts as destroyed.
		if vm.State == "Destroyed" {
			continue
		}
		return fmt.Errorf("zstack_instance %s still exists with state %s", id, vm.State)
	}
	return nil
}

var testAccCheckVmCdRomDestroy = testAccCheckResourceDestroyByQuery("zstack_vm_cdrom", func(cli *client.ZSClient, q *param.QueryParam) ([]view.VmCdRomInventoryView, error) {
	return cli.QueryVmCdRom(q)
})

var testAccCheckVmNicDestroy = testAccCheckResourceDestroyByQuery("zstack_vm_nic", func(cli *client.ZSClient, q *param.QueryParam) ([]view.VmNicInventoryView, error) {
	return cli.QueryVmNic(q)
})

func testAccCheckGuestToolsAttachmentDestroy(s *terraform.State) error {
	// Delete is a no-op for guest_tools_attachment — the ISO remains attached.
	// There is no "destroyed" state to check.
	return nil
}

func testAccCheckScriptExecutionDestroy(s *terraform.State) error {
	// Delete only removes the record from Terraform state; no API call is made.
	// Nothing to verify on the ZStack side.
	return nil
}

var testAccCheckVolumeSnapshotDestroy = testAccCheckResourceDestroyByGet("zstack_volume_snapshot", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetVolumeSnapshot(id)
	return err
})

var testAccCheckVolumeBackupDestroy = testAccCheckResourceDestroyByGet("zstack_volume_backup", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetVolumeBackup(id)
	return err
})

// Story-13: Monitoring & Ops destroy check functions

var testAccCheckMonitorGroupDestroy = testAccCheckResourceDestroyByQuery("zstack_monitor_group", func(cli *client.ZSClient, q *param.QueryParam) ([]view.MonitorGroupInventoryView, error) {
	return cli.QueryMonitorGroup(q)
})

var testAccCheckCdpTaskDestroy = testAccCheckResourceDestroyByQuery("zstack_cdp_task", func(cli *client.ZSClient, q *param.QueryParam) ([]view.CdpTaskInventoryView, error) {
	return cli.QueryCdpTask(q)
})

var testAccCheckBackupStorageDestroy = testAccCheckResourceDestroyByGet("zstack_backup_storage", func(cli *client.ZSClient, id string) error {
	_, err := cli.GetBackupStorage(id)
	return err
})
