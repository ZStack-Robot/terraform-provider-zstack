# Copyright (c) ZStack.io, Inc.

resource "zstack_resource_stack" "example" {
  name = "example-resource-stack"
}

output "zstack_resource_stack" {
  value = zstack_resource_stack.example
}
