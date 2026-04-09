# Copyright (c) ZStack.io, Inc.

resource "zstack_access_key" "example" {
  account_uuid = "example-account-uuid"
  user_uuid    = "example-user-uuid"
  description  = "API access key"
}

output "access_key" {
  value = zstack_access_key.example
}
