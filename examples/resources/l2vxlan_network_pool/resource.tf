# Copyright (c) ZStack.io, Inc.

resource "zstack_l2vxlan_network_pool" "example" {
  name               = "example-vxlan-pool"
  description        = "Example L2 VXLAN network pool"
  zone_uuid          = "a1b2c3d4e5f6789a0b1c2d3e4f5a6b7c"
  physical_interface = "bond0"
  vswitch_type       = "LinuxBridge"
}

output "l2vxlan_network_pool" {
  value = zstack_l2vxlan_network_pool.example
}
