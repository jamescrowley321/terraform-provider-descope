**Community Fork** — This is an independently maintained fork of [descope/terraform-provider-descope](https://github.com/descope/terraform-provider-descope).
> It is **not** an official Descope product. For the official provider, see the [upstream repository](https://github.com/descope/terraform-provider-descope).

<div align="center">
  <h3 align="center">Descope Terraform Provider (Community Fork)</h3>

  <p align="center">
    A community-maintained Terraform provider for managing Descope projects with extended resource coverage
  </p>

  <p align="center">
    <a href="https://securityscorecards.dev/viewer/?uri=github.com/jamescrowley321/terraform-provider-descope"><img src="https://api.securityscorecards.dev/projects/github.com/jamescrowley321/terraform-provider-descope/badge" alt="OpenSSF Scorecard"></a>
  </p>
</div>

<br />

## Why This Fork?

The official Descope Terraform provider supports project-level configuration through the `descope_project` resource. However, the Descope Management API exposes many additional entities that are not yet available as Terraform resources.

This fork aims to close that gap by adding standalone resources for the full Management API surface:

| Resource | Upstream | This Fork |
|----------|:--------:|:---------:|
| Project settings, auth methods, connectors, flows | Yes | Yes |
| Roles & permissions (nested in project) | Yes | Yes |
| Applications (nested in project) | Yes | Yes |
| **Tenants** (standalone CRUD) | No | Planned |
| **Access Keys** (M2M) | No | Planned |
| **SSO Configuration** (per-tenant) | No | Planned |
| **SSO Applications** (standalone) | No | Planned |
| **Fine-Grained Authorization** (FGA) | No | Planned |
| **Third-Party Applications** | No | Planned |
| **Outbound Applications** | No | Planned |
| **Password Settings** (standalone) | No | Planned |
| **Standalone Roles & Permissions** | No | Planned |
| **Standalone Flows** | No | Planned |
| **Project Export/Import** | No | Planned |

See the [open issues](https://github.com/jamescrowley321/terraform-provider-descope/issues) for the full roadmap.

<br/>

## About

Use this Terraform provider to manage your [Descope](https://www.descope.com) project using Terraform configuration files.

* Modify project settings and authentication methods.
* Create connectors, roles, permissions, applications and other entities.
* Use custom themes and flows created in the Descope console.
* Ensure dependencies between entities are satisfied.
* **Manage tenants, access keys, SSO, and more as standalone resources** (coming soon).

> **Note:** Users are intentionally excluded from this provider. Users are runtime entities created through authentication flows, not infrastructure — managing them via Terraform would cause perpetual state drift and is an anti-pattern. Use the [Descope SDK](https://docs.descope.com) or console for user management.

<br/>

## Getting Started

### Requirements

-   [Terraform](https://www.terraform.io/downloads.html) >= 1.0
-   [Go](https://golang.org/doc/install) >= 1.21 (to build the provider)

### Installation

To use this provider in your Terraform configuration, add the following block:

```hcl
terraform {
  required_providers {
    descope = {
      source = "jamescrowley321/descope"
    }
  }
}

provider "descope" {
  project_id = "YOUR_PROJECT_ID"
  management_key = "YOUR_MANAGEMENT_KEY"
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change. Please ensure tests are updated as appropriate.