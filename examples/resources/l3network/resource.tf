# Copyright (c) ZStack.io, Inc.

resource "zstack_l3network" "example" {
  name            = "l3-network-1"
  description     = "L3 network example"
  l2_network_uuid = "example-l2-network-uuid"
  type            = "L3BasicNetwork"
  category        = "Private"
}

output "l3network" {
  value = zstack_l3network.example
}
