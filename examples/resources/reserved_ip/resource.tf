# Copyright (c) ZStack.io, Inc.

resource "zstack_reserved_ip" "test" {
  l3_network_uuid = "a5e77b2972e64316878993af7a695400"
  start_ip        = "172.26.111.250"
  end_ip          = "172.26.111.253"
}

output "zstack_reserved_ip" {
  value = zstack_reserved_ip.test
}


