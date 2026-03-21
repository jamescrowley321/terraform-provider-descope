variable "name" { type = string }

resource "descope_flow" "test" {
  flow_id    = var.name
  definition = jsonencode({
    metadata = {
      name = var.name
    }
    contents = {
      startTask = "0"
      tasks = {
        "0" = {
          action = "logged-in"
          id     = "0"
          name   = "End"
          type   = "automated"
          next   = {}
          view   = { x = 0, y = 0 }
        }
      }
    }
    screens = {}
  })
}
