---
page_title: "descope_access_key Resource - descope"
subcategory: ""
description: |-
  Manages a Descope access key for programmatic API authentication.
---

# descope_access_key (Resource)

Manages a Descope access key for programmatic API authentication. Access keys allow server-to-server communication and can be scoped to specific tenants, roles, and permissions.

~> **Important:** The `cleartext` attribute (the raw key value) is only available immediately after creation and **cannot be retrieved later** through the API. Store it securely using a secrets manager (e.g., AWS Secrets Manager, HashiCorp Vault) immediately after `terraform apply`.

## Example Usage

### Basic Access Key

```hcl
resource "descope_access_key" "api_key" {
  name       = "my-api-key"
  role_names = ["Tenant Admin"]
}

# Store the cleartext key in your secrets manager
output "access_key_secret" {
  value     = descope_access_key.api_key.cleartext
  sensitive = true
}
```

### Access Key with Expiration and IP Restrictions

```hcl
resource "descope_access_key" "restricted_key" {
  name        = "ci-cd-key"
  description = "Limited access key for CI/CD pipelines"

  expire_time   = 1893456000
  permitted_ips = ["203.0.113.0/24", "198.51.100.10"]

  role_names = ["Editor"]

  custom_claims = {
    environment = "production"
  }
}
```

### Multi-Tenant Access Key

```hcl
resource "descope_access_key" "multi_tenant_key" {
  name = "multi-tenant-key"

  key_tenants = [
    {
      tenant_id = descope_tenant.primary.id
      roles     = ["Admin"]
    },
    {
      tenant_id = descope_tenant.secondary.id
      roles     = ["Viewer"]
    }
  ]
}
```

## Schema

### Required

- `name` (String) The name of the access key.

### Optional

- `custom_attributes` (Map of String) Custom attributes attached to the access key.
- `custom_claims` (Map of String) Custom claims to include in tokens issued with this key.
- `description` (String) A description for the access key.
- `expire_time` (Number) The expiration time of the access key as a Unix timestamp. If not set, the key will not expire. Changing this value after creation will require the access key to be replaced.
- `key_tenants` (Attributes List) Tenant associations for the access key. (see [below for nested schema](#nestedatt--key_tenants))
- `permitted_ips` (List of String) A list of IP addresses or CIDR ranges that are allowed to use this access key. If not set, the key can be used from any IP address.
- `role_names` (Set of String) A set of role names to assign to the access key.
- `status` (String) The status of the access key. Must be either `active` or `inactive`. Cannot be set to `inactive` on creation.
- `user_id` (String) Associates the access key with a specific user. Changing this value after creation will require the access key to be replaced.

### Read-Only

- `cleartext` (String, Sensitive) The plaintext value of the access key. This is only available after the key is created and cannot be retrieved later. Store this value securely as it is required to authenticate API requests.
- `client_id` (String) The client ID associated with this access key.
- `created_by` (String) The ID of the user who created the access key.
- `created_time` (Number) Unix timestamp when the access key was created.
- `id` (String) The ID of this resource.

<a id="nestedatt--key_tenants"></a>
### Nested Schema for `key_tenants`

Required:

- `tenant_id` (String) The ID of the tenant.

Optional:

- `roles` (Set of String) Roles to assign within this tenant.

Read-Only:

- `tenant_name` (String) The name of the tenant (populated from API).

## Known Limitations

### Token Issuer Mismatch

Tokens obtained by exchanging an access key via `POST /v1/auth/accesskey/exchange` use a different `iss` (issuer) claim than the one advertised in Descope's OIDC discovery document:

| Source | Issuer Format |
|---|---|
| OIDC Discovery (`/.well-known/openid-configuration`) | `https://api.descope.com/{project_id}` |
| Access key exchange (`/v1/auth/accesskey/exchange`) | `https://api.descope.com/v1/apps/{project_id}` |

Per the OIDC specifications (OpenID Connect Core 3.1.3.7, Discovery 1.0 Section 3), the `iss` claim in tokens **MUST exactly match** the discovery document's `issuer` value. Descope session tokens from the access key exchange API do not satisfy this requirement.

If you validate these tokens using OIDC discovery, you must either disable issuer verification or validate against the session issuer format directly. Tokens obtained via standard OAuth 2.0 flows (client credentials, authorization code) use the correct OIDC issuer.

## Import

Access keys can be imported using the access key ID:

```shell
terraform import descope_access_key.example <access_key_id>
```
