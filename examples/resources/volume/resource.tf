# Copyright (c) ZStack.io, Inc.

data "zstack_disk_offerings" "example" {
  name = "name of disk offering"
}

resource "zstack_volume" "example" {
  name               = "my-data-volume"
  description        = "A data volume created via Terraform"
  disk_offering_uuid = data.zstack_disk_offerings.example.disk_offers.0.uuid
  # vm_instance_uuid = "optional-vm-uuid-to-attach"
}

output "zstack_volume" {
  value = zstack_volume.example
}
