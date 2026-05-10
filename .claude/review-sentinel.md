## Review: Security (Sentinel)

### BLOCK (must fix before merge)
No BLOCK findings.

### WARN (should fix)
No WARN findings.

### INFO (acceptable risk)
- `internal/resources/inboundapp.go:16` — The `inboundAppEntity` constant is dual-purposed: it is both the Terraform resource type suffix (`descope_inbound_app`) and the `entity` discriminator sent to the Descope `/v1/mgmt/infra` API (see `internal/infra/client.go:43-134`). Renaming `inbound_application` to `inbound_app` therefore changes the on-the-wire API call, not just the Terraform schema. This is a backend-contract concern (the Descope infra API must accept `inbound_app`), but it is not a security regression: authorization is enforced by the management key (`c.managementKey`) and the per-project API client, neither of which is influenced by the entity string. The entity value is a routing discriminator on an endpoint already gated by a management key — there is no tenant- or principal-scoped check that the rename could weaken.
- `internal/provider/provider.go:50-54` — The `project_id` provider attribute gained a `DeprecationMessage` and a longer `Description`. Both are static strings rendered by Terraform; no user input flows into them and no new code path was added. The `Configure` logic at lines 87-90 and 109 that consumes `project_id` is unchanged, so the deprecation does not alter how the project ID is resolved, validated, or passed to the SDK client. No new attack surface.
- `internal/provider/provider.go:52` — The new description references `descope_access_key` and the `DESCOPE_PROJECT_ID` env var — the same information already present in the previous description. It does not disclose any secret, internal endpoint, or credential-handling detail. Descope project IDs are non-sensitive identifiers (they appear publicly in JWT issuer URLs).
- `tools/testacc/resource.go:34` — Test harness rename mirrors the production constant. Test-only; not reachable from a deployed provider.

### Summary
- BLOCK: 0 | WARN: 0 | INFO: 4
- Overall: PASS
