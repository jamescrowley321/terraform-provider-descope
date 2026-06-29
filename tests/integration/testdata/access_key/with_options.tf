variable "name" {
  type = string
}

resource "descope_project" "test" {
  name = "${var.name}-proj"

  authorization = {
    roles = [
      { name = "Viewer" }
    ]
  }
}

resource "descope_access_key" "test" {
  project_id  = descope_project.test.id
  name        = var.name
  description = "Test access key"
  roles       = ["Viewer"]

  permitted_ips = ["192.168.1.0/24"]

  custom_claims = jsonencode({ claim1 = "value1" })
}

output "id" {
  value = descope_access_key.test.id
}

output "status" {
  value = descope_access_key.test.status
}

output "project_id" {
  value = descope_access_key.test.project_id
}
