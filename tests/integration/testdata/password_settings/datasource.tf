variable "name" { type = string }
resource "descope_project" "test" {
  name                = var.name
  deletion_protection = false
}
resource "descope_password_settings" "test" {
  project_id       = descope_project.test.id
  min_length       = 10
  non_alphanumeric = false
}
data "descope_password_settings" "test" {
  project_id = descope_project.test.id
  depends_on = [descope_password_settings.test]
}
