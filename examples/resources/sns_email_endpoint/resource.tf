# Copyright (c) ZStack.io, Inc.

resource "zstack_sns_email_endpoint" "example" {
  name          = "example-email-endpoint"
  email         = "user@example.com"
  description   = "Example Email Endpoint"
  platform_uuid = "5d09ee200ebf450b9c962ce2082a64f8"
}

output "sns_email_endpoint" {
  value = zstack_sns_email_endpoint.example
}
