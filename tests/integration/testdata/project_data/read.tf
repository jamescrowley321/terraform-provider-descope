variable "project_id" { type = string }

data "descope_project" "test" {
  id = var.project_id
}
