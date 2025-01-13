data "zstack_primary_storages" "primary_storages_test" {
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
    values = ["7937400246271"]
  }
}

output "zstack_primary_storages_test" {
  value = data.zstack_primary_storages.primary_storages_test
}