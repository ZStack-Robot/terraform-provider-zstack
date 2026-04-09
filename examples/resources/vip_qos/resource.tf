# Copyright (c) ZStack.io, Inc.

resource "zstack_vip_qos" "example" {
  vip_uuid           = "5d09ee200ebf450b9c962ce2082a64f8"
  port               = 80
  outbound_bandwidth = 10000000
  inbound_bandwidth  = 10000000
}

output "vip_qos" {
  value = zstack_vip_qos.example
}
