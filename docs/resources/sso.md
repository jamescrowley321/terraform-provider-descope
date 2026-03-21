---
page_title: "descope_sso Resource - descope"
subcategory: ""
description: |-
  Manages SSO configuration for a Descope tenant.
---

# descope_sso (Resource)

Manages SSO configuration for a Descope tenant. Each SSO configuration is scoped to a tenant and supports either OIDC or SAML authentication. Multiple SSO configurations per tenant are supported via `sso_id`.

~> **Note:** SSO configurations managed by this resource may conflict with SSO settings defined in `descope_project.authentication.sso`. Avoid managing the same tenant's SSO in both places.

## Example Usage

### OIDC SSO

```terraform
resource "descope_tenant" "acme" {
  name = "Acme Corp"
}

resource "descope_sso" "acme_oidc" {
  tenant_id    = descope_tenant.acme.id
  display_name = "Acme OIDC"

  oidc {
    name          = "Acme SSO"
    client_id     = "my-client-id"
    client_secret = var.oidc_client_secret
    auth_url      = "https://idp.acme.com/authorize"
    token_url     = "https://idp.acme.com/token"
    user_data_url = "https://idp.acme.com/userinfo"
    scope         = ["openid", "profile", "email"]
  }
}
```

### SAML SSO (manual)

```terraform
resource "descope_sso" "partner_saml" {
  tenant_id    = descope_tenant.partner.id
  display_name = "Partner SAML"

  saml {
    idp_url       = "https://idp.partner.com/sso"
    idp_entity_id = "https://idp.partner.com"
    idp_cert      = file("idp-cert.pem")

    attribute_mapping {
      email = "user.email"
      name  = "user.name"
    }
  }
}
```

### SAML SSO (metadata URL)

```terraform
resource "descope_sso" "vendor_saml" {
  tenant_id    = descope_tenant.vendor.id
  display_name = "Vendor SAML"

  saml_metadata {
    idp_metadata_url = "https://idp.vendor.com/metadata"
  }
}
```

## Schema

### Required

- `tenant_id` (String) The tenant ID this SSO configuration belongs to. Changing this forces a new resource.

### Optional

- `display_name` (String) Display name for this SSO configuration. Defaults to `""`.
- `domains` (Set of String) Domains used to map users to this tenant when authenticating via SSO.
- `oidc` (Block) OIDC SSO configuration. Conflicts with `saml` and `saml_metadata`. See [OIDC](#oidc) below.
- `saml` (Block) SAML SSO configuration (manual). Conflicts with `oidc` and `saml_metadata`. See [SAML](#saml) below.
- `saml_metadata` (Block) SAML SSO configuration via metadata URL. Conflicts with `oidc` and `saml`. See [SAML Metadata](#saml-metadata) below.
- `sso_id` (String) Custom SSO configuration ID. Auto-generated if omitted. Changing this forces a new resource.

### Read-Only

- `id` (String) Composite identifier in the format `tenantID/ssoID`.

### OIDC

- `name` (String, Required) Display name for the OIDC provider.
- `client_id` (String, Required) OAuth client ID.
- `client_secret` (String, Optional, Sensitive) OAuth client secret.
- `auth_url` (String) Authorization endpoint URL.
- `callback_domain` (String) Callback domain override.
- `grant_type` (String) OAuth grant type.
- `issuer` (String) Token issuer.
- `jwks_url` (String) JWKS endpoint URL.
- `manage_provider_tokens` (Boolean) Whether to manage provider tokens. Defaults to `false`.
- `redirect_url` (String) Redirect URL after authentication.
- `scope` (Set of String) OAuth scopes to request.
- `token_url` (String) Token endpoint URL.
- `user_data_url` (String) User info endpoint URL.

### SAML

- `idp_url` (String, Required) Identity provider SSO URL.
- `idp_entity_id` (String, Required) Identity provider entity ID.
- `idp_cert` (String, Required) Identity provider certificate (PEM).
- `redirect_url` (String) Redirect URL after authentication.
- `sp_entity_id` (String, Read-Only) Service provider entity ID (computed by Descope).
- `sp_acs_url` (String, Read-Only) Service provider ACS URL (computed by Descope).
- `attribute_mapping` (Block) Maps IDP attributes to Descope user fields. All fields are optional strings: `name`, `given_name`, `middle_name`, `family_name`, `picture`, `email`, `phone_number`, `group`.

### SAML Metadata

- `idp_metadata_url` (String, Required) URL to the IDP's SAML metadata.
- `redirect_url` (String) Redirect URL after authentication.
- `sp_entity_id` (String, Read-Only) Service provider entity ID (computed by Descope).
- `sp_acs_url` (String, Read-Only) Service provider ACS URL (computed by Descope).
- `attribute_mapping` (Block) Same as SAML attribute_mapping above.

## Import

SSO configurations can be imported using `tenantID/ssoID`:

```shell
terraform import descope_sso.example "T2VxG0YjHhLW5K3iFqMo/default"
```
