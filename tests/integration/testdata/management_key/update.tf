variable "name" {
  type = string
}

resource "descope_management_key" "test" {
  name        = var.name
  status      = "inactive"
  description = "Updated via integration test"

  rebac {
    company_roles = ["Company Admin"]
  }
}

output "id" {
  value = descope_management_key.test.id
}

output "status" {
  value = descope_management_key.test.status
}

output "description" {
  value = descope_management_key.test.description
}
