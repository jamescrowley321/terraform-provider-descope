variable "name" { type = string }

resource "descope_permission" "test" {
  name        = var.name
  description = "Updated test permission"
}
