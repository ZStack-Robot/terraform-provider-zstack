# Copyright (c) ZStack.io, Inc.

resource "zstack_preconfiguration_template" "example" {
  name         = "example-preconfiguration-template"
  distribution = "CentOS7"
  type         = "kickstart"
  content      = <<-EOT
    # Base ZStack system variables required by the API:
    # REPO_URL
    # USERNAME
    # PASSWORD
    # NETWORK_CFGS
    # FORCE_INSTALL
    # PRE_SCRIPTS
    # POST_SCRIPTS
    install
    text
  EOT
}

output "zstack_preconfiguration_template" {
  value = zstack_preconfiguration_template.example
}
