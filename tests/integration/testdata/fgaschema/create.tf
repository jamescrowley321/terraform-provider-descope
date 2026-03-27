variable "name" { type = string }

resource "descope_fga_schema" "test" {
  schema = "model AuthZ 1.0\n\ntype user\n\ntype document\n  relation owner: user"
}
