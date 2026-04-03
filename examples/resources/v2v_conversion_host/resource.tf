# Copyright (c) ZStack.io, Inc.

resource "zstack_v2v_conversion_host" "example" {
  name         = "example-v2v-conversion-host"
  type         = "V2V"
  host_uuid    = "example-host-uuid"
  storage_path = "/data/v2v-conversion"
}

output "zstack_v2v_conversion_host" {
  value = zstack_v2v_conversion_host.example
}
