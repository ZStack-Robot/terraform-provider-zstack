# Copyright (c) ZStack.io, Inc.

data "zstack_global_configs" "vm_deletion_policy" {
  category = "vm"
  name     = "deletionPolicy"
}

output "vm_deletion_policy" {
  value = data.zstack_global_configs.vm_deletion_policy.global_configs[0]
}
