---
page_title: "zstack_vxlan_preflight Data Source - terraform-provider-zstack"
subcategory: ""
description: |-
  Validates VXLAN host interface, VTEP reachability, and VNI range readiness without changing ZStack resources.
---

# zstack_vxlan_preflight (Data Source)

Validates VXLAN host interface, VTEP reachability, and VNI range readiness without changing ZStack resources.

All input values must be known for Terraform to run the checks during planning. The VXLAN pool must exist and its zone and physical interface must match the requested values.

For a target `zstack_vni_range`, preallocate a UUID and pass the same value to its `resource_uuid` and this data source's `exclude_vni_range_uuid`. This prevents the target range from conflicting with itself on later plans. Make the target range depend on the preflight data source so creation cannot race the check.

The check fails when the pool does not exist or has a different zone or physical interface, the cluster is outside the requested zone, a host lacks the requested NIC or bonding, VTEP address selection is missing or ambiguous, a VTEP address is duplicated, any non-self directed address pair is missing or not `Connected`, or the requested VNI range overlaps another range in the same pool.

## Example Usage

```terraform
# Copyright (c) ZStack.io, Inc.

locals {
  target_vni_range_uuid = "d1e2f3a4b5c6789d0e1f2a3b4c5d6e7f"
}

data "zstack_vxlan_preflight" "example" {
  zone_uuid              = "a1b2c3d4e5f6789a0b1c2d3e4f5a6b7c"
  cluster_uuid           = "b1c2d3e4f5a6789b0c1d2e3f4a5b6c7d"
  pool_uuid              = "c1d2e3f4a5b6789c0d1e2f3a4b5c6d7e"
  physical_interface     = "bond0"
  vtep_cidr              = "172.25.0.0/16"
  start_vni              = 1000
  end_vni                = 2000
  exclude_vni_range_uuid = local.target_vni_range_uuid
}

resource "zstack_vni_range" "target" {
  name          = "example-vni-range"
  pool_uuid     = "c1d2e3f4a5b6789c0d1e2f3a4b5c6d7e"
  start_vni     = 1000
  end_vni       = 2000
  resource_uuid = local.target_vni_range_uuid

  depends_on = [data.zstack_vxlan_preflight.example]
}

output "vxlan_readiness" {
  value = {
    ready           = data.zstack_vxlan_preflight.example.ready
    hosts           = data.zstack_vxlan_preflight.example.hosts
    connectivity    = data.zstack_vxlan_preflight.example.connectivity
    evidence_digest = data.zstack_vxlan_preflight.example.evidence_digest
  }
}
```

## Schema

### Required

- `cluster_uuid` (String) UUID of the cluster whose compute hosts are validated.
- `end_vni` (Number) Last VNI to validate.
- `physical_interface` (String) Physical NIC or bonding name that must exist on every cluster host.
- `pool_uuid` (String) UUID of the VXLAN pool used to scope VNI overlap checks.
- `start_vni` (Number) First VNI to validate.
- `vtep_cidr` (String) IPv4 CIDR containing exactly one unambiguous VTEP address on every host.
- `zone_uuid` (String) UUID of the zone expected to contain the cluster.

### Optional

- `exclude_vni_range_uuid` (String) UUID of the target VNI range to exclude from overlap detection. Use the same preallocated UUID as `zstack_vni_range.resource_uuid` for a stable pre-create check.

### Read-Only

- `connectivity` (Attributes List) Connected non-self directed VTEP address pairs, sorted by source and target IP. (see [below for nested schema](#nestedatt--connectivity))
- `evidence_digest` (String) SHA-256 digest of versioned canonical readiness evidence. Credentials and raw API responses are excluded.
- `hosts` (Attributes List) Canonical host-to-VTEP mappings sorted by host UUID. (see [below for nested schema](#nestedatt--hosts))
- `ready` (Boolean) True when every preflight check succeeds.

<a id="nestedatt--connectivity"></a>
### Nested Schema for `connectivity`

Read-Only:

- `source_ip` (String)
- `status` (String)
- `target_ip` (String)

<a id="nestedatt--hosts"></a>
### Nested Schema for `hosts`

Read-Only:

- `host_name` (String)
- `host_uuid` (String)
- `vtep_ip` (String)
