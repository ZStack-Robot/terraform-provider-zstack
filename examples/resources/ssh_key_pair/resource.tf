# Copyright (c) ZStack.io, Inc.

resource "zstack_ssh_key_pair" "example" {
  name        = "my-ssh-key"
  public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB... user@host"
  description = "An example SSH key pair"
}

output "zstack_ssh_key_pair" {
  value = zstack_ssh_key_pair.example
}
