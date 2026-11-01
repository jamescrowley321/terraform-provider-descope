variable "project_id" { type = string }
variable "name" { type = string }
variable "tenant_id" { type = string }

resource "descope_tenant" "test" {
  project_id = var.project_id
  tenant_id  = var.tenant_id
  name       = var.name
}
