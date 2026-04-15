// Copyright (c) ZStack.io, Inc.

package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"testing"
)

// EnvData holds all environment resource data loaded from testdata/env.json
type EnvData struct {
	// Infrastructure
	Zones           []map[string]interface{} `json:"zones"`
	Clusters        []map[string]interface{} `json:"clusters"`
	Hosts           []map[string]interface{} `json:"hosts"`
	PrimaryStorages []map[string]interface{} `json:"primary_storages"`
	BackupStorages  []map[string]interface{} `json:"backup_storages"`

	// Compute
	Images             []map[string]interface{} `json:"images"`
	InstanceOfferings  []map[string]interface{} `json:"instance_offerings"`
	DiskOfferings      []map[string]interface{} `json:"disk_offerings"`
	VmInstances        []map[string]interface{} `json:"vm_instances"`
	GpuDevices         []map[string]interface{} `json:"gpu_devices"`
	AutoScalingGroups  []map[string]interface{} `json:"auto_scaling_groups"`
	PciDeviceOfferings []map[string]interface{} `json:"pci_device_offerings"`
	VCenters           []map[string]interface{} `json:"vcenters"`

	// Storage
	Volumes                  []map[string]interface{} `json:"volumes"`
	VolumeSnapshots          []map[string]interface{} `json:"volume_snapshots"`
	CephPrimaryStorages      []map[string]interface{} `json:"ceph_primary_storages"`
	CephBackupStorages       []map[string]interface{} `json:"ceph_backup_storages"`
	CephPools                []map[string]interface{} `json:"ceph_pools"`
	ImageStoreBackupStorages []map[string]interface{} `json:"image_store_backup_storages"`
	VolumeBackups            []map[string]interface{} `json:"volume_backups"`
	DatabaseBackups          []map[string]interface{} `json:"database_backups"`
	NvmeServers              []map[string]interface{} `json:"nvme_servers"`
	IscsiServers             []map[string]interface{} `json:"iscsi_servers"`

	// Network
	L2Networks           []map[string]interface{} `json:"l2_networks"`
	L2VlanNetworks       []map[string]interface{} `json:"l2_vlan_networks"`
	L3Networks           []map[string]interface{} `json:"l3_networks"`
	IpRanges             []map[string]interface{} `json:"ip_ranges"`
	Vips                 []map[string]interface{} `json:"vips"`
	Eips                 []map[string]interface{} `json:"eips"`
	PortForwardingRules  []map[string]interface{} `json:"port_forwarding_rules"`
	LoadBalancers        []map[string]interface{} `json:"load_balancers"`
	LoadBalancerListeners []map[string]interface{} `json:"load_balancer_listeners"`
	SecurityGroups       []map[string]interface{} `json:"security_groups"`
	SecurityGroupRules   []map[string]interface{} `json:"security_group_rules"`
	VmNics               []map[string]interface{} `json:"vm_nics"`
	AccessControlLists   []map[string]interface{} `json:"access_control_lists"`
	Certificates         []map[string]interface{} `json:"certificates"`
	FlowCollectors       []map[string]interface{} `json:"flow_collectors"`
	FlowMeters           []map[string]interface{} `json:"flow_meters"`
	IPSecConnections     []map[string]interface{} `json:"ipsec_connections"`
	MulticastRouters     []map[string]interface{} `json:"multicast_routers"`
	PolicyRouteRuleSets  []map[string]interface{} `json:"policy_route_rule_sets"`
	PolicyRouteRules     []map[string]interface{} `json:"policy_route_rules"`
	VpcFirewalls         []map[string]interface{} `json:"vpc_firewalls"`
	VpcHaGroups          []map[string]interface{} `json:"vpc_ha_groups"`
	VpcSharedQos         []map[string]interface{} `json:"vpc_shared_qos"`
	VRouterRouteTables   []map[string]interface{} `json:"vrouter_route_tables"`
	VRouterRouteEntries  []map[string]interface{} `json:"vrouter_route_entries"`

	// Virtual Router
	VirtualRouterOfferings []map[string]interface{} `json:"virtual_router_offerings"`
	VirtualRouters         []map[string]interface{} `json:"virtual_routers"`

	// System / IAM
	Accounts         []map[string]interface{} `json:"accounts"`
	IAM2Projects     []map[string]interface{} `json:"iam2_projects"`
	AccessKeys       []map[string]interface{} `json:"access_keys"`
	IAM2Organizations []map[string]interface{} `json:"iam2_organizations"`
	IAM2VirtualIDs   []map[string]interface{} `json:"iam2_virtual_ids"`
	Roles            []map[string]interface{} `json:"roles"`
	Users            []map[string]interface{} `json:"users"`

	// Monitoring
	Alarms            []map[string]interface{} `json:"alarms"`
	EmailMedia        []map[string]interface{} `json:"email_media"`
	LogServers        []map[string]interface{} `json:"log_servers"`
	MonitorGroups     []map[string]interface{} `json:"monitor_groups"`
	MonitorTemplates  []map[string]interface{} `json:"monitor_templates"`
	SNSEmailEndpoints []map[string]interface{} `json:"sns_email_endpoints"`
	SNSTopics         []map[string]interface{} `json:"sns_topics"`
	SnmpAgents        []map[string]interface{} `json:"snmp_agents"`
	Webhooks          []map[string]interface{} `json:"webhooks"`

	// Auxiliary
	AffinityGroups []map[string]interface{} `json:"affinity_groups"`
	SshKeyPairs    []map[string]interface{} `json:"ssh_key_pairs"`
	UserTags       []map[string]interface{} `json:"user_tags"`
	SystemTags     []map[string]interface{} `json:"system_tags"`

	// Operations
	SdnControllers   []map[string]interface{} `json:"sdn_controllers"`
	InstanceScripts  []map[string]interface{} `json:"instance_scripts"`
	ScriptExecutions []map[string]interface{} `json:"script_executions"`
	MnNodes          []map[string]interface{} `json:"mn_nodes"`

	// Automation
	LdapServers       []map[string]interface{} `json:"ldap_servers"`
	PortMirrors       []map[string]interface{} `json:"port_mirrors"`
	PortMirrorSessions []map[string]interface{} `json:"port_mirror_sessions"`
	PriceTables       []map[string]interface{} `json:"price_tables"`
	SchedulerJobs     []map[string]interface{} `json:"scheduler_jobs"`
	SchedulerTriggers []map[string]interface{} `json:"scheduler_triggers"`
}

