# Copyright (c) ZStack.io, Inc.



data "zstack_images" "example" {
    #  name = "imageName"
    #   name_pattern = "hostname%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
}



output "zstack_images" {
  value = data.zstack_images.example
}


