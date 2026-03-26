// Copyright (c) ZStack.io, Inc.
//
// This program queries a real ZStack environment and writes the results
// to testdata/env.json for use by acceptance tests.
//
// Usage:
//   source .env.test && go run ./zstack/provider/testdata/generate_env.go

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

func q() *param.QueryParam {
	p := param.NewQueryParam()
	return &p
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// EnvData is the top-level structure written to env.json
type EnvData struct {
	// Infrastructure
	Zones           []map[string]interface{} `json:"zones"`
	Clusters        []map[string]interface{} `json:"clusters"`
	Hosts           []map[string]interface{} `json:"hosts"`
	PrimaryStorages []map[string]interface{} `json:"primary_storages"`
	BackupStorages  []map[string]interface{} `json:"backup_storages"`

	// Compute
	Images            []map[string]interface{} `json:"images"`
	InstanceOfferings []map[string]interface{} `json:"instance_offerings"`
	DiskOfferings     []map[string]interface{} `json:"disk_offerings"`
	VmInstances       []map[string]interface{} `json:"vm_instances"`
	GpuDevices        []map[string]interface{} `json:"gpu_devices"`
	AutoScalingGroups []map[string]interface{} `json:"auto_scaling_groups"`

	// Storage
	Volumes         []map[string]interface{} `json:"volumes"`
	VolumeSnapshots []map[string]interface{} `json:"volume_snapshots"`

	// Network
	L2Networks            []map[string]interface{} `json:"l2_networks"`
	L2VlanNetworks        []map[string]interface{} `json:"l2_vlan_networks"`
	L3Networks            []map[string]interface{} `json:"l3_networks"`
	IpRanges              []map[string]interface{} `json:"ip_ranges"`
	Vips                  []map[string]interface{} `json:"vips"`
	Eips                  []map[string]interface{} `json:"eips"`
	PortForwardingRules   []map[string]interface{} `json:"port_forwarding_rules"`
	LoadBalancers         []map[string]interface{} `json:"load_balancers"`
	LoadBalancerListeners []map[string]interface{} `json:"load_balancer_listeners"`
	SecurityGroups        []map[string]interface{} `json:"security_groups"`
	SecurityGroupRules    []map[string]interface{} `json:"security_group_rules"`
	VmNics                []map[string]interface{} `json:"vm_nics"`

	// Virtual Router
	VirtualRouterOfferings []map[string]interface{} `json:"virtual_router_offerings"`
	VirtualRouters         []map[string]interface{} `json:"virtual_routers"`

	// System / IAM
	Accounts     []map[string]interface{} `json:"accounts"`
	IAM2Projects []map[string]interface{} `json:"iam2_projects"`

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
}

func main() {
	host := getEnvOrDefault("ZSTACK_HOST", "172.24.227.46")
	portStr := getEnvOrDefault("ZSTACK_PORT", "8080")
	port, _ := strconv.Atoi(portStr)
	akID := getEnvOrDefault("ZSTACK_ACCESS_KEY_ID", "")
	akSecret := getEnvOrDefault("ZSTACK_ACCESS_KEY_SECRET", "")

	if akID == "" || akSecret == "" {
		fmt.Fprintln(os.Stderr, "ZSTACK_ACCESS_KEY_ID and ZSTACK_ACCESS_KEY_SECRET must be set")
		os.Exit(1)
	}

	cli := client.NewZSClient(client.NewZSConfig(host, port, "zstack").AccessKey(akID, akSecret).ReadOnly(true).Debug(false))

	data := EnvData{}

	fmt.Printf("Querying ZStack environment at %s:%d ...\n", host, port)

	// Zones
	if zones, err := cli.QueryZone(q()); err == nil {
		for _, z := range zones {
			data.Zones = append(data.Zones, map[string]interface{}{
				"name":  z.Name,
				"uuid":  z.UUID,
				"state": z.State,
				"type":  z.Type,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryZone error: %v\n", err)
	}

	// Clusters
	if clusters, err := cli.QueryCluster(q()); err == nil {
		for _, c := range clusters {
			data.Clusters = append(data.Clusters, map[string]interface{}{
				"name":            c.Name,
				"uuid":            c.UUID,
				"state":           c.State,
				"hypervisor_type": c.HypervisorType,
				"zone_uuid":      c.ZoneUuid,
				"type":            c.Type,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryCluster error: %v\n", err)
	}

	// Hosts
	if hosts, err := cli.QueryHost(q()); err == nil {
		for _, h := range hosts {
			data.Hosts = append(data.Hosts, map[string]interface{}{
				"name":         h.Name,
				"uuid":         h.UUID,
				"state":        h.State,
				"status":       h.Status,
				"architecture": h.Architecture,
				"cluster_uuid": h.ClusterUuid,
				"zone_uuid":    h.ZoneUuid,
				"type":         "KVM",
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryHost error: %v\n", err)
	}

	// Images
	if images, err := cli.QueryImage(q()); err == nil {
		for _, img := range images {
			data.Images = append(data.Images, map[string]interface{}{
				"name":          img.Name,
				"uuid":          img.UUID,
				"state":         img.State,
				"status":        img.Status,
				"format":        img.Format,
				"platform":      img.Platform,
				"architecture":  img.Architecture,
				"guest_os_type": img.GuestOsType,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryImage error: %v\n", err)
	}

	// Backup Storages
	if bss, err := cli.QueryBackupStorage(q()); err == nil {
		for _, bs := range bss {
			data.BackupStorages = append(data.BackupStorages, map[string]interface{}{
				"name":               bs.Name,
				"uuid":               bs.UUID,
				"state":              bs.State,
				"status":             bs.Status,
				"type":               bs.Type,
				"total_capacity":     bs.TotalCapacity,
				"available_capacity": bs.AvailableCapacity,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryBackupStorage error: %v\n", err)
	}

	// Primary Storages
	if pss, err := cli.QueryPrimaryStorage(q()); err == nil {
		for _, ps := range pss {
			data.PrimaryStorages = append(data.PrimaryStorages, map[string]interface{}{
				"name":               ps.Name,
				"uuid":               ps.UUID,
				"state":              ps.State,
				"status":             ps.Status,
				"type":               ps.Type,
				"total_capacity":     ps.TotalCapacity,
				"available_capacity": ps.AvailableCapacity,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryPrimaryStorage error: %v\n", err)
	}

	// Instance Offerings (only UserVm type, matching Terraform data source filter)
	if ios, err := cli.QueryInstanceOffering(q()); err == nil {
		for _, io := range ios {
			if io.Type != "UserVm" {
				continue
			}
			data.InstanceOfferings = append(data.InstanceOfferings, map[string]interface{}{
				"name":               io.Name,
				"uuid":               io.UUID,
				"state":              io.State,
				"cpu_num":            io.CpuNum,
				"memory_size":        io.MemorySize / (1024 * 1024), // bytes -> MB (matches Terraform)
				"allocator_strategy": io.AllocatorStrategy,
				"type":               io.Type,
				"sort_key":           io.SortKey,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryInstanceOffering error: %v\n", err)
	}

	// Disk Offerings (disk_size converted to GB, matching Terraform)
	if dos, err := cli.QueryDiskOffering(q()); err == nil {
		for _, do := range dos {
			data.DiskOfferings = append(data.DiskOfferings, map[string]interface{}{
				"name":               do.Name,
				"uuid":               do.UUID,
				"state":              do.State,
				"disk_size":          do.DiskSize / (1024 * 1024 * 1024), // bytes -> GB (matches Terraform)
				"allocator_strategy": do.AllocatorStrategy,
				"type":               do.Type,
				"sort_key":           do.SortKey,
				"description":        do.Description,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryDiskOffering error: %v\n", err)
	}

	// L2 Networks
	if l2s, err := cli.QueryL2Network(q()); err == nil {
		for _, l2 := range l2s {
			m := map[string]interface{}{
				"name":      l2.Name,
				"uuid":      l2.UUID,
				"type":      l2.Type,
				"zone_uuid": l2.ZoneUuid,
				"physical_interface": l2.PhysicalInterface,
			}
			if l2.AttachedClusterUuids != nil {
				m["attached_cluster_uuids"] = l2.AttachedClusterUuids
			}
			data.L2Networks = append(data.L2Networks, m)
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryL2Network error: %v\n", err)
	}

	// L3 Networks
	if l3s, err := cli.QueryL3Network(q()); err == nil {
		for _, l3 := range l3s {
			data.L3Networks = append(data.L3Networks, map[string]interface{}{
				"name":           l3.Name,
				"uuid":           l3.UUID,
				"type":           l3.Type,
				"category":       l3.Category,
				"system":         l3.System,
				"l2_network_uuid": l3.L2NetworkUuid,
				"zone_uuid":     l3.ZoneUuid,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryL3Network error: %v\n", err)
	}

	// VM Instances (memory_size converted to MB, matching Terraform)
	if vms, err := cli.QueryVmInstance(q()); err == nil {
		for _, vm := range vms {
			data.VmInstances = append(data.VmInstances, map[string]interface{}{
				"name":            vm.Name,
				"uuid":            vm.UUID,
				"state":           vm.State,
				"type":            vm.Type,
				"platform":        vm.Platform,
				"architecture":    vm.Architecture,
				"cluster_uuid":    vm.ClusterUuid,
				"host_uuid":       vm.HostUuid,
				"image_uuid":      vm.ImageUuid,
				"hypervisor_type": vm.HypervisorType,
				"cpu_num":         vm.CpuNum,
				"memory_size":     vm.MemorySize / (1024 * 1024), // bytes -> MB
				"zone_uuid":       vm.ZoneUuid,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryVmInstance error: %v\n", err)
	}

	// Security Groups
	if sgs, err := cli.QuerySecurityGroup(q()); err == nil {
		for _, sg := range sgs {
			data.SecurityGroups = append(data.SecurityGroups, map[string]interface{}{
				"name":  sg.Name,
				"uuid":  sg.UUID,
				"state": sg.State,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QuerySecurityGroup error: %v\n", err)
	}

	// Security Group Rules
	if sgrs, err := cli.QuerySecurityGroupRule(q()); err == nil {
		for _, sgr := range sgrs {
			data.SecurityGroupRules = append(data.SecurityGroupRules, map[string]interface{}{
				"uuid":                sgr.UUID,
				"type":                sgr.Type,
				"protocol":            sgr.Protocol,
				"start_port":          sgr.StartPort,
				"end_port":            sgr.EndPort,
				"security_group_uuid": sgr.SecurityGroupUuid,
				"state":               sgr.State,
				"action":              sgr.Action,
				"src_ip_range":        sgr.SrcIpRange,
				"dst_port_range":      sgr.DstPortRange,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QuerySecurityGroupRule error: %v\n", err)
	}

	// Virtual Router Offerings (memory_size converted to MB, matching Terraform)
	if vros, err := cli.QueryVirtualRouterOffering(q()); err == nil {
		for _, vro := range vros {
			data.VirtualRouterOfferings = append(data.VirtualRouterOfferings, map[string]interface{}{
				"name":                    vro.Name,
				"uuid":                    vro.UUID,
				"state":                   vro.State,
				"cpu_num":                 vro.CpuNum,
				"memory_size":             vro.MemorySize / (1024 * 1024), // bytes -> MB
				"allocator_strategy":      vro.AllocatorStrategy,
				"type":                    vro.Type,
				"sort_key":                vro.SortKey,
				"image_uuid":              vro.ImageUuid,
				"management_network_uuid": vro.ManagementNetworkUuid,
				"public_network_uuid":     vro.PublicNetworkUuid,
				"zone_uuid":               vro.ZoneUuid,
				"is_default":              vro.IsDefault,
				"reserved_memory_size":    vro.ReservedMemorySize,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryVirtualRouterOffering error: %v\n", err)
	}

	// Virtual Routers
	if vrs, err := cli.QueryVirtualRouterVm(q()); err == nil {
		for _, vr := range vrs {
			vrMap := map[string]interface{}{
				"name":                    vr.Name,
				"uuid":                    vr.UUID,
				"state":                   vr.State,
				"status":                  vr.Status,
				"hypervisor_type":         vr.HypervisorType,
				"appliance_vm_type":       vr.ApplianceVmType,
				"agent_port":              vr.AgentPort,
				"type":                    vr.Type,
				"ha_status":               vr.HaStatus,
				"zone_uuid":               vr.ZoneUuid,
				"cluster_uuid":            vr.ClusterUuid,
				"management_network_uuid": vr.ManagementNetworkUuid,
				"image_uuid":              vr.ImageUuid,
				"host_uuid":               vr.HostUuid,
				"instance_offering_uuid":  vr.InstanceOfferingUuid,
				"platform":                vr.Platform,
				"architecture":            vr.Architecture,
				"cpu_num":                 vr.CpuNum,
				"memory_size":             vr.MemorySize / (1024 * 1024), // bytes -> MB
			}

			var nics []map[string]interface{}
			for _, nic := range vr.VmNics {
				nics = append(nics, map[string]interface{}{
					"ip":      nic.Ip,
					"mac":     nic.Mac,
					"netmask": nic.Netmask,
					"gateway": nic.Gateway,
				})
			}
			vrMap["vm_nics"] = nics

			data.VirtualRouters = append(data.VirtualRouters, vrMap)
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryVirtualRouter error: %v\n", err)
	}

	// SDN Controllers
	if sdns, err := cli.QuerySdnController(q()); err == nil {
		for _, sdn := range sdns {
			data.SdnControllers = append(data.SdnControllers, map[string]interface{}{
				"name":        sdn.Name,
				"uuid":        sdn.UUID,
				"ip":          sdn.Ip,
				"status":      sdn.Status,
				"vendor_type": sdn.VendorType,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QuerySdnController error: %v\n", err)
	}

	// Instance Scripts
	if scripts, err := cli.QueryGuestVmScript(q()); err == nil {
		for _, s := range scripts {
			data.InstanceScripts = append(data.InstanceScripts, map[string]interface{}{
				"name":           s.Name,
				"uuid":           s.UUID,
				"script_type":    s.ScriptType,
				"platform":       s.Platform,
				"script_timeout": s.ScriptTimeout,
				"render_params":  s.RenderParams,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryGuestVmScript error: %v\n", err)
	}

	// Script Executions
	if executions, err := cli.QueryGuestVmScriptExecutedRecordDetail(q()); err == nil {
		for _, e := range executions {
			data.ScriptExecutions = append(data.ScriptExecutions, map[string]interface{}{
				"uuid":             e.UUID,
				"vm_instance_uuid": e.VmInstanceUuid,
				"record_uuid":     e.RecordUuid,
				"status":          e.Status,
				"exit_code":       e.ExitCode,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryGuestVmScriptExecutedRecordDetail error: %v\n", err)
	}

	// IP Ranges
	if ipRanges, err := cli.QueryIpRange(q()); err == nil {
		for _, ipr := range ipRanges {
			data.IpRanges = append(data.IpRanges, map[string]interface{}{
				"name":            ipr.Name,
				"uuid":            ipr.UUID,
				"l3_network_uuid": ipr.L3NetworkUuid,
				"start_ip":        ipr.StartIp,
				"end_ip":          ipr.EndIp,
				"netmask":         ipr.Netmask,
				"gateway":         ipr.Gateway,
				"ip_version":      ipr.IpVersion,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryIpRange error: %v\n", err)
	}

	// VM NICs
	if vmNics, err := cli.QueryVmNic(q()); err == nil {
		for _, nic := range vmNics {
			data.VmNics = append(data.VmNics, map[string]interface{}{
				"uuid":             nic.UUID,
				"vm_instance_uuid": nic.VmInstanceUuid,
				"l3_network_uuid":  nic.L3NetworkUuid,
				"ip":               nic.Ip,
				"mac":              nic.Mac,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryVmNic error: %v\n", err)
	}

	// MN Nodes
	if mns, err := cli.QueryManagementNode(q()); err == nil {
		for _, mn := range mns {
			data.MnNodes = append(data.MnNodes, map[string]interface{}{
				"uuid":      mn.UUID,
				"host_name": mn.HostName,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryManagementNode error: %v\n", err)
	}

	// Accounts
	if accounts, err := cli.QueryAccount(q()); err == nil {
		for _, a := range accounts {
			data.Accounts = append(data.Accounts, map[string]interface{}{
				"name":        a.Name,
				"uuid":        a.UUID,
				"type":        a.Type,
				"description": a.Description,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryAccount error: %v\n", err)
	}

	// IAM2 Projects
	if projects, err := cli.QueryIAM2Project(q()); err == nil {
		for _, p := range projects {
			data.IAM2Projects = append(data.IAM2Projects, map[string]interface{}{
				"name":        p.Name,
				"uuid":        p.UUID,
				"state":       p.State,
				"description": p.Description,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryIAM2Project error: %v\n", err)
	}

	// Affinity Groups
	if ags, err := cli.QueryAffinityGroup(q()); err == nil {
		for _, ag := range ags {
			data.AffinityGroups = append(data.AffinityGroups, map[string]interface{}{
				"name":        ag.Name,
				"uuid":        ag.UUID,
				"policy":      ag.Policy,
				"type":        ag.Type,
				"zone_uuid":   ag.ZoneUuid,
				"state":       ag.State,
				"description": ag.Description,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryAffinityGroup error: %v\n", err)
	}

	// SSH Key Pairs
	if skps, err := cli.QuerySshKeyPair(q()); err == nil {
		for _, skp := range skps {
			data.SshKeyPairs = append(data.SshKeyPairs, map[string]interface{}{
				"name":        skp.Name,
				"uuid":        skp.UUID,
				"public_key":  skp.PublicKey,
				"description": skp.Description,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QuerySshKeyPair error: %v\n", err)
	}

	// Volumes
	if vols, err := cli.QueryVolume(q()); err == nil {
		for _, v := range vols {
			data.Volumes = append(data.Volumes, map[string]interface{}{
				"name":                 v.Name,
				"uuid":                 v.UUID,
				"type":                 v.Type,
				"state":                v.State,
				"status":               v.Status,
				"size":                 v.Size,
				"actual_size":          v.ActualSize,
				"format":               v.Format,
				"vm_instance_uuid":     v.VmInstanceUuid,
				"primary_storage_uuid": v.PrimaryStorageUuid,
				"disk_offering_uuid":   v.DiskOfferingUuid,
				"is_shareable":         v.IsShareable,
				"device_id":            v.DeviceId,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryVolume error: %v\n", err)
	}

	// Volume Snapshots
	if snaps, err := cli.QueryVolumeSnapshot(q()); err == nil {
		for _, s := range snaps {
			data.VolumeSnapshots = append(data.VolumeSnapshots, map[string]interface{}{
				"name":                 s.Name,
				"uuid":                 s.UUID,
				"type":                 s.Type,
				"state":                s.State,
				"status":               s.Status,
				"size":                 s.Size,
				"volume_uuid":          s.VolumeUuid,
				"volume_type":          s.VolumeType,
				"tree_uuid":            s.TreeUuid,
				"parent_uuid":          s.ParentUuid,
				"primary_storage_uuid": s.PrimaryStorageUuid,
				"format":               s.Format,
				"latest":               s.Latest,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryVolumeSnapshot error: %v\n", err)
	}

	// VIPs
	if vips, err := cli.QueryVip(q()); err == nil {
		for _, v := range vips {
			data.Vips = append(data.Vips, map[string]interface{}{
				"name":            v.Name,
				"uuid":            v.UUID,
				"ip":              v.Ip,
				"state":           v.State,
				"gateway":         v.Gateway,
				"netmask":         v.Netmask,
				"l3_network_uuid": v.L3NetworkUuid,
				"use_for":         v.UseFor,
				"system":          v.System,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryVip error: %v\n", err)
	}

	// EIPs
	if eips, err := cli.QueryEip(q()); err == nil {
		for _, e := range eips {
			data.Eips = append(data.Eips, map[string]interface{}{
				"name":       e.Name,
				"uuid":       e.UUID,
				"state":      e.State,
				"vip_uuid":   e.VipUuid,
				"vip_ip":     e.VipIp,
				"guest_ip":   e.GuestIp,
				"vm_nic_uuid": e.VmNicUuid,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryEip error: %v\n", err)
	}

	// L2 VLAN Networks (with vlan_id)
	if l2vlans, err := cli.QueryL2VlanNetwork(q()); err == nil {
		for _, l2v := range l2vlans {
			m := map[string]interface{}{
				"name":               l2v.Name,
				"uuid":               l2v.UUID,
				"vlan":               l2v.Vlan,
				"zone_uuid":          l2v.ZoneUuid,
				"physical_interface": l2v.PhysicalInterface,
				"type":               l2v.Type,
			}
			if l2v.AttachedClusterUuids != nil {
				m["attached_cluster_uuids"] = l2v.AttachedClusterUuids
			}
			data.L2VlanNetworks = append(data.L2VlanNetworks, m)
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryL2VlanNetwork error: %v\n", err)
	}

	// Port Forwarding Rules
	if pfrs, err := cli.QueryPortForwardingRule(q()); err == nil {
		for _, pf := range pfrs {
			data.PortForwardingRules = append(data.PortForwardingRules, map[string]interface{}{
				"name":               pf.Name,
				"uuid":               pf.UUID,
				"state":              pf.State,
				"vip_uuid":           pf.VipUuid,
				"vip_ip":             pf.VipIp,
				"vip_port_start":     pf.VipPortStart,
				"vip_port_end":       pf.VipPortEnd,
				"private_port_start": pf.PrivatePortStart,
				"private_port_end":   pf.PrivatePortEnd,
				"vm_nic_uuid":        pf.VmNicUuid,
				"protocol_type":      pf.ProtocolType,
				"guest_ip":           pf.GuestIp,
				"allowed_cidr":       pf.AllowedCidr,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryPortForwardingRule error: %v\n", err)
	}

	// Load Balancers
	if lbs, err := cli.QueryLoadBalancer(q()); err == nil {
		for _, lb := range lbs {
			data.LoadBalancers = append(data.LoadBalancers, map[string]interface{}{
				"name":              lb.Name,
				"uuid":              lb.UUID,
				"state":             lb.State,
				"type":              lb.Type,
				"vip_uuid":          lb.VipUuid,
				"server_group_uuid": lb.ServerGroupUuid,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryLoadBalancer error: %v\n", err)
	}

	// Load Balancer Listeners
	if lbls, err := cli.QueryLoadBalancerListener(q()); err == nil {
		for _, lbl := range lbls {
			data.LoadBalancerListeners = append(data.LoadBalancerListeners, map[string]interface{}{
				"name":                 lbl.Name,
				"uuid":                 lbl.UUID,
				"load_balancer_uuid":   lbl.LoadBalancerUuid,
				"protocol":             lbl.Protocol,
				"load_balancer_port":   lbl.LoadBalancerPort,
				"instance_port":        lbl.InstancePort,
				"security_policy_type": lbl.SecurityPolicyType,
				"server_group_uuid":    lbl.ServerGroupUuid,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryLoadBalancerListener error: %v\n", err)
	}

	// GPU Devices
	if gpus, err := cli.QueryGpuDevice(q()); err == nil {
		for _, g := range gpus {
			data.GpuDevices = append(data.GpuDevices, map[string]interface{}{
				"name":               g.Name,
				"uuid":               g.UUID,
				"gpu_type":           g.GpuType,
				"gpu_status":         g.GpuStatus,
				"allocate_status":    g.AllocateStatus,
				"host_uuid":          g.HostUuid,
				"vm_instance_uuid":   g.VmInstanceUuid,
				"vendor":             g.Vendor,
				"device":             g.Device,
				"memory":             g.Memory,
				"pci_device_address": g.PciDeviceAddress,
				"state":              g.State,
				"status":             g.Status,
				"type":               g.Type,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryGpuDevice error: %v\n", err)
	}

	// Auto Scaling Groups
	if asgs, err := cli.QueryAutoScalingGroup(q()); err == nil {
		for _, asg := range asgs {
			data.AutoScalingGroups = append(data.AutoScalingGroups, map[string]interface{}{
				"name":                  asg.Name,
				"uuid":                  asg.UUID,
				"state":                 asg.State,
				"scaling_resource_type": asg.ScalingResourceType,
				"default_cooldown":      asg.DefaultCooldown,
				"min_resource_size":     asg.MinResourceSize,
				"max_resource_size":     asg.MaxResourceSize,
				"removal_policy":        asg.RemovalPolicy,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryAutoScalingGroup error: %v\n", err)
	}

	// User Tags
	if tags, err := cli.QueryUserTag(q()); err == nil {
		for _, t := range tags {
			data.UserTags = append(data.UserTags, map[string]interface{}{
				"uuid":          t.UUID,
				"tag":           t.Tag,
				"type":          t.Type,
				"resource_uuid": t.ResourceUuid,
				"resource_type": t.ResourceType,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QueryUserTag error: %v\n", err)
	}

	// System Tags
	if tags, err := cli.QuerySystemTag(q()); err == nil {
		for _, t := range tags {
			data.SystemTags = append(data.SystemTags, map[string]interface{}{
				"uuid":          t.UUID,
				"tag":           t.Tag,
				"type":          t.Type,
				"inherent":      t.Inherent,
				"resource_uuid": t.ResourceUuid,
				"resource_type": t.ResourceType,
			})
		}
	} else {
		fmt.Fprintf(os.Stderr, "QuerySystemTag error: %v\n", err)
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON marshal error: %v\n", err)
		os.Exit(1)
	}

	outPath := "zstack/provider/testdata/env.json"
	if err := os.WriteFile(outPath, jsonBytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Write file error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Environment data written to %s\n", outPath)
	fmt.Printf("  Zones: %d, Clusters: %d, Hosts: %d, Images: %d\n",
		len(data.Zones), len(data.Clusters), len(data.Hosts), len(data.Images))
	fmt.Printf("  BackupStorages: %d, PrimaryStorages: %d\n",
		len(data.BackupStorages), len(data.PrimaryStorages))
	fmt.Printf("  InstanceOfferings: %d, DiskOfferings: %d\n",
		len(data.InstanceOfferings), len(data.DiskOfferings))
	fmt.Printf("  L2Networks: %d, L3Networks: %d, VmInstances: %d\n",
		len(data.L2Networks), len(data.L3Networks), len(data.VmInstances))
	fmt.Printf("  SecurityGroups: %d, SecurityGroupRules: %d\n",
		len(data.SecurityGroups), len(data.SecurityGroupRules))
	fmt.Printf("  VirtualRouterOfferings: %d, VirtualRouters: %d\n",
		len(data.VirtualRouterOfferings), len(data.VirtualRouters))
	fmt.Printf("  SdnControllers: %d, InstanceScripts: %d, MnNodes: %d\n",
		len(data.SdnControllers), len(data.InstanceScripts), len(data.MnNodes))
	fmt.Printf("  IpRanges: %d, VmNics: %d\n",
		len(data.IpRanges), len(data.VmNics))
	fmt.Printf("  Accounts: %d, IAM2Projects: %d, AffinityGroups: %d, SshKeyPairs: %d\n",
		len(data.Accounts), len(data.IAM2Projects), len(data.AffinityGroups), len(data.SshKeyPairs))
	fmt.Printf("  Volumes: %d, VolumeSnapshots: %d\n",
		len(data.Volumes), len(data.VolumeSnapshots))
	fmt.Printf("  L2VlanNetworks: %d, Vips: %d, Eips: %d, PortForwardingRules: %d\n",
		len(data.L2VlanNetworks), len(data.Vips), len(data.Eips), len(data.PortForwardingRules))
	fmt.Printf("  LoadBalancers: %d, LoadBalancerListeners: %d\n",
		len(data.LoadBalancers), len(data.LoadBalancerListeners))
	fmt.Printf("  GpuDevices: %d, AutoScalingGroups: %d\n",
		len(data.GpuDevices), len(data.AutoScalingGroups))
	fmt.Printf("  UserTags: %d, SystemTags: %d\n",
		len(data.UserTags), len(data.SystemTags))
}
