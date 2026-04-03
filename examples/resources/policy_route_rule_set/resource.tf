# Copyright (c) ZStack.io, Inc.

resource "zstack_policy_route_rule_set" "example" {
  name         = "example-policy-route-rule-set"
  vrouter_uuid = "5d09ee200ebf450b9c962ce2082a64f8"
  description  = "Example Policy Route Rule Set"
  type         = "vpc"
}

output "policy_route_rule_set" {
  value = zstack_policy_route_rule_set.example
}
