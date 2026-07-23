// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/netip"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

const (
	vxlanPreflightEvidenceVersion = "v1"
	vxlanPreflightPageSize        = 500
)

var (
	_ datasource.DataSource                   = &vxlanPreflightDataSource{}
	_ datasource.DataSourceWithConfigure      = &vxlanPreflightDataSource{}
	_ datasource.DataSourceWithValidateConfig = &vxlanPreflightDataSource{}
)

type vxlanPreflightDataSource struct {
	api vxlanPreflightAPI
}

type vxlanPreflightAPI interface {
	vniRangePageQuery
	GetCluster(context.Context, string) (*view.ClusterInventoryView, error)
	PageHosts(context.Context, *param.QueryParam) ([]view.HostInventoryView, int, error)
	GetClusterHostNetworkFacts(context.Context, string) (*view.GetClusterHostNetworkFactsView, error)
	CheckNetworkReachable(context.Context, []string, []string) (*view.CheckNetworkReachableView, error)
}

type zstackVxlanPreflightAPI struct {
	client *client.ZSClient
}

func (a zstackVxlanPreflightAPI) GetCluster(ctx context.Context, uuid string) (*view.ClusterInventoryView, error) {
	var cluster view.ClusterInventoryView
	if err := a.client.ZSHttpClient.Get(ctx, "v1/clusters", uuid, nil, &cluster); err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (a zstackVxlanPreflightAPI) PageHosts(ctx context.Context, query *param.QueryParam) ([]view.HostInventoryView, int, error) {
	var hosts []view.HostInventoryView
	total, err := a.client.ZSHttpClient.Page(ctx, "v1/hosts", query, &hosts)
	return hosts, total, err
}

func (a zstackVxlanPreflightAPI) GetClusterHostNetworkFacts(ctx context.Context, clusterUuid string) (*view.GetClusterHostNetworkFactsView, error) {
	var facts view.GetClusterHostNetworkFactsView
	if err := a.client.ZSHttpClient.GetWithRespKey(
		ctx,
		"v1/cluster/hosts-network-facts",
		clusterUuid,
		"",
		nil,
		&facts,
	); err != nil {
		return nil, err
	}
	return &facts, nil
}

func (a zstackVxlanPreflightAPI) CheckNetworkReachable(
	ctx context.Context,
	sourceHostnames, targetHostnames []string,
) (*view.CheckNetworkReachableView, error) {
	query := struct {
		SourceHostnames []string `json:"sourceHostnames,omitempty"`
		TargetHostnames []string `json:"targetHostnames"`
	}{
		SourceHostnames: sourceHostnames,
		TargetHostnames: targetHostnames,
	}

	var result view.CheckNetworkReachableView
	if err := a.client.ZSHttpClient.GetWithRespKey(
		ctx,
		"v1/zops/check/network",
		"",
		"",
		&query,
		&result,
	); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a zstackVxlanPreflightAPI) PageVniRanges(
	ctx context.Context,
	query *param.QueryParam,
) ([]view.VniRangeInventoryView, int, error) {
	return zstackVniRangePageQuery{client: a.client}.PageVniRanges(ctx, query)
}

type vxlanPreflightDataSourceModel struct {
	ZoneUuid          types.String                      `tfsdk:"zone_uuid"`
	ClusterUuid       types.String                      `tfsdk:"cluster_uuid"`
	PoolUuid          types.String                      `tfsdk:"pool_uuid"`
	PhysicalInterface types.String                      `tfsdk:"physical_interface"`
	VtepCidr          types.String                      `tfsdk:"vtep_cidr"`
	StartVni          types.Int64                       `tfsdk:"start_vni"`
	EndVni            types.Int64                       `tfsdk:"end_vni"`
	Ready             types.Bool                        `tfsdk:"ready"`
	Hosts             []vxlanPreflightHostModel         `tfsdk:"hosts"`
	Connectivity      []vxlanPreflightConnectivityModel `tfsdk:"connectivity"`
	EvidenceDigest    types.String                      `tfsdk:"evidence_digest"`
}

type vxlanPreflightHostModel struct {
	HostUuid types.String `tfsdk:"host_uuid"`
	HostName types.String `tfsdk:"host_name"`
	VtepIp   types.String `tfsdk:"vtep_ip"`
}

type vxlanPreflightConnectivityModel struct {
	SourceIp types.String `tfsdk:"source_ip"`
	TargetIp types.String `tfsdk:"target_ip"`
	Status   types.String `tfsdk:"status"`
}

type vxlanPreflightHostEvidence struct {
	HostUuid string `json:"hostUuid"`
	HostName string `json:"hostName"`
	VtepIp   string `json:"vtepIp"`
}

type vxlanPreflightConnectivityEvidence struct {
	SourceIp string `json:"sourceIp"`
	TargetIp string `json:"targetIp"`
	Status   string `json:"status"`
}

type vxlanPreflightEvidence struct {
	Version           string                               `json:"version"`
	ZoneUuid          string                               `json:"zoneUuid"`
	ClusterUuid       string                               `json:"clusterUuid"`
	PoolUuid          string                               `json:"poolUuid"`
	PhysicalInterface string                               `json:"physicalInterface"`
	VtepCidr          string                               `json:"vtepCidr"`
	StartVni          int64                                `json:"startVni"`
	EndVni            int64                                `json:"endVni"`
	Hosts             []vxlanPreflightHostEvidence         `json:"hosts"`
	Connectivity      []vxlanPreflightConnectivityEvidence `json:"connectivity"`
}

func ZStackVxlanPreflightDataSource() datasource.DataSource {
	return &vxlanPreflightDataSource{}
}

func (d *vxlanPreflightDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	configuredClient, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer.", req.ProviderData),
		)
		return
	}
	d.api = zstackVxlanPreflightAPI{client: configuredClient}
}

