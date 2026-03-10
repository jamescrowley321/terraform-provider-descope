variable "name" {
  type = string
}

resource "descope_project" "test" {
  name = var.name
  tags = ["integration-test"]
}

output "id" {
  value = descope_project.test.id
}

output "name" {
  value = descope_project.test.name
}
