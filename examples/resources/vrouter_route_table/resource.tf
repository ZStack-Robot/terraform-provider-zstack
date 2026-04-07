# Copyright (c) ZStack.io, Inc.

resource "zstack_vrouter_route_table" "example" {
  name                = "route-table-1"
  description         = "VRouter route table example"
  virtual_router_uuid = "example-vrouter-uuid"
}

output "vrouter_route_table" {
  value = zstack_vrouter_route_table.example
}
