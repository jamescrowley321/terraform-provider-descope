variable "name" {
  type = string
}

variable "email" {
  type = string
}

resource "descope_project" "test" {
  name = "${var.name}-proj"
  tags = ["production", "staging"]
}

resource "descope_descoper" "test" {
  email = var.email
  name  = var.name

  rbac = {
    tag_roles = [
      {
        tags = ["production", "staging"]
        role = "admin"
      }
    ]
  }
}

output "id" {
  value = descope_descoper.test.id
}
