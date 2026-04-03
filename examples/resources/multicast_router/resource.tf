# Copyright (c) ZStack.io, Inc.

resource "zstack_multicast_router" "example" {
  vpc_router_vm_uuid = "example-uuid-placeholder"
}

output "zstack_multicast_router" {
  value = zstack_multicast_router.example
}
