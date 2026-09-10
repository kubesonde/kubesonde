#!/usr/bin/env bash
set -euo pipefail
set -x

# Run from the frontend directory (where package.json lives), regardless of CWD.
cd "$(dirname "$0")"

# `npm version patch` bumps package.json/package-lock, creates the release
# commit, and tags it (e.g. v0.8.33) — no separate add/commit/tag needed.
VERSION=$(npm version patch)

git push origin "${VERSION}"
git push origin dev

git checkout main
git merge --no-edit origin/dev
git push origin main

git checkout dev
