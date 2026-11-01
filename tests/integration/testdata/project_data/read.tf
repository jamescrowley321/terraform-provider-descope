variable "name" { type = string }

resource "descope_project" "test" {
  deletion_protection = false
  name                = var.name
}

data "descope_project" "test" {
  id = descope_project.test.id
}
