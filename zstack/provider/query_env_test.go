package provider

import (
	"fmt"
	"testing"

	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

func q() *param.QueryParam {
	p := param.NewQueryParam()
	return &p
}

// TestQueryEnvironment is a helper test that queries the ZStack environment
// and prints all resource data needed to update test assertions.
// Run with: source .env.test && go test ./zstack/provider/ -run TestQueryEnvironment -v -count=1
func TestQueryEnvironment(t *testing.T) {
	host := getEnvOrDefault("ZSTACK_HOST", "172.24.227.46")
	port := 8080
	akID := getEnvOrDefault("ZSTACK_ACCESS_KEY_ID", "")
	akSecret := getEnvOrDefault("ZSTACK_ACCESS_KEY_SECRET", "")

	if akID == "" || akSecret == "" {
		t.Fatal("ZSTACK_ACCESS_KEY_ID and ZSTACK_ACCESS_KEY_SECRET must be set")
	}

	cli := client.NewZSClient(client.NewZSConfig(host, port, "zstack").AccessKey(akID, akSecret).ReadOnly(true).Debug(false))

	fmt.Println("\n========================================")
	fmt.Println("  ZStack Environment Resource Summary")
	fmt.Printf("  Host: %s:%d\n", host, port)
	fmt.Println("========================================")

	// Zones
	fmt.Println("\n--- ZONES ---")
	zones, err := cli.QueryZone(q())
	if err != nil {
		t.Logf("QueryZone error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(zones))
		for i, z := range zones {
			fmt.Printf("[%d] name=%q uuid=%q state=%q type=%q\n", i, z.Name, z.UUID, z.State, z.Type)
		}
	}

	// Clusters
	fmt.Println("\n--- CLUSTERS ---")
	clusters, err := cli.QueryCluster(q())
	if err != nil {
		t.Logf("QueryCluster error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(clusters))
		for i, c := range clusters {
			fmt.Printf("[%d] name=%q uuid=%q state=%q hypervisorType=%q zoneUuid=%q\n", i, c.Name, c.UUID, c.State, c.HypervisorType, c.ZoneUuid)
		}
	}

	// Hosts
	fmt.Println("\n--- HOSTS ---")
	hosts, err := cli.QueryHost(q())
	if err != nil {
		t.Logf("QueryHost error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(hosts))
		for i, h := range hosts {
			fmt.Printf("[%d] name=%q uuid=%q state=%q status=%q architecture=%q clusterUuid=%q zoneUuid=%q\n",
				i, h.Name, h.UUID, h.State, h.Status, h.Architecture, h.ClusterUuid, h.ZoneUuid)
		}
	}

	// Images
	fmt.Println("\n--- IMAGES ---")
	images, err := cli.QueryImage(q())
	if err != nil {
		t.Logf("QueryImage error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(images))
		for i, img := range images {
			fmt.Printf("[%d] name=%q uuid=%q state=%q status=%q format=%q platform=%q\n",
				i, img.Name, img.UUID, img.State, img.Status, img.Format, img.Platform)
		}
	}

	// Backup Storages
	fmt.Println("\n--- BACKUP STORAGES ---")
	bss, err := cli.QueryBackupStorage(q())
	if err != nil {
		t.Logf("QueryBackupStorage error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(bss))
		for i, bs := range bss {
			fmt.Printf("[%d] name=%q uuid=%q state=%q status=%q type=%q totalCapacity=%d availableCapacity=%d\n",
				i, bs.Name, bs.UUID, bs.State, bs.Status, bs.Type, bs.TotalCapacity, bs.AvailableCapacity)
		}
	}

	// Primary Storages
	fmt.Println("\n--- PRIMARY STORAGES ---")
	pss, err := cli.QueryPrimaryStorage(q())
	if err != nil {
		t.Logf("QueryPrimaryStorage error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(pss))
		for i, ps := range pss {
			fmt.Printf("[%d] name=%q uuid=%q state=%q status=%q type=%q totalCapacity=%d availableCapacity=%d\n",
				i, ps.Name, ps.UUID, ps.State, ps.Status, ps.Type, ps.TotalCapacity, ps.AvailableCapacity)
		}
	}

	// Instance Offerings
	fmt.Println("\n--- INSTANCE OFFERINGS ---")
	ios, err := cli.QueryInstanceOffering(q())
	if err != nil {
		t.Logf("QueryInstanceOffering error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(ios))
		for i, io := range ios {
			fmt.Printf("[%d] name=%q uuid=%q state=%q cpuNum=%d memorySize=%d\n",
				i, io.Name, io.UUID, io.State, io.CpuNum, io.MemorySize)
		}
	}

	// Disk Offerings
	fmt.Println("\n--- DISK OFFERINGS ---")
	dos, err := cli.QueryDiskOffering(q())
	if err != nil {
		t.Logf("QueryDiskOffering error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(dos))
		for i, do := range dos {
			fmt.Printf("[%d] name=%q uuid=%q state=%q diskSize=%d\n",
				i, do.Name, do.UUID, do.State, do.DiskSize)
		}
	}

	// L2 Networks
	fmt.Println("\n--- L2 NETWORKS ---")
	l2s, err := cli.QueryL2Network(q())
	if err != nil {
		t.Logf("QueryL2Network error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(l2s))
		for i, l2 := range l2s {
			fmt.Printf("[%d] name=%q uuid=%q type=%q zoneUuid=%q\n",
				i, l2.Name, l2.UUID, l2.Type, l2.ZoneUuid)
		}
	}

	// L3 Networks
	fmt.Println("\n--- L3 NETWORKS ---")
	l3s, err := cli.QueryL3Network(q())
	if err != nil {
		t.Logf("QueryL3Network error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(l3s))
		for i, l3 := range l3s {
			fmt.Printf("[%d] name=%q uuid=%q type=%q system=%v category=%q l2NetworkUuid=%q zoneUuid=%q\n",
				i, l3.Name, l3.UUID, l3.Type, l3.System, l3.Category, l3.L2NetworkUuid, l3.ZoneUuid)
		}
	}

	// VM Instances
	fmt.Println("\n--- VM INSTANCES ---")
	vms, err := cli.QueryVmInstance(q())
	if err != nil {
		t.Logf("QueryVmInstance error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(vms))
		for i, vm := range vms {
			fmt.Printf("[%d] name=%q uuid=%q state=%q type=%q platform=%q\n",
				i, vm.Name, vm.UUID, vm.State, vm.Type, vm.Platform)
		}
	}

	// Security Groups
	fmt.Println("\n--- SECURITY GROUPS ---")
	sgs, err := cli.QuerySecurityGroup(q())
	if err != nil {
		t.Logf("QuerySecurityGroup error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(sgs))
		for i, sg := range sgs {
			fmt.Printf("[%d] name=%q uuid=%q state=%q\n", i, sg.Name, sg.UUID, sg.State)
		}
	}

	// Security Group Rules
	fmt.Println("\n--- SECURITY GROUP RULES ---")
	sgrs, err := cli.QuerySecurityGroupRule(q())
	if err != nil {
		t.Logf("QuerySecurityGroupRule error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(sgrs))
		for i, sgr := range sgrs {
			fmt.Printf("[%d] uuid=%q type=%q protocol=%q startPort=%d endPort=%d securityGroupUuid=%q\n",
				i, sgr.UUID, sgr.Type, sgr.Protocol, sgr.StartPort, sgr.EndPort, sgr.SecurityGroupUuid)
		}
	}

	// Virtual Router Offerings
	fmt.Println("\n--- VIRTUAL ROUTER OFFERINGS ---")
	vros, err := cli.QueryVirtualRouterOffering(q())
	if err != nil {
		t.Logf("QueryVirtualRouterOffering error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(vros))
		for i, vro := range vros {
			fmt.Printf("[%d] name=%q uuid=%q state=%q cpuNum=%d memorySize=%d\n",
				i, vro.Name, vro.UUID, vro.State, vro.CpuNum, vro.MemorySize)
		}
	}

	// Virtual Routers
	fmt.Println("\n--- VIRTUAL ROUTERS ---")
	vrs, err := cli.QueryVirtualRouterVm(q())
	if err != nil {
		t.Logf("QueryVirtualRouter error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(vrs))
		for i, vr := range vrs {
			fmt.Printf("[%d] name=%q uuid=%q state=%q\n", i, vr.Name, vr.UUID, vr.State)
		}
	}

	// SDN Controllers
	fmt.Println("\n--- SDN CONTROLLERS ---")
	sdns, err := cli.QuerySdnController(q())
	if err != nil {
		t.Logf("QuerySdnController error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(sdns))
		for i, sdn := range sdns {
			fmt.Printf("[%d] name=%q uuid=%q ip=%q status=%q vendorType=%q\n",
				i, sdn.Name, sdn.UUID, sdn.Ip, sdn.Status, sdn.VendorType)
		}
	}

	// Instance Scripts
	fmt.Println("\n--- INSTANCE SCRIPTS ---")
	scripts, err := cli.QueryGuestVmScript(q())
	if err != nil {
		t.Logf("QueryGuestVmScript error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(scripts))
		for i, s := range scripts {
			fmt.Printf("[%d] name=%q uuid=%q scriptType=%q platform=%q\n",
				i, s.Name, s.UUID, s.ScriptType, s.Platform)
		}
	}

	// Script Executions
	fmt.Println("\n--- SCRIPT EXECUTIONS ---")
	executions, err := cli.QueryGuestVmScriptExecutedRecordDetail(q())
	if err != nil {
		t.Logf("QueryGuestVmScriptExecutedRecordDetail error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(executions))
		for i, e := range executions {
			fmt.Printf("[%d] uuid=%q vmInstanceUuid=%q recordUuid=%q status=%q exitCode=%d\n",
				i, e.UUID, e.VmInstanceUuid, e.RecordUuid, e.Status, e.ExitCode)
		}
	}

	// MN Nodes
	fmt.Println("\n--- MN NODES ---")
	mns, err := cli.QueryManagementNode(q())
	if err != nil {
		t.Logf("QueryManagementNode error: %v", err)
	} else {
		fmt.Printf("Count: %d\n", len(mns))
		for i, mn := range mns {
			fmt.Printf("[%d] uuid=%q hostName=%q\n", i, mn.UUID, mn.HostName)
		}
	}

	fmt.Println("\n========================================")
	fmt.Println("  Query Complete")
	fmt.Println("========================================")
}