func (d *vxlanPreflightDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vxlan_preflight"
}

func (d *vxlanPreflightDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	nonEmptyString := []validator.String{stringvalidator.LengthAtLeast(1)}
	resp.Schema = schema.Schema{
		Description: "Validates VXLAN host interface, VTEP reachability, and VNI range readiness without changing ZStack resources.",
		Attributes: map[string]schema.Attribute{
			"zone_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the zone expected to contain the cluster.",
				Validators:  nonEmptyString,
			},
			"cluster_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the cluster whose compute hosts are validated.",
				Validators:  nonEmptyString,
			},
			"pool_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the VXLAN pool used to scope VNI overlap checks.",
				Validators:  nonEmptyString,
			},
			"physical_interface": schema.StringAttribute{
				Required:    true,
				Description: "Physical NIC or bonding name that must exist on every cluster host.",
				Validators:  nonEmptyString,
			},
			"vtep_cidr": schema.StringAttribute{
				Required:    true,
				Description: "IPv4 CIDR containing exactly one unambiguous VTEP address on every host.",
				Validators:  nonEmptyString,
			},
			"start_vni": schema.Int64Attribute{
				Required:    true,
				Description: "First VNI to validate.",
				Validators: []validator.Int64{
					int64validator.Between(1, maxVni),
				},
			},
			"end_vni": schema.Int64Attribute{
				Required:    true,
				Description: "Last VNI to validate.",
				Validators: []validator.Int64{
					int64validator.Between(1, maxVni),
				},
			},
			"ready": schema.BoolAttribute{
				Computed:    true,
				Description: "True when every preflight check succeeds.",
			},
			"hosts": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Canonical host-to-VTEP mappings sorted by host UUID.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host_uuid": schema.StringAttribute{
							Computed: true,
						},
						"host_name": schema.StringAttribute{
							Computed: true,
						},
						"vtep_ip": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"connectivity": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Connected non-self directed VTEP address pairs, sorted by source and target IP.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"source_ip": schema.StringAttribute{
							Computed: true,
						},
						"target_ip": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"evidence_digest": schema.StringAttribute{
				Computed:    true,
				Description: "SHA-256 digest of versioned canonical readiness evidence. Credentials and raw API responses are excluded.",
			},
		},
	}
}

