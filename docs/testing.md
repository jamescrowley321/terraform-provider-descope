# Testing the provider

Run `go test ./...` for unit and Terraform lifecycle tests. The lifecycle regression test uses a local mock
API and a Terraform CLI; it needs no Descope credentials. It checks precise numeric claims, JSON formatting,
drift repair, clearing, secret retention, and deletion.

Run `go test -tags=integration -count=1 -p 1 ./tests/integration` with `DESCOPE_MANAGEMENT_KEY`,
`DESCOPE_BASE_URL`, and `DESCOPE_PROJECT_ID` for the retained-feature suite. The project ID supplies test
fixtures and SDK verification; provider resources and data sources use explicit project IDs.

`TestM2MTokenLifecycle` creates a disposable project and verifies both SDK access-key exchange and OAuth
client credentials, signed JWTs, custom claims, project and tenant authorization, drift repair, clearing,
deactivation, and reactivation. Password-settings and FGA-schema tests use disposable projects too.
Project export is read-only and skips if the account lacks the required snapshot-export license.

The canonical upstream acceptance tests remain under `internal/models` and `tools/tfexport`. Enable them
with `TF_ACC=1` and `DESCOPE_TESTACC_PROJECT_ID`; some tests configure singletons in that project, so use
an isolated test project. CI runs upstream project acceptance coverage and the retained fork suite.

Before committing, run `go build ./...`, `gofmt`, the appropriate tests, `make lint`, and generated documentation
checks. Regenerate model documentation with `make terragen` using the required connector templates, then
registry documentation with `make docs`. The Go toolchain must match the lint tool's supported Go version.
