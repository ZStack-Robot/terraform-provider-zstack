# Copyright (c) ZStack.io, Inc.

resource "zstack_vip" "test" {    
  name = "vipfromtf"
  description = "vip desc"
  l3_network_uuid  = "0f5ce0fd3d074462bb752c70ee88eca2"
#  vip = "static virtual ip"  
}




output "zstack_vip" {
  value = zstack_vip.test
}