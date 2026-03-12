variable "name" {
  type = string
}

resource "descope_management_key" "test" {
  name          = var.name
  description   = "With permitted IPs"
  permitted_ips = ["192.168.1.0/24", "10.0.0.1"]

  rebac = {
    company_roles = ["company-full-access"]
  }
}

output "id" {
  value = descope_management_key.test.id
}

output "status" {
  value = descope_management_key.test.status
}
