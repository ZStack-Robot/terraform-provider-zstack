# Copyright (c) ZStack.io, Inc.

data "zstack_images" "images" {
  name = "image name"
}

data "zstack_l3network" "example" {
  name = "network name"
}

resource "zstack_vm" "vm" {
  count            = 2
  name             = "disk-test-${count.index + 1}"
  description      = "test"
  image_uuid       = data.zstack_images.images.images.0.uuid #"${data.zstack_images.images.images[0].uuid}" #"9b26312501614ec0b6dc731e6977dfb2"
  l3_network_uuids = [data.zstack_l3network.example.l3networks.0.uuid]
  root_disk = {
    offering_uuid = "e002a969bb3041c087e355c77e997be5"
  }
  memory_size = 1147483640
  cpu_num     = 1
}

output "zstack_vm" {
  value = zstack_vm.vm
}


