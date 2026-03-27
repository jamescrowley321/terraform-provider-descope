variable "name" { type = string }

resource "descope_list" "test" {
  name        = var.name
  description = "Test IP list"
  type        = "ips"
  data        = ["192.0.2.1", "198.51.100.0/24"]
}
