# Copyright (c) ZStack.io, Inc.

resource "zstack_webhook" "example" {
  name = "example-webhook"
  url  = "https://example.com/webhook"
  type = "alert"
}

output "zstack_webhook" {
  value = zstack_webhook.example
}
