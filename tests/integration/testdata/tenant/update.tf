variable "project_id" { type = string }
variable "name" { type = string }

resource "descope_tenant" "test" {
  project_id                = var.project_id
  name                      = var.name
  self_provisioning_domains = ["${var.name}.example.com"]
  enforce_sso               = true
}
