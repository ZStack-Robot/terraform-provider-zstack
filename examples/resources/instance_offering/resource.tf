# Copyright (c) ZStack.io, Inc.

resource "zstack_instance_offering" "example" {
  name        = "instanceoffertest"
  description = "An example instance offer"
  cpu_num     = 1
  memory_size = 1024 # in megabytes, MB
}

output "zstack_instance_offering" {
  value = zstack_instance_offering.example
}