# Copyright (c) ZStack.io, Inc.

resource "zstack_iam2_organization" "example" {
  name        = "example-organization"
  type        = "organization"
  description = "Example IAM2 Organization"
  parent_uuid = "5d09ee200ebf450b9c962ce2082a64f8"
}

output "iam2_organization" {
  value = zstack_iam2_organization.example
}
