data "zstack_disks" "example" {
  #   name = "DATA-for-Zaku-1"
  name_pattern = "%Zaku%" # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter {
    name   = "type"
    values = ["Data"]
  }
  filter {
    name   = "is_shareable"
    values = [true]
  }
}

output "zstack_disks" {
  value = data.zstack_disks.example
}