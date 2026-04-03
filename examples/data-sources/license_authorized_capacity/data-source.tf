# Copyright (c) ZStack.io, Inc.

data "zstack_license_authorized_capacity" "current" {
}

output "license_authorized_capacity" {
  value = data.zstack_license_authorized_capacity.current
}
