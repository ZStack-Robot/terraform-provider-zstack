# Copyright (c) ZStack.io, Inc.

resource "zstack_sns_http_endpoint" "example" {
  name          = "example-http-endpoint"
  url           = "https://example.com/webhook"
  username      = "admin"
  password      = "example-password"
  description   = "Example HTTP Endpoint"
  platform_uuid = "5d09ee200ebf450b9c962ce2082a64f8"
}

output "sns_http_endpoint" {
  value = zstack_sns_http_endpoint.example
}
