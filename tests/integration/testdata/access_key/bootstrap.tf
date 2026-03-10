variable "name" {
  type = string
}

resource "descope_management_key" "bootstrap" {
  name = var.name

  rebac = {
    company_roles = ["company-full-access"]
  }
}

output "cleartext" {
  value     = descope_management_key.bootstrap.cleartext
  sensitive = true
}
