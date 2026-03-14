variable "name" { type = string }

resource "descope_project" "test" {
  name = var.name
}

data "descope_project" "test" {
  id = descope_project.test.id
}
