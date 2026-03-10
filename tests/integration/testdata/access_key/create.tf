variable "name" {
  type = string
}

resource "descope_access_key" "test" {
  name       = var.name
  role_names = ["Tenant Admin"]
}

output "id" {
  value = descope_access_key.test.id
}

output "status" {
  value = descope_access_key.test.status
}
