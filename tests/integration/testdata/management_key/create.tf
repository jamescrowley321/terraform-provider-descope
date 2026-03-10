variable "name" {
  type = string
}

resource "descope_management_key" "test" {
  name = var.name

  rebac = {
    company_roles = ["Company Admin"]
  }
}

output "id" {
  value = descope_management_key.test.id
}

output "status" {
  value = descope_management_key.test.status
}
