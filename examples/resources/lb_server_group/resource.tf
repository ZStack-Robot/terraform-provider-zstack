# Copyright (c) ZStack.io, Inc.

resource "zstack_lb_server_group" "example" {
  name               = "example-lb-server-group"
  load_balancer_uuid = "5d09ee200ebf450b9c962ce2082a64f8"
  description        = "Example Load Balancer Server Group"
  ip_version         = "ipv4"
}

output "lb_server_group" {
  value = zstack_lb_server_group.example
}
