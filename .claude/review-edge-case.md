## Review: Edge Case Hunter

### Findings

No unhandled edge cases found.

The diff is a pure rename refactor plus a schema metadata addition:

- `README.md`, `docs/architecture/auth-flows.md`, `docs/architecture/descope-model.md`, `docs/architecture/oauth2-oidc-mapping.md`, `docs/index.md`, `docs/resources/inbound_app.md` (renamed from `inbound_application.md`): documentation-only string replacements. No runtime code paths.
- `internal/provider/provider.go` (lines 50-54): adds a `DeprecationMessage` field on the existing `project_id` schema attribute and rewords its `Description`. The `Configure` function body is unchanged — no new branches, no new validation, no behavior change.
- `internal/resources/inboundapp.go` (renamed from `inbound_application.go`, line 16): changes the const `inboundAppEntity` from `"inbound_application"` to `"inbound_app"`. The Create/Read/Update/Delete/ImportState bodies are byte-identical to the prior file; all pre-existing error paths (`infra.AsValidationError`, generic `err != nil`, `entity.Diagnostics.HasError`, the malformed import-ID rejection at line 160) are preserved.
- `tools/testacc/resource.go` (line 34): test helper string literal rename only.

No new branching paths, no new boundary conditions, no new async boundaries, no new type boundaries, no new integration boundaries were introduced by this diff.

Out-of-scope observation: the entity-string change is a behavioral change against the Descope management API contract at `/v1/mgmt/infra` in `internal/infra/client.go:43-134`, where `entity` is passed verbatim to the backend. If the upstream service still expects `"inbound_application"`, every CRUD call will be rejected by the API — but that rejection flows through the existing `err != nil` handlers, so there is still no unhandled path locally.

### Summary
- Unhandled paths found: 0
- Critical (crash/data loss): 0
- Non-critical (wrong result/degraded): 0