func (d *vxlanPreflightDataSource) ValidateConfig(
	ctx context.Context,
	req datasource.ValidateConfigRequest,
	resp *datasource.ValidateConfigResponse,
) {
	var config vxlanPreflightDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.VtepCidr.IsNull() && !config.VtepCidr.IsUnknown() {
		prefix, err := netip.ParsePrefix(strings.TrimSpace(config.VtepCidr.ValueString()))
		if err != nil || !prefix.Addr().Is4() {
			resp.Diagnostics.AddAttributeError(
				path.Root("vtep_cidr"),
				"Invalid VTEP CIDR",
				"vtep_cidr must be a valid IPv4 prefix.",
			)
		}
	}

	if !config.StartVni.IsNull() && !config.StartVni.IsUnknown() &&
		!config.EndVni.IsNull() && !config.EndVni.IsUnknown() &&
		config.StartVni.ValueInt64() > config.EndVni.ValueInt64() {
		resp.Diagnostics.AddAttributeError(
			path.Root("start_vni"),
			"Invalid VNI range",
			"start_vni must be less than or equal to end_vni.",
		)
	}
}

func (d *vxlanPreflightDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vxlanPreflightDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if d.api == nil {
		resp.Diagnostics.AddError("VXLAN preflight client is not configured", "The ZStack client was not provided to the data source.")
		return
	}

	vtepPrefix, err := netip.ParsePrefix(strings.TrimSpace(state.VtepCidr.ValueString()))
	if err != nil || !vtepPrefix.Addr().Is4() {
		resp.Diagnostics.AddAttributeError(path.Root("vtep_cidr"), "Invalid VTEP CIDR", "vtep_cidr must be a valid IPv4 prefix.")
		return
	}
	vtepPrefix = vtepPrefix.Masked()

	cluster, err := d.api.GetCluster(ctx, state.ClusterUuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("VXLAN preflight cluster query failed", err.Error())
		return
	}
	if cluster.ZoneUuid != state.ZoneUuid.ValueString() {
		resp.Diagnostics.AddAttributeError(
			path.Root("zone_uuid"),
			"Cluster is not in the requested zone",
			fmt.Sprintf("Cluster %s belongs to zone %s, not %s.", cluster.UUID, cluster.ZoneUuid, state.ZoneUuid.ValueString()),
		)
		return
	}

	hosts, err := queryAllVxlanPreflightHosts(ctx, d.api, state.ClusterUuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("VXLAN preflight host query failed", err.Error())
		return
	}
	if len(hosts) == 0 {
		resp.Diagnostics.AddError("VXLAN preflight found no hosts", fmt.Sprintf("Cluster %s contains no compute hosts.", state.ClusterUuid.ValueString()))
		return
	}

	facts, err := d.api.GetClusterHostNetworkFacts(ctx, state.ClusterUuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("VXLAN preflight network facts query failed", err.Error())
		return
	}
	if !facts.Success {
		resp.Diagnostics.AddError("VXLAN preflight network facts query failed", "ZStack returned success=false for cluster host network facts.")
		return
	}

	hostEvidence, vtepIps, err := buildVxlanPreflightHostEvidence(
		hosts,
		facts,
		state.PhysicalInterface.ValueString(),
		vtepPrefix,
	)
	if err != nil {
		resp.Diagnostics.AddError("VXLAN preflight host validation failed", err.Error())
		return
	}

	ranges, err := queryAllVniRanges(ctx, d.api, state.PoolUuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("VXLAN preflight VNI query failed", err.Error())
		return
	}
	if conflict := findVniRangeOverlap(ranges, state.StartVni.ValueInt64(), state.EndVni.ValueInt64(), ""); conflict != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("start_vni"),
			"VXLAN preflight found an overlapping VNI range",
			fmt.Sprintf(
				"Requested range [%d, %d] overlaps %q (%s) [%d, %d] in pool %s.",
				state.StartVni.ValueInt64(),
				state.EndVni.ValueInt64(),
				conflict.Name,
				conflict.UUID,
				conflict.StartVni,
				conflict.EndVni,
				state.PoolUuid.ValueString(),
			),
		)
		return
	}

	reachability, err := d.api.CheckNetworkReachable(ctx, vtepIps, vtepIps)
	if err != nil {
		resp.Diagnostics.AddError("VXLAN preflight reachability query failed", err.Error())
		return
	}
	connectivityEvidence, err := validateVxlanPreflightConnectivity(vtepIps, reachability)
	if err != nil {
		resp.Diagnostics.AddError("VXLAN preflight reachability validation failed", err.Error())
		return
	}

	evidence := vxlanPreflightEvidence{
		Version:           vxlanPreflightEvidenceVersion,
		ZoneUuid:          state.ZoneUuid.ValueString(),
		ClusterUuid:       state.ClusterUuid.ValueString(),
		PoolUuid:          state.PoolUuid.ValueString(),
		PhysicalInterface: state.PhysicalInterface.ValueString(),
		VtepCidr:          vtepPrefix.String(),
		StartVni:          state.StartVni.ValueInt64(),
		EndVni:            state.EndVni.ValueInt64(),
		Hosts:             hostEvidence,
		Connectivity:      connectivityEvidence,
	}
	digest, err := vxlanPreflightEvidenceDigest(evidence)
	if err != nil {
		resp.Diagnostics.AddError("VXLAN preflight evidence generation failed", err.Error())
		return
	}

	state.Ready = types.BoolValue(true)
	state.Hosts = make([]vxlanPreflightHostModel, 0, len(hostEvidence))
	for _, host := range hostEvidence {
		state.Hosts = append(state.Hosts, vxlanPreflightHostModel{
			HostUuid: types.StringValue(host.HostUuid),
			HostName: types.StringValue(host.HostName),
			VtepIp:   types.StringValue(host.VtepIp),
		})
	}
	state.Connectivity = make([]vxlanPreflightConnectivityModel, 0, len(connectivityEvidence))
	for _, pair := range connectivityEvidence {
		state.Connectivity = append(state.Connectivity, vxlanPreflightConnectivityModel{
			SourceIp: types.StringValue(pair.SourceIp),
			TargetIp: types.StringValue(pair.TargetIp),
			Status:   types.StringValue(pair.Status),
		})
	}
	state.EvidenceDigest = types.StringValue(digest)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func queryAllVxlanPreflightHosts(
	ctx context.Context,
	api vxlanPreflightAPI,
	clusterUuid string,
) ([]view.HostInventoryView, error) {
	hosts := make([]view.HostInventoryView, 0)
	seenUuids := make(map[string]struct{})
	expectedTotal := -1

	for start := 0; ; start = len(hosts) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		query := param.NewQueryParam()
		query.AddQ(fmt.Sprintf("clusterUuid=%s", clusterUuid))
		query.Start(start).Limit(vxlanPreflightPageSize).ReplyWithCount(true)
		page, total, err := api.PageHosts(ctx, &query)
		if err != nil {
			return nil, err
		}
		if total < 0 {
			return nil, fmt.Errorf("host query returned invalid total %d", total)
		}
		if expectedTotal == -1 {
			expectedTotal = total
		} else if total != expectedTotal {
			return nil, fmt.Errorf("host query total changed during pagination: expected %d, got %d", expectedTotal, total)
		}

		if len(page) == 0 {
			if len(hosts) != expectedTotal {
				return nil, fmt.Errorf("incomplete host query: received %d of %d records", len(hosts), expectedTotal)
			}
			return hosts, nil
		}

		for _, host := range page {
			if host.ClusterUuid != clusterUuid {
				return nil, fmt.Errorf("host query for cluster %s returned host %s from cluster %s", clusterUuid, host.UUID, host.ClusterUuid)
			}
			if host.UUID == "" {
				return nil, fmt.Errorf("host query returned a host without a UUID")
			}
			if _, exists := seenUuids[host.UUID]; exists {
				return nil, fmt.Errorf("host query returned duplicate host %s", host.UUID)
			}
			seenUuids[host.UUID] = struct{}{}
			hosts = append(hosts, host)
		}

		if len(hosts) > expectedTotal {
			return nil, fmt.Errorf("host query returned %d records, exceeding total %d", len(hosts), expectedTotal)
		}
		if len(hosts) == expectedTotal {
			sort.Slice(hosts, func(i, j int) bool { return hosts[i].UUID < hosts[j].UUID })
			return hosts, nil
		}
	}
}

