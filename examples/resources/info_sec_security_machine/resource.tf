# Copyright (c) ZStack.io, Inc.

resource "zstack_info_sec_security_machine" "example" {
  name                  = "example-info-sec-security-machine"
  security_machine_type = "info_sec"
  url                   = "http://192.168.1.102:8080"
}

output "zstack_info_sec_security_machine" {
  value = zstack_info_sec_security_machine.example
}
