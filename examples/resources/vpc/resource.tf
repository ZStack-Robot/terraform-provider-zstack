# Copyright (c) ZStack.io, Inc.

resource "zstack_vpc" "test" {
  l2_network_uuid     = "l2_network_uuid data source"
  name                = "example"
  description         = "vpc network description"
  enable_ipam         = true
  dns                 = "dns ip address"
  virtual_router_uuid = "Attach virtual router  for this VPC network."
  subnet_cidr = {
    name         = "example-subnet"
    network_cidr = "192.168.110.0/24" # subnet cidr block
    gateway      = "192.168.110.1"    # gateway of subnet
  }
}

output "zstack_vpc" {
  value = zstack_vpc.test
}
