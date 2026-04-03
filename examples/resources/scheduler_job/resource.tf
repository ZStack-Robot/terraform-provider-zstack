# Copyright (c) ZStack.io, Inc.

resource "zstack_scheduler_job" "example" {
  name                 = "example-scheduler-job"
  target_resource_uuid = "5d09ee200ebf450b9c962ce2082a64f8"
  type                 = "scheduler_job"
  description          = "Example Scheduler Job"
}

output "scheduler_job" {
  value = zstack_scheduler_job.example
}
