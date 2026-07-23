// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

func TestVxlanPreflightDataSourceSchema(t *testing.T) {
	var dataSource vxlanPreflightDataSource
	resp := &datasource.SchemaResponse{}
	dataSource.Schema(context.Background(), datasource.SchemaRequest{}, resp)

	for _, name := range []string{
		"zone_uuid",
		"cluster_uuid",
		"pool_uuid",
		"physical_interface",
		"vtep_cidr",
		"start_vni",
		"end_vni",
	} {
		attribute, ok := resp.Schema.Attributes[name]
		if !ok || !attribute.IsRequired() {
			t.Errorf("attribute %q should be required", name)
		}
	}
	exclude, ok := resp.Schema.Attributes["exclude_vni_range_uuid"]
	if !ok || !exclude.IsOptional() {
		t.Error("exclude_vni_range_uuid should be optional")
	}
	for _, name := range []string{"ready", "hosts", "connectivity", "evidence_digest"} {
		attribute, ok := resp.Schema.Attributes[name]
		if !ok || !attribute.IsComputed() {
			t.Errorf("attribute %q should be computed", name)
		}
	}
}

func TestBuildVxlanPreflightHostEvidence(t *testing.T) {
	hosts := []view.HostInventoryView{
		{
			BaseInfoView: view.BaseInfoView{UUID: "host-2", Name: "second"},
			ClusterUuid:  "cluster-uuid",
		},
		{
			BaseInfoView: view.BaseInfoView{UUID: "host-1", Name: "first"},
			ClusterUuid:  "cluster-uuid",
		},
	}
	facts := &view.GetClusterHostNetworkFactsView{
		Success: true,
		Nics: []view.HostNetworkInterfaceInventoryView{{
			HostUuid:      "host-1",
			InterfaceName: "ens3",
			IpAddresses:   []string{"172.26.0.11/16", "192.168.1.11/24"},
		}},
		Bondings: []view.HostNetworkBondingInventoryView{{
			HostUuid:    "host-2",
			BondingName: "ens3",
			IpAddresses: []string{"172.26.0.12"},
		}},
	}

	evidence, ips, err := buildVxlanPreflightHostEvidence(
		hosts,
		facts,
		"ens3",
		netip.MustParsePrefix("172.26.0.0/16"),
	)
	if err != nil {
		t.Fatalf("build host evidence: %v", err)
	}
	if len(evidence) != 2 || evidence[0].HostUuid != "host-1" || evidence[0].VtepIp != "172.26.0.11" {
		t.Fatalf("unexpected canonical host evidence: %#v", evidence)
	}
	if len(ips) != 2 || ips[0] != "172.26.0.11" || ips[1] != "172.26.0.12" {
		t.Fatalf("unexpected canonical VTEP IPs: %#v", ips)
	}
}

