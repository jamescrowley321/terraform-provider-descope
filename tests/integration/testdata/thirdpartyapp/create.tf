variable "name" { type = string }

resource "descope_third_party_application" "test" {
  name        = var.name
  description = "Test third-party application"
}
