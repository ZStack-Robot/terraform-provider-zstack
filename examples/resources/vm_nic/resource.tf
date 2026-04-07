# Copyright (c) ZStack.io, Inc.

resource "zstack_vm_nic" "example" {
  vm_instance_uuid = "example-vm-instance-uuid"
  l3_network_uuid  = "example-l3-network-uuid"
}

output "vm_nic" {
  value = zstack_vm_nic.example
}
