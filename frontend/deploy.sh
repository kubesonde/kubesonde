#!/usr/bin/env bash
set -euo pipefail
set -x

# Run from the frontend directory (where package.json lives), regardless of CWD.
cd "$(dirname "$0")"

# Deploys must originate from dev.
BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [[ "${BRANCH}" != "dev" ]]; then
  echo "deploy.sh: must be on 'dev' to deploy (currently on '${BRANCH}')." >&2
  exit 1
fi

# The working tree must be clean, otherwise the version bump can't be committed
# cleanly (this is what previously left a tag on the wrong commit).
if [[ -n "$(git status --porcelain)" ]]; then
  echo "deploy.sh: working tree is dirty. Commit or stash changes first:" >&2
  git status --short >&2
  exit 1
fi

git pull --ff-only origin dev

# Bump the version only (no auto commit/tag) so we control the git steps and
# they land reliably even though package.json lives in a subdirectory.
VERSION=$(npm version patch --no-git-tag-version)  # e.g. v0.8.35

git add package.json package-lock.json
git commit -m "release ${VERSION}"
git tag "${VERSION}"

git push origin dev
git push origin "${VERSION}"

# Create the GitHub release for this tag on every deploy.
gh release create "${VERSION}" --title "${VERSION}" --generate-notes

# Fast-forward main to the released dev state.
git checkout main
git merge --no-edit origin/dev
git push origin main

git checkout dev
