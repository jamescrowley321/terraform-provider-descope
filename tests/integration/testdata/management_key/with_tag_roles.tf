variable "name" {
  type = string
}

resource "descope_management_key" "test" {
  name        = var.name
  description = "With tag roles"

  rebac = {
    tag_roles = [
      {
        tags  = ["production", "staging"]
        roles = ["tag-infra-read-write"]
      }
    ]
  }
}

output "id" {
  value = descope_management_key.test.id
}

output "status" {
  value = descope_management_key.test.status
}
