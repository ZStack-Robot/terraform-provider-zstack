#  Copyright (c) ZStack.io, Inc.

data "zstack_images" "example" {
  #  name = "imageName"
  #   name_pattern = "hostname%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter = {
    Status       = "Ready"
    State        = "Enabled"
    Platform     = "Linux"
    Architecture = "x86_64"
  }
}

output "zstack_images" {
  value = data.zstack_images.example
}


