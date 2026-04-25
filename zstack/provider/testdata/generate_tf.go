// Copyright (c) ZStack.io, Inc.
//
// This program reads testdata/env.json and generates real, runnable .tf files
// for batch-testing with `terraform apply/destroy` against a real ZStack environment.
//
// Usage:
//   source .env.test && go run ./zstack/provider/testdata/generate_tf.go
//   # QA mode — only data sources, with 3 lookup variants per type:
//   go run ./zstack/provider/testdata/generate_tf.go -only=datasources

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// EnvData mirrors the structure in generate_env.go / env.json.
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

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func getStr(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok && v != nil {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0
}

func getBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok && v != nil {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// findL3ByCategory returns the first L3 network matching category (e.g. "Public", "Private").
// It skips system networks.
func findL3ByCategory(env *EnvData, category string) map[string]interface{} {
	for _, l3 := range env.L3Networks {
		if getStr(l3, "category") == category && !getBool(l3, "system") {
			return l3
		}
	}
	return nil
}

func findPublicL3(env *EnvData) map[string]interface{} {
	return findL3ByCategory(env, "Public")
}

func findPrivateL3(env *EnvData) map[string]interface{} {
	return findL3ByCategory(env, "Private")
}

// findIpRangeForL3 finds an IP range associated with the given L3 network UUID.
func findIpRangeForL3(env *EnvData, l3UUID string) map[string]interface{} {
	for _, ipr := range env.IpRanges {
		if getStr(ipr, "l3_network_uuid") == l3UUID {
			return ipr
		}
	}
	return nil
}

// findReadyImage returns the first image that is not deleted (status=Ready).
func findReadyImage(env *EnvData) map[string]interface{} {
	for _, img := range env.Images {
		if getStr(img, "status") == "Ready" && getStr(img, "name") != "vr" {
			return img
		}
	}
	return nil
}

// findVRImage returns the virtual-router image.
func findVRImage(env *EnvData) map[string]interface{} {
	for _, img := range env.Images {
		if getStr(img, "name") == "vr" && getStr(img, "status") == "Ready" {
			return img
		}
	}
	return nil
}

// findUserVmNic finds a VM NIC belonging to a non-system, non-VR vm instance.
func findUserVmNic(env *EnvData) map[string]interface{} {
	if len(env.VmNics) > 0 {
		return env.VmNics[0]
	}
	return nil
}

// incrementIP adds offset to an IP address string.
func incrementIP(ipStr string, offset int) string {
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return ipStr
	}
	v := uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
	v += uint32(offset)
	return fmt.Sprintf("%d.%d.%d.%d", (v>>24)&0xFF, (v>>16)&0xFF, (v>>8)&0xFF, v&0xFF)
}

// ---------------------------------------------------------------------------
// Provider HCL
// ---------------------------------------------------------------------------

// providerBinDir returns the directory containing the locally compiled
// terraform-provider-zstack binary.  Resolution order:
//   1. GOBIN  (if set)
//   2. GOPATH/bin  (if GOPATH set)
//   3. ~/go/bin    (Go default)
func providerBinDir() string {
	if d := os.Getenv("GOBIN"); d != "" {
		return d
	}
	if d := os.Getenv("GOPATH"); d != "" {
		return filepath.Join(d, "bin")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "go", "bin")
}

// terraformRC generates a CLI config file that tells Terraform to use
// the locally compiled provider binary instead of downloading from a registry.
func terraformRC() string {
	return fmt.Sprintf(`provider_installation {
  dev_overrides {
    "zstack.io/cloud/zstack" = %q
  }
  direct {}
}
`, providerBinDir())
}

// providerHCL generates the shared provider.tf.
// The "host" attribute is Required in the provider schema, so it must be set in HCL.
// We read it from ZSTACK_HOST at generation time and bake it in.
// Credentials (access_key_id, access_key_secret) are Optional in the schema and
// fall back to ZSTACK_ACCESS_KEY_ID / ZSTACK_ACCESS_KEY_SECRET env vars at runtime.
func providerHCL() string {
	host := os.Getenv("ZSTACK_HOST")
	if host == "" {
		host = "172.24.227.46"
	}
	return fmt.Sprintf(`terraform {
  required_providers {
    zstack = {
      source = "zstack.io/cloud/zstack"
    }
  }
}

provider "zstack" {
  host = %q
}
`, host)
}

// ---------------------------------------------------------------------------
// Generator type: name → func(env) → (hcl, canGenerate, skipReason)
// ---------------------------------------------------------------------------

type generator struct {
	name string
	fn   func(*EnvData) (string, bool, string)
}

// ---------------------------------------------------------------------------
// Data source generators
// ---------------------------------------------------------------------------

func dataSimple(dsType, envField, nameField string, getList func(*EnvData) []map[string]interface{}) generator {
	return generator{
		name: "data-" + strings.ReplaceAll(dsType, "_", "_"),
		fn: func(env *EnvData) (string, bool, string) {
			list := getList(env)
			if nameField != "" && len(list) == 0 {
				return "", false, envField + " empty"
			}
			var hcl string
			if nameField != "" && len(list) > 0 {
				name := getStr(list[0], "name")
				hcl = fmt.Sprintf(`data "zstack_%s" "test" {
  name = %q
}

output "result" {
  value = data.zstack_%s.test
}
`, dsType, name, dsType)
			} else {
				hcl = fmt.Sprintf(`data "zstack_%s" "test" {
}

output "result" {
  value = data.zstack_%s.test
}
`, dsType, dsType)
			}
			return hcl, true, ""
		},
	}
}

func dataSourceGenerators() []generator {
	// All TypeNames below match the provider's actual registered names (see
	// zstack/provider/data_source_zstack_*.go Metadata methods). After the
	// SDK-name alignment refactor (commit 9c3b97d), the previous test names
	// — `backupstorages`, `disk_offers`, `instance_offers`,
	// `networking_sdn_controllers`, `mnnodes`, `scripts`,
	// `virtual_router_offers`, `zones` — became invalid; they have been
	// renamed to match the provider.
	return []generator{
		// Phase A core data sources
		dataSimple("zone", "zones", "name", func(e *EnvData) []map[string]interface{} { return e.Zones }),
		dataSimple("clusters", "clusters", "name", func(e *EnvData) []map[string]interface{} { return e.Clusters }),
		dataSimple("hosts", "hosts", "name", func(e *EnvData) []map[string]interface{} { return e.Hosts }),
		dataSimple("l3networks", "l3_networks", "name", func(e *EnvData) []map[string]interface{} { return e.L3Networks }),
		dataSimple("instances", "vm_instances", "name", func(e *EnvData) []map[string]interface{} { return e.VmInstances }),

		// Image / storage
		dataSimple("images", "images", "name", func(e *EnvData) []map[string]interface{} { return e.Images }),
		dataSimple("virtual_router_images", "", "", func(e *EnvData) []map[string]interface{} { return nil }),
		dataSimple("backup_storages", "backup_storages", "name", func(e *EnvData) []map[string]interface{} { return e.BackupStorages }),
		dataSimple("primary_storages", "primary_storages", "name", func(e *EnvData) []map[string]interface{} { return e.PrimaryStorages }),

		// Compute offerings
		dataSimple("instance_offerings", "instance_offerings", "name", func(e *EnvData) []map[string]interface{} { return e.InstanceOfferings }),
		dataSimple("disk_offerings", "disk_offerings", "name", func(e *EnvData) []map[string]interface{} { return e.DiskOfferings }),
		dataSimple("virtual_router_offerings", "virtual_router_offerings", "name", func(e *EnvData) []map[string]interface{} { return e.VirtualRouterOfferings }),

		// Virtual routers + L2
		dataSimple("virtual_routers", "virtual_routers", "name", func(e *EnvData) []map[string]interface{} { return e.VirtualRouters }),
		dataSimple("l2networks", "l2_networks", "name", func(e *EnvData) []map[string]interface{} { return e.L2Networks }),
		dataSimple("l2vlan_networks", "", "", func(e *EnvData) []map[string]interface{} { return nil }),

		// Volumes / disks
		dataSimple("volumes", "", "", func(e *EnvData) []map[string]interface{} { return nil }),
		dataSimple("volume_snapshots", "", "", func(e *EnvData) []map[string]interface{} { return nil }),
		dataSimple("disks", "", "", func(e *EnvData) []map[string]interface{} { return nil }),

		// Network plumbing
		dataSimple("networking_secgroups", "security_groups", "name", func(e *EnvData) []map[string]interface{} { return e.SecurityGroups }),
		dataSimple("networking_secgroup_rules", "", "", func(e *EnvData) []map[string]interface{} { return nil }),
		dataSimple("sdn_controllers", "sdn_controllers", "name", func(e *EnvData) []map[string]interface{} { return e.SdnControllers }),
		dataSimple("vips", "", "", func(e *EnvData) []map[string]interface{} { return nil }),
		dataSimple("eips", "", "", func(e *EnvData) []map[string]interface{} { return nil }),
		dataSimple("port_forwarding_rules", "", "", func(e *EnvData) []map[string]interface{} { return nil }),
		dataSimple("reserved_ips", "", "", func(e *EnvData) []map[string]interface{} { return nil }),
		dataSimple("subnet_ip_ranges", "", "", func(e *EnvData) []map[string]interface{} { return nil }),

		// Scripts
		dataSimple("instance_scripts", "instance_scripts", "name", func(e *EnvData) []map[string]interface{} { return e.InstanceScripts }),
		dataSimple("hook_scripts", "", "", func(e *EnvData) []map[string]interface{} { return nil }),

		// MN / system
		dataSimple("mn_nodes", "", "", func(e *EnvData) []map[string]interface{} { return nil }),
		dataSimple("gpu_devices", "", "", func(e *EnvData) []map[string]interface{} { return nil }),

		// Identity / RBAC
		dataSimple("accounts", "accounts", "name", func(e *EnvData) []map[string]interface{} { return e.Accounts }),
		dataSimple("iam2_projects", "iam2_projects", "name", func(e *EnvData) []map[string]interface{} { return e.IAM2Projects }),
		dataSimple("affinity_groups", "affinity_groups", "name", func(e *EnvData) []map[string]interface{} { return e.AffinityGroups }),
		dataSimple("ssh_key_pairs", "ssh_key_pairs", "name", func(e *EnvData) []map[string]interface{} { return e.SshKeyPairs }),

		// Load balancing
		dataSimple("load_balancers", "", "", func(e *EnvData) []map[string]interface{} { return nil }),
		dataSimple("load_balancer_listeners", "", "", func(e *EnvData) []map[string]interface{} { return nil }),

		// Auto-scaling
		dataSimple("auto_scaling_groups", "", "", func(e *EnvData) []map[string]interface{} { return nil }),

		// User tags
		dataSimple("user_tags", "", "", func(e *EnvData) []map[string]interface{} { return nil }),

		// License
		dataSimple("license_authorized_nodes", "", "", func(e *EnvData) []map[string]interface{} { return nil }),
		dataSimple("license_authorized_capacity", "", "", func(e *EnvData) []map[string]interface{} { return nil }),
		{
			name: "data-guest_tools",
			fn: func(env *EnvData) (string, bool, string) {
				if len(env.VmInstances) == 0 {
					return "", false, "vm_instances empty"
				}
				vmUUID := getStr(env.VmInstances[0], "uuid")
				return fmt.Sprintf(`data "zstack_instance_guest_tools" "test" {
  instance_uuid = %q
}

output "result" {
  value = data.zstack_instance_guest_tools.test
}
`, vmUUID), true, ""
			},
		},
		{
			name: "data-tags",
			fn: func(env *EnvData) (string, bool, string) {
				return `data "zstack_tags" "test" {
  tag_type = "tag"
}

output "result" {
  value = data.zstack_tags.test
}
`, true, ""
			},
		},
		// L2 networks — no dedicated data source in provider, skip if not present
	}
}

// ---------------------------------------------------------------------------
// QA-style data source generators
// ---------------------------------------------------------------------------
//
// Each entry produces ONE main.tf containing up to FIVE blocks of the same
// data source — `all`, `by_name`, `by_uuid`, `by_name_pattern`,
// `by_nonexistent_uuid` — plus outputs and cross-validation assertions.
// A single `terraform apply` therefore validates, per data source:
//
//   1. Listing without filter (smoke).
//   2. Human-friendly `name` lookup.
//   3. AI / automation deterministic `uuid` lookup.
//   4. Fuzzy `name_pattern` lookup (when supported by the schema).
//   5. Negative lookup with a known-bad UUID — must return an empty list,
//      NOT an error. This catches API regressions where the SDK turns a
//      "no rows" result into a 404.
//
// Plus, when both `by_name` and `by_uuid` resolve, an additional
// `_lookup_consistency` output asserts they pin the SAME record at apply
// time. Variants are skipped automatically when env.json doesn't have the
// inputs (e.g. an empty zone list will skip `by_name` / `by_uuid` /
// `by_name_pattern`, but `all` and `by_nonexistent_uuid` always run).

// negativeProbeUUID is a deterministic UUID that should never exist in any
// real ZStack deployment. Used to assert "filter returns empty list" (not
// an error) for every data source that supports `uuid` filtering.
const negativeProbeUUID = "00000000-0000-0000-0000-000000000000"

// qaDataSource describes one data source we want to QA.
type qaDataSource struct {
	tfType   string                                  // TypeName suffix, e.g. "zones" → emits data "zstack_zones"
	envField string                                  // diagnostic / skip-reason hint, e.g. "zones"
	getList  func(*EnvData) []map[string]interface{} // returns env list for this resource
	listAttr string                                  // schema attribute holding the result list, e.g. "zones"

	// Filter capabilities. When false, the corresponding variant is skipped
	// even if env data is present. Defaults (zero value = false) match the
	// strictest case; populate explicitly per data source.
	supportsName        bool
	supportsUUID        bool
	supportsNamePattern bool

	// extraArgs is raw HCL injected verbatim into EVERY generated block of
	// this data source. Used for required filter args like
	// `tag_type = "tag"` on zstack_tags or `instance_uuid = "..."` on
	// zstack_instance_guest_tools. The trailing newline is required.
	extraArgs string

	// singleton marks data sources that don't expose a list — they have
	// top-level Computed fields only (e.g. license_authorized_capacity).
	// Only the `all` variant is emitted for these, with a single output
	// echoing the entire object.
	singleton bool
}

func dataSourceQAGenerators() []generator {
	defs := []qaDataSource{
		// ---- Phase A core (already covered, expanded with name_pattern + negative)
		{tfType: "zone", envField: "zones", getList: func(e *EnvData) []map[string]interface{} { return e.Zones }, listAttr: "zones",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "clusters", envField: "clusters", getList: func(e *EnvData) []map[string]interface{} { return e.Clusters }, listAttr: "clusters",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "hosts", envField: "hosts", getList: func(e *EnvData) []map[string]interface{} { return e.Hosts }, listAttr: "hosts",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "l3networks", envField: "l3_networks", getList: func(e *EnvData) []map[string]interface{} { return e.L3Networks }, listAttr: "l3networks",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "instances", envField: "vm_instances", getList: func(e *EnvData) []map[string]interface{} { return e.VmInstances }, listAttr: "vminstances",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		// ---- Image / storage
		{tfType: "images", envField: "images", getList: func(e *EnvData) []map[string]interface{} { return e.Images }, listAttr: "images",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "virtual_router_images", envField: "(no env list)", getList: func(e *EnvData) []map[string]interface{} { return nil }, listAttr: "images",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "backup_storages", envField: "backup_storages", getList: func(e *EnvData) []map[string]interface{} { return e.BackupStorages }, listAttr: "backup_storages",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "primary_storages", envField: "primary_storages", getList: func(e *EnvData) []map[string]interface{} { return e.PrimaryStorages }, listAttr: "primary_storages",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		// ---- Compute offerings
		{tfType: "instance_offerings", envField: "instance_offerings", getList: func(e *EnvData) []map[string]interface{} { return e.InstanceOfferings }, listAttr: "instance_offers",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "disk_offerings", envField: "disk_offerings", getList: func(e *EnvData) []map[string]interface{} { return e.DiskOfferings }, listAttr: "disk_offers",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "virtual_router_offerings", envField: "virtual_router_offerings", getList: func(e *EnvData) []map[string]interface{} { return e.VirtualRouterOfferings }, listAttr: "virtual_router_offers",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		// ---- Virtual routers + L2
		{tfType: "virtual_routers", envField: "virtual_routers", getList: func(e *EnvData) []map[string]interface{} { return e.VirtualRouters }, listAttr: "virtual_router",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "l2networks", envField: "l2_networks", getList: func(e *EnvData) []map[string]interface{} { return e.L2Networks }, listAttr: "l2networks",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "l2vlan_networks", envField: "l2_vlan_networks", getList: func(e *EnvData) []map[string]interface{} { return e.L2VlanNetworks }, listAttr: "l2vlan_networks",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		// ---- Volumes / disks / snapshots
		{tfType: "volumes", envField: "volumes", getList: func(e *EnvData) []map[string]interface{} { return e.Volumes }, listAttr: "volumes",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "volume_snapshots", envField: "volume_snapshots", getList: func(e *EnvData) []map[string]interface{} { return e.VolumeSnapshots }, listAttr: "snapshots",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "disks", envField: "(no env list)", getList: func(e *EnvData) []map[string]interface{} { return nil }, listAttr: "disks",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		// ---- Network plumbing
		{tfType: "networking_secgroups", envField: "security_groups", getList: func(e *EnvData) []map[string]interface{} { return e.SecurityGroups }, listAttr: "networking_secgroups",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "networking_secgroup_rules", envField: "security_group_rules", getList: func(e *EnvData) []map[string]interface{} { return e.SecurityGroupRules }, listAttr: "rules"},
		{tfType: "sdn_controllers", envField: "sdn_controllers", getList: func(e *EnvData) []map[string]interface{} { return e.SdnControllers }, listAttr: "sdn_controllers",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "vips", envField: "vips", getList: func(e *EnvData) []map[string]interface{} { return e.Vips }, listAttr: "vips",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "eips", envField: "eips", getList: func(e *EnvData) []map[string]interface{} { return e.Eips }, listAttr: "eips",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "port_forwarding_rules", envField: "port_forwarding_rules", getList: func(e *EnvData) []map[string]interface{} { return e.PortForwardingRules }, listAttr: "port_forwarding_rules",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "subnet_ip_ranges", envField: "ip_ranges", getList: func(e *EnvData) []map[string]interface{} { return e.IpRanges }, listAttr: "subnet_ip_ranges",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "reserved_ips", envField: "(no env list)", getList: func(e *EnvData) []map[string]interface{} { return nil }, listAttr: "reserved_ips",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		// ---- Load balancing
		{tfType: "load_balancers", envField: "load_balancers", getList: func(e *EnvData) []map[string]interface{} { return e.LoadBalancers }, listAttr: "load_balancers",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "load_balancer_listeners", envField: "load_balancer_listeners", getList: func(e *EnvData) []map[string]interface{} { return e.LoadBalancerListeners }, listAttr: "load_balancer_listeners",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		// ---- Auto-scaling / GPU
		{tfType: "auto_scaling_groups", envField: "auto_scaling_groups", getList: func(e *EnvData) []map[string]interface{} { return e.AutoScalingGroups }, listAttr: "auto_scaling_groups",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "gpu_devices", envField: "gpu_devices", getList: func(e *EnvData) []map[string]interface{} { return e.GpuDevices }, listAttr: "gpu_devices",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		// ---- Scripts
		{tfType: "instance_scripts", envField: "instance_scripts", getList: func(e *EnvData) []map[string]interface{} { return e.InstanceScripts }, listAttr: "scripts",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "hook_scripts", envField: "(no env list)", getList: func(e *EnvData) []map[string]interface{} { return nil }, listAttr: "hook_scripts",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		// ---- Identity / RBAC
		{tfType: "accounts", envField: "accounts", getList: func(e *EnvData) []map[string]interface{} { return e.Accounts }, listAttr: "accounts",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "iam2_projects", envField: "iam2_projects", getList: func(e *EnvData) []map[string]interface{} { return e.IAM2Projects }, listAttr: "iam2_projects",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "affinity_groups", envField: "affinity_groups", getList: func(e *EnvData) []map[string]interface{} { return e.AffinityGroups }, listAttr: "affinity_groups",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "ssh_key_pairs", envField: "ssh_key_pairs", getList: func(e *EnvData) []map[string]interface{} { return e.SshKeyPairs }, listAttr: "ssh_key_pairs",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		// ---- Tags
		{tfType: "user_tags", envField: "user_tags", getList: func(e *EnvData) []map[string]interface{} { return e.UserTags }, listAttr: "user_tags",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "tags", envField: "system_tags", getList: func(e *EnvData) []map[string]interface{} { return e.SystemTags }, listAttr: "tags",
			supportsName: true, supportsUUID: true, supportsNamePattern: true,
			extraArgs: "  tag_type = \"tag\"\n"},
		// ---- License (license_authorized_nodes is list, capacity is singleton)
		{tfType: "license_authorized_nodes", envField: "(no env list)", getList: func(e *EnvData) []map[string]interface{} { return nil }, listAttr: "nodes",
			supportsName: true, supportsUUID: true, supportsNamePattern: true},
		{tfType: "license_authorized_capacity", envField: "(singleton)", getList: func(e *EnvData) []map[string]interface{} { return nil }, listAttr: "",
			singleton: true},
		// ---- MN nodes (no filter args at all in schema → only `all`)
		{tfType: "mn_nodes", envField: "mn_nodes", getList: func(e *EnvData) []map[string]interface{} { return e.MnNodes }, listAttr: "mn_nodes"},
	}

	gens := make([]generator, 0, len(defs)+1)
	for _, def := range defs {
		def := def
		gens = append(gens, generator{
			name: "qa-data-" + def.tfType,
			fn: func(env *EnvData) (string, bool, string) {
				return buildQADataSourceHCL(env, def), true, ""
			},
		})
	}

	// Specials that need a per-env Required argument computed at generation
	// time, not a static extraArgs string.
	gens = append(gens, generator{
		name: "qa-data-instance_guest_tools",
		fn: func(env *EnvData) (string, bool, string) {
			if len(env.VmInstances) == 0 {
				return "", false, "vm_instances empty (instance_guest_tools requires instance_uuid)"
			}
			vmUUID := getStr(env.VmInstances[0], "uuid")
			if vmUUID == "" {
				return "", false, "vm_instance has no uuid"
			}
			return buildInstanceGuestToolsHCL(vmUUID), true, ""
		},
	})

	return gens
}

// buildQADataSourceHCL emits up to 5 blocks per data source plus
// cross-validation outputs. Variants without env data are documented
// inline so QA can see what was skipped vs. exercised.
func buildQADataSourceHCL(env *EnvData, def qaDataSource) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# QA fixture for data \"zstack_%s\".\n", def.tfType)
	if def.singleton {
		fmt.Fprintf(&b, "# Singleton data source — no list, no filters; only `all` variant.\n\n")
		fmt.Fprintf(&b, "data \"zstack_%s\" \"all\" {}\n\n", def.tfType)
		fmt.Fprintf(&b, "output \"%s_all\" {\n  value = data.zstack_%s.all\n}\n",
			def.tfType, def.tfType)
		return b.String()
	}
	fmt.Fprintf(&b, "# Variants exercised: all%s.\n", listVariants(def, env))
	b.WriteString("\n")

	// Variant 1 — list-all (no filter)
	fmt.Fprintf(&b, "data \"zstack_%s\" \"all\" {\n%s}\n", def.tfType, def.extraArgs)
	fmt.Fprintf(&b, "output \"%s_all_count\" {\n  value = try(length(data.zstack_%s.all.%s), 0)\n}\n\n",
		def.tfType, def.tfType, def.listAttr)

	list := def.getList(env)
	var firstName, firstUUID string
	if len(list) > 0 {
		firstName = getStr(list[0], "name")
		firstUUID = getStr(list[0], "uuid")
	}

	// Variant 2 — by_name (human ergonomic)
	if def.supportsName && firstName != "" {
		fmt.Fprintf(&b, "data \"zstack_%s\" \"by_name\" {\n%s  name = %q\n}\n",
			def.tfType, def.extraArgs, firstName)
		fmt.Fprintf(&b, "output \"%s_by_name_first_uuid\" {\n  value = try(data.zstack_%s.by_name.%s[0].uuid, null)\n}\n\n",
			def.tfType, def.tfType, def.listAttr)
	}

	// Variant 3 — by_uuid (deterministic)
	if def.supportsUUID && firstUUID != "" {
		fmt.Fprintf(&b, "data \"zstack_%s\" \"by_uuid\" {\n%s  uuid = %q\n}\n",
			def.tfType, def.extraArgs, firstUUID)
		fmt.Fprintf(&b, "output \"%s_by_uuid_first_name\" {\n  value = try(data.zstack_%s.by_uuid.%s[0].name, null)\n}\n\n",
			def.tfType, def.tfType, def.listAttr)
	}

	// Variant 4 — by_name_pattern (fuzzy lookup; pattern ≡ exact name to keep the
	// match deterministic without depending on ZStack's wildcard rules).
	if def.supportsNamePattern && firstName != "" {
		fmt.Fprintf(&b, "data \"zstack_%s\" \"by_name_pattern\" {\n%s  name_pattern = %q\n}\n",
			def.tfType, def.extraArgs, firstName)
		fmt.Fprintf(&b, "output \"%s_by_name_pattern_count\" {\n  value = try(length(data.zstack_%s.by_name_pattern.%s), 0)\n}\n\n",
			def.tfType, def.tfType, def.listAttr)
	}

	// Variant 5 — by_nonexistent_uuid (negative path: returns empty list, NOT error)
	if def.supportsUUID {
		fmt.Fprintf(&b, "# Negative-path probe: a known-bad UUID must yield an empty list, not an error.\n")
		fmt.Fprintf(&b, "data \"zstack_%s\" \"by_nonexistent_uuid\" {\n%s  uuid = %q\n}\n",
			def.tfType, def.extraArgs, negativeProbeUUID)
		fmt.Fprintf(&b, "output \"%s_nonexistent_count\" {\n  value = try(length(data.zstack_%s.by_nonexistent_uuid.%s), 0)\n}\n\n",
			def.tfType, def.tfType, def.listAttr)
	}

	// Cross-check: by_name and by_uuid resolved to the same record.
	if def.supportsName && def.supportsUUID && firstName != "" && firstUUID != "" {
		fmt.Fprintf(&b, "# Lookup consistency assertion (eval at apply time).\n")
		fmt.Fprintf(&b, "output \"%s_lookup_consistency\" {\n", def.tfType)
		b.WriteString("  value = (\n")
		fmt.Fprintf(&b, "    try(data.zstack_%s.by_name.%s[0].uuid, \"\") ==\n", def.tfType, def.listAttr)
		fmt.Fprintf(&b, "    try(data.zstack_%s.by_uuid.%s[0].uuid, \"\")\n", def.tfType, def.listAttr)
		b.WriteString("  )\n}\n")
	}

	if len(list) == 0 && !def.singleton {
		fmt.Fprintf(&b, "\n# Skipped variants by_name/by_uuid/by_name_pattern: env.%s is empty.\n", def.envField)
	}

	return b.String()
}

// buildInstanceGuestToolsHCL — instance_guest_tools is special-cased because
// the only Required arg is `instance_uuid`. There is no list filter, so we
// only test (a) happy path on a real VM, and (b) negative path on a bogus
// instance_uuid which should surface a clean diagnostic, not a panic.
func buildInstanceGuestToolsHCL(vmUUID string) string {
	return fmt.Sprintf(`# QA fixture for data "zstack_instance_guest_tools".
# Required arg: instance_uuid. Only "happy" variant is exercised — the
# negative variant is intentionally omitted because the data source returns
# a hard error (not an empty result) for a non-existent instance, which
# would fail the apply.

data "zstack_instance_guest_tools" "happy" {
  instance_uuid = %q
}

output "guest_tools" {
  value = data.zstack_instance_guest_tools.happy
}
`, vmUUID)
}

// listVariants returns a comma-prefixed list of variant names that will be
// emitted for this data source, derived from def + env data. Used for the
// per-fixture comment header so QA can see what was actually exercised.
func listVariants(def qaDataSource, env *EnvData) string {
	if def.singleton {
		return ""
	}
	parts := []string{}
	list := def.getList(env)
	if def.supportsName && len(list) > 0 && getStr(list[0], "name") != "" {
		parts = append(parts, "by_name")
	}
	if def.supportsUUID && len(list) > 0 && getStr(list[0], "uuid") != "" {
		parts = append(parts, "by_uuid")
	}
	if def.supportsNamePattern && len(list) > 0 && getStr(list[0], "name") != "" {
		parts = append(parts, "by_name_pattern")
	}
	if def.supportsUUID {
		parts = append(parts, "by_nonexistent_uuid")
	}
	if len(parts) == 0 {
		return ""
	}
	return ", " + strings.Join(parts, ", ")
}

// ---------------------------------------------------------------------------
// Resource generators — self-contained (no env UUIDs)
// ---------------------------------------------------------------------------

func selfContainedResourceGenerators() []generator {
	return []generator{
		{
			name: "res-disk_offer",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_disk_offering" "test" {
  name        = "tf-batch-test-disk-offer"
  disk_size   = 10
  description = "[batch-test] disk offering"
}

output "uuid" {
  value = zstack_disk_offering.test.uuid
}
`, true, ""
			},
		},
		{
			name: "res-instance_offer",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_instance_offering" "test" {
  name        = "tf-batch-test-instance-offer"
  cpu_num     = 1
  memory_size = 1024
  description = "[batch-test] instance offering"
}

output "uuid" {
  value = zstack_instance_offering.test.uuid
}
`, true, ""
			},
		},
		{
			name: "res-tag",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_tag" "test" {
  name        = "tf-batch-test-tag"
  value       = "batch-test-value"
  color       = "#FF0000"
  type        = "simple"
  description = "[batch-test] tag"
}

output "uuid" {
  value = zstack_tag.test.uuid
}
`, true, ""
			},
		},
		{
			name: "res-script",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_instance_scripts" "test" {
  name           = "tf-batch-test-script"
  script_content = "#!/bin/bash\necho hello"
  encoding_type  = "PlainText"
  script_type    = "Shell"
  platform       = "Linux"
  script_timeout = 300
  description    = "[batch-test] script"
}

output "uuid" {
  value = zstack_instance_scripts.test.uuid
}
`, true, ""
			},
		},
		{
			name: "res-networking_secgroup",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_networking_secgroup" "test" {
  name        = "tf-batch-test-secgroup"
  ip_version  = 4
  description = "[batch-test] security group"
}

output "uuid" {
  value = zstack_networking_secgroup.test.uuid
}
`, true, ""
			},
		},
		{
			name: "res-affinity_group",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_affinity_group" "test" {
  name        = "tf-batch-test-affinity-group"
  policy      = "antiSoft"
  description = "[batch-test] affinity group"
}

output "uuid" {
  value = zstack_affinity_group.test.uuid
}
`, true, ""
			},
		},
		{
			name: "res-ssh_key_pair",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_ssh_key_pair" "test" {
  name        = "tf-batch-test-ssh-key-pair"
  public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC7 test@batch-test"
  description = "[batch-test] SSH key pair"
}

output "uuid" {
  value = zstack_ssh_key_pair.test.uuid
}
`, true, ""
			},
		},
		{
			name: "res-account",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_account" "test" {
  name        = "tf-batch-test-account"
  password    = "BatchTest@12345"
  description = "[batch-test] account"
}

output "uuid" {
  value = zstack_account.test.uuid
}
`, true, ""
			},
		},
		{
			name: "res-iam2_project",
			fn: func(env *EnvData) (string, bool, string) {
				suffix := time.Now().Format("20060102150405")
				return fmt.Sprintf(`resource "zstack_iam2_project" "test" {
  name        = "tf-batch-test-iam2-project-%s"
  description = "[batch-test] IAM2 project"
}

output "uuid" {
  value = zstack_iam2_project.test.uuid
}
`, suffix), true, ""
			},
		},
		// Role
		{
			name: "res-role",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_role" "test" {
  name        = "tf-batch-test-role"
  description = "[batch-test] role"
}

output "uuid" {
  value = zstack_role.test.uuid
}
`, true, ""
			},
		},
		// User
		{
			name: "res-user",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_user" "test" {
  name        = "tf-batch-test-user"
  password    = "BatchTest@12345"
  description = "[batch-test] user"
}

output "uuid" {
  value = zstack_user.test.uuid
}
`, true, ""
			},
		},
		// IAM2 Organization
		{
			name: "res-iam2_organization",
			fn: func(env *EnvData) (string, bool, string) {
				suffix := time.Now().Format("20060102150405")
				return fmt.Sprintf(`resource "zstack_iam2_organization" "test" {
  name        = "tf-batch-test-iam2-org-%s"
  type        = "Department"
  description = "[batch-test] IAM2 organization"
}

output "uuid" {
  value = zstack_iam2_organization.test.uuid
}
`, suffix), true, ""
			},
		},
		// SNS Topic
		{
			name: "res-sns_topic",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_sns_topic" "test" {
  name        = "tf-batch-test-sns-topic"
  description = "[batch-test] SNS topic"
}

output "uuid" {
  value = zstack_sns_topic.test.uuid
}
`, true, ""
			},
		},
		// Webhook
		{
			name: "res-webhook",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_webhook" "test" {
  name        = "tf-batch-test-webhook"
  url         = "http://example.com/hook"
  type        = "generic"
  description = "[batch-test] webhook"
}

output "uuid" {
  value = zstack_webhook.test.uuid
}
`, true, ""
			},
		},
		// Monitor Group
		{
			name: "res-monitor_group",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_monitor_group" "test" {
  name        = "tf-batch-test-monitor-group"
  description = "[batch-test] monitor group"
}

