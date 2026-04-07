# Copyright (c) ZStack.io, Inc.

data "zstack_eips" "all" {
}

data "zstack_eips" "filtered" {
  name_pattern = "my-eip%"
}

output "eips" {
  value = data.zstack_eips.all
}
