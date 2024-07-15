# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    zstack = {
      source = "zstack.io/terraform-provider-zstack/zstack"
    }
  }
}

provider "zstack" {
  host            = "test-qa.zstack.io"
  accesskeyid     = "46Ca208N3yLLn1JfkQR3"
  accesskeysecret = "JjdhscQnoN1uqS7VmZKZGuUmYQz6rFu35fh3hvcu"
}


data "zstack_zsclusters" "clusters" {
   name_regex = "Cluster-1"
}



output "zstack_clusters" {
  value = data.zstack_zsclusters.clusters
}


