# Copyright (c) ZStack.io, Inc.

resource "zstack_flow_meter" "example" {
  type = "traffic"
}

output "zstack_flow_meter" {
  value = zstack_flow_meter.example
}
