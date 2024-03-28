# Copyright (c) HashiCorp, Inc.

variable "zstack_mn_db_host" {
  type        = string
  description = "ZStacK MN Mariadb HOST for example 1.1.1.1:3306"
}

variable "zstack_mn_db_database" {
  type    = string
  default = "keycloak"
}

variable "zstack_mn_db_username" {
  type    = string
  default = "keycloak"
}

variable "zstack_mn_db_password" {
  type      = string
  default   = "ZStack@123"
  sensitive = true
}

variable "zstack_idp_admin" {
  type    = string
  default = "admin"
}

variable "zstack_idp_password" {
  type      = string
  default   = "password"
  sensitive = true
}

