---
page_title: "descope_third_party_application Resource - descope"
subcategory: ""
description: |-
  Manages a Descope third-party application.
---

# descope_third_party_application (Resource)

Manages a Descope third-party application that authenticates against Descope as an OAuth/OIDC provider. Third-party applications represent external services that use Descope for user authentication and authorization.

## Example Usage

```hcl
resource "descope_third_party_application" "partner_app" {
  name                   = "Partner Portal"
  description            = "OAuth client for partner system"
  login_page_url         = "https://partner.example.com/login"
  approved_callback_urls = ["https://partner.example.com/callback"]
}
```

## Schema

### Required

- `name` (String) The application name.

### Optional

- `approved_callback_urls` (Set of String) Approved OAuth callback URLs for this application.
- `description` (String) Application description.
- `login_page_url` (String) URL where the login page is hosted.
- `logo` (String) Application logo URL.

### Read-Only

- `client_id` (String) The OAuth client ID assigned by Descope.
- `client_secret` (String, Sensitive) The OAuth client secret, only available at creation time. Not returned by subsequent API reads.
- `id` (String) The application identifier.

~> **Note:** `client_secret` is generated server-side and only returned when the resource is created. If Terraform state is lost, the secret cannot be recovered. Use `terraform import` to re-import the resource, but the secret will not be available.

## Import

```shell
terraform import descope_third_party_application.example "app-id-here"
```
