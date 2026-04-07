# Copyright (c) ZStack.io, Inc.

resource "zstack_policy_route_rule" "example" {
  rule_set_uuid = "5d09ee200ebf450b9c962ce2082a64f8"
  table_uuid    = "a1b2c3d4e5f6789a0b1c2d3e4f5a6b7c"
  rule_number   = 1
  dest_ip       = "192.168.1.0/24"
  source_ip     = "10.0.0.0/8"
  protocol      = "tcp"
}

output "policy_route_rule" {
  value = zstack_policy_route_rule.example
}
