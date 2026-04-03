# Copyright (c) ZStack.io, Inc.

resource "zstack_monitor_group" "example" {
  name = "example-monitor-group"
}

output "zstack_monitor_group" {
  value = zstack_monitor_group.example
}
