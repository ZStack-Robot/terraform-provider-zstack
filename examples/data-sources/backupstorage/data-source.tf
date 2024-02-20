
data "zstack_backupstorage" "imagestorage" {
  name_regex = "image"
}



output "zstack_imagestorage" {
   value =  data.zstack_backupstorage.imagestorage
}


