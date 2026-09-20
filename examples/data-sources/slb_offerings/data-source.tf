# Copyright (c) ZStack.io, Inc.

# For a deterministic lookup, set uuid instead of name or name_pattern.
data "zstack_slb_offerings" "example" {
  name_pattern = "slb-%"
  filter {
    name   = "state"
    values = ["Enabled"]
  }
  filter {
    name   = "memory_size"
    values = ["8192"] # MiB
  }
}

output "slb_offerings" {
  value = data.zstack_slb_offerings.example.slb_offers
}

output "slb_offering_count" {
  value = length(data.zstack_slb_offerings.example.slb_offers)
}
