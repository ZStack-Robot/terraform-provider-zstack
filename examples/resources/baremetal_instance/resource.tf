# Copyright (c) ZStack.io, Inc.

resource "zstack_baremetal_instance" "example" {
  name                   = "example-baremetal-instance"
  chassis_uuid           = "example-uuid-placeholder"
  image_uuid             = "example-uuid-placeholder"
  instance_offering_uuid = "example-uuid-placeholder"
}

output "zstack_baremetal_instance" {
  value = zstack_baremetal_instance.example
}
