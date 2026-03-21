variable "name" { type = string }

resource "descope_tenant" "test" {
  name = var.name
}

resource "descope_sso" "test" {
  tenant_id    = descope_tenant.test.id
  display_name = "Test OIDC SSO"

  oidc {
    name      = "Test OIDC"
    client_id = "test-client-id"
    auth_url  = "https://idp.example.com/authorize"
    token_url = "https://idp.example.com/token"
  }
}
