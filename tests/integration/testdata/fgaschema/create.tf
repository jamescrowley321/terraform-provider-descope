variable "project_id" { type = string }
variable "name" { type = string }

resource "descope_fga_schema" "test" {
  project_id = var.project_id
  schema     = "model AuthZ 1.0\n\ntype user\n\ntype document\n  relation owner: user"
}


data "descope_fga_check" "owner" {
  project_id    = var.project_id
  resource      = "document-1"
  resource_type = "document"
  relation      = "owner"
  target        = "user-1"
  target_type   = "user"
  depends_on    = [descope_fga_schema.test]
}
