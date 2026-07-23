# Copyright (c) ZStack.io, Inc.

resource "zstack_l2vxlan_network_pool" "example" {
  name               = "example-vxlan-pool"
  zone_uuid          = "a1b2c3d4e5f6789a0b1c2d3e4f5a6b7c"
  physical_interface = "bond0"
}

resource "zstack_vni_range" "example" {
  name        = "example-vni-range"
  description = "Example VNI range"
  start_vni   = 1000
  end_vni     = 2000
  pool_uuid   = zstack_l2vxlan_network_pool.example.uuid
}

output "vni_range" {
  value = zstack_vni_range.example
}
