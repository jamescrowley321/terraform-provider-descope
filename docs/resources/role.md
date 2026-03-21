---
page_title: "descope_role Resource - descope"
subcategory: ""
description: |-
  Manages a standalone Descope role.
---

# descope_role (Resource)

Manages a standalone Descope role. Roles group permissions together and can be assigned to users. Roles can be global or scoped to a specific tenant.

~> **Note:** Roles managed by this resource may conflict with roles defined in `descope_project.authorization`. Avoid managing the same role in both places.

## Example Usage

### Global Role

```terraform
resource "descope_permission" "read_data" {
  name        = "read:data"
  description = "Allows reading data resources"
}

resource "descope_permission" "write_data" {
  name        = "write:data"
  description = "Allows writing data resources"
}

resource "descope_role" "editor" {
  name        = "editor"
  description = "Can read and write data"

  permission_names = [
    descope_permission.read_data.name,
    descope_permission.write_data.name,
  ]
}
```

### Tenant-Scoped Role

```terraform
resource "descope_role" "tenant_admin" {
  name        = "admin"
  description = "Tenant administrator"
  tenant_id   = descope_tenant.acme.id

  permission_names = [
    descope_permission.read_data.name,
    descope_permission.write_data.name,
  ]
}
```

## Schema

### Required

- `name` (String) The unique name of the role. Must be unique within its scope (global or per-tenant).

### Optional

- `default_role` (Boolean) Whether this role is automatically assigned to new users. Defaults to `false`.
- `description` (String) A brief description of what this role allows. Defaults to `""`.
- `permission_names` (Set of String) The names of permissions included in this role.
- `private` (Boolean) Whether this role is hidden from user-facing APIs. Defaults to `false`.
- `tenant_id` (String) The tenant ID to scope this role to. If omitted, the role is global. Changing this forces a new resource.

### Read-Only

- `id` (String) The role identifier. For global roles: the role name. For tenant-scoped roles: `tenantID/name`.

## Import

Global roles can be imported by name:

```shell
terraform import descope_role.example "editor"
```

Tenant-scoped roles use the format `tenantID/name`:

```shell
terraform import descope_role.example "T2VxG0YjHhLW5K3iFqMo/admin"
```
