# Copyright (c) ZStack.io, Inc.

resource "zstack_san_sec_security_machine" "example" {
  name                  = "example-san-sec-security-machine"
  security_machine_type = "san_sec"
  url                   = "http://192.168.1.101:8080"
}

output "zstack_san_sec_security_machine" {
  value = zstack_san_sec_security_machine.example
}
