# Copyright (c) ZStack.io, Inc.

# Create a load balancer listener that forwards HTTP traffic
resource "zstack_load_balancer_listener" "example" {
  name               = "http-listener"
  description        = "Forward port 80 traffic to backend instances on port 8080"
  load_balancer_uuid = "lb-uuid-placeholder"
  protocol           = "tcp"
  load_balancer_port = 80   # Frontend port receiving traffic
  instance_port      = 8080 # Backend port on instances
}

output "zstack_load_balancer_listener" {
  value = zstack_load_balancer_listener.example
}
