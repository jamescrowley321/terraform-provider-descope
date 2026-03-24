variable "name" { type = string }

resource "descope_fga_schema" "test" {
  schema = "model\n  schema 1.1\n\ntype user\n\ntype document\n  relations\n    define owner: [user]"
}
