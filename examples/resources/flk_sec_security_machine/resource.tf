# Copyright (c) ZStack.io, Inc.

resource "zstack_flk_sec_security_machine" "example" {
  name                  = "example-flk-sec-security-machine"
  security_machine_type = "flk_sec"
  url                   = "http://192.168.1.104:8080"
}

output "zstack_flk_sec_security_machine" {
  value = zstack_flk_sec_security_machine.example
}
