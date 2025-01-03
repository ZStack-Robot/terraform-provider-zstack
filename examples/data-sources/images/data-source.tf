#  Copyright (c) ZStack.io, Inc.

data "zstack_images" "example" {
  #  name = "imageName"
  #   name_pattern = "hostname%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter {
    name   = "architecture"
    values = ["aarch64", "x86_64"]
  }
  filter {
    name   = "state"
    values = ["Enabled"]
  }
  filter {
    name   = "status"
    values = ["Ready", "Deleted"]
  }
  filter {
    name   = "guest_os_type"
    values = ["Linux"]
  }
}

output "zstack_images" {
  value = data.zstack_images.example
}


