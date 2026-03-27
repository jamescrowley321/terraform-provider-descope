
> **Community Fork** — This is an independently maintained fork of [descope/terraform-provider-descope](https://github.com/descope/terraform-provider-descope).
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

The upstream Descope Terraform provider supports project-level configuration through `descope_project`, plus a handful of company-level resources. This fork extends coverage to the full Management API surface, adding 11 new resources and 4 data sources.

### Resource Coverage

| Resource | Type | Upstream | This Fork |
|----------|------|:--------:|:---------:|
| `descope_project` | Resource | Yes | Yes |
| `descope_descoper` | Resource | Yes | Yes |
| `descope_management_key` | Resource | Yes | Yes |
| `descope_inbound_application` | Resource | Yes | Yes |
| `descope_tenant` | Resource | - | **Yes** |
| `descope_access_key` | Resource | - | **Yes** |
| `descope_role` | Resource | - | **Yes** |
| `descope_permission` | Resource | - | **Yes** |
| `descope_sso` | Resource | - | **Yes** |
| `descope_sso_application` | Resource | - | **Yes** |
| `descope_third_party_application` | Resource | - | **Yes** |
| `descope_outbound_application` | Resource | - | **Yes** |
| `descope_fga_schema` | Resource | - | **Yes** |
| `descope_list` | Resource | - | **Yes** |
| `descope_password_settings` | Resource | - | **Yes** |
| `data.descope_project` | Data Source | - | **Yes** |
| `data.descope_password_settings` | Data Source | - | **Yes** |
| `data.descope_project_export` | Data Source | - | **Yes** |
| `data.descope_fga_check` | Data Source | - | **Yes** |

See the [open issues](https://github.com/jamescrowley321/terraform-provider-descope/issues) for the full roadmap.

<br/>

## About

Use this Terraform provider to manage your [Descope](https://www.descope.com) project using Terraform configuration files.

- Create and manage projects with settings, auth methods, connectors, and flows.
- Manage tenants, access keys, roles, permissions, and SSO configuration as standalone resources.
- Configure fine-grained authorization (FGA) schemas and IP/text allow/deny lists.
- Create and manage SSO, third-party, outbound, and inbound applications.
- Read project data, password settings, and FGA authorization checks via data sources.

> **Note:** Users are intentionally excluded from this provider. Users are runtime entities created through authentication flows, not infrastructure — managing them via Terraform would cause perpetual state drift and is an anti-pattern. Use the [Descope SDK](https://docs.descope.com) or console for user management.

<br/>

## Getting Started

### Requirements

- [Terraform CLI](https://developer.hashicorp.com/terraform/install)
- A Descope account (free tier works for most resources; pro/enterprise required for SSO applications and project export)
- A [management key](https://app.descope.com/settings/company) for your Descope company

### Usage

> **Note:** Until this fork is published to a Terraform registry, you must build and install the provider locally.
> See [Development](#development) below.

Configure the provider and declare resources. The management key and other provider settings can also be set via environment variables (`DESCOPE_MANAGEMENT_KEY`, `DESCOPE_PROJECT_ID`, `DESCOPE_BASE_URL`).

```hcl
provider "descope" {
  project_id     = "P..."
  management_key = "K..."
}
```

<br/>

## Examples

### Tenants

Create and manage tenants with self-provisioning domains and SSO enforcement:

```hcl
resource "descope_tenant" "production" {
  name                      = "Acme Corp"
  self_provisioning_domains = ["acme.com"]
  enforce_sso               = true

  settings = {
    session_settings_enabled      = true
    refresh_token_expiration      = 30
    refresh_token_expiration_unit = "days"
  }
}
```

### Roles and Permissions

Manage authorization as standalone resources with explicit dependencies:

```hcl
resource "descope_permission" "build_apps" {
  name        = "build-apps"
  description = "Allowed to build and sign applications"
}

resource "descope_permission" "deploy" {
  name        = "deploy"
  description = "Allowed to deploy to production"
}

resource "descope_role" "developer" {
  name             = "Developer"
  description      = "Builds and deploys applications"
  permission_names = [
    descope_permission.build_apps.name,
    descope_permission.deploy.name,
  ]
}
```

### Access Keys

Create machine-to-machine access keys with IP restrictions and custom claims:

```hcl
resource "descope_access_key" "ci_deploy" {
  name          = "CI Deploy Key"
  description   = "Used by GitHub Actions for deployments"
  role_names    = ["Tenant Admin"]
  permitted_ips = ["192.168.1.0/24"]

  custom_claims = {
    environment = "production"
  }
}
```

### SSO Configuration

Configure OIDC SSO for a tenant:

```hcl
resource "descope_sso" "okta" {
  tenant_id    = descope_tenant.production.id
  display_name = "Okta SSO"

  oidc = {
    name          = "Okta"
    client_id     = "0oa..."
    client_secret = var.okta_client_secret
    auth_url      = "https://company.okta.com/oauth2/v1/authorize"
    token_url     = "https://company.okta.com/oauth2/v1/token"
    user_data_url = "https://company.okta.com/oauth2/v1/userinfo"

    attribute_mapping = {
      login_id = "sub"
      email    = "email"
      name     = "name"
    }
  }
}
```

### Fine-Grained Authorization

Define an FGA schema and check authorization:

```hcl
resource "descope_fga_schema" "authz" {
  schema = <<-EOT
    model AuthZ 1.0

    type user

    type document
      relation owner: user
      relation viewer: user
  EOT
}

data "descope_fga_check" "can_view" {
  resource    = "document:readme"
  relation    = "viewer"
  target      = "user:alice"
  depends_on  = [descope_fga_schema.authz]
}
```

### IP/Text Allow and Deny Lists

Manage lists for IP filtering or text-based rules:

```hcl
resource "descope_list" "blocked_ips" {
  name        = "Blocked IPs"
  description = "IPs blocked from authentication"
  type        = "ips"
  data        = ["192.0.2.1", "198.51.100.0/24"]
}
```

### Password Settings

Configure project-wide password policy:

```hcl
resource "descope_password_settings" "policy" {
  enabled          = true
  min_length       = 12
  lowercase        = true
  uppercase        = true
  number           = true
  non_alphanumeric = true
  expiration       = true
  expiration_weeks = 26
  lock             = true
  lock_attempts    = 5
}
```

### Project Settings

Create a project with custom session configuration and authorization:

```hcl
resource "descope_project" "my_project" {
  name = "My Project"

  project_settings = {
    refresh_token_expiration = "3 weeks"
    session_token_expiration = "1 hour"
    refresh_token_rotation   = true
  }
}
```

<br/>

## Development

See the [README](internal/README.md) file in the `internal` directory for architecture details.

### Setup

```bash
git clone https://github.com/jamescrowley321/terraform-provider-descope
cd terraform-provider-descope
make dev
```

This builds the provider binary, installs it to `$GOPATH/bin`, and creates a `~/.terraformrc` with `dev_overrides` so Terraform uses your local build.

### Build and Test

```bash
make install            # Rebuild and install the provider
make test               # Unit tests (no credentials needed)
make testintegration    # Integration tests against live Descope API
make testacc            # Acceptance + unit tests
make lint               # Linting + secret detection
make testcleanup        # Delete leftover testacc-* resources
```

Integration tests require a `.env` file with credentials (see `.env.example`):

```bash
source .env && make testintegration
```

### Testing

This project uses a multi-layered testing strategy with adversarial SDK verification. Integration tests don't just check Terraform state — they call the Descope API directly to verify resources actually exist with the correct field values. See [docs/testing.md](docs/testing.md) for the full breakdown of every validation gate.

<br/>

## Contributing

Contributions are welcome, including AI-assisted submissions. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

<br/>

## Relationship to Upstream

This fork tracks the upstream [descope/terraform-provider-descope](https://github.com/descope/terraform-provider-descope) repository. Upstream changes are merged periodically to stay current. New resources developed here may be proposed back to the upstream project via pull requests.

### Attribution

This project is based on the work of [Descope](https://www.descope.com) and the original provider authors. The original source code is licensed under the [MIT License](LICENSE).

<br/>

## Support

#### Learn more

To learn more please see the [Descope documentation](https://docs.descope.com).

#### Issues

If anything is missing or not working correctly please [open an issue](https://github.com/jamescrowley321/terraform-provider-descope/issues).

#### Descope Community

For general Descope questions (not specific to this fork) you can use the [Slack community](https://www.descope.com/community) or contact [Descope support](mailto:support@descope.com).
