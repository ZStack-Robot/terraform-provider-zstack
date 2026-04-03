# Copyright (c) ZStack.io, Inc.

resource "zstack_vpc_shared_qos" "example" {
  name            = "example-vpc-shared-qos"
  vpc_uuid        = "5d09ee200ebf450b9c962ce2082a64f8"
  l3_network_uuid = "a1b2c3d4e5f6789a0b1c2d3e4f5a6b7c"
  description     = "Example VPC Shared QoS"
  bandwidth       = 100000000
}

output "vpc_shared_qos" {
  value = zstack_vpc_shared_qos.example
}
