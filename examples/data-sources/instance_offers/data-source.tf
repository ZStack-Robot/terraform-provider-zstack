# Copyright (c) ZStack.io, Inc.

data "zstack_instance_offers" "example" {

}

output "zstack_instance_offers" {
  value = data.zstack_instance_offers.example
}


