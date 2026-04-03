# Copyright (c) ZStack.io, Inc.

resource "zstack_jit_security_machine" "example" {
  name                  = "example-jit-security-machine"
  security_machine_type = "jit"
  url                   = "http://192.168.1.100:8080"
}

output "zstack_jit_security_machine" {
  value = zstack_jit_security_machine.example
}
