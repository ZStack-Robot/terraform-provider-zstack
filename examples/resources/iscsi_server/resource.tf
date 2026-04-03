# Copyright (c) ZStack.io, Inc.

resource "zstack_iscsi_server" "example" {
  name = "example-iscsi-server"
  ip   = "192.168.1.105"
}

output "zstack_iscsi_server" {
  value = zstack_iscsi_server.example
}
