# Claude Code Instructions

## Before Committing

**Always run these checks before committing any code changes:**

1. `gofmt -l ./...` — fix any formatting issues with `gofmt -w <file>`
2. `go build ./...` — must compile cleanly
3. `source .env && go test -v -count=1 -p 1 ./internal/models/accesskey/` — run relevant acceptance tests

If any check fails, fix the issue before committing. Do not commit broken code.

## Git & GitHub

- This is a fork: `jamescrowley321/terraform-provider-descope`
- **Always create PRs against the fork**, not upstream `descope/terraform-provider-descope`
- If `gh repo set-default` gets reset, run: `gh repo set-default jamescrowley321/terraform-provider-descope`
