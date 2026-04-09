# Copyright (c) ZStack.io, Inc.

resource "zstack_policy" "example" {
  name        = "example-policy"
  description = "Example IAM Policy"
}

output "policy" {
  value = zstack_policy.example
}