func buildVxlanPreflightHostEvidence(
	hosts []view.HostInventoryView,
	facts *view.GetClusterHostNetworkFactsView,
	physicalInterface string,
	vtepPrefix netip.Prefix,
) ([]vxlanPreflightHostEvidence, []string, error) {
	hostByUuid := make(map[string]view.HostInventoryView, len(hosts))
	for _, host := range hosts {
		if _, exists := hostByUuid[host.UUID]; exists {
			return nil, nil, fmt.Errorf("duplicate host UUID %s", host.UUID)
		}
		hostByUuid[host.UUID] = host
	}

	interfaceFound := make(map[string]bool, len(hosts))
	candidates := make(map[string]map[string]struct{}, len(hosts))
	addAddresses := func(hostUuid, interfaceName string, addresses []string) error {
		if interfaceName != physicalInterface {
			return nil
		}
		if _, exists := hostByUuid[hostUuid]; !exists {
			return nil
		}
		interfaceFound[hostUuid] = true
		if candidates[hostUuid] == nil {
			candidates[hostUuid] = make(map[string]struct{})
		}
		for _, rawAddress := range addresses {
			address, err := parseVxlanPreflightAddress(rawAddress)
			if err != nil {
				return fmt.Errorf("host %s interface %s returned invalid IP address %q: %w", hostUuid, physicalInterface, rawAddress, err)
			}
			if address.Is4() && vtepPrefix.Contains(address) {
				candidates[hostUuid][address.String()] = struct{}{}
			}
		}
		return nil
	}

	for _, nic := range facts.Nics {
		if err := addAddresses(nic.HostUuid, nic.InterfaceName, nic.IpAddresses); err != nil {
			return nil, nil, err
		}
	}
	for _, bonding := range facts.Bondings {
		if err := addAddresses(bonding.HostUuid, bonding.BondingName, bonding.IpAddresses); err != nil {
			return nil, nil, err
		}
	}

	evidence := make([]vxlanPreflightHostEvidence, 0, len(hosts))
	vtepIps := make([]string, 0, len(hosts))
	vtepOwner := make(map[string]string, len(hosts))
	for _, host := range hosts {
		if !interfaceFound[host.UUID] {
			return nil, nil, fmt.Errorf("host %s does not expose interface or bonding %q", host.UUID, physicalInterface)
		}
		if len(candidates[host.UUID]) == 0 {
			return nil, nil, fmt.Errorf("host %s interface %q has no IPv4 address in %s", host.UUID, physicalInterface, vtepPrefix.String())
		}
		if len(candidates[host.UUID]) > 1 {
			return nil, nil, fmt.Errorf(
				"host %s interface %q has multiple IPv4 addresses in %s; VTEP selection is ambiguous",
				host.UUID,
				physicalInterface,
				vtepPrefix.String(),
			)
		}

		var vtepIp string
		for candidate := range candidates[host.UUID] {
			vtepIp = candidate
		}
		if owner, exists := vtepOwner[vtepIp]; exists {
			return nil, nil, fmt.Errorf("VTEP IP %s is duplicated on hosts %s and %s", vtepIp, owner, host.UUID)
		}
		vtepOwner[vtepIp] = host.UUID
		evidence = append(evidence, vxlanPreflightHostEvidence{
			HostUuid: host.UUID,
			HostName: host.Name,
			VtepIp:   vtepIp,
		})
		vtepIps = append(vtepIps, vtepIp)
	}

	sort.Slice(evidence, func(i, j int) bool { return evidence[i].HostUuid < evidence[j].HostUuid })
	sort.Strings(vtepIps)
	return evidence, vtepIps, nil
}

