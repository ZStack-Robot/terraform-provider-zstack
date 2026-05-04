# Copyright (c) ZStack.io, Inc.

resource "zstack_resource_stack" "example" {
  name = "example-resource-stack"
  template_content = jsonencode({
    ZStackTemplateFormatVersion = "2018-06-18"
    Resources                   = {}
  })
}

output "zstack_resource_stack" {
  value = zstack_resource_stack.example
}
