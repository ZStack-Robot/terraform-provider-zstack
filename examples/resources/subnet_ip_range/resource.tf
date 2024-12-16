# Copyright (c) ZStack.io, Inc.

resource "zstack_subnet" "test" {

  l3_network_uuid = "6a7c9dd9d6e449f992a59df8c102b3ba"
  name            = "net1"
  start_ip        = "192.168.1.1"
  end_ip          = "192.168.1.100"
  netmask         = "255.255.0.0"
  gateway         = "192.168.100.1"
}

output "zstack_subnet" {
  value = zstack_subnet.test
}

