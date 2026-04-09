# Copyright (c) ZStack.io, Inc.

resource "zstack_l2vxlan_network" "example" {
  name        = "example-l2-vxlan-network"
  pool_uuid   = "5d09ee200ebf450b9c962ce2082a64f8"
  description = "Example L2 VXLAN Network"
  vni         = 1000
  zone_uuid   = "a1b2c3d4e5f6789a0b1c2d3e4f5a6b7c"
}

output "l2vxlan_network" {
  value = zstack_l2vxlan_network.example
}
