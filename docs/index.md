---
page_title: "descope Provider"
description: |-
  Use the Descope Terraform Provider to manage your Descope project's authentication methods, flows, roles, permissions, connectors, and more as infrastructure-as-code.
---

# Descope Provider

The [Descope](https://www.descope.com) Terraform Provider lets you manage your Descope project configuration as infrastructure-as-code. Configure authentication methods, define roles and permissions, set up third-party connectors, manage flows, and more—all declaratively in Terraform.

Descope is an authentication and user management platform. The Terraform provider manages _project configuration_ (how your project behaves), not users or tenants (use the [Descope Management API](https://docs.descope.com/api/openapi) or [SDKs](https://docs.descope.com) for those).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/install) 1.0 or later
- A Descope account (Community tier works; some features require a **Pro or Enterprise** plan)
- A **Management Key** from [Company Settings](https://app.descope.com/settings/company) in the Descope console

## Authentication

The provider authenticates with the Descope API using a management key. You can create management keys in the [Company Settings](https://app.descope.com/settings/company) section of the Descope console.

~> **Security Note:** Never hardcode your management key in Terraform configuration files, as this risks exposing it in version control. Use environment variables or a secrets manager instead.

### Environment Variables

| Variable                | Description                                              |
|-------------------------|----------------------------------------------------------|
| `DESCOPE_MANAGEMENT_KEY` | A valid management key for your Descope company         |
| `DESCOPE_BASE_URL`       | Override the Descope API base URL (optional, for testing) |

```shell
export DESCOPE_MANAGEMENT_KEY="K2..."
terraform plan
```

## Example Usage

### Minimal Configuration

```hcl
terraform {
  required_providers {
    descope = {
      source  = "jamescrowley321/descope"
      version = "~> 1.0"
    }
  }
}

# Credentials are read from the DESCOPE_MANAGEMENT_KEY environment variable
provider "descope" {}
```

### Creating a Project

```hcl
resource "descope_project" "example" {
  name        = "my-app"
  environment = "production"
  tags        = ["prod"]
}
```

### Full Provider Configuration

If you need to explicitly configure the provider (e.g. in a module or when not using environment variables):

```hcl
terraform {
  required_providers {
    descope = {
      source  = "jamescrowley321/descope"
      version = "~> 1.0"
    }
  }
}

variable "descope_management_key" {
  type      = string
  sensitive = true
}

provider "descope" {
  management_key = var.descope_management_key
}
```

## Resources

| Resource | Description |
|----------|-------------|
| [`descope_project`](resources/project) | Full project configuration: authentication, RBAC, connectors, flows, and more |
| [`descope_access_key`](resources/access_key) | API access keys for programmatic authentication |
| [`descope_management_key`](resources/management_key) | Programmatic management keys with scoped permissions |
| [`descope_tenant`](resources/tenant) | Multi-tenant environments with session and SSO configuration |
| [`descope_permission`](resources/permission) | Standalone permissions for role-based access control |
| [`descope_role`](resources/role) | Roles that group permissions for RBAC |
| [`descope_descoper`](resources/descoper) | Descope console user accounts with role assignments |
| [`descope_sso`](resources/sso) | SSO configuration for a project |
| [`descope_sso_application`](resources/sso_application) | SSO application definitions |
| [`descope_inbound_application`](resources/inbound_application) | Inbound application integrations |
| [`descope_outbound_application`](resources/outbound_application) | Outbound application integrations |
| [`descope_third_party_application`](resources/third_party_application) | Third-party application integrations |
| [`descope_password_settings`](resources/password_settings) | Password authentication policy configuration |
| [`descope_fga_schema`](resources/fga_schema) | Fine-grained authorization schema |
| [`descope_list`](resources/list) | Custom lists (IP allowlists, text lists, JSON data) |

## Data Sources

| Data Source | Description |
|-------------|-------------|
| [`descope_project`](data-sources/project) | Read an existing project's full configuration |
| [`descope_password_settings`](data-sources/password_settings) | Read current password authentication settings |
| [`descope_project_export`](data-sources/project_export) | Export a project's configuration |
| [`descope_fga_check`](data-sources/fga_check) | Check fine-grained authorization permissions |

## Guides

- [Quickstart](guides/quickstart) – Set up the provider and manage your first project

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `base_url` (String) An optional base URL for the Descope API
- `management_key` (String, Sensitive) A valid management key for your Descope company
- `project_id` (String, Deprecated)


