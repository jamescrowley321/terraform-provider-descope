variable "name" { type = string }

resource "descope_tenant" "test" {
  name = var.name

  settings = {
    session_settings_enabled      = true
    refresh_token_expiration      = 30
    refresh_token_expiration_unit = "days"
    session_token_expiration      = 10
    session_token_expiration_unit = "minutes"
  }
}
