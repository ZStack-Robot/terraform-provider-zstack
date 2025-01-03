# Copyright (c) ZStack.io, Inc.

resource "zstack_eip" "test" {
  name        = "eip"
  description = "eip desc"
  vip_uuid    = "5d09ee200ebf450b9c962ce2082a64f8"
  vm_nic_uuid = "910ea22533ba41f0b037e063eb207c2e"
}

output "zstack_eip" {
  value = zstack_eip.test
}


