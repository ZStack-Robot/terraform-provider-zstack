# Copyright (c) ZStack.io, Inc.

resource "zstack_volume_snapshot" "example" {
  name        = "my-snapshot"
  description = "A volume snapshot created via Terraform"
  volume_uuid = "uuid-of-volume-to-snapshot"
}

output "zstack_volume_snapshot" {
  value = zstack_volume_snapshot.example
}
