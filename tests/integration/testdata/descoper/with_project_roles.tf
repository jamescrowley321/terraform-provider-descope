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
    project_roles = [
      {
        project_ids = [descope_project.test.id]
        role        = "developer"
      }
    ]
  }
}

output "id" {
  value = descope_descoper.test.id
}