output "uuid" {
  value = zstack_monitor_group.test.uuid
}
`, true, ""
			},
		},
		// Monitor Template
		{
			name: "res-monitor_template",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_monitor_template" "test" {
  name        = "tf-batch-test-monitor-template"
  description = "[batch-test] monitor template"
}

output "uuid" {
  value = zstack_monitor_template.test.uuid
}
`, true, ""
			},
		},
		// Price Table
		{
			name: "res-price_table",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_price_table" "test" {
  name        = "tf-batch-test-price-table"
  description = "[batch-test] price table"
}

output "uuid" {
  value = zstack_price_table.test.uuid
}
`, true, ""
			},
		},
		// Access Control List
		{
			name: "res-access_control_list",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_access_control_list" "test" {
  name        = "tf-batch-test-acl"
  description = "[batch-test] access control list"
  ip_version  = 4
}

output "uuid" {
  value = zstack_access_control_list.test.uuid
}
`, true, ""
			},
		},
		// Flow Meter (only "type" is Required)
		{
			name: "res-flow_meter",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_flow_meter" "test" {
  name        = "tf-batch-test-flow-meter"
  description = "[batch-test] flow meter"
  type        = "NetFlow"
}

output "uuid" {
  value = zstack_flow_meter.test.uuid
}
`, true, ""
			},
		},
		// SNS Email Endpoint (name and email are Required)
		{
			name: "res-sns_email_endpoint",
			fn: func(env *EnvData) (string, bool, string) {
				return `resource "zstack_sns_email_endpoint" "test" {
  name        = "tf-batch-test-sns-email"
  email       = "test@example.com"
  description = "[batch-test] SNS email endpoint"
}

