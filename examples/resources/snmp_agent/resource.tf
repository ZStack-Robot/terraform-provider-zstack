# Copyright (c) ZStack.io, Inc.

resource "zstack_snmp_agent" "example" {
  version = "2c"
  port    = 161
}

output "zstack_snmp_agent" {
  value = zstack_snmp_agent.example
}
