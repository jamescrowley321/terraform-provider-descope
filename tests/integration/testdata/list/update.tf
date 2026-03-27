variable "name" { type = string }

resource "descope_list" "test" {
  name        = var.name
  description = "Updated IP list"
  type        = "ips"
  data        = ["192.168.1.1", "10.0.0.0/8", "172.16.0.0/12"]
}
