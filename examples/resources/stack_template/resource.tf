# Copyright (c) ZStack.io, Inc.

resource "zstack_stack_template" "example" {
  name = "example-stack-template"
}

output "zstack_stack_template" {
  value = zstack_stack_template.example
}
