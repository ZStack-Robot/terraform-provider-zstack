# Copyright (c) ZStack.io, Inc.

# Create an auto scaling group for VM instances
resource "zstack_auto_scaling_group" "example" {
  name                  = "web-tier-scaling"
  description           = "Auto scaling group for web tier VM instances"
  scaling_resource_type = "VmInstance"
  default_cooldown      = 60               # Cooldown period in seconds between scaling activities
  min_resource_size     = 1                # Minimum number of instances
  max_resource_size     = 10               # Maximum number of instances
  removal_policy        = "OldestInstance" # OldestInstance, NewestInstance, OldestScalingConfiguration
}

output "zstack_auto_scaling_group" {
  value = zstack_auto_scaling_group.example
}