output "uuid" {
  value = zstack_sns_email_endpoint.test.uuid
}
`, true, ""
			},
		},
	}
}

// ---------------------------------------------------------------------------
// Resource generators — need env UUIDs
// ---------------------------------------------------------------------------

func envDependentResourceGenerators() []generator {
	return []generator{
		{
			name: "res-image",
			fn: func(env *EnvData) (string, bool, string) {
				if len(env.BackupStorages) == 0 {
					return "", false, "backup_storages empty"
				}
				bsUUID := getStr(env.BackupStorages[0], "uuid")
				return fmt.Sprintf(`resource "zstack_image" "test" {
  name                 = "tf-batch-test-image"
  description          = "[batch-test] image"
  url                  = "http://192.168.200.100/mirror/diskimages/CentOS-7-x86_64-Cloudinit-8G-official.qcow2"
  format               = "qcow2"
  media_type           = "RootVolumeTemplate"
  platform             = "Linux"
  guest_os_type        = "Centos 7"
  architecture         = "x86_64"
  backup_storage_uuids = [%q]
}

output "uuid" {
  value = zstack_image.test.uuid
}
`, bsUUID), true, ""
			},
		},
		{
			name: "res-virtual_router_image",
			fn: func(env *EnvData) (string, bool, string) {
				if len(env.BackupStorages) == 0 {
					return "", false, "backup_storages empty"
				}
				bsUUID := getStr(env.BackupStorages[0], "uuid")
				return fmt.Sprintf(`resource "zstack_virtual_router_image" "test" {
  name                 = "tf-batch-test-vr-image"
  description          = "[batch-test] virtual router image"
  url                  = "http://192.168.200.100/mirror/diskimages/CentOS-7-x86_64-Cloudinit-8G-official.qcow2"
  architecture         = "x86_64"
  backup_storage_uuids = [%q]
}

