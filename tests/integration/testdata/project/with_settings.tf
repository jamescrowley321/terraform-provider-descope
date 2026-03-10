variable "name" {
  type = string
}

resource "descope_project" "test" {
  name = var.name

  project_settings = {
    refresh_token_expiration = "3 weeks"
    session_token_expiration = "1 hour"
    refresh_token_rotation   = true
  }
}

output "id" {
  value = descope_project.test.id
}

output "name" {
  value = descope_project.test.name
}
