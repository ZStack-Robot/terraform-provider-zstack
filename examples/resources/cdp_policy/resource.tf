# Copyright (c) ZStack.io, Inc.

resource "zstack_cdp_policy" "example" {
  name                      = "example-cdp-policy"
  recovery_point_per_second = 60
}

output "zstack_cdp_policy" {
  value = zstack_cdp_policy.example
}
