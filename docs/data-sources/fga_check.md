---
page_title: "descope_fga_check Data Source - descope"
subcategory: ""
description: |-
  Checks whether a given FGA relation is authorized.
---

# descope_fga_check (Data Source)

Checks whether a target has the specified relation to a resource in the Descope FGA system. Returns an `allowed` boolean indicating whether the authorization check passed.

## Example Usage

```terraform
data "descope_fga_check" "can_alice_view" {
  resource_type = "document"
  resource      = "doc1"
  relation      = "viewer"
  target_type   = "user"
  target        = "alice"
}

output "alice_can_view" {
  value = data.descope_fga_check.can_alice_view.allowed
}
```

## Schema

### Required

- `relation` (String) The relation to check (e.g., `viewer`, `editor`, `owner`).
- `resource` (String) The resource identifier to check against.
- `resource_type` (String) The type of the resource (e.g., `document`, `folder`).
- `target` (String) The target identifier (e.g., user ID).
- `target_type` (String) The type of the target (e.g., `user`, `group`).

### Read-Only

- `allowed` (Boolean) Whether the target has the specified relation to the resource.
- `id` (String) A synthetic identifier derived from the relation tuple.
