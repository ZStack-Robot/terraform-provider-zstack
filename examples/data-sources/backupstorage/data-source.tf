# Copyright (c) ZStack.io, Inc.


data "zstack_backupstorages" "example" {
 #   name  = "backupstorage name"
 #   name_pattern = "image%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
}   

output "zstack_imagestorages" {
  value = data.zstack_backupstorage.example
}


