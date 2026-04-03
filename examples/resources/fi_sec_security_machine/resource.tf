# Copyright (c) ZStack.io, Inc.

resource "zstack_fi_sec_security_machine" "example" {
  name                  = "example-fi-sec-security-machine"
  security_machine_type = "fi_sec"
  url                   = "http://192.168.1.103:8080"
}

output "zstack_fi_sec_security_machine" {
  value = zstack_fi_sec_security_machine.example
}
