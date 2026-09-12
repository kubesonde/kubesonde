# How to install and uninstall Kubesonde

Task-oriented recipes for getting Kubesonde onto (and off) a cluster. For a guided first run, see the [Quickstart](/guide/quickstart).

## Install from a release

```bash
kubectl apply -f kubesonde.yaml   # downloaded from the releases page
```

Confirm the controller is up:

```bash
kubectl -n kubesonde-system get pods
kubectl -n kubesonde-system get deploy kubesonde-controller-manager
```

## Fix RBAC errors during install

If `kubectl apply` fails with a forbidden/RBAC error, your user lacks the permissions to create the CRDs or cluster roles. Grant yourself admin (evaluation clusters only) or ask a cluster admin to apply the manifest:

```bash
kubectl create clusterrolebinding my-admin \
  --clusterrole=cluster-admin --user="$(kubectl config view -o jsonpath='{.users[0].name}')"
```

## Building from source

Clone the repository — it contains the controller, the CRDs, and the frontend:

```bash
git clone https://github.com/kubesonde/kubesonde.git
cd kubesonde
```

Repository layout:

| Path | Contents |
| --- | --- |
| `crd/` | Backend controller and the Kubesonde CRD |
| `frontend/` | The UI for analyzing probe outputs |
| `examples/` | Sample output from Kubesonde |
| `docs/` | Documentation and design notes |

### Build the controller

The controller uses a Kubebuilder-style `Makefile`:

```bash
cd crd
make build   # build the manager binary
make test    # run the tests
make run     # run the controller against your current kubecontext
```

Regenerate the consolidated installer manifest:

```bash
make build-installer
```

### Build the frontend

```bash
cd frontend
nvm use      # switch to the Node version pinned in .nvmrc
npm install
npm start    # dev server
npm run build
npm test
```

## Uninstall

Delete your `Kubesonde` resources first, then the controller and CRDs:

```bash
kubectl delete -f probe.yaml      # your Kubesonde resource(s)
kubectl delete -f kubesonde.yaml  # controller, frontend, and CRDs
```

Deleting the CRD removes any remaining `Kubesonde` objects with it.
