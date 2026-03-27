variable "name" { type = string }

resource "descope_list" "test" {
  name        = var.name
  description = "Updated IP list"
  type        = "ips"
  data        = ["192.0.2.1", "198.51.100.0/24", "203.0.113.0/24"]
}
