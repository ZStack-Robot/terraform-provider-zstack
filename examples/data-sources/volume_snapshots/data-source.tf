# Copyright (c) ZStack.io, Inc.

data "zstack_volume_snapshots" "example" {
  # name = "my-snapshot"
  # name_pattern = "my-%"
}

output "zstack_volume_snapshots" {
  value = data.zstack_volume_snapshots.example
}
