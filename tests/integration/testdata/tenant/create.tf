variable "name" { type = string }

resource "descope_tenant" "test" {
  name = var.name
}
