# Copyright (c) ZStack.io, Inc.

resource "zstack_sns_topic" "example" {
  name        = "example-sns-topic"
  description = "Example SNS Topic"
}

output "sns_topic" {
  value = zstack_sns_topic.example
}
