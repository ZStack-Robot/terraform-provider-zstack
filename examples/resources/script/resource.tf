# Copyright (c) ZStack.io, Inc.


resource "zstack_script" "example" {
  name           = "test-script"
  description    = "test-script-desc"
  script_content = "echo Hello World"
  platform       = "Linux"
  script_type    = "Shell"
  script_timeout = 50
  encoding_type  = "Base64"
}


output "zstack_script_out" {
  value = zstack_script.example
}


