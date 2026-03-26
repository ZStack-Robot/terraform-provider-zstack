# Copyright (c) ZStack.io, Inc.

# Create a load balancer bound to a VIP address
resource "zstack_load_balancer" "example" {
  name        = "my-load-balancer"
  description = "Production load balancer for web traffic"
  vip_uuid    = "vip-uuid-placeholder"
}

output "zstack_load_balancer" {
  value = zstack_load_balancer.example
}
