# Copyright (c) ZStack.io, Inc.

resource "zstack_l2vxlan_network_pool" "example" {
  name         = "example-vxlan-pool"
  description  = "Example L2 VXLAN network pool"
  zone_uuid    = "a1b2c3d4e5f6789a0b1c2d3e4f5a6b7c"
  vswitch_type = "LinuxBridge"

  # Optional. Omit this attribute to let ZStack discover the host interface
  # from the VTEP CIDR when the pool is attached to a cluster.
  # physical_interface = "bond0"
}

output "l2vxlan_network_pool" {
  value = zstack_l2vxlan_network_pool.example
}
