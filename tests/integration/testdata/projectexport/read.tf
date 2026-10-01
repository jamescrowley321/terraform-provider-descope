variable "project_id" { type = string }
data "descope_project_export" "test" {
  project_id = var.project_id
}
