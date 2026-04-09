# Copyright (c) ZStack.io, Inc.

resource "zstack_email_media" "example" {
  name        = "example-email-media"
  smtp_server = "smtp.example.com"
  smtp_port   = 587
}

output "zstack_email_media" {
  value = zstack_email_media.example
}
