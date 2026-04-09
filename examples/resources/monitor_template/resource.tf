# Copyright (c) ZStack.io, Inc.

resource "zstack_monitor_template" "example" {
  name = "example-monitor-template"
}

output "zstack_monitor_template" {
  value = zstack_monitor_template.example
}
