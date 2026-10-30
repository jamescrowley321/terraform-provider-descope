variable "project_id" { type = string }
variable "name" { type = string }

resource "descope_tenant" "test" {
  project_id = var.project_id
  name       = var.name
}
