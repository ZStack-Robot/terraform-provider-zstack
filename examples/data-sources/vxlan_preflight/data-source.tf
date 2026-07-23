# Copyright (c) ZStack.io, Inc.

data "zstack_vxlan_preflight" "example" {
  zone_uuid          = "a1b2c3d4e5f6789a0b1c2d3e4f5a6b7c"
  cluster_uuid       = "b1c2d3e4f5a6789b0c1d2e3f4a5b6c7d"
  pool_uuid          = "c1d2e3f4a5b6789c0d1e2f3a4b5c6d7e"
  physical_interface = "bond0"
  vtep_cidr          = "172.25.0.0/16"
  start_vni          = 1000
  end_vni            = 2000
}

output "vxlan_readiness" {
  value = {
    ready           = data.zstack_vxlan_preflight.example.ready
    hosts           = data.zstack_vxlan_preflight.example.hosts
    connectivity    = data.zstack_vxlan_preflight.example.connectivity
    evidence_digest = data.zstack_vxlan_preflight.example.evidence_digest
  }
}
