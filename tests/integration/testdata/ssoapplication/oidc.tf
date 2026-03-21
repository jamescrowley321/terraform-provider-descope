variable "name" { type = string }

resource "descope_sso_application" "test" {
  name        = var.name
  description = "Test OIDC SSO Application"

  oidc = {
    login_page_url = "https://app.example.com/login"
  }
}
