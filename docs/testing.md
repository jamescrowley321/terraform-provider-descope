# Testing and Validation

This document describes every gate a code change must pass before it reaches production. The goal is defense in depth: no single layer is trusted to catch everything.

## Validation Gates

A change passes through these layers, roughly in order:

```
Local commit
  │
  ├─ pre-commit hooks ─── gitleaks, gofmt, go build, go test
  │
  ▼
Push / Pull Request
  │
  ├─ CI: build ────────── go mod tidy (no drift), go build ./...
  ├─ CI: unit tests ───── go test ./... (33 test files, no credentials)
  ├─ CI: lint ─────────── golangci-lint (24 linters + 24 revive rules)
  ├─ CI: leaks ────────── gitleaks (custom Descope key patterns + 50 built-in rules)
  ├─ CI: security ─────── gosec, govulncheck, Snyk SCA + SAST
  ├─ CI: CodeQL ───────── GitHub semantic analysis (weekly + on PR)
  ├─ CI: integration ──── terraform apply against live Descope API
  │   ├─ upstream ─────── TestProject{CRUD,Settings,Authorization}
  │   └─ fork ─────────── all fork-tagged tests (14 test files)
  ├─ CI: SDK verify ───── every integration test calls Descope SDK after apply
  └─ CI: cleanup ──────── tools/testcleanup deletes testacc-* resources
  │
  ▼
Release (v* tag)
  │
  ├─ GoReleaser ───────── multi-platform binaries
  ├─ Syft ─────────────── SBOM generation
  ├─ Cosign ───────────── keyless artifact signing (Sigstore)
  ├─ GPG ──────────────── checksums signature
  └─ SLSA ─────────────── v3 build provenance attestation
```

## Layer 1: Pre-Commit Hooks

Configured in `.pre-commit-config.yaml`. Runs automatically on every `git commit`.

| Hook | What it catches |
|---|---|
| **gitleaks** v8.30.0 | Secrets in staged changes (management keys, API tokens, private keys) |
| **go-fmt** | Formatting violations |
| **go-build-repo-mod** | Compilation failures |
| **go-test-repo-mod** | Unit test regressions |

A commit that leaks a secret or breaks the build is rejected before it exists in history.

## Layer 2: Unit Tests

**33 test files** across `internal/models/`, `internal/infra/`, and `tools/terragen/`. No credentials or external services required.

What they cover:

- **Model validation** — schema version checks, field constraints, type conversions
- **Conversion logic** — Terraform state ↔ SDK type mapping (`*_test.go` in each model package)
- **Retry and error handling** — `internal/infra/retry_test.go`, `errors_test.go`
- **Code generation utilities** — JSON parsing, string helpers for terragen

Run locally:

```bash
make test                              # all unit tests
go test -v ./internal/models/tenant/   # single package
```

## Layer 3: Linting

**golangci-lint** with 24 linters and 24 revive rules, configured in `.github/actions/ci/lint/golangci.yml`.

Key linters:

| Category | Linters |
|---|---|
| **Correctness** | govet, staticcheck, errcheck, nilerr, forcetypeassert |
| **Security** | (delegated to gosec/govulncheck — see Layer 5) |
| **Style** | revive (24 rules), gofmt, goimports, misspell |
| **Performance** | ineffassign, unconvert, wastedassign, unused, unparam |
| **Maintainability** | exhaustive, copyloopvar, contextcheck, predeclared |

Run locally:

```bash
make lint   # golangci-lint + gitleaks protect + detect
```

## Layer 4: Secret Detection

**gitleaks** runs in two modes:

1. **Pre-commit** — scans staged changes (`.pre-commit-config.yaml`)
2. **CI** — scans full diff (`.github/workflows/ci.yml`)

Custom rules in `.github/actions/ci/leaks/gitleaks.toml`:

- `descope-management-key` — matches `K` followed by 200+ alphanumeric chars
- `gitlab-pat` — matches `glpat-` tokens
- Plus 50+ built-in patterns (AWS, GitHub, SSH keys, Stripe, etc.)

## Layer 5: Security Scanning

Three scanners run on every push and PR (`.github/workflows/security.yml`):

| Scanner | Focus | Blocks PR? |
|---|---|---|
| **gosec** v2.24.7 | Go-specific vulnerabilities (SQL injection, crypto misuse, etc.) | No (advisory) |
| **govulncheck** v1.1.4 | Known CVEs in Go dependencies | Yes |
| **Snyk** (SCA + SAST) | Dependency vulns (high severity) + code-level analysis | Conditional |

