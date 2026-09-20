# Copyright (c) ZStack.io, Inc.

resource "zstack_slb_offering" "example" {
  name        = "slb-8c8g"
  description = "Dedicated load balancer offering"
  cpu_num     = 8
  memory_size = 8192 # MiB (8 GiB)
  # Replace these UUIDs with an existing zone, management L3 network,
  # and an image suitable for SLB instances in your environment.
  zone_uuid               = "d29f4847a99f4dea83bc446c8fe6e64c"
  management_network_uuid = "50e8c0d69681447fbe347c8dae2b1bef"
  image_uuid              = "93005c8a2a314a489635eca8c30794d4"
}

output "slb_offering_uuid" {
  value = zstack_slb_offering.example.uuid
}
