## Review: Acceptance Auditor

### PASS
- [AC-1] Rename `descope_inbound_application` resource back to `descope_inbound_app` — `inboundAppEntity` constant set to `"inbound_app"` at `internal/resources/inboundapp.go:16`; resource type name derived via `req.ProviderTypeName + "_" + inboundAppEntity` at `internal/resources/inboundapp.go:40` (yields `descope_inbound_app`); source file renamed from `internal/resources/inbound_application.go` to `internal/resources/inboundapp.go`; docs file renamed from `docs/resources/inbound_application.md` to `docs/resources/inbound_app.md`. Tested at `internal/models/inboundapp/inboundapp_test.go:10-12` via `testacc.InboundApp(t)` helper which now resolves to `inbound_app` (tier: unit; full build/test pass per AC-5 confirms wiring).
- [AC-2] Fork-only docs updated — all five files cleanly switched from `descope_inbound_application` to `descope_inbound_app`: `README.md:30`, `docs/index.md:102`, `docs/architecture/auth-flows.md:162`, `docs/architecture/descope-model.md:51`, and `docs/architecture/oauth2-oidc-mapping.md` (lines 9-10, 40-41, 63, 70). Repo-wide grep for `inbound_application` returns zero hits in source/docs.
- [AC-3] Internal tooling string updated — `tools/testacc/resource.go:34` now calls `newResource(t, "inbound_app")` instead of `"inbound_application"`. Tested at `internal/models/inboundapp/inboundapp_test.go:12` which calls `testacc.InboundApp(t)`; unit test suite passes (tier: unit).
- [AC-4] `DeprecationMessage` re-added on provider-level `project_id` — `internal/provider/provider.go:53` adds `DeprecationMessage`; field remains `Optional: true` (line 51) and behavior is preserved: Configure still reads it at lines 87-90 and passes to `infra.NewProviderData(..., projectID)` at line 109. Field is not removed, no behavioral change — signaling only.
- [AC-5] No build/test regression — `go build ./...` exits clean with no output; `go test ./...` shows all packages with tests passing.

### FAIL
(none)

### PARTIAL
(none)

### SCOPE CREEP
(none) — every changed line traces to one of AC-1 through AC-4. The provider attribute description at `internal/provider/provider.go:52` was rewritten in the same hunk as the `DeprecationMessage` addition; this is consistent with AC-4 since the description elaborates the same deprecation message and adds no behavior.

### Summary
- Total ACs: 5
- Pass: 5 | Fail: 0 | Partial: 0
