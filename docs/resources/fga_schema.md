---
page_title: "descope_fga_schema Resource - descope"
subcategory: ""
description: |-
  Manages the FGA schema for a Descope project.
---

# descope_fga_schema (Resource)

Manages the Fine-Grained Authorization (FGA) schema for a Descope project. The schema defines object types, relations, and permissions for relationship-based access control (ReBAC).

This is a singleton resource — only one FGA schema exists per project. Creating or updating the resource saves the schema to Descope. Destroying the resource removes it from Terraform state only — the schema remains active on the project because the Descope API does not support deleting an FGA schema.

## Example Usage

```hcl
resource "descope_fga_schema" "main" {
  schema = <<-EOT
model AuthZ 1.0

type user

type document
  relation owner: user
  relation editor: user
  permission can_edit: editor | owner
  permission can_view: editor | owner
EOT
}
```

## Schema

### Required

- `schema` (String) The FGA authorization model in [Descope AuthZ DSL](https://docs.descope.com/authorization/rebac/define-schema) format. Use a heredoc (`<<-EOT`) to define multi-line schemas inline. The DSL must start with `model AuthZ 1.0`.

### Read-Only

- `id` (String) The resource identifier (always `fga_schema`).

~> **Note:** Destroying this resource only removes it from Terraform state. The FGA schema remains active on the Descope project because the API does not support schema deletion. To re-manage the schema, import it with `terraform import`.