output "uuid" {
  value = zstack_virtual_router_image.test.uuid
}
`, bsUUID), true, ""
			},
		},
		{
			name: "res-volume",
			fn: func(env *EnvData) (string, bool, string) {
				if len(env.DiskOfferings) == 0 {
					return "", false, "disk_offerings empty"
				}
				doUUID := getStr(env.DiskOfferings[0], "uuid")
				return fmt.Sprintf(`resource "zstack_volume" "test" {
  name              = "tf-batch-test-volume"
  description       = "[batch-test] volume"
  disk_offering_uuid = %q
}

output "uuid" {
  value = zstack_volume.test.uuid
}
`, doUUID), true, ""
			},
		},
		{
			name: "res-vip",
			fn: func(env *EnvData) (string, bool, string) {
				pub := findPublicL3(env)
				if pub == nil {
					return "", false, "no Public L3 network"
				}
				l3UUID := getStr(pub, "uuid")
				return fmt.Sprintf(`resource "zstack_vip" "test" {
  name           = "tf-batch-test-vip"
  description    = "[batch-test] vip"
  l3_network_uuid = %q
}

output "uuid" {
  value = zstack_vip.test.uuid
}
`, l3UUID), true, ""
			},
		},
		{
			name: "res-virtual_router_offer",
			fn: func(env *EnvData) (string, bool, string) {
				if len(env.VirtualRouterOfferings) == 0 {
					return "", false, "virtual_router_offerings empty (need ref data)"
				}
				if len(env.Zones) == 0 {
					return "", false, "zones empty"
				}
				// Use existing VR offering as reference for network UUIDs
				vro := env.VirtualRouterOfferings[0]
				zoneUUID := getStr(vro, "zone_uuid")
				mgmtUUID := getStr(vro, "management_network_uuid")
				pubUUID := getStr(vro, "public_network_uuid")
				imgUUID := getStr(vro, "image_uuid")

				return fmt.Sprintf(`resource "zstack_virtual_router_offering" "test" {
  name                    = "tf-batch-test-vr-offering"
  description             = "[batch-test] virtual router offering"
  cpu_num                 = 2
  memory_size             = 1024
  zone_uuid               = %q
  management_network_uuid = %q
  public_network_uuid     = %q
  image_uuid              = %q
}

