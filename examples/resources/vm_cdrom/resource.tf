# Copyright (c) ZStack.io, Inc.

resource "zstack_vm_cdrom" "example" {
  name             = "example-vm-cdrom"
  vm_instance_uuid = "5d09ee200ebf450b9c962ce2082a64f8"
  iso_uuid         = "a1b2c3d4e5f6789a0b1c2d3e4f5a6b7c"
  description      = "Example VM CDROM"
}

output "vm_cdrom" {
  value = zstack_vm_cdrom.example
}
