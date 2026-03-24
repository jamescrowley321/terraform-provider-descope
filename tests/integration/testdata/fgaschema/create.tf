variable "name" { type = string }

resource "descope_fga_schema" "test" {
  schema = "model\n  schema 1.1\ntype user\ntype document\n  relations\n    define owner: [user]"
}