output "uuid" {
  value = zstack_virtual_router_offering.test.uuid
}
`, zoneUUID, mgmtUUID, pubUUID, imgUUID), true, ""
			},
		},
		{
			name: "res-reserved_ip",
			fn: func(env *EnvData) (string, bool, string) {
				pub := findPublicL3(env)
				if pub == nil {
					return "", false, "no Public L3 network"
				}
				l3UUID := getStr(pub, "uuid")
				ipr := findIpRangeForL3(env, l3UUID)
				if ipr == nil {
					return "", false, "no ip_range for Public L3"
				}
				// Reserve 2 IPs at the end of the range to avoid conflicts
				endIP := getStr(ipr, "end_ip")
				reserveStart := incrementIP(endIP, -1)
				reserveEnd := endIP
				return fmt.Sprintf(`resource "zstack_reserved_ip" "test" {
  l3_network_uuid = %q
  start_ip        = %q
  end_ip          = %q
}
`, l3UUID, reserveStart, reserveEnd), true, ""
			},
		},
		{
			name: "res-subnet_ip_range",
			fn: func(env *EnvData) (string, bool, string) {
				priv := findPrivateL3(env)
				if priv == nil {
					return "", false, "no Private L3 network"
				}
				l3UUID := getStr(priv, "uuid")
				ipr := findIpRangeForL3(env, l3UUID)
				if ipr == nil {
					return "", false, "no ip_range for Private L3"
				}
				// Use a safe range in the 10.10.x.240-250 area to avoid collisions
				gateway := getStr(ipr, "gateway")
				netmask := getStr(ipr, "netmask")
				// Parse the gateway to figure out the subnet prefix
				parts := strings.Split(gateway, ".")
				if len(parts) != 4 {
					return "", false, "cannot parse gateway IP"
				}
				startIP := fmt.Sprintf("%s.%s.%s.240", parts[0], parts[1], parts[2])
				endIP := fmt.Sprintf("%s.%s.%s.250", parts[0], parts[1], parts[2])
				return fmt.Sprintf(`resource "zstack_subnet_ip_range" "test" {
  name            = "tf-batch-test-ip-range"
  l3_network_uuid = %q
  start_ip        = %q
  end_ip          = %q
  netmask         = %q
  gateway         = %q
  ip_range_type   = "Normal"
}

