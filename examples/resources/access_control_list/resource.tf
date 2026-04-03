# Copyright (c) ZStack.io, Inc.

resource "zstack_access_control_list" "example" {
  name        = "example-acl"
  description = "Example Access Control List"
  ip_version  = "ipv4"
}

output "access_control_list" {
  value = zstack_access_control_list.example
}
