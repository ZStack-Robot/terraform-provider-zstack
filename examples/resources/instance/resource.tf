data "zstack_images" "centos" {
  name = "centos8"
}

data "zstack_l3networks" "l3networks" {
  name = "public-net"
}

data "zstack_instance_offers" "offer" {
  name = "min"
}

resource "zstack_instance" "example_vm" {
  name             = "example-v"
  image_uuid       = data.zstack_images.centos.images[0].uuid
  l3_network_uuids = [data.zstack_l3networks.l3networks.l3networks[0].uuid]
  description      = "jumper server"
  #  instance_offering_uuid = data.zstack_instance_offers.offer.instance_offers[0].uuid #using Instance offering uuid or custom cpu and memory 
  memory_size = 4096
  cpu_num     = 4
  expunge     = true
}

output "zstack_instance" {
  value = zstack_instance.example_vm
}