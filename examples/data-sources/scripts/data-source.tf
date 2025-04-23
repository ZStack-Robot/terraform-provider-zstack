data "zstack_scripts" "all_scripts" {

}

data "zstack_scripts" "script_by_name" {
  name = "test-script"
}

data "zstack_scripts" "scripts_by_pattern" {
  name_pattern = "test%"
}

data "zstack_scripts" "filtered_scripts" {
  filter {
    name   = "platform"
    values = ["Windows"]
  }

  filter {
    name   = "script_type"
    values = ["Powershell"]
  }
}


# Output examples
output "all_scripts" {
  value = data.zstack_scripts.all_scripts.scripts
}

output "script_by_name" {
  value = data.zstack_scripts.script_by_name.scripts
}

output "linux_scripts_count" {
  value = length(data.zstack_scripts.filtered_scripts.scripts)
}

# Output a specific script's content
output "first_script_content" {
  value = length(data.zstack_scripts.all_scripts.scripts) > 0 ? data.zstack_scripts.all_scripts.scripts[0].script_content : "No scripts found"
}