func TestBuildVxlanPreflightHostEvidenceRejectsInvalidCandidates(t *testing.T) {
	hosts := []view.HostInventoryView{
		{BaseInfoView: view.BaseInfoView{UUID: "host-1"}, ClusterUuid: "cluster-uuid"},
		{BaseInfoView: view.BaseInfoView{UUID: "host-2"}, ClusterUuid: "cluster-uuid"},
	}
	prefix := netip.MustParsePrefix("172.26.0.0/16")

	testCases := map[string]struct {
		facts       *view.GetClusterHostNetworkFactsView
		errorSubstr string
	}{
		"missing interface": {
			facts: &view.GetClusterHostNetworkFactsView{
				Nics: []view.HostNetworkInterfaceInventoryView{
					{HostUuid: "host-1", InterfaceName: "ens3", IpAddresses: []string{"172.26.0.11"}},
				},
			},
			errorSubstr: "host-2 does not expose interface",
		},
		"no address in CIDR": {
			facts: &view.GetClusterHostNetworkFactsView{
				Nics: []view.HostNetworkInterfaceInventoryView{
					{HostUuid: "host-1", InterfaceName: "ens3", IpAddresses: []string{"192.168.1.11"}},
					{HostUuid: "host-2", InterfaceName: "ens3", IpAddresses: []string{"172.26.0.12"}},
				},
			},
			errorSubstr: "has no IPv4 address",
		},
		"ambiguous addresses": {
			facts: &view.GetClusterHostNetworkFactsView{
				Nics: []view.HostNetworkInterfaceInventoryView{
					{HostUuid: "host-1", InterfaceName: "ens3", IpAddresses: []string{"172.26.0.11", "172.26.0.21"}},
					{HostUuid: "host-2", InterfaceName: "ens3", IpAddresses: []string{"172.26.0.12"}},
				},
			},
			errorSubstr: "multiple IPv4 addresses",
		},
		"duplicate address": {
			facts: &view.GetClusterHostNetworkFactsView{
				Nics: []view.HostNetworkInterfaceInventoryView{
					{HostUuid: "host-1", InterfaceName: "ens3", IpAddresses: []string{"172.26.0.11"}},
					{HostUuid: "host-2", InterfaceName: "ens3", IpAddresses: []string{"172.26.0.11"}},
				},
			},
			errorSubstr: "is duplicated",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			_, _, err := buildVxlanPreflightHostEvidence(hosts, testCase.facts, "ens3", prefix)
			if err == nil || !strings.Contains(err.Error(), testCase.errorSubstr) {
				t.Fatalf("expected error containing %q, got %v", testCase.errorSubstr, err)
			}
		})
	}
}

func TestValidateVxlanPreflightConnectivity(t *testing.T) {
	ips := []string{"172.26.0.11", "172.26.0.12"}
	connected := &view.CheckNetworkReachableView{
		Success: true,
		Results: []view.NetworkReachablePairView{
			{SourceHostname: "172.26.0.12", TargetHostname: "172.26.0.11", Status: "Connected"},
			{SourceHostname: "172.26.0.11", TargetHostname: "172.26.0.12", Status: "Connected"},
		},
	}

	evidence, err := validateVxlanPreflightConnectivity(ips, connected)
	if err != nil {
		t.Fatalf("validate connectivity: %v", err)
	}
	if len(evidence) != 2 || evidence[0].SourceIp != "172.26.0.11" {
		t.Fatalf("unexpected connectivity evidence: %#v", evidence)
	}

	testCases := map[string]struct {
		result      *view.CheckNetworkReachableView
		errorSubstr string
	}{
		"success false": {
			result:      &view.CheckNetworkReachableView{},
			errorSubstr: "success=false",
		},
		"missing pair": {
			result: &view.CheckNetworkReachableView{
				Success: true,
				Results: connected.Results[:1],
			},
			errorSubstr: "missing pair",
		},
		"not connected": {
			result: &view.CheckNetworkReachableView{
				Success: true,
				Results: []view.NetworkReachablePairView{
					{SourceHostname: "172.26.0.11", TargetHostname: "172.26.0.12", Status: "Disconnected"},
					{SourceHostname: "172.26.0.12", TargetHostname: "172.26.0.11", Status: "Connected"},
				},
			},
			errorSubstr: "expected Connected",
		},
		"duplicate pair": {
			result: &view.CheckNetworkReachableView{
				Success: true,
				Results: append(append([]view.NetworkReachablePairView{}, connected.Results...), connected.Results[0]),
			},
			errorSubstr: "duplicate pair",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := validateVxlanPreflightConnectivity(ips, testCase.result)
			if err == nil || !strings.Contains(err.Error(), testCase.errorSubstr) {
				t.Fatalf("expected error containing %q, got %v", testCase.errorSubstr, err)
			}
		})
	}
}

