# Copyright (c) ZStack.io, Inc.

resource "zstack_stack_template" "example" {
  name = "example-stack-template"
  template_content = jsonencode({
    ZStackTemplateFormatVersion = "2018-06-18"
    Resources                   = {}
  })
}

output "zstack_stack_template" {
  value = zstack_stack_template.example
}
