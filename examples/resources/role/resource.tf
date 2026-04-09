# Copyright (c) ZStack.io, Inc.

resource "zstack_role" "example" {
  name        = "example-role"
  description = "Example IAM Role"
  identity    = "system"
}

output "role" {
  value = zstack_role.example
}
