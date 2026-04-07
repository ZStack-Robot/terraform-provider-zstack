# Copyright (c) ZStack.io, Inc.

resource "zstack_cdp_task" "example" {
  name                = "example-cdp-task"
  task_type           = "incremental"
  policy_uuid         = "example-uuid-placeholder"
  backup_storage_uuid = "example-uuid-placeholder"
  resource_uuids      = ["example-uuid-placeholder-1", "example-uuid-placeholder-2"]
}

output "zstack_cdp_task" {
  value = zstack_cdp_task.example
}
