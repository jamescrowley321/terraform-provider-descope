---
page_title: "descope_sso_application Resource - descope"
subcategory: ""
description: |-
  Manages a Descope SSO application (OIDC or SAML).
---

# descope_sso_application (Resource)

Manages a standalone Descope SSO application. SSO applications define how external service providers authenticate users via your Descope project, using either OIDC or SAML protocols.

~> **Note:** SSO applications managed by this resource may conflict with applications defined in `descope_project.applications`. Avoid managing the same application in both places.

## Example Usage

### OIDC Application

```hcl
resource "descope_sso_application" "my_app" {
  name        = "My OIDC App"
  description = "OIDC SSO application for my service"

  oidc = {
    login_page_url = "https://app.example.com/login"
  }
}
```

### SAML Application

```hcl
resource "descope_sso_application" "saml_app" {
  name        = "My SAML App"
  description = "SAML SSO application"

  saml = {
    login_page_url  = "https://app.example.com/login"
    use_metadata_info = true
    metadata_url    = "https://app.example.com/saml/metadata"
  }
}
```

## Schema

### Required

- `name` (String) The SSO application name. Must be unique.

### Optional

- `description` (String) A description of the SSO application. Defaults to `""`.
- `enabled` (Boolean) Whether the application is enabled. Defaults to `true`.
- `logo` (String) URL to the application logo. Defaults to `""`.
- `oidc` (Block) OIDC application settings. Conflicts with `saml`. See below.
- `saml` (Block) SAML application settings. Conflicts with `oidc`. See below.

### Read-Only

- `app_type` (String) The application type (`oidc` or `saml`).
- `id` (String) The application identifier.

### OIDC Settings

- `back_channel_logout_url` (String) Back-channel logout URL. Defaults to `""`.
- `force_authentication` (Boolean) Force re-authentication. Defaults to `false`.
- `login_page_url` (String) URL of the login page. Defaults to `""`.

### SAML Settings

- `acs_url` (String) SP Assertion Consumer Service URL.
- `certificate` (String) SP certificate for signed requests.
- `default_relay_state` (String) Default relay state value.
- `entity_id` (String) SP entity ID.
- `force_authentication` (Boolean) Force re-authentication. Defaults to `false`.
- `login_page_url` (String) URL of the login page.
- `logout_redirect_url` (String) Redirect URL after logout.
- `metadata_url` (String) SP metadata URL (when `use_metadata_info` is true).
- `use_metadata_info` (Boolean) Fetch SP info from metadata URL. Defaults to `false`.

## Import

SSO applications can be imported by ID:

```shell
terraform import descope_sso_application.example "app-id-here"
```
