# Copyright (c) ZStack.io, Inc.

resource "zstack_virtual_router_offer" "test" {
  name                    = "vroffer"
  description             = "An example virtual router offer"
  cpu_num                 = 1
  memory_size             = 1073741824
  zone_uuid               = "d29f4847a99f4dea83bc446c8fe6e64c"
  management_network_uuid = "50e8c0d69681447fbe347c8dae2b1bef"
  image_uuid              = "93005c8a2a314a489635eca8c30794d4"

}

output "zstack_vroffer" {
  value = zstack_virtual_router_offer.test
}