# Copyright (c) ZStack.io, Inc.

resource "zstack_vpc_firewall" "example" {
  name            = "vpc-firewall-1"
  description     = "VPC firewall example"
  l3_network_uuid = "example-l3-network-uuid"
}

output "vpc_firewall" {
  value = zstack_vpc_firewall.example
}
