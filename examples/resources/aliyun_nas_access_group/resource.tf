# Copyright (c) ZStack.io, Inc.

resource "zstack_aliyun_nas_access_group" "example" {
  name             = "example-aliyun-nas-access-group"
  data_center_uuid = "example-uuid-placeholder"
}

output "zstack_aliyun_nas_access_group" {
  value = zstack_aliyun_nas_access_group.example
}
