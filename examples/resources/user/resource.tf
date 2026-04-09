# Copyright (c) ZStack.io, Inc.

resource "zstack_user" "example" {
  name        = "example-user"
  password    = "example-password-123"
  description = "Example User Account"
}

output "user" {
  value = zstack_user.example
}
