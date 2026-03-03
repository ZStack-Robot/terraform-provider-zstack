data "zstack_hook_scripts" "test" {

}

data "zstack_hook_scripts" "pattern_match" {
  name_pattern = "es%"
}

output "hook_scripts_test" {
  value = data.zstack_hook_scripts.test
}

output "hook_scripts_pattern" {
  value = data.zstack_hook_scripts.pattern_match
}