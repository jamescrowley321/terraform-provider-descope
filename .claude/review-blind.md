## Review: Blind Hunter

### MUST FIX
No findings.

### SHOULD FIX
- [`internal/resources/inboundapp.go:16`] The constant `inboundAppEntity` is changed from `"inbound_application"` to `"inbound_app"`. The diff shows no corresponding state migration / `MoveResourceState` / upgrade-state handler. If this constant feeds the Terraform resource type name, every existing user state file referencing `descope_inbound_application` will break on next plan with no automated migration path — users will be forced to manually `terraform state mv` every instance. This is a breaking change to a public resource name without a visible migration shim.
- [`internal/resources/inboundapp.go` (file rename)] The file is renamed `inbound_application.go` -> `inboundapp.go`, dropping the underscore, while sibling resources still keep `outbound_application.go`/`third_party_application.go` style names (per the docs in this same diff). Inconsistent file naming will confuse grep-based navigation and future maintainers.
- [`internal/provider/provider.go:51-52`] The new `Description` and the `DeprecationMessage` overlap and partially contradict: Description says the field is a usable "fallback", DeprecationMessage says it will be removed in a future major. Pick one canonical phrasing — duplicate guidance in two surfaces tends to drift over time. Also: the diff shows only this one attribute; verify no other attribute flags (e.g., `Sensitive`) were inadvertently dropped from the surrounding lines that are not visible in this hunk.

### NITPICK
- [`docs/resources/inbound_app.md`] Page title is now `descope_inbound_app` but the body prose still says "inbound OIDC/OAuth 2.0 application" — cosmetic inconsistency.
- [`README.md:30`] The support table is updated to the new name but the diff contains no CHANGELOG / migration-note entry for users discovering their state no longer matches the resource name.
- [`internal/provider/provider.go:52`] `DeprecationMessage` text duplicates the trailing sentence of `Description`. Could deduplicate.
