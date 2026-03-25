---
page_title: "descope_list Resource - descope"
subcategory: ""
description: |-
  Manages a Descope list for IP or text-based filtering.
---

# descope_list (Resource)

Manages a Descope list for IP allowlisting/denylisting or text-based filtering. Lists can contain IP addresses/ranges or text entries used in authentication flows and access policies.

## Example Usage

### IP Allowlist

```terraform
resource "descope_list" "ip_allowlist" {
  name        = "Production IP Allowlist"
  description = "Allowed IPs for production access"
  type        = "ips"
  data        = ["192.168.1.0/24", "10.0.0.1"]
}
```

### Text Denylist

```terraform
resource "descope_list" "blocked_domains" {
  name        = "Blocked Email Domains"
  description = "Email domains to block during signup"
  type        = "texts"
  data        = ["disposable.email", "tempmail.com"]
}
```

## Schema

### Required

- `name` (String) The list name. Must be unique per project.
- `type` (String) The list type. One of `ips` or `texts`. Cannot be changed after creation.

### Optional

- `data` (Set of String) The list entries. IP addresses/CIDR ranges for `ips` type, or text values for `texts` type.
- `description` (String) List description.

### Read-Only

- `id` (String) The list identifier.

## Import

```shell
terraform import descope_list.example "list-id-here"
```
