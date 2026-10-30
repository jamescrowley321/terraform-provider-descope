variable "name" { type = string }
variable "service" { default = "worker" }
variable "status" { default = "active" }
variable "clear" { default = false }

resource "descope_project" "test" {
  name                = var.name
  deletion_protection = false
}
resource "descope_permission" "read" {
  project_id = descope_project.test.id
  name       = "read:jobs"
}
resource "descope_role" "service" {
  project_id  = descope_project.test.id
  name        = "service"
  permissions = [descope_permission.read.name]
}
resource "descope_tenant" "test" {
  project_id = descope_project.test.id
  name       = "customer"
}
resource "descope_role" "tenant_service" {
  project_id  = descope_project.test.id
  tenant_id   = descope_tenant.test.id
  name        = "tenant-service"
  permissions = [descope_permission.read.name]
}
resource "descope_jwt_template" "service" {
  project_id = descope_project.test.id
  name       = "service"
  type       = "key"
  template = jsonencode({
    purpose = "terraform-m2m-test"
    version = 1
  })
}
resource "descope_session_settings" "service" {
  project_id              = descope_project.test.id
  access_key_jwt_template = descope_jwt_template.service.id
}
resource "descope_access_key_attribute" "owner" {
  project_id = descope_project.test.id
  id         = "owner"
  name       = "Owner"
  type       = "string"
}
resource "descope_access_key" "service" {
  project_id        = descope_project.test.id
  name              = "service"
  status            = var.status
  roles             = var.clear ? [] : [descope_role.service.name]
  custom_attributes = var.clear ? jsonencode({}) : jsonencode({ (descope_access_key_attribute.owner.id) = var.service })
  custom_claims = var.clear ? jsonencode({}) : jsonencode({
    service = var.service
    enabled = true
    retries = 3
    nested  = { scopes = ["read", "write"] }
  })
  depends_on = [descope_session_settings.service]
}
resource "descope_access_key" "tenant_service" {
  project_id = descope_project.test.id
  name       = "tenant-service"
  tenants = [{
    tenant_id = descope_tenant.test.id
    roles     = [descope_role.tenant_service.name]
  }]
  depends_on = [descope_session_settings.service]
}
