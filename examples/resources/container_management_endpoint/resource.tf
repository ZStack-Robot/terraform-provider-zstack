# Copyright (c) ZStack.io, Inc.

resource "zstack_container_management_endpoint" "example" {
  name                        = "example-container-management-endpoint"
  management_ip               = "192.168.1.107"
  vendor                      = "docker"
  management_port             = 8080
  container_access_key_id     = "example-key-id"
  container_access_key_secret = "example-secret-key"
}

output "zstack_container_management_endpoint" {
  value = zstack_container_management_endpoint.example
}
