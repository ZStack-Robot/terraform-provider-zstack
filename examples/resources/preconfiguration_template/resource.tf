# Copyright (c) ZStack.io, Inc.

resource "zstack_preconfiguration_template" "example" {
  name         = "example-preconfiguration-template"
  distribution = "ubuntu"
  type         = "cloudinit"
  content      = "#!/bin/bash\necho 'Preconfiguration template'"
}

output "zstack_preconfiguration_template" {
  value = zstack_preconfiguration_template.example
}
