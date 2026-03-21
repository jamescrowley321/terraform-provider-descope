variable "name" { type = string }

resource "descope_permission" "read" {
  name        = "${var.name}-read"
  description = "Read permission"
}

resource "descope_permission" "write" {
  name        = "${var.name}-write"
  description = "Write permission"
}

resource "descope_role" "test" {
  name        = var.name
  description = "Test role"

  permission_names = [
    descope_permission.read.name,
    descope_permission.write.name,
  ]
}
