# Copyright (c) ZStack.io, Inc.

resource "zstack_log_server" "example" {
  name          = "example-log-server"
  category      = "ManagementNodeLog"
  type          = "Log4j2"
  level         = "INFO"
  appender_type = "Syslog"
  appender_configuration = {
    hostname = "logs.example.com"
    port     = "514"
    protocol = "UDP"
    facility = "LOCAL5"
  }
}

output "zstack_log_server" {
  value = zstack_log_server.example
}