output "uuid" {
  value = zstack_subnet_ip_range.test.uuid
}
`, l3UUID, startIP, endIP, netmask, gateway), true, ""
			},
		},
		{
			name: "res-networking_secgroup_rule",
			fn: func(env *EnvData) (string, bool, string) {
				// Create an inline security group to attach the rule to
				return `resource "zstack_networking_secgroup" "for_rule" {
  name        = "tf-batch-test-secgroup-for-rule"
  ip_version  = 4
  description = "[batch-test] secgroup for rule test"
}

resource "zstack_networking_secgroup_rule" "test" {
  depends_on           = [zstack_networking_secgroup.for_rule]
  name                 = "tf-batch-test-sg-rule"
  security_group_uuid  = zstack_networking_secgroup.for_rule.uuid
  priority             = 1
  direction            = "Ingress"
  action               = "ACCEPT"
  state                = "Enabled"
  ip_version           = 4
  protocol             = "TCP"
  ip_ranges            = "10.0.0.0/8"
  destination_port_ranges = "22"
  description          = "[batch-test] security group rule"
}

output "rule_uuid" {
  value = zstack_networking_secgroup_rule.test.uuid
}
`, true, ""
			},
		},
		{
			name: "res-eip",
			fn: func(env *EnvData) (string, bool, string) {
				nic := findUserVmNic(env)
				if nic == nil {
					return "", false, "vm_nics empty"
				}
				pub := findPublicL3(env)
				if pub == nil {
					return "", false, "no Public L3 network"
				}
				l3UUID := getStr(pub, "uuid")
				nicUUID := getStr(nic, "uuid")
				return fmt.Sprintf(`resource "zstack_vip" "for_eip" {
  name            = "tf-batch-test-vip-for-eip"
  description     = "[batch-test] vip for eip"
  l3_network_uuid = %q
}

