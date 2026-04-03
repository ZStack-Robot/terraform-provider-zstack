# Copyright (c) ZStack.io, Inc.

resource "zstack_global_config" "example" {
  category = "vm"
  name     = "vm.cpuNum"
  value    = "128"
}

output "global_config" {
  value = zstack_global_config.example
}
