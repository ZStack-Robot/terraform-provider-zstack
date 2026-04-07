# Copyright (c) ZStack.io, Inc.

resource "zstack_sdn_controller" "example" {
  name        = "sdn-controller-1"
  description = "SDN Controller for VXLAN"
  vendor_type = "H3C_VCFC"
  ip          = "192.168.1.100"
  username    = "admin"
  password    = "password"
}

output "sdn_controller" {
  value = zstack_sdn_controller.example
}