func TestVxlanPreflightEvidenceDigestDeterministic(t *testing.T) {
	evidence := vxlanPreflightEvidence{
		Version:             vxlanPreflightEvidenceVersion,
		ZoneUuid:            "zone-uuid",
		ClusterUuid:         "cluster-uuid",
		PoolUuid:            "pool-uuid",
		PhysicalInterface:   "ens3",
		VtepCidr:            "172.26.0.0/16",
		StartVni:            100,
		EndVni:              200,
		ExcludeVniRangeUuid: "target-range-uuid",
		Hosts: []vxlanPreflightHostEvidence{
			{HostUuid: "host-1", HostName: "first", VtepIp: "172.26.0.11"},
		},
	}

	first, err := vxlanPreflightEvidenceDigest(evidence)
	if err != nil {
		t.Fatalf("first digest: %v", err)
	}
	second, err := vxlanPreflightEvidenceDigest(evidence)
	if err != nil {
		t.Fatalf("second digest: %v", err)
	}
	if first != second || len(first) != 64 {
		t.Fatalf("unexpected evidence digests: %q %q", first, second)
	}
	evidence.ExcludeVniRangeUuid = ""
	withoutExclusion, err := vxlanPreflightEvidenceDigest(evidence)
	if err != nil {
		t.Fatalf("digest without exclusion: %v", err)
	}
	if first == withoutExclusion {
		t.Fatal("exclude_vni_range_uuid must affect the evidence digest")
	}
}

func TestVxlanPreflightDataSourceRead(t *testing.T) {
	server := newVxlanPreflightServer(t, vxlanPreflightServerOptions{})
	defer server.Close()

	tfresource.UnitTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config: vxlanPreflightTestConfig(t, server),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_vxlan_preflight.test", tfjsonpath.New("ready"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.zstack_vxlan_preflight.test", tfjsonpath.New("hosts"), knownvalue.ListSizeExact(2)),
					statecheck.ExpectKnownValue(
						"data.zstack_vxlan_preflight.test",
						tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("host_uuid"),
						knownvalue.StringExact("host-1"),
					),
					statecheck.ExpectKnownValue("data.zstack_vxlan_preflight.test", tfjsonpath.New("connectivity"), knownvalue.ListSizeExact(2)),
					statecheck.ExpectKnownValue("data.zstack_vxlan_preflight.test", tfjsonpath.New("evidence_digest"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestVxlanPreflightDataSourceBlocksAPIFailure(t *testing.T) {
	server := newVxlanPreflightServer(t, vxlanPreflightServerOptions{failReachability: true})
	defer server.Close()

	tfresource.UnitTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config:      vxlanPreflightTestConfig(t, server),
				ExpectError: regexp.MustCompile(`(?s)VXLAN preflight reachability query failed.*mock.*reachability.*failure`),
			},
		},
	})
}

func TestVxlanPreflightDataSourceRejectsInvalidPool(t *testing.T) {
	testCases := map[string]struct {
		options vxlanPreflightServerOptions
		error   *regexp.Regexp
	}{
		"not found": {
			options: vxlanPreflightServerOptions{poolStatus: http.StatusNotFound},
			error:   regexp.MustCompile(`(?s)VXLAN preflight pool query failed.*pool-uuid`),
		},
		"identity mismatch": {
			options: vxlanPreflightServerOptions{poolUuid: "different-pool-uuid"},
			error:   regexp.MustCompile(`(?s)VXLAN preflight pool identity mismatch.*different-pool-uuid`),
		},
		"zone mismatch": {
			options: vxlanPreflightServerOptions{poolZoneUuid: "different-zone-uuid"},
			error:   regexp.MustCompile(`(?s)VXLAN pool is not in the requested zone.*different-zone-uuid`),
		},
		"physical interface mismatch": {
			options: vxlanPreflightServerOptions{poolPhysicalInterface: "bond0"},
			error:   regexp.MustCompile(`(?s)VXLAN pool uses a different physical interface.*bond0.*ens3`),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			server := newVxlanPreflightServer(t, testCase.options)
			defer server.Close()

			tfresource.UnitTest(t, tfresource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []tfresource.TestStep{
					{
						Config:      vxlanPreflightTestConfig(t, server),
						ExpectError: testCase.error,
					},
				},
			})
		})
	}
}

func TestVxlanPreflightDataSourceExcludesTargetRangeAndRemainsNoop(t *testing.T) {
	server := newVxlanPreflightServer(t, vxlanPreflightServerOptions{
		ranges: []view.VniRangeInventoryView{{
			BaseInfoView:  view.BaseInfoView{UUID: "target-range-uuid", Name: "target"},
			StartVni:      100,
			EndVni:        200,
			L2NetworkUuid: "pool-uuid",
		}},
	})
	defer server.Close()

	config := vxlanPreflightTestConfigWithExclude(t, server, "target-range-uuid")
	tfresource.UnitTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config: config,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_vxlan_preflight.test", tfjsonpath.New("ready"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(
						"data.zstack_vxlan_preflight.test",
						tfjsonpath.New("exclude_vni_range_uuid"),
						knownvalue.StringExact("target-range-uuid"),
					),
				},
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

func TestVxlanPreflightDataSourceDoesNotExcludeOtherRange(t *testing.T) {
	server := newVxlanPreflightServer(t, vxlanPreflightServerOptions{
		ranges: []view.VniRangeInventoryView{{
			BaseInfoView:  view.BaseInfoView{UUID: "other-range-uuid", Name: "other"},
			StartVni:      100,
			EndVni:        200,
			L2NetworkUuid: "pool-uuid",
		}},
	})
	defer server.Close()

	tfresource.UnitTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config:      vxlanPreflightTestConfigWithExclude(t, server, "target-range-uuid"),
				ExpectError: regexp.MustCompile(`(?s)VXLAN preflight found an overlapping VNI range.*other-range-uuid`),
			},
		},
	})
}

type vxlanPreflightServerOptions struct {
	failReachability      bool
	poolStatus            int
	poolUuid              string
	poolZoneUuid          string
	poolPhysicalInterface string
	ranges                []view.VniRangeInventoryView
}

