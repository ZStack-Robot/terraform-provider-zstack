# Copyright (c) ZStack.io, Inc.

resource "zstack_affinity_group" "example" {
  name        = "my-affinity-group"
  policy      = "antiSoft" # antiSoft, antiHard, proSoft, proHard
  description = "An example affinity group"
  # type      = "host"     # optional, defaults to host
  # zone_uuid = "..."      # optional, zone UUID
}

output "zstack_affinity_group" {
  value = zstack_affinity_group.example
}
