# Copyright (c) ZStack.io, Inc.

resource "zstack_baremetal_pxe_server" "example" {
  name               = "example-baremetal-pxe-server"
  dhcp_interface     = "eth0"
  dhcp_range_begin   = "192.168.1.100"
  dhcp_range_end     = "192.168.1.200"
  dhcp_range_netmask = "255.255.255.0"
}

output "zstack_baremetal_pxe_server" {
  value = zstack_baremetal_pxe_server.example
}
