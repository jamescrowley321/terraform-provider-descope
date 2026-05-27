#!/usr/bin/env bash
#
# upstream-sync-check — surface merge-risk signals before syncing upstream/main.
#
# Run BEFORE `git merge upstream/main` to catch silent shadows and naming
# divergences that won't show up as git conflicts. Three checks:
#
#   1. COLLISION   — files added independently on both fork and upstream since
#                    the merge-base. These will arrive as content conflicts but
#                    listing them up front lets you decide whose version wins.
#   2. CONVENTION  — fork-only files in shared-with-upstream directories whose
#                    basenames use underscores. Upstream's convention is no
#                    underscores between words (e.g. inboundapp.go). A fork
#                    file at internal/resources/access_key.go silently shadows
#                    whatever upstream would add at internal/resources/accesskey.go.
#   3. INCOMING    — Go files present in upstream but not in fork. Informational:
#                    these are arriving on merge.
#
# Usage:
#   ./scripts/upstream-sync-check.sh
#
# Requires: `git` with `upstream` and `origin` remotes configured.
# Exit code: 0 if clean; 1 if COLLISION or CONVENTION findings.

set -euo pipefail

REMOTE_FORK="${REMOTE_FORK:-origin}"
REMOTE_UPSTREAM="${REMOTE_UPSTREAM:-upstream}"
BASE_BRANCH="${BASE_BRANCH:-main}"

FORK_REF="$REMOTE_FORK/$BASE_BRANCH"
UPSTREAM_REF="$REMOTE_UPSTREAM/$BASE_BRANCH"

git fetch "$REMOTE_UPSTREAM" "$BASE_BRANCH" --quiet
git fetch "$REMOTE_FORK" "$BASE_BRANCH" --quiet

MERGE_BASE=$(git merge-base "$FORK_REF" "$UPSTREAM_REF")
echo "=== upstream-sync-check ==="
echo "merge-base:    $MERGE_BASE"
echo "fork ref:      $FORK_REF"
echo "upstream ref:  $UPSTREAM_REF"
echo ""

exit_code=0

# 1. COLLISION — files added on both sides since merge-base.
fork_added=$(git diff --name-only --diff-filter=A "$MERGE_BASE..$FORK_REF" | sort)
upstream_added=$(git diff --name-only --diff-filter=A "$MERGE_BASE..$UPSTREAM_REF" | sort)
collisions=$(comm -12 <(echo "$fork_added") <(echo "$upstream_added"))
if [ -n "$collisions" ]; then
  echo "COLLISION: files added independently on both sides — decide whose version wins:"
  echo "$collisions" | sed 's/^/  /'
  echo ""
  exit_code=1
fi

# 2. CONVENTION — fork-only files in shared dirs with underscored basenames.
#    Compare directories present in both upstream and fork to find "shared dirs".
shared_dirs=$(comm -12 \
  <(git ls-tree -rd --name-only "$UPSTREAM_REF" | sort) \
  <(git ls-tree -rd --name-only "$FORK_REF" | sort))

convention_warnings=""
while IFS= read -r dir; do
  [ -z "$dir" ] && continue
  fork_files=$(git ls-tree --name-only "$FORK_REF" "$dir/" | grep -E '\.go$' | grep -v '_test\.go$' || true)
  upstream_files=$(git ls-tree --name-only "$UPSTREAM_REF" "$dir/" | grep -E '\.go$' | grep -v '_test\.go$' || true)
  fork_only=$(comm -23 <(echo "$fork_files" | sort) <(echo "$upstream_files" | sort))
  underscored=$(echo "$fork_only" | grep -E '[a-z]_[a-z]' || true)
  if [ -n "$underscored" ]; then
    while IFS= read -r f; do
      [ -z "$f" ] && continue
      convention_warnings+="  $f"$'\n'
    done <<< "$underscored"
  fi
done <<< "$shared_dirs"

if [ -n "$convention_warnings" ]; then
  echo "CONVENTION: fork-only Go files in shared-with-upstream dirs use underscores"
  echo "            (upstream convention is no underscores between words):"
  echo -n "$convention_warnings"
  echo ""
  exit_code=1
fi

# 3. INCOMING — Go files in upstream not yet in fork (informational).
incoming=$(comm -23 \
  <(git ls-tree -r --name-only "$UPSTREAM_REF" | grep -E '\.go$' | sort) \
  <(git ls-tree -r --name-only "$FORK_REF" | grep -E '\.go$' | sort))
if [ -n "$incoming" ]; then
  echo "INCOMING: Go files in upstream not yet in fork (arriving on merge):"
  echo "$incoming" | sed 's/^/  /'
  echo ""
fi

[ $exit_code -eq 0 ] && echo "OK — no merge-risk signals detected."
exit $exit_code