Additionally:

- **CodeQL** (`.github/workflows/codeql.yml`) — GitHub's semantic analysis, runs weekly + on PRs
- **OpenSSF Scorecard** (`.github/workflows/scorecard.yml`) — supply chain security assessment, runs weekly

## Layer 6: Integration Tests

Integration tests run `terraform apply` against a live Descope account, then verify the results.

### Architecture

```
tests/integration/
├── harness.go          # Test harness — builds provider, manages workspaces
├── verify.go           # SDK verification — loads resources directly from API
├── testdata/
│   ├── provider.tf     # Shared provider config
│   ├── tenant/         # Fixtures: create.tf, update.tf, with_settings.tf, ...
│   ├── access_key/
│   ├── role/
│   └── ...             # One directory per resource type
├── tenant_test.go      # Test files — one per resource
├── access_key_test.go
└── ...
```

### Harness (`harness.go`)

The `Harness` struct manages a complete Terraform workspace per test:

1. **Builds the provider binary once** — `sync.Once` compiles the provider from source
2. **Creates an isolated workspace** — temp directory with `.terraformrc` pointing to the local binary via `dev_overrides`
3. **Exposes Terraform operations** — `Apply`, `Plan`, `Destroy`, `Import`, `StateResource`, `Output`
4. **Retries transient failures** — up to 2 retries with 5s backoff for API flakiness
5. **Cleans up on test end** — `t.Cleanup` runs `terraform destroy`; waits for async project deletion

Key methods:

| Method | Purpose |
|---|---|
| `ApplyFixture(fixture, address, vars...)` | Load HCL fixture → apply → return state attributes |
| `ReimportResource(fixture, address, id, vars...)` | Remove from state → import → verify round-trip |
| `TryApply(vars...)` | Like Apply, but returns error instead of fataling (for license checks) |
| `StateResource(address)` | Parse `terraform show -json` and return attribute map |
| `GenerateName(t)` | `testacc-{TestName}-{MMddHHmm}-{UUID8}` — unique, cleanup-friendly |

### SDK Verification (`verify.go`)

After every `terraform apply`, tests call the Descope Management SDK directly to verify the API matches Terraform's view. This catches the class of bugs where the provider silently drops or zeros fields.

**Why this exists:** PRs #80, #86, and #94 all introduced bugs that passed the existing integration tests because those tests only checked Terraform state, not the actual API.

Client factories:

| Factory | Scope | Used for |
|---|---|---|
| `newProjectSDKClient(t)` | `DESCOPE_PROJECT_ID` | Project-scoped resources (tenant, role, access key, etc.) |
| `newCompanySDKClient(t)` | No project ID | Company-level resources (descoper, management key) |
| `newSDKClientWithProject(t, id)` | Specific project | Verifying resources in newly-created projects |

Load functions — one per resource type:

| Function | SDK method |
|---|---|
| `LoadTenantViaSDK(t, id)` | `Tenant().Load(id)` |
| `LoadTenantSettingsViaSDK(t, id)` | `Tenant().GetSettings(id)` |
| `LoadSSOApplicationViaSDK(t, id)` | `SSOApplication().Load(id)` |
| `LoadThirdPartyAppViaSDK(t, id)` | `ThirdPartyApplication().LoadApplication(id)` |
| `LoadFGASchemaViaSDK(t)` | `FGA().LoadSchema()` |
| `LoadListViaSDK(t, id)` | `List().Load(id)` |
| `LoadAccessKeyViaSDK(t, id)` | `AccessKey().Load(id)` |
| `FindPermissionViaSDK(t, name)` | `Permission().LoadAll()` → find by name |
| `FindRoleViaSDK(t, name)` | `Role().LoadAll()` → find by name |
| `LoadSSOSettingsViaSDK(t, tenantID, ssoID)` | `SSO().LoadSettings(tenantID, ssoID)` |
| `LoadOutboundAppViaSDK(t, id)` | `OutboundApplication().LoadApplication(id)` |
| `LoadPasswordSettingsViaSDK(t, tenantID)` | `Password().GetSettings(tenantID)` |
| `LoadManagementKeyViaSDK(t, id)` | `ManagementKey().Get(id)` |
| `LoadDescoperViaSDK(t, id)` | `Descoper().Get(id)` |
| `ProjectExistsViaSDK(t, id)` | `Project().ListProjects()` → find by ID |

