# Copyright (c) ZStack.io, Inc.

resource "zstack_account" "example" {
  name        = "my-account"
  password    = "SecurePassword123"
  description = "An example account"
}

output "zstack_account" {
  value = zstack_account.example
}