func newVxlanPreflightServer(t *testing.T, options vxlanPreflightServerOptions) *httptest.Server {
	t.Helper()
	poolUuid := options.poolUuid
	if poolUuid == "" {
		poolUuid = "pool-uuid"
	}
	poolZoneUuid := options.poolZoneUuid
	if poolZoneUuid == "" {
		poolZoneUuid = "zone-uuid"
	}
	poolPhysicalInterface := options.poolPhysicalInterface
	if poolPhysicalInterface == "" {
		poolPhysicalInterface = "ens3"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/zstack/v1/l2-networks/vxlan-pool/pool-uuid", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Fatalf("unexpected VXLAN pool method: %s", req.Method)
		}
		if options.poolStatus != 0 {
			http.Error(w, "mock VXLAN pool failure", options.poolStatus)
			return
		}
		writeVxlanPreflightResponse(t, w, map[string]any{
			"inventories": []view.L2VxlanNetworkPoolInventoryView{{
				BaseInfoView:      view.BaseInfoView{UUID: poolUuid, Name: "pool"},
				ZoneUuid:          poolZoneUuid,
				PhysicalInterface: poolPhysicalInterface,
			}},
		})
	})
	mux.HandleFunc("/zstack/v1/clusters/cluster-uuid", func(w http.ResponseWriter, req *http.Request) {
		writeVxlanPreflightResponse(t, w, map[string]any{
			"inventories": []view.ClusterInventoryView{{
				BaseInfoView: view.BaseInfoView{UUID: "cluster-uuid", Name: "cluster"},
				ZoneUuid:     "zone-uuid",
			}},
		})
	})
	mux.HandleFunc("/zstack/v1/hosts", func(w http.ResponseWriter, req *http.Request) {
		if got := req.URL.Query().Get("q"); got != "clusterUuid=cluster-uuid" {
			t.Fatalf("unexpected host query: %q", got)
		}
		hosts := []view.HostInventoryView{
			{BaseInfoView: view.BaseInfoView{UUID: "host-2", Name: "second"}, ClusterUuid: "cluster-uuid"},
			{BaseInfoView: view.BaseInfoView{UUID: "host-1", Name: "first"}, ClusterUuid: "cluster-uuid"},
		}
		writeVxlanPreflightResponse(t, w, map[string]any{"inventories": hosts, "total": len(hosts)})
	})
	mux.HandleFunc("/zstack/v1/cluster/hosts-network-facts/cluster-uuid", func(w http.ResponseWriter, req *http.Request) {
		writeVxlanPreflightResponse(t, w, view.GetClusterHostNetworkFactsView{
			Success: true,
			Nics: []view.HostNetworkInterfaceInventoryView{
				{HostUuid: "host-1", InterfaceName: "ens3", IpAddresses: []string{"172.26.0.11/16"}},
				{HostUuid: "host-2", InterfaceName: "ens3", IpAddresses: []string{"172.26.0.12/16"}},
			},
		})
	})
	mux.HandleFunc("/zstack/v1/l2-networks/vxlan-pool/vni-range", func(w http.ResponseWriter, req *http.Request) {
		if got := req.URL.Query().Get("q"); got != "l2NetworkUuid=pool-uuid" {
			t.Fatalf("unexpected VNI query: %q", got)
		}
		writeVxlanPreflightResponse(t, w, map[string]any{"inventories": options.ranges, "total": len(options.ranges)})
	})
	mux.HandleFunc("/zstack/v1/zops/check/network", func(w http.ResponseWriter, req *http.Request) {
		if options.failReachability {
			http.Error(w, "mock reachability failure", http.StatusServiceUnavailable)
			return
		}
		query := req.URL.Query()
		wantIps := []string{"172.26.0.11", "172.26.0.12"}
		if got := query["sourceHostnames"]; !equalStrings(got, wantIps) {
			t.Fatalf("unexpected sourceHostnames: %#v", got)
		}
		if got := query["targetHostnames"]; !equalStrings(got, wantIps) {
			t.Fatalf("unexpected targetHostnames: %#v", got)
		}
		writeVxlanPreflightResponse(t, w, view.CheckNetworkReachableView{
			Success: true,
			Results: []view.NetworkReachablePairView{
				{SourceHostname: "172.26.0.11", TargetHostname: "172.26.0.12", Status: "Connected"},
				{SourceHostname: "172.26.0.12", TargetHostname: "172.26.0.11", Status: "Connected"},
			},
		})
	})
	return httptest.NewServer(mux)
}

func vxlanPreflightTestConfig(t *testing.T, server *httptest.Server) string {
	return vxlanPreflightTestConfigWithExclude(t, server, "")
}

func vxlanPreflightTestConfigWithExclude(t *testing.T, server *httptest.Server, excludeVniRangeUuid string) string {
	t.Helper()
	host, port, err := net.SplitHostPort(server.Listener.Addr().String())
	if err != nil {
		t.Fatalf("split mock server address: %v", err)
	}
	excludeConfig := ""
	if excludeVniRangeUuid != "" {
		excludeConfig = fmt.Sprintf("  exclude_vni_range_uuid = %q\n", excludeVniRangeUuid)
	}
	return fmt.Sprintf(`
provider "zstack" {
  host              = %q
  port              = %s
  access_key_id     = "test-access-key"
  access_key_secret = "test-access-key-secret"
}

data "zstack_vxlan_preflight" "test" {
  zone_uuid          = "zone-uuid"
  cluster_uuid       = "cluster-uuid"
  pool_uuid          = "pool-uuid"
  physical_interface = "ens3"
  vtep_cidr          = "172.26.0.0/16"
  start_vni          = 100
  end_vni            = 200
%s
}
`, host, port, excludeConfig)
}

func writeVxlanPreflightResponse(t *testing.T, w http.ResponseWriter, response any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		t.Fatalf("write mock response: %v", err)
	}
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
