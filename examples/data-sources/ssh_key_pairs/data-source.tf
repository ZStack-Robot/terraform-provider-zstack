# Copyright (c) ZStack.io, Inc.

data "zstack_ssh_key_pairs" "example" {
  name = "my-ssh-key"
  # name_pattern = "my-%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
}

output "zstack_ssh_key_pairs" {
  value = data.zstack_ssh_key_pairs.example
}