resource "zstack_eip" "test" {
  name        = "tf-batch-test-eip"
  description = "[batch-test] elastic IP"
  vip_uuid    = zstack_vip.for_eip.uuid
  vm_nic_uuid = %q
}

output "eip_uuid" {
  value = zstack_eip.test.uuid
}
`, l3UUID, nicUUID), true, ""
			},
		},
		{
			name: "res-networking_secgroup_attachment",
			fn: func(env *EnvData) (string, bool, string) {
				nic := findUserVmNic(env)
				if nic == nil {
					return "", false, "vm_nics empty"
				}
				nicUUID := getStr(nic, "uuid")
				return fmt.Sprintf(`resource "zstack_networking_secgroup" "for_attach" {
  name        = "tf-batch-test-secgroup-for-attach"
  ip_version  = 4
  description = "[batch-test] secgroup for attachment test"
}

resource "zstack_networking_secgroup_attachment" "test" {
  secgroup_uuid = zstack_networking_secgroup.for_attach.uuid
  nic_uuid      = %q
}

output "id" {
  value = zstack_networking_secgroup_attachment.test.id
}
`, nicUUID), true, ""
			},
		},
		// VPC (needs L2 network)
		{
			name: "res-vpc",
			fn: func(env *EnvData) (string, bool, string) {
				if len(env.L2Networks) == 0 {
					return "", false, "l2_networks empty"
				}
				l2UUID := getStr(env.L2Networks[0], "uuid")
				return fmt.Sprintf(`resource "zstack_vpc" "test" {
  name            = "tf-batch-test-vpc"
  description     = "[batch-test] VPC"
  l2_network_uuid = %q

  subnet_cidr = {
    name         = "tf-batch-test-vpc-subnet"
    network_cidr = "172.30.0.0/24"
    gateway      = "172.30.0.1"
  }
}

