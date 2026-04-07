# Copyright (c) ZStack.io, Inc.

resource "zstack_log_server" "example" {
  name          = "example-log-server"
  category      = "syslog"
  type          = "syslog"
  configuration = "{\"hostname\": \"logs.example.com\", \"port\": 514}"
}

output "zstack_log_server" {
  value = zstack_log_server.example
}
