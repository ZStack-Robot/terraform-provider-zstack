terraform {
  required_providers {
    local = {
      source = "hashicorp/local"
      version = "2.4.1"
    }
    zstack = {
      source = "zstack.io/terraform-provider-zstack/zstack"
    }
    template = {
      source = "hashicorp/template"
      version = "2.2.0"
    }
  }
}

provider "zstack" {
  host  =  "172.25.16.104"
  accountname = "admin"
  accountpassword = "password"
  accesskeyid = "mO6W9gzCQxsfK6OsE7dg"
  accesskeysecret = "Z1B3KVQlGqeaxcpeP55M3WxpRPDUyLqsppp1aLms"
}

provider "local" {

}

provider "template" {
  # Configuration options
}
data "zstack_images" "images" {
    name_regex = "RDS-3.13.10"
}

resource "zstack_vm" "vm" {
  count = 3
  name = "RDS-${count.index + 1}"
  description = "chi test"
  imageuuid = data.zstack_images.images.images.0.uuid #"${data.zstack_images.images.images[0].uuid}" #"9b26312501614ec0b6dc731e6977dfb2"
  l3networkuuids = "de7f26a7304d45aea9e9871a1ba7dbae"
  rootdiskofferinguuid = "e6ed934030244c7c8465975f7a23ae79"
  rootdisksize = 102400
  memorysize = 21474836480
  cupnum = 16
}

#data "zstack_images" "images" {}

resource "local_file" "hosts_cfg" {
  content = templatefile("./inventory.tpl",
    {
      master = [zstack_vm.vm.0.ip]
      mastershadow = [zstack_vm.vm.1.ip, zstack_vm.vm.2.ip]
    }
  )
  filename = "./inventory.ini"
}

resource "terraform_data" "remote-exec" {

  depends_on = [zstack_vm.vm]
  connection {
    type     = "ssh"
    user     = "root"
    password = "Cljslrl0620!"
    host     = zstack_vm.vm.0.ip
  }

  provisioner "file" {
    source = "./inventory.ini"
    destination = "/gaea/installer/inventory.ini"
  }

  provisioner "remote-exec" {
   inline = [
     "cd /gaea/installer/",
     "ansible-playbook -i inventory.ini install.yml"
    ]
  }
}

output "zstack_vm" {
   value =  zstack_vm.vm.*.ip
}


