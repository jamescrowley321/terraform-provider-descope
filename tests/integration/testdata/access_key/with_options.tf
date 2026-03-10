variable "name" {
  type = string
}

resource "descope_access_key" "test" {
  name        = var.name
  description = "Test access key"
  role_names  = ["Tenant Admin"]

  permitted_ips = ["192.168.1.0/24"]

  custom_claims = {
    "claim1" = "value1"
  }
}

output "id" {
  value = descope_access_key.test.id
}

output "status" {
  value = descope_access_key.test.status
}
