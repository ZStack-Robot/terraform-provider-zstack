# Copyright (c) ZStack.io, Inc.

data "zstack_volumes" "example" {
  # name = "my-data-volume"
  # name_pattern = "my-%"
}

output "zstack_volumes" {
  value = data.zstack_volumes.example
}
