# Copyright (c) ZStack.io, Inc.

resource "zstack_vcenter" "example" {
  name        = "example-vcenter"
  domain_name = "vcenter.example.com"
  https_port  = 443
  user_name   = "administrator@vsphere.local"
  password    = "password123"
  zone_uuid   = "example-uuid-placeholder"
}

output "zstack_vcenter" {
  value = zstack_vcenter.example
}
