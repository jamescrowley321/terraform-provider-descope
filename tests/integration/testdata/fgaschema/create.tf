variable "name" { type = string }

resource "descope_fga_schema" "test" {
  schema = jsonencode({
    schema_version = "1.1"
    type_definitions = [
      { type = "user" },
      {
        type = "document"
        relations = {
          owner = {
            this = {}
          }
        }
        metadata = {
          relations = {
            owner = {
              directly_related_user_types = [
                { type = "user" }
              ]
            }
          }
        }
      }
    ]
  })
}
