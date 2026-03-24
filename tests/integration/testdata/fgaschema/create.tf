variable "name" { type = string }

resource "descope_fga_schema" "test" {
  schema = <<-EOT
    model
      schema 1.1
    type user
    type ${var.name}
      relations
        define owner: [user]
  EOT
}
