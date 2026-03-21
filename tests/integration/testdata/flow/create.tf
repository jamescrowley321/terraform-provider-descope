variable "name" { type = string }

resource "descope_flow" "test" {
  flow_id    = var.name
  definition = jsonencode({
    flow = {
      id   = var.name
      name = var.name
      type = "custom"
    }
    screens = []
  })
}
