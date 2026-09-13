# Contributing to Kubesonde

Thanks for your interest in improving Kubesonde! Contributions of all kinds are
welcome — code, docs, bug reports, and ideas.

By participating you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## Ways to contribute

- **Report a bug** or **request a feature** via the
  [issue tracker](https://github.com/kubesonde/kubesonde/issues/new/choose).
- **Improve the docs** — the site lives in [`docs/`](docs/) (VitePress).
- **Fix a bug or build a feature** — open a PR (see below).

## Repository layout

| Path | Contents |
| --- | --- |
| `crd/` | Controller and the `Kubesonde` CRD (Go, Kubebuilder-style) |
| `frontend/` | The UI for analyzing probe outputs (Vite + React) |
| `docker/` | The `gonetstat` monitor image |
| `docs/` | Documentation site (VitePress) |
| `examples/` | Sample Kubesonde output |

## Development

### Prerequisites

- Go (see `crd/go.mod` for the version), Docker, and `kubectl`.
- A Kubernetes cluster for manual testing (a local [kind](https://kind.sigs.k8s.io/)
  or [minikube](https://minikube.sigs.k8s.io/) cluster is fine).
- Node (see `frontend/.nvmrc`) for the frontend.

### Controller (`crd/`)

```bash
cd crd
make build      # build the manager binary
make test       # unit tests (envtest)
make lint       # golangci-lint
make run        # run the controller against your current kubecontext
make test-e2e   # e2e tests against a kind cluster
```

### Frontend (`frontend/`)

```bash
cd frontend
nvm use
npm install
npm start       # dev server
npm run build
npm test
```

### Docs (`docs/`)

```bash
cd docs
npm install
npm run docs:dev     # http://localhost:5173/kubesonde/
npm run docs:build
```

## Pull requests

1. **Fork and branch** off `main` with a descriptive branch name.
2. **Keep PRs focused** — one logical change per PR.
3. **Add tests** for new behavior and make sure `make test` (controller) and
   `npm test` (frontend) pass.
4. **Write a clear description** of what changed and why.
5. **Label your PR** so it lands in the right release-notes section
   (see below).

Maintainers review PRs and may request changes. Once approved and green, a
maintainer merges.

### Labels and release notes

Release notes are generated automatically from merged PRs, grouped by label
(see [`.github/release.yml`](.github/release.yml)):

| Label | Section |
| --- | --- |
| `enhancement` / `feature` | 🚀 Features |
| `bug` / `fix` | 🐛 Bug Fixes |
| `documentation` | 📖 Documentation |
| `dependencies` | ⬆️ Dependencies |

Label your PR accordingly. Use `ignore-for-release` to exclude a PR from the
changelog.

## Releases

Releases are automated and triggered by pushing a `v*` git tag — see
[RELEASING.md](RELEASING.md). Only maintainers cut releases.

## Reporting security issues

Please do **not** open a public issue for security vulnerabilities. See
[SECURITY.md](SECURITY.md) for private reporting.
