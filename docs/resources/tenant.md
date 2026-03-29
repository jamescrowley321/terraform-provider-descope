---
page_title: "descope_tenant Resource - descope"
subcategory: ""
description: |-
  Manages a Descope tenant for multi-tenant authentication.
---

# descope_tenant (Resource)

Manages a Descope tenant. Tenants are isolated environments within a Descope project for managing users, authentication methods, and configurations in multi-tenant applications.

~> **Note:** Tenants managed by this resource may conflict with tenants defined in `descope_project`. Avoid managing the same tenant in both places.

## Example Usage

### Basic Tenant

```hcl
resource "descope_tenant" "example" {
  name = "Acme Corp"
}
```

### Tenant with Custom ID and Self-Provisioning

```hcl
resource "descope_tenant" "example" {
  tenant_id                 = "acme-corp"
  name                      = "Acme Corp"
  self_provisioning_domains = ["acme.com", "acme.io"]
  enforce_sso               = true
  default_roles             = ["user"]
}
```

### Tenant with Session Settings

```hcl
resource "descope_tenant" "example" {
  name = "Acme Corp"

  settings = {
    session_settings_enabled      = true
    refresh_token_expiration      = 30
    refresh_token_expiration_unit = "days"
    session_token_expiration      = 10
    session_token_expiration_unit = "minutes"
    enable_inactivity             = true
    inactivity_time               = 30
    inactivity_time_unit          = "minutes"
  }
}
```

### Child Tenant with Role Inheritance

```hcl
resource "descope_tenant" "parent" {
  name = "Parent Org"
}

resource "descope_tenant" "child" {
  name             = "Child Team"
  parent_tenant_id = descope_tenant.parent.id
  role_inheritance  = "userOnly"
  default_roles     = ["developer"]
}
```

## Schema

### Required

- `name` (String) The human-readable name of the tenant.

### Optional

- `cascade_delete` (Boolean) When `true`, deleting the tenant also deletes all associated users and data. When `false` (default), the delete fails if the tenant has associated resources.
- `custom_attributes` (Map of String) Custom key-value attributes for the tenant.
- `default_roles` (Set of String) Roles automatically assigned to new users created in this tenant.
- `disabled` (Boolean) Whether the tenant is disabled. Disabled tenants cannot be accessed.
- `enforce_sso` (Boolean) Whether to enforce single sign-on (SSO) for all users in the tenant.
- `enforce_sso_exclusions` (Set of String) Email addresses or patterns excluded from SSO enforcement.
- `parent_tenant_id` (String) ID of a parent tenant for hierarchical multi-tenancy. Cannot be changed after creation.
- `role_inheritance` (String) Controls role inheritance in hierarchical tenants. Valid values: `"none"`, `"userOnly"`.
- `self_provisioning_domains` (Set of String) Email domains that allow self-provisioning. Users with matching email addresses can automatically create accounts.
- `settings` (Attributes) Session and token configuration for the tenant. (see [below for nested schema](#nestedatt--settings))
- `tenant_id` (String) Custom tenant ID. If not set, the system generates one automatically. Cannot be changed after creation.

### Read-Only

- `auth_type` (String) The authentication type configured for the tenant.
- `created_time` (Number) Unix timestamp when the tenant was created.
- `domains` (Set of String) Domains associated with the tenant.
- `id` (String) The ID of this resource.

<a id="nestedatt--settings"></a>
### Nested Schema for `settings`

Optional:

- `enable_inactivity` (Boolean) Enable inactivity timeout enforcement.
- `inactivity_time` (Number) Inactivity timeout duration.
- `inactivity_time_unit` (String) Unit for inactivity timeout. Valid values: `"seconds"`, `"minutes"`, `"hours"`, `"days"`, `"weeks"`.
- `jit_disabled` (Boolean) Disable just-in-time user provisioning.
- `refresh_token_expiration` (Number) Refresh token TTL value.
- `refresh_token_expiration_unit` (String) Unit for refresh token expiration. Valid values: `"seconds"`, `"minutes"`, `"hours"`, `"days"`, `"weeks"`.
- `session_settings_enabled` (Boolean) Enable custom session management settings for this tenant.
- `session_token_expiration` (Number) Session token TTL value.
- `session_token_expiration_unit` (String) Unit for session token expiration. Valid values: `"seconds"`, `"minutes"`, `"hours"`, `"days"`, `"weeks"`.
- `stepup_token_expiration` (Number) Step-up token TTL value.
- `stepup_token_expiration_unit` (String) Unit for step-up token expiration. Valid values: `"seconds"`, `"minutes"`, `"hours"`, `"days"`, `"weeks"`.

## Import

Tenants can be imported using the tenant ID:

```shell
terraform import descope_tenant.example <tenant_id>
```
