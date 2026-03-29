---
page_title: "descope_permission Resource - descope"
subcategory: ""
description: |-
  Manages a standalone Descope permission.
---

# descope_permission (Resource)

Manages a standalone Descope permission. Permissions can be referenced by `descope_role` resources to build role-based access control.

~> **Note:** Permissions managed by this resource may conflict with permissions defined in `descope_project.authorization`. Avoid managing the same permission in both places.

## Example Usage

```hcl
resource "descope_permission" "read_data" {
  name        = "read:data"
  description = "Allows reading data resources"
}

resource "descope_permission" "write_data" {
  name        = "write:data"
  description = "Allows writing data resources"
}
```

## Schema

### Required

- `name` (String) The unique name of the permission.

### Optional

- `description` (String) A brief description of what this permission allows. Defaults to `""`.

### Read-Only

- `id` (String) The permission identifier (same as `name`).

## Import

Permissions can be imported using the permission name:

```shell
terraform import descope_permission.example "read:data"
```
