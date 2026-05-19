# Copyright (c) ZStack.io, Inc.

resource "zstack_virtual_router_instance" "example" {
  name                         = "vpc-router-1"
  description                  = "VPC router example"
  virtual_router_offering_uuid = "example-virtual-router-offering-uuid"
}

resource "zstack_vpc_firewall" "example" {
  name        = "vpc-firewall-1"
  description = "VPC firewall example"
  vpc_uuid    = zstack_virtual_router_instance.example.uuid
}

output "vpc_firewall" {
  value = zstack_vpc_firewall.example
}
