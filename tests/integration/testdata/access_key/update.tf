variable "name" {
  type = string
}

resource "descope_access_key" "test" {
  name        = var.name
  status      = "inactive"
  description = "Updated via integration test"
  role_names  = ["Tenant Admin"]
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
