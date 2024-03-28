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
    name_regex = "idp2"
}

resource "zstack_vm" "vm" {
  count = 2
  name = "IDP-${count.index + 1}"
  description = "IDP test-${count.index + 1}"
  imageuuid = data.zstack_images.images.images.0.uuid #"${data.zstack_images.images.images[0].uuid}" #"9b26312501614ec0b6dc731e6977dfb2"
  l3networkuuids = "de7f26a7304d45aea9e9871a1ba7dbae"
  rootdiskofferinguuid = "e6ed934030244c7c8465975f7a23ae79"
  rootdisksize = 102400
  memorysize = 2147483648
  cupnum = 8
}

resource "local_file" "hosts_cfg" {
  content = templatefile("./docker-compose.tpl",
    {
      cluster = [zstack_vm.vm.0.ip, zstack_vm.vm.1.ip]
      db_url = var.zstack_mn_db_host
      db_name = var.zstack_mn_db_database
      db_user = var.zstack_mn_db_username
      db_pwd = var.zstack_mn_db_password
      idp_admin = var.zstack_idp_admin
      idp_password = var.zstack_idp_password
    }
  )
  filename = "./docker-compose.yaml"
}

resource "terraform_data" "remote-exec" {

  depends_on = [zstack_vm.vm]
  connection {
    type     = "ssh"
    user     = "root"
    password = "password"
    host     = zstack_vm.vm.0.ip
  }

  provisioner "file" {
    source = "./docker-compose.yaml"
    destination = "/root/docker/docker-compose.yaml"
  }

  provisioner "remote-exec" {
   inline = [
     "cd /root/docker/",
     "sh install.sh"
    ]
  }
}

output "zstack_vm" {
   value =  zstack_vm.vm.*.ip
}


