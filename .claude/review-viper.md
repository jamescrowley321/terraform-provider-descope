## Review: Red Team (Viper)

### Attack Surface

This PR is a cosmetic rename of the `descope_inbound_application` resource type to `descope_inbound_app`, plus a deprecation annotation on the provider-level `project_id` attribute. The change set is overwhelmingly documentation:

- `README.md`, `docs/index.md`, `docs/architecture/*.md`, `docs/resources/inbound_application.md` -> `inbound_app.md`: text-only renames. No code paths affected.
- `internal/resources/inboundapp.go` (renamed from `inbound_application.go`): the only behavioral change is one constant value: `inboundAppEntity = "inbound_application"` -> `"inbound_app"`. This constant is the `TypeName` suffix (`resp.TypeName = req.ProviderTypeName + "_" + inboundAppEntity`) and is also passed as the entity string to `r.client.Create/Read/Update/Delete`. No auth, token, middleware, RBAC, or session handling is touched. No new endpoints. No change to request signing, key handling, or transport.
- `internal/provider/provider.go`: adds a `DeprecationMessage` string and rewords the `project_id` description. The runtime resolution logic for `projectID` (env var fallback + explicit value) is unchanged; the field remains `Optional: true`, still not `Sensitive` (consistent with prior state — it is a non-secret identifier). Deprecation messages are surfaced by Terraform Core as plan-time warnings; they do not alter the auth model.
- `tools/testacc/resource.go`: test helper string updated to match the new entity name. Test-only.

No middleware ordering changes. No new endpoints. No token-handling changes. No new dependencies. No infrastructure / Docker / network changes. No changes to `infra.Client` (which carries the management key) or to how the management key (`Sensitive: true`) is read. No issuer-format handling touched. No RBAC / tenant claim logic touched.

The rename is a breaking change for consumers, but that is a migration/stability concern, not a security concern. It does not weaken authn/authz or expand attacker capability.

### Findings

No exploitable findings.

Considered and dismissed:

- **Resource-type shadowing / typosquatting**: resource type names are namespaced under the configured provider source; a third party cannot register `descope_inbound_application` for the `descope` source. Upstream Descope republishing the old name is a supply-chain concern outside this PR.
- **State drift -> ghost resource**: orphaned state for the old type does not change server-side credentials' trust boundary. Operational risk only.
- **`project_id` deprecation message info disclosure**: text contains only field name and migration guidance. `project_id` is a public identifier (appears in JWT `iss`). No info leak.
- **Behavioral change to `project_id` resolution**: `Optional`, absence of `Sensitive`, and the `Configure()` resolution logic at `provider.go:50-54` and `provider.go:87-90` are identical pre/post-patch. `DeprecationMessage` only emits a plan-time warning diagnostic.
- **Test helper rename**: acceptance-test tooling only, not compiled into release binary.

### Summary
- Attack surface elements: 0 meaningful (1 entity-name constant + 1 cosmetic schema description/deprecation; remainder is docs)
- Findings: 0 critical, 0 high, 0 medium, 0 low
- Overall: PASS
