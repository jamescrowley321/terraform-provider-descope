variable "name" {
  type = string
}

variable "email" {
  type = string
}

resource "descope_descoper" "test" {
  email = var.email
  name  = var.name
  phone = "+15551234567"

  rbac = {
    is_company_admin = true
  }
}

output "id" {
  value = descope_descoper.test.id
}

output "email" {
  value = descope_descoper.test.email
}

output "name" {
  value = descope_descoper.test.name
}

output "phone" {
  value = descope_descoper.test.phone
}
