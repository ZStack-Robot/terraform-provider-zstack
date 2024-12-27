#  Copyright (c) ZStack.io, Inc.

data "zstack_backupstorages" "example" {
  #   name  = "backupstorage name"
  #   name_pattern = "image%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  # optional support schema  filter
  filter {
    name   = "status"
    values = ["Connected"]
  }
  filter {
    name   = "state"
    values = ["Enabled"]
  }
  filter {
    name   = "total_capacity"
    values = ["7999424823296"]
  }
}

output "zstack_imagestorages" {
  value = data.zstack_backupstorage.example
}