var (
	envData     *EnvData
	envDataOnce sync.Once
	envDataErr  error
)

// loadEnvData loads testdata/env.json. It is safe for concurrent use.
// If the file does not exist, it calls t.Skip.
func loadEnvData(t *testing.T) *EnvData {
	t.Helper()
	envDataOnce.Do(func() {
		_, filename, _, _ := runtime.Caller(0)
		dir := filepath.Dir(filename)
		path := filepath.Join(dir, "testdata", "env.json")

		data, err := os.ReadFile(path)
		if err != nil {
			envDataErr = fmt.Errorf("cannot read testdata/env.json: %w (run: source .env.test && go run ./zstack/provider/testdata/generate_env.go)", err)
			return
		}
		var env EnvData
		if err := json.Unmarshal(data, &env); err != nil {
			envDataErr = fmt.Errorf("cannot parse testdata/env.json: %w", err)
			return
		}
		envData = &env
	})
	if envDataErr != nil {
		t.Skip(envDataErr.Error())
	}
	return envData
}

// Helper functions to extract string values from map[string]interface{}

func envStr(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case string:
			return val
		case float64:
			if val == float64(int64(val)) {
				return strconv.FormatInt(int64(val), 10)
			}
			return strconv.FormatFloat(val, 'f', -1, 64)
		case bool:
			return strconv.FormatBool(val)
		case nil:
			return ""
		default:
			return fmt.Sprintf("%v", val)
		}
	}
	return ""
}

func envInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val)
		}
	}
	return 0
}
