# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    zstack = {
      source = "zstack.io/terraform-provider-zstack/zstack"
    }
  }
}

provider "zstack" {
  host            = "ip"
  accesskeyid     = "accesskeyid"
  accesskeysecret = "accesskeysecret"
}


data "zstack_zsclusters" "clusters" {
#   name_regex = "Cluster-1"
}



output "zstack_clusters" {
  value = data.zstack_zsclusters.clusters.0.uuid
}


