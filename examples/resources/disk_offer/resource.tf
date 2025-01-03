# Copyright (c) ZStack.io, Inc.

resource "zstack_disk_offer" "test" {
  name        = "largeDiskOffering-test"
  description = "An example disk offer"
  disk_size   = 1073741824
}

output "zstack_disk_offer" {
  value = zstack_disk_offer.test
}


