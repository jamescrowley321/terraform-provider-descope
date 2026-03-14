variable "name" { type = string }
variable "tenant_id" { type = string }

resource "descope_tenant" "test" {
  tenant_id = var.tenant_id
  name      = var.name
}
