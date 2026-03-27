variable "name" { type = string }

resource "descope_tenant" "test" {
  name                      = var.name
  self_provisioning_domains = ["${var.name}.example.com"]
  enforce_sso               = true
}
