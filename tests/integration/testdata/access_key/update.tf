variable "name" {
  type = string
}

resource "descope_project" "test" {
  name = "${var.name}-proj"
}

resource "descope_access_key" "test" {
  project_id  = descope_project.test.id
  name        = var.name
  status      = "inactive"
  description = "Updated via integration test"
}

output "id" {
  value = descope_access_key.test.id
}

output "status" {
  value = descope_access_key.test.status
}

output "description" {
  value = descope_access_key.test.description
}

output "project_id" {
  value = descope_access_key.test.project_id
}