func parseVxlanPreflightAddress(value string) (netip.Addr, error) {
	value = strings.TrimSpace(value)
	if strings.Contains(value, "/") {
		prefix, err := netip.ParsePrefix(value)
		if err != nil {
			return netip.Addr{}, err
		}
		return prefix.Addr().Unmap(), nil
	}
	address, err := netip.ParseAddr(value)
	if err != nil {
		return netip.Addr{}, err
	}
	return address.Unmap(), nil
}

func validateVxlanPreflightConnectivity(
	vtepIps []string,
	result *view.CheckNetworkReachableView,
) ([]vxlanPreflightConnectivityEvidence, error) {
	if result == nil {
		return nil, fmt.Errorf("ZStack returned no reachability result")
	}
	if !result.Success {
		return nil, fmt.Errorf("ZStack returned success=false for network reachability")
	}

	knownIps := make(map[string]struct{}, len(vtepIps))
	for _, ip := range vtepIps {
		if _, exists := knownIps[ip]; exists {
			return nil, fmt.Errorf("duplicate VTEP IP %s", ip)
		}
		knownIps[ip] = struct{}{}
	}

	type pairKey struct {
		source string
		target string
	}
	pairs := make(map[pairKey]view.NetworkReachablePairView)
	for _, pair := range result.Results {
		if _, exists := knownIps[pair.SourceHostname]; !exists {
			return nil, fmt.Errorf("reachability result contains unknown source %s", pair.SourceHostname)
		}
		if _, exists := knownIps[pair.TargetHostname]; !exists {
			return nil, fmt.Errorf("reachability result contains unknown target %s", pair.TargetHostname)
		}
		if pair.SourceHostname == pair.TargetHostname {
			continue
		}
		key := pairKey{source: pair.SourceHostname, target: pair.TargetHostname}
		if _, exists := pairs[key]; exists {
			return nil, fmt.Errorf("reachability result contains duplicate pair %s -> %s", key.source, key.target)
		}
		pairs[key] = pair
	}

	evidence := make([]vxlanPreflightConnectivityEvidence, 0, len(vtepIps)*(len(vtepIps)-1))
	for _, source := range vtepIps {
		for _, target := range vtepIps {
			if source == target {
				continue
			}
			key := pairKey{source: source, target: target}
			pair, exists := pairs[key]
			if !exists {
				return nil, fmt.Errorf("reachability result is missing pair %s -> %s", source, target)
			}
			if pair.Status != "Connected" {
				return nil, fmt.Errorf("VTEP pair %s -> %s is %q, expected Connected", source, target, pair.Status)
			}
			evidence = append(evidence, vxlanPreflightConnectivityEvidence{
				SourceIp: source,
				TargetIp: target,
				Status:   pair.Status,
			})
		}
	}
	sort.Slice(evidence, func(i, j int) bool {
		if evidence[i].SourceIp == evidence[j].SourceIp {
			return evidence[i].TargetIp < evidence[j].TargetIp
		}
		return evidence[i].SourceIp < evidence[j].SourceIp
	})
	return evidence, nil
}

func vxlanPreflightEvidenceDigest(evidence vxlanPreflightEvidence) (string, error) {
	encoded, err := json.Marshal(evidence)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:]), nil
}
