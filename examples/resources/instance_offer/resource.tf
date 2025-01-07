# Copyright (c) ZStack.io, Inc.

resource "zstack_instance_offer" "example" {
  name        = "instanceoffertest"
  description = "An example instance offer"
  cpu_num     = 1
  memory_size = 1024 # in megabytes, MB
}

output "zstack_instance_offer" {
  value = zstack_instance_offer.example
}