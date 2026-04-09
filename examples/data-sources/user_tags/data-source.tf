# Copyright (c) ZStack.io, Inc.

data "zstack_user_tags" "all" {
}

output "user_tags" {
  value = data.zstack_user_tags.all
}
