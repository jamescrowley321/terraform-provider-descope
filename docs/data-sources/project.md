---
page_title: "descope_project Data Source - descope"
subcategory: ""
description: |-
  Reads a Descope project's configuration.
---

# descope_project (Data Source)

Reads the full configuration of a Descope project, including authentication methods, roles, permissions, connectors, flows, and more.

This data source is useful for referencing an existing project's settings without managing them, or for exporting configuration details for use in other resources or modules.

## Example Usage

```terraform
data "descope_project" "current" {
  id = "P2abc123def456"
}

output "project_name" {
  value = data.descope_project.current.name
}

output "project_environment" {
  value = data.descope_project.current.environment
}
```

## Schema

### Required

- `id` (String) The Descope project ID to look up.

### Read-Only

All attributes from the [`descope_project`](../resources/project) resource are available as computed read-only values, including:

- `name` (String) The name of the project.
- `environment` (String) The project environment (`production` or unset).
- `tags` (Set of String) Tags assigned to the project.
- `project_settings` (Object) General project settings including token expiration and session behavior.
- `authentication` (Object) Authentication method configuration.
- `authorization` (Object) Roles and permissions configuration.
- `connectors` (Object) Third-party service integrations.
- `applications` (Object) OIDC and SAML application configuration.
- `flows` (Map) Authentication flow definitions.

See the [`descope_project` resource](../resources/project) for the full schema reference.
