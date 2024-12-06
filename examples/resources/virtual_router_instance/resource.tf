# Copyright (c) ZStack.io, Inc.

resource "zstack_virtual_router_instance" "test" {
  name                         = "example-vr"
  description                  = "Example Virtual router instance"
  virtual_router_offering_uuid = "uuid of virtual router offering"
}

output "zstack_virtual_router_instance" {
  value = zstack_virtual_router_instance.test
}