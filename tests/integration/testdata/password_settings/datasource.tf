resource "descope_password_settings" "test" {
  enabled          = true
  min_length       = 12
  lowercase        = true
  uppercase        = true
  number           = true
  non_alphanumeric = true
  expiration       = true
  expiration_weeks = 26
  reuse            = true
  reuse_amount     = 5
  lock             = true
  lock_attempts    = 5
}

data "descope_password_settings" "test" {
  depends_on = [descope_password_settings.test]
}
