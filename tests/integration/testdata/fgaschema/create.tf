variable "name" { type = string }

resource "descope_fga_schema" "test" {
  schema = jsonencode({
    types = {
      (var.name) = {
        relations = {
          owner = {
            this = {}
          }
        }
      }
    }
  })
}
