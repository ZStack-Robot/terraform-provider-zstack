# Copyright (c) ZStack.io, Inc.

resource "zstack_flow_collector" "example" {
  flow_meter_uuid = "example-uuid-placeholder"
}

output "zstack_flow_collector" {
  value = zstack_flow_collector.example
}