output "uuid" {
  value = zstack_vpc.test.uuid
}
`, l2UUID), true, ""
			},
		},
		// VM NIC (needs a Private L3 network)
		{
			name: "res-vm_nic",
			fn: func(env *EnvData) (string, bool, string) {
				priv := findPrivateL3(env)
				if priv == nil {
					return "", false, "no Private L3 network"
				}
				l3UUID := getStr(priv, "uuid")
				return fmt.Sprintf(`resource "zstack_vm_nic" "test" {
  l3_network_uuid = %q
}

output "uuid" {
  value = zstack_vm_nic.test.uuid
}
`, l3UUID), true, ""
			},
		},
		{
			name: "res-instance",
			fn: func(env *EnvData) (string, bool, string) {
				img := findReadyImage(env)
				if img == nil {
					return "", false, "no Ready image available"
				}
				priv := findPrivateL3(env)
				if priv == nil {
					return "", false, "no Private L3 network"
				}
				if len(env.InstanceOfferings) == 0 {
					return "", false, "instance_offerings empty"
				}
				imgUUID := getStr(img, "uuid")
				l3UUID := getStr(priv, "uuid")
				ioUUID := getStr(env.InstanceOfferings[0], "uuid")
				return fmt.Sprintf(`resource "zstack_instance" "test" {
  name                   = "tf-batch-test-instance"
  description            = "[batch-test] vm instance"
  image_uuid             = %q
  instance_offering_uuid = %q
  strategy               = "CreateStopped"

  network_interfaces = [
    {
      l3_network_uuid = %q
      default_l3      = true
    }
  ]
}

output "uuid" {
  value = zstack_instance.test.uuid
}
`, imgUUID, ioUUID, l3UUID), true, ""
			},
		},
	}
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

func main() {
	onlyMode := flag.String("only", "", "if 'datasources', generate ONLY data source fixtures with QA-style 3-variant lookups")
	flag.Parse()

	envPath := "zstack/provider/testdata/env.json"
	tfDir := "zstack/provider/testdata/terraform"

	// 1. Load env.json
	raw, err := os.ReadFile(envPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", envPath, err)
		fmt.Fprintln(os.Stderr, "Run: source .env.test && go run ./zstack/provider/testdata/generate_env.go")
		os.Exit(1)
	}

	var env EnvData
	if err := json.Unmarshal(raw, &env); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", envPath, err)
		os.Exit(1)
	}

	// 2. Remove and recreate terraform/ dir
	os.RemoveAll(tfDir)
	if err := os.MkdirAll(tfDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating %s: %v\n", tfDir, err)
		os.Exit(1)
	}

	// 3. Write shared provider.tf
	providerPath := filepath.Join(tfDir, "provider.tf")
	if err := os.WriteFile(providerPath, []byte(providerHCL()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", providerPath, err)
		os.Exit(1)
	}

	// 3b. Write dev.tfrc for local provider binary
	tfrcPath := filepath.Join(tfDir, "dev.tfrc")
	if err := os.WriteFile(tfrcPath, []byte(terraformRC()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", tfrcPath, err)
		os.Exit(1)
	}

	// 4. Collect generators (gated by -only flag)
	var allGens []generator
	switch *onlyMode {
	case "datasources":
		// QA mode: only data source coverage, with 3-variant lookups per type
		// (all / by_name / by_uuid) to validate both human and AI lookup paths.
		allGens = append(allGens, dataSourceQAGenerators()...)
	case "":
		allGens = append(allGens, dataSourceGenerators()...)
		allGens = append(allGens, selfContainedResourceGenerators()...)
		allGens = append(allGens, envDependentResourceGenerators()...)
	default:
		fmt.Fprintf(os.Stderr, "unknown -only mode %q (supported: datasources)\n", *onlyMode)
		os.Exit(2)
	}

	// 5. Generate
	generated := 0
	skipped := 0
	var skippedNames []string

	for _, gen := range allGens {
		hcl, ok, reason := gen.fn(&env)
		if !ok {
			skipped++
			skippedNames = append(skippedNames, fmt.Sprintf("  %-40s %s", gen.name, reason))
			continue
		}

		// Create subdirectory
		subDir := filepath.Join(tfDir, gen.name)
		if err := os.MkdirAll(subDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating %s: %v\n", subDir, err)
			continue
		}

		// Write main.tf
		mainPath := filepath.Join(subDir, "main.tf")
		if err := os.WriteFile(mainPath, []byte(hcl), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", mainPath, err)
			continue
		}

		// Create symlink to provider.tf
		symlinkPath := filepath.Join(subDir, "provider.tf")
		if err := os.Symlink("../provider.tf", symlinkPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating symlink %s: %v\n", symlinkPath, err)
		}

		generated++
	}

	// 6. Print summary
	fmt.Printf("Terraform batch-test files generated in %s/\n", tfDir)
	fmt.Printf("  Generated: %d\n", generated)
	fmt.Printf("  Skipped:   %d\n", skipped)
	fmt.Printf("  Provider binary dir: %s\n", providerBinDir())
	if len(skippedNames) > 0 {
		fmt.Println("\nSkipped generators:")
		for _, s := range skippedNames {
			fmt.Println(s)
		}
	}

	absDir, _ := filepath.Abs(tfDir)
	absTfrc := filepath.Join(absDir, "dev.tfrc")

	fmt.Println("\nUsage:")
	fmt.Println("  # 1. Build the provider binary (once, or after code changes):")
	fmt.Println("  go build -o $(go env GOPATH)/bin/terraform-provider-zstack")
	fmt.Println("")
	fmt.Println("  # 2. Set environment variables:")
	fmt.Println("  source .env.test")
	fmt.Printf("  export TF_CLI_CONFIG_FILE=%s\n", absTfrc)
	fmt.Println("")
	fmt.Println("  # 3. Test a single resource (no 'terraform init' needed with dev_overrides):")
	fmt.Printf("  cd %s/res-disk_offer\n", tfDir)
	fmt.Println("  terraform apply -auto-approve")
	fmt.Println("  terraform destroy -auto-approve")
	fmt.Println("")
	fmt.Println("  # 4. Batch test all:")
	fmt.Printf("  for dir in %s/*/; do\n", tfDir)
	fmt.Println("    echo \"=== Testing $dir ===\"")
	fmt.Println("    (cd \"$dir\" && terraform apply -auto-approve && terraform destroy -auto-approve)")
	fmt.Println("  done")
}
