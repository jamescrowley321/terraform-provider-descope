---
page_title: "descope_flow Resource - descope"
subcategory: ""
description: |-
  Manages a Descope authentication flow.
---

# descope_flow (Resource)

Manages a Descope authentication flow by importing and exporting flow definitions as JSON. Flows define the authentication experience for users (sign-up, sign-in, MFA, etc.).

~> **Note:** Flows managed by this resource may conflict with flows defined in `descope_project`. Avoid managing the same flow in both places.

## Example Usage

```terraform
resource "descope_flow" "sign_up" {
  flow_id    = "sign-up"
  definition = file("flows/sign-up.json")
}
```

### Using jsonencode

```terraform
resource "descope_flow" "custom" {
  flow_id    = "custom-flow"
  definition = jsonencode({
    flow = {
      id   = "custom-flow"
      name = "Custom Flow"
      type = "custom"
    }
    screens = []
  })
}
```

## Schema

### Required

- `definition` (String) The flow definition as a JSON string. Use `file()` to load from a file or `jsonencode()` to define inline.
- `flow_id` (String) The flow identifier. Changing this forces a new resource.

### Read-Only

- `id` (String) The flow identifier (same as `flow_id`).

## Import

Flows can be imported by flow ID:

```shell
terraform import descope_flow.example "sign-up"
```

## Notes

- On Read, the `definition` attribute is populated with the server-exported JSON, which may differ from the original input due to server normalization (additional fields, key ordering).
- Subsequent plans may show diffs if the server adds or reorders fields in the flow definition.
