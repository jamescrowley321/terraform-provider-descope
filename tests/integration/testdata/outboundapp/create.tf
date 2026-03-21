variable "name" { type = string }

resource "descope_outbound_application" "test" {
  name              = var.name
  description       = "Test outbound application"
  authorization_url = "https://external.example.com/authorize"
  token_url         = "https://external.example.com/token"
}
