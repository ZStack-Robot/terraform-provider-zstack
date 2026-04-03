# Copyright (c) ZStack.io, Inc.

resource "zstack_port_mirror" "example" {
  mirror_network_uuid = "example-uuid-placeholder"
}

output "zstack_port_mirror" {
  value = zstack_port_mirror.example
}