### Test Pattern

Every resource test follows this structure:

```go
func TestTenantCRUD(t *testing.T) {
    h := NewHarness(t)
    name := GenerateName(t)

    // 1. Create via Terraform
    attrs := h.ApplyFixture("tenant/create.tf", "descope_tenant.test", "name="+name)
    assert.Equal(t, name, attrs["name"])
    id := StringAttr(attrs, "id")

    // 2. Verify create via SDK (adversarial check)
    sdkTenant := LoadTenantViaSDK(t, id)
    assert.Equal(t, name, sdkTenant.Name)

    // 3. Update via Terraform
    attrs = h.ApplyFixture("tenant/update.tf", "descope_tenant.test", "name="+name)
    assert.Equal(t, true, attrs["enforce_sso"])

    // 4. Verify update via SDK (adversarial check)
    sdkTenant = LoadTenantViaSDK(t, id)
    assert.True(t, sdkTenant.EnforceSSO)

    // 5. Import round-trip
    attrs = h.ReimportResource("tenant/update.tf", "descope_tenant.test", id, "name="+name)
    assert.Equal(t, id, StringAttr(attrs, "id"))

    // 6. Destroy and verify cleanup
    h.Destroy("name=" + name)
    assert.False(t, h.HasState())
}
```

Steps 2 and 4 are the adversarial layer — they bypass Terraform entirely and query the API to confirm the resource actually has the values Terraform claims.

### Build Tags

Tests use build tags to control which tests run in each CI job:

| Tag | Meaning | CI job |
|---|---|---|
| `integration` | Runs in upstream CI | `integration_test` (subset: project tests only) |
| `fork` | Runs in fork CI | `fork_test` (all fork-tagged tests) |
| `integration \|\| fork` | Runs in both | Most resource tests |

Tests requiring paid Descope licenses (`TestSSOApplicationOIDC`, `TestProjectExportDataSource`) detect the license error and call `t.Skip()` instead of failing.

### CI Execution

From `.github/workflows/ci.yml`:

```yaml
# Job: integration_test — runs first
go test -v -count=1 -tags=integration -p 1 -timeout 30m \
  -run 'TestProject(CRUD|Settings|Authorization)' ./tests/integration/

# Job: fork_test — runs after integration_test passes
go test -v -count=1 -tags=fork -p 1 -timeout 30m ./tests/integration/
```

Both jobs use `-p 1` (no parallel packages) to avoid resource conflicts on the shared Descope account.

### Resource Cleanup

After each CI job, `tools/testcleanup` deletes all resources with the `testacc-` prefix:

| Resource type | Discovery method |
|---|---|
| Access keys | `AccessKey().SearchAll()` |
| Management keys | `ManagementKey().Search()` |
| Descopers | `Descoper().List()` |
| Projects | `Project().ListProjects()` |

This runs unconditionally (`if: always()`) so leaked resources don't accumulate even when tests fail.

Run locally:

```bash
source .env && go run ./tools/testcleanup
```

## Layer 7: Release Signing and Provenance

When a `v*` tag is pushed (`.github/workflows/release.yml`):

1. **GoReleaser** — builds binaries for all platforms
2. **Syft** — generates SBOM (Software Bill of Materials)
3. **Cosign** — signs artifacts with Sigstore (keyless)
4. **GPG** — signs SHA256SUMS with GPG key
5. **SLSA** — generates v3 build provenance via `slsa-github-generator`

This ensures consumers can verify the binary they downloaded was built from the claimed source commit.

## Running Tests Locally

```bash
# Unit tests only (no credentials needed)
make test

# Integration tests (requires .env with credentials)
source .env && make testintegration

# Single integration test
source .env && go test -v -count=1 -p 1 -tags=integration \
  -run TestTenantCRUD ./tests/integration/

# Full lint + secret scan
make lint

# Clean up leaked test resources
source .env && make testcleanup
```

Required environment variables for integration tests (see `.env.example`):

```
DESCOPE_MANAGEMENT_KEY=K...
DESCOPE_PROJECT_ID=P...
DESCOPE_BASE_URL=https://api.descope.com
```

## Dependency Management

**Dependabot** (`.github/dependabot.yml`) opens PRs weekly for:

- Go module updates (max 10 open PRs)
- GitHub Actions version bumps (max 5 open PRs, pinned by SHA)

All action references in workflows are pinned to commit SHAs, not tags, to prevent supply chain attacks via tag rewriting.
