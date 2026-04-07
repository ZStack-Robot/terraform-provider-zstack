# Copyright (c) ZStack.io, Inc.

resource "zstack_iam2_virtual_id" "example" {
  name        = "test-user"
  password    = "SecureP@ssw0rd"
  description = "Test IAM2 virtual ID"
}

output "iam2_virtual_id" {
  value = zstack_iam2_virtual_id.example
}
