# Contributing

Thank you for your interest in contributing to this community fork of the Descope Terraform Provider.

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) (see `go.mod` for the required version)
- [Terraform CLI](https://developer.hashicorp.com/terraform/install)
- A Descope management key (required for acceptance tests)

### Setup

```bash
git clone https://github.com/jamescrowley321/terraform-provider-descope
cd terraform-provider-descope
make dev
```

This builds the provider binary and configures a local `~/.terraformrc` override so Terraform uses your local build.

### Build and Test

```bash
make install       # Rebuild and install the provider
make testacc       # Run acceptance tests
make testcoverage  # Run tests with coverage report
make lint          # Run linting and security checks
```

## Workflow

1. **Open an issue first** — Describe the feature, bug, or change you want to make. This avoids duplicate work and lets us align on the approach.
2. **Fork and branch** — Create a feature branch from `main`:
   ```bash
   git checkout -b feat/my-feature main
   ```
3. **Make your changes** — Follow the architecture patterns described in `internal/README.md`.
4. **Test** — Add or update unit tests and acceptance tests. All tests must pass.
5. **Lint** — Run `make lint` before committing.
6. **Commit** — Use [Conventional Commits](https://www.conventionalcommits.org/) (see below).
7. **Push and open a PR** — Push your branch and open a pull request against `main`.

### Conventional Commits

All commit messages must follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<optional scope>): <description>

<optional body>

<optional footer(s)>
```

**Types:** `feat`, `fix`, `docs`, `test`, `chore`, `refactor`, `ci`, `build`, `perf`, `style`

**Examples:**

```
feat(tenant): add descope_tenant resource with full CRUD
fix(access-key): handle missing expiration in API response
docs: update README with new resource table
test(user): add acceptance tests for user creation
```

**Breaking changes:** Add `!` after the type/scope and include a `BREAKING CHANGE:` footer:

```
feat(api)!: rename descope_project settings block

BREAKING CHANGE: The `project_settings` block has been renamed to `settings`.
```

### Branch Naming

Use the same type prefixes for branch names:

```
feat/tenant-resource
fix/access-key-expiration
docs/update-readme
```

## AI-Assisted Contributions

Contributions that use AI tools (GitHub Copilot, Claude, ChatGPT, Cursor, etc.) are welcome. We apply the same quality standards to all contributions regardless of how they were authored.

### Requirements for AI-Assisted PRs

- **All CI checks must pass** — lint, tests, security scans. No exceptions.
- **Tests are required** — AI-generated code must include unit tests and acceptance tests with adequate coverage.
- **Disclose tooling** — Note which AI tools were used in the PR description. A simple line like "Co-authored with Claude Code" is sufficient.
- **Human review is required** — All PRs require at least one human reviewer approval before merge.
- **You are responsible** — The submitter is accountable for the correctness, security, and quality of the code, regardless of whether it was AI-generated.

### Requesting a review from Claude

Comment `@claude review this PR` on a pull request and Claude runs the
[blind-peer-review](https://github.com/jamescrowley321/blind-peer-review) lenses
against the diff — one fresh reviewer per lens (Cold Read, Edge Cases, Acceptance
Criteria, Security Review, Red Team), adjudicated fail-closed into a PASS/BLOCK
verdict. It is the same lens library the maintainers run locally, so CI and the
laptop do not drift apart.

It runs **only when asked**. There is no automatic review on every pull request;
#215 removed that deliberately.

**Only maintainers can trigger it.** The workflow runs with repository secrets on
a public repo, so the job's `if:` requires the commenter's `author_association`
to be `OWNER`, `MEMBER` or `COLLABORATOR`. GitHub evaluates that before
scheduling anything, so a mention from anyone else never starts a job at all — it
is not merely rejected after the fact. GitHub offers no "approve this run" prompt
for comment-triggered workflows the way it does for pull requests from forks, so
ask a maintainer to comment `@claude review this PR` on your pull request — a
maintainer asking is the approval.

**The verdict is advice, not approval.** Two limitations, stated plainly because
a clean PASS is exactly when they matter most:

- The lenses run as subagents of one session. Each starts from a clean context,
  but they share a process, a model and a host — weaker than the upstream CI
  adapter, which runs one job per lens. Upstream measured that batching lenses
  into a single agent's context turned a MUST FIX finding into a clean PASS.
- Claude reviews code that Claude often helped write, so the reviewer shares the
  author's blind spots. Every run prints this caveat; leave it in the posted
  findings rather than trimming it.

**Human review is required** regardless, per the rule above: every PR needs at
least one human reviewer's approval before merge. A lens verdict does not count
as that approval and is not a required check.

### What We Look For

- Code follows the existing architecture patterns (see `internal/README.md`)
- No hallucinated APIs or invented SDK methods — verify against the [Descope Go SDK](https://github.com/descope/go-sdk)
- Tests actually run and cover the new functionality
- Documentation is accurate and complete

## Adding New Resources

When adding a new Terraform resource (e.g., `descope_tenant`), follow this checklist:

1. [ ] Define the model in `internal/models/` following the existing patterns
2. [ ] Implement the entity in `internal/entities/`
3. [ ] Add the resource to the provider in `internal/resources/`
4. [ ] Add unit tests
5. [ ] Add acceptance tests
6. [ ] Generate documentation with `make docs`
7. [ ] Update the resource table in `README.md`

See `internal/README.md` for detailed architecture documentation.

## Code of Conduct

Be respectful, constructive, and collaborative. We're building something useful together.

## Questions?

Open a [discussion or issue](https://github.com/jamescrowley321/terraform-provider-descope/issues) and we'll help you get started.
