resource "descope_password_settings" "test" {
  enabled         = true
  min_length      = 10
  lowercase       = true
  uppercase       = true
  number          = true
  non_alphanumeric = false
  expiration      = false
  lock            = false
}
