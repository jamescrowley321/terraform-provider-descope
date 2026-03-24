---
page_title: "descope_fga_schema Resource - descope"
subcategory: ""
description: |-
  Manages the FGA schema for a Descope project.
---

# descope_fga_schema (Resource)

Manages the Fine-Grained Authorization (FGA) schema for a Descope project. The schema defines object types and their relations for relationship-based access control (ReBAC), similar to Google Zanzibar.

This is a singleton resource — only one FGA schema exists per project. Creating the resource saves the schema; destroying it clears the schema from the project.

## Example Usage

```terraform
resource "descope_fga_schema" "main" {
  schema = jsonencode({
    types = {
      document = {
        relations = {
          owner  = { this = {} }
          editor = { union = [{ this = {} }, { computedUserset = { relation = "owner" } }] }
          viewer = { union = [{ this = {} }, { computedUserset = { relation = "editor" } }] }
        }
      }
    }
  })
}
```

## Schema

### Required

- `schema` (String) The FGA schema as a JSON string. Use `jsonencode()` to define the schema in HCL. The schema defines object types and their relations.

### Read-Only

- `id` (String) The resource identifier (always `fga_schema`).

~> **Warning:** Destroying this resource clears the FGA schema from the project. This may break authorization checks that depend on the schema.
