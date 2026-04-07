# Copyright (c) ZStack.io, Inc.

resource "zstack_vrouter_route_entry" "example" {
  route_table_uuid = "example-route-table-uuid"
  destination      = "10.0.0.0/24"
  target           = "192.168.1.1"
  type             = "User"
}

output "vrouter_route_entry" {
  value = zstack_vrouter_route_entry.example
}
