# Copyright (c) ZStack.io, Inc.

resource "zstack_alarm" "example" {
  name                = "example-alarm"
  comparison_operator = "greater_than"
  namespace           = "zstack.vm"
  metric_name         = "cpu_usage"
  threshold           = 80.0
  description         = "Example CPU Usage Alarm"
  period              = 300
  repeat_interval     = 600
}

output "alarm" {
  value = zstack_alarm.example
}
