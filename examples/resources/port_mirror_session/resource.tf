# Copyright (c) ZStack.io, Inc.

resource "zstack_port_mirror_session" "example" {
  name             = "example-port-mirror-session"
  port_mirror_uuid = "example-uuid-placeholder"
  type             = "INGRESS"
  src_end_point    = "example-source-endpoint"
  dst_end_point    = "example-dest-endpoint"
}

output "zstack_port_mirror_session" {
  value = zstack_port_mirror_session.example
}
