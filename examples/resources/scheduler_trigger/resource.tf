# Copyright (c) ZStack.io, Inc.

resource "zstack_scheduler_trigger" "example" {
  name               = "example-scheduler-trigger"
  scheduler_type     = "simple"
  description        = "Example Scheduler Trigger"
  scheduler_interval = 300
  repeat_count       = 10
  start_time         = "2024-01-01T00:00:00Z"
  cron               = "0 0 * * *"
}

output "scheduler_trigger" {
  value = zstack_scheduler_trigger.example
}
