# Copyright (c) ZStack.io, Inc.

resource "zstack_nvme_server" "example" {
  name      = "example-nvme-server"
  ip        = "192.168.1.106"
  transport = "tcp"
}

output "zstack_nvme_server" {
  value = zstack_nvme_server.example
}
