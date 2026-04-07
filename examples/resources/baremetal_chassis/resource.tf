# Copyright (c) ZStack.io, Inc.

resource "zstack_baremetal_chassis" "example" {
  name          = "example-baremetal-chassis"
  cluster_uuid  = "example-uuid-placeholder"
  ipmi_address  = "192.168.1.100"
  ipmi_username = "admin"
  ipmi_password = "password123"
}

output "zstack_baremetal_chassis" {
  value = zstack_baremetal_chassis.example
}
