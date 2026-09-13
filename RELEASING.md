# Releasing Kubesonde

Releases are fully automated and triggered by pushing a Git tag that starts with
`v` (for example `v0.9.0`). The [`Release` workflow](.github/workflows/go_main.yaml)
then:

1. Builds and tests the CRD (via the reusable `go_base.yaml` workflow).
2. Generates the installer manifest `kubesonde.yaml`, pinned to the released
   tag for **both** the controller and the frontend.
3. Builds and pushes multi-arch images to GHCR — all at the **same tag**:
   - `ghcr.io/kubesonde/controller:<tag>` and `:latest`
   - `ghcr.io/kubesonde/gonetstat:<tag>` and `:latest`
   - `ghcr.io/kubesonde/frontend:<tag>` and `:latest`
4. Creates a GitHub Release with **auto-generated release notes** (categorized
   via [`.github/release.yml`](.github/release.yml)) and attaches
   `kubesonde.yaml` as a downloadable asset.

Images are pushed **before** the Release is created, so a published release
always points at images that exist in the registry.

## Versioning

The **git tag is the single source of truth**. Controller, gonetstat, and
frontend always share it — there is nothing to bump by hand:

- The controller/gonetstat image tags come straight from the pushed tag.
- The frontend version banner is injected into its Docker build via the
  `VERSION` build arg (the tag). `frontend/package.json`'s `version` field is
  **not** used for releases and is intentionally frozen at `0.0.0`; do not bump
  it. Local `npm run build` shows `dev`.

## Cutting a release

```bash
# Make sure main is green and up to date, then:
git checkout main && git pull

# Tag and push. Use a semver tag prefixed with v.
git tag v0.9.0
git push origin v0.9.0
```

That's it — the workflow does the rest. Watch it under the repository's
**Actions** tab.

## Pre-releases

Any tag containing a hyphen is automatically marked as a **pre-release** on
GitHub (it still builds and pushes images):

```bash
git tag v0.9.0-rc.1
git push origin v0.9.0-rc.1
```

## Nicer release notes

Release notes are grouped by pull-request labels (Features, Bug Fixes,
Documentation, Dependencies, …). Label your PRs accordingly to get a clean
changelog. Dependabot PRs and anything labelled `ignore-for-release` are
excluded.
