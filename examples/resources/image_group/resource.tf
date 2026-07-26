# Copyright (c) ZStack.io, Inc.

resource "zstack_image_group" "example" {
  name                      = "example-image-group"
  description               = "Image group managed by Terraform"
  root_volume_template_uuid = "a1b2c3d4e5f6789a0b1c2d3e4f5a6b7c"

  data_volume_template_uuids = [
    "b1c2d3e4f5a6789a0b1c2d3e4f5a6b7c",
  ]
}

output "image_group_uuid" {
  value = zstack_image_group.example.uuid
}

output "image_group_status" {
  value = zstack_image_group.example.status
}
