# Copyright (c) ZStack.io, Inc.

# Create a zone for organizing cloud resources
resource "zstack_zone" "example" {
  name        = "production-zone"
  description = "Production availability zone for cloud resources"
}

output "zstack_zone" {
  value = zstack_zone.example
}
