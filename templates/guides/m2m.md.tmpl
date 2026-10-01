---
page_title: "Machine-to-machine authentication"
description: |-
  Manage service access keys, authorization, and custom token claims.
---

# Machine-to-machine authentication

Use the [complete example](https://github.com/jamescrowley321/terraform-provider-descope/tree/main/examples/m2m)
to create a project, permission, role, access-key JWT template, session settings, and service key.
Set `DESCOPE_MANAGEMENT_KEY` before applying. Each project-scoped resource requires its own `project_id`.
The provider-level `project_id` and `DESCOPE_PROJECT_ID` are ignored.

`descope_project` manages the project container. Roles, permissions, templates, settings, and applications
are standalone resources. `descope_access_key.roles` takes role names. Use Terraform references to order
their creation. `custom_claims`, `custom_attributes`, and JWT `template` are JSON object strings: use
`jsonencode({ ... })`, including for booleans, numbers, arrays, and nested objects.

A JWT template for machine identities must use `type = "key"`. Select it with
`descope_session_settings.access_key_jwt_template`. Key-specific claims live in
`descope_access_key.custom_claims`; template claims apply to every key exchanged under that template.
The example orders key creation after session configuration.

Reads refresh claims and custom attributes from the API. A real API change causes a Terraform plan;
whitespace and key ordering do not. Setting either JSON attribute to `jsonencode({})` clears it.
Updating claims, roles, or status preserves the key ID and original `cleartext`.

## Exchange a key

Exchange `descope_access_key.worker.cleartext` through the SDK's `ExchangeAccessKey` method, or use OAuth
client credentials at `/oauth2/v1/token`. For the default OIDC application, use the access key's `client_id`
and its `cleartext` as the OAuth client secret. The live integration test covers both exchanges, including
`descope.claims` and `descope.custom_claims` scopes, signature verification, claim changes, and revocation.

A newly created access key starts active. Update `status = "inactive"` to prevent further token exchanges;
set it back to `"active"` to reactivate. Deactivation does not prove that previously issued JWTs have expired.
`cleartext` is sensitive and stored in Terraform state; Descope does not return it on later reads or imports.

## Tenant authorization

For a tenant role, set `descope_role.tenant_id = descope_tenant.example.id`. Assign tenant-scoped keys with:

```hcl
resource "descope_access_key" "customer_worker" {
  project_id = descope_project.service.id
  name       = "customer-worker"
  tenants = [{
    tenant_id = descope_tenant.example.id
    roles     = [descope_role.customer_worker.name]
  }]
}
```

Use either project `roles` or tenant assignments on a key. The live integration test verifies tenant roles
and permissions with `SelectedTenant` when exchanging the tenant key.

## Current API limits

The API used for live validation rejected `aud` overrides with `E113801`, including direct SDK requests and
a copy of the built-in access-key audience template. Custom string, boolean, numeric, and nested claims passed.
Audience customization needs validation against your Descope environment; this provider does not bypass
server template validation. The example uses the default issuer and audience behavior.

Project snapshot export excludes tenants and tenant SSO. Manage those resources explicitly.

See Descope's [access-key documentation](https://docs.descope.com/management/m2m-access-keys),
[JWT template documentation](https://docs.descope.com/management/token/jwt-templates), and
[OIDC client credentials documentation](https://docs.descope.com/identity-federation/applications/oidc-apps/oidc-endpoints).
