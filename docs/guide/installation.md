# Installation

Kubesonde ships as a single consolidated installer manifest, `kubesonde.yaml`, that contains the CRDs, RBAC, the controller deployment, and the frontend.

## Requirements

| Requirement | Version |
| --- | --- |
| Kubernetes cluster | v1.29+ |
| `kubectl` | v1.29+ |
| Cluster permissions | Ability to create CRDs, namespaces, RBAC, and deployments |

The controller runs probing workloads (a debugger and a monitor container) against the pods in the target namespace, so the cluster must allow those pods to be scheduled.

## Install from a release (recommended)

1. Download `kubesonde.yaml` from the [releases page](https://github.com/kubesonde/kubesonde/releases).
2. Apply it:

   ```bash
   kubectl apply -f kubesonde.yaml
   ```

3. Confirm the controller is running:

   ```bash
   kubectl -n kubesonde-system get pods
   ```

This creates the `kubesonde-system` namespace and installs the `Kubesonde` CRD (API group `security.kubesonde.io/v1`).

## Verify the CRD is registered

```bash
kubectl get crd kubesondes.security.kubesonde.io
kubectl explain kubesonde.spec
```

## Build the installer yourself

If you want to pin a specific image or build from source, regenerate the consolidated manifest from the `crd/` directory:

```bash
git clone https://github.com/kubesonde/kubesonde.git
cd kubesonde/crd
make build-installer
```

See [Contributing / building locally](/how-to/install#building-from-source) for building the controller and frontend images.

## Uninstall

Remove any `Kubesonde` resources first, then the controller and CRDs:

```bash
kubectl delete -f probe.yaml      # your Kubesonde resource(s)
kubectl delete -f kubesonde.yaml  # controller, frontend, and CRDs
```

