# Copyright (c) ZStack.io, Inc.

resource "zstack_disk_offering" "test" {
  name        = "largeDiskOffering-test"
  description = "An example disk offer"
  disk_size   = 100 # in gigabytes (GB)
}

output "zstack_disk_offering" {
  value = zstack_disk_offering.test
}


