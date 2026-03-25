variable "name" { type = string }

resource "descope_list" "test" {
  name        = var.name
  description = "Test IP list"
  type        = "ips"
  data        = ["192.168.1.1", "10.0.0.0/8"]
}
