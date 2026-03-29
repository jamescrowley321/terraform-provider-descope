---
page_title: "descope_outbound_application Resource - descope"
subcategory: ""
description: |-
  Manages a Descope outbound application.
---

# descope_outbound_application (Resource)

Manages a Descope outbound application for OAuth integrations with external services. Outbound applications allow Descope to act as an OAuth client, enabling users to connect their accounts with third-party services.

## Example Usage

```hcl
resource "descope_outbound_application" "github" {
  name              = "GitHub Integration"
  description       = "OAuth integration with GitHub"
  client_id         = "github-client-id"
  client_secret     = var.github_client_secret
  authorization_url = "https://github.com/login/oauth/authorize"
  token_url         = "https://github.com/login/oauth/access_token"
  default_scopes    = ["read:user", "user:email"]
}
```

## Schema

### Required

- `name` (String) The application name.

### Optional

- `authorization_url` (String) OAuth authorization endpoint URL.
- `callback_domain` (String) Callback domain override.
- `client_id` (String) OAuth client ID.
- `client_secret` (String, Sensitive) OAuth client secret. Write-only — not returned by the API.
- `default_redirect_url` (String) Default redirect URL after authorization.
- `default_scopes` (Set of String) Default OAuth scopes to request.
- `description` (String) Application description.
- `discovery_url` (String) OIDC discovery URL for auto-configuration.
- `logo` (String) Application logo URL.
- `pkce` (Boolean) Enable PKCE for the OAuth flow. Defaults to `false`.
- `revocation_url` (String) Token revocation endpoint URL.
- `token_url` (String) OAuth token endpoint URL.

### Read-Only

- `id` (String) The application identifier.

## Import

```shell
terraform import descope_outbound_application.example "app-id-here"
```
