terraform {
  required_providers {
    descope = { source = "jamescrowley321/descope" }
  }
}
provider "descope" {}

resource "descope_project" "service" {
  name = "Service identities"
}
resource "descope_permission" "jobs_read" {
  project_id = descope_project.service.id
  name       = "jobs:read"
}
resource "descope_role" "worker" {
  project_id  = descope_project.service.id
  name        = "worker"
  permissions = [descope_permission.jobs_read.name]
}
resource "descope_jwt_template" "worker" {
  project_id = descope_project.service.id
  name       = "Worker tokens"
  type       = "key"
  template = jsonencode({
    purpose = "jobs-api"
    version = 1
  })
}
resource "descope_session_settings" "service" {
  project_id              = descope_project.service.id
  access_key_jwt_template = descope_jwt_template.worker.id
}
resource "descope_access_key" "worker" {
  project_id = descope_project.service.id
  name       = "jobs-worker"
  roles      = [descope_role.worker.name]
  custom_claims = jsonencode({
    service = "jobs-worker"
    enabled = true
    retries = 3
    scopes  = ["jobs:read"]
  })
  depends_on = [descope_session_settings.service]
}
output "client_id" {
  value = descope_access_key.worker.client_id
}
output "client_secret" {
  value     = descope_access_key.worker.cleartext
  sensitive = true
}
