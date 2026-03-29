---
page_title: "descope_project_export Data Source - descope"
subcategory: ""
description: |-
  Exports a snapshot of the current Descope project configuration.
---

# descope_project_export (Data Source)

Exports a snapshot of all settings and configurations for the current Descope project. The snapshot is returned as a JSON string containing all project files.

This data source is useful for backing up project configurations, comparing environments, or feeding into project import operations.

## Example Usage

```hcl
data "descope_project_export" "current" {}

output "project_config" {
  value     = data.descope_project_export.current.files
  sensitive = true
}

# Access a specific configuration file
locals {
  config = jsondecode(data.descope_project_export.current.files)
}
```

## Schema

### Read-Only

- `files` (String) The exported project configuration as a JSON string. Use `jsondecode()` to access individual configuration files. Note: secret values (tokens, keys) are left blank in exported snapshots.
- `id` (String) The data source identifier (always `project_export`).
