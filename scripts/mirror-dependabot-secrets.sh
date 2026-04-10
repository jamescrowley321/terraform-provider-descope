#!/usr/bin/env bash
#
# Mirror the repo secrets that CI references into the Dependabot secret scope
# so dependabot-authored PRs can run the full integration suite. Run once per
# repo; re-run on value rotation or when CI references a new secret.
#
# GitHub does not let any principal read an existing secret value, so this
# script reads required values from the current shell environment. Populate
# the shell first:
#
#     source .env
#
# Alternatively export each variable by hand or pull from a secret store.
#
# Usage:
#   ./scripts/mirror-dependabot-secrets.sh
#
# Requires: `gh` authenticated as a repo admin on the target repo.

set -euo pipefail

REPO="${REPO:-jamescrowley321/terraform-provider-descope}"

SECRETS=(
  DESCOPE_MANAGEMENT_KEY
  DESCOPE_PROJECT_ID
)

missing=()
for name in "${SECRETS[@]}"; do
  if [[ -z "${!name:-}" ]]; then
    missing+=("$name")
  fi
done

if (( ${#missing[@]} > 0 )); then
  printf 'ERROR: the following env vars are not set in the current shell:\n' >&2
  printf '  %s\n' "${missing[@]}" >&2
  printf '\nExport them (e.g. `source .env`) then re-run this script.\n' >&2
  exit 1
fi

printf 'Mirroring %d secrets to Dependabot scope on %s...\n' "${#SECRETS[@]}" "$REPO"
for name in "${SECRETS[@]}"; do
  printf '%s' "${!name}" | gh secret set "$name" --app dependabot --repo "$REPO"
  printf '  set %s\n' "$name"
done

printf '\nDone. Verify with:\n  gh api repos/%s/dependabot/secrets\n' "$REPO"
