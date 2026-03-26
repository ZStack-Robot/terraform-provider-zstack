# Copyright (c) ZStack.io, Inc.

resource "zstack_iam2_project" "example" {
  name        = "my-project"
  description = "An example IAM2 project"
}

output "zstack_iam2_project" {
  value = zstack_iam2_project.example
}
