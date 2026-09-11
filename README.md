![workflow](https://github.com/kubesonde/kubesonde/actions/workflows/go_main.yaml/badge.svg)
![frontend_test](https://github.com/kubesonde/kubesonde/actions/workflows/frontend_dev.yaml/badge.svg)
![frontend_deployment](https://github.com/kubesonde/kubesonde/actions/workflows/deploy_frontend.yaml/badge.svg)
[![Netlify Status](https://api.netlify.com/api/v1/badges/df3643ab-e317-4b96-b5c2-de937837b375/deploy-status)](https://app.netlify.com/sites/testksonde/deploys)

# Kubesonde

Kubesonde is a tool that probes and visualizes the *actual* network connectivity of applications running in a Kubernetes cluster, so you can compare it against the network policies you meant to enforce. It works by instrumenting live pods at runtime and reporting every connection attempt observed, rather than relying only on what NetworkPolicy manifests declare.

Kubesonde is **not** a network policy engine, admission controller, or firewall: it does not block or modify traffic, and it does not replace tools like Cilium, Calico, or Kubernetes NetworkPolicies. It is a diagnostic and auditing tool for finding gaps between the connectivity you think you have and the connectivity you actually have.

## Why Kubesonde

Kubernetes network security is declarative: you write NetworkPolicies and trust that the cluster enforces them. In practice, the gap between *intended* and *actual* connectivity is where misconfigurations hide:

- **NetworkPolicies are hard to get right.** A missing selector, an overly broad `podSelector`, or a forgotten default-deny policy can silently leave pods reachable when they should be isolated. Reading the manifests does not tell you what traffic is really allowed.
- **The enforced state can differ from the declared state.** CNI plugins, admission controllers, and overlapping policies interact in non-obvious ways. The only way to know what a pod can actually reach is to try.
- **Auditing connectivity by hand does not scale.** In a namespace with dozens of microservices, manually reasoning about which pod can talk to which — and on which ports — is error-prone and quickly becomes impossible.

Kubesonde answers a concrete question that manifests alone cannot: *from this pod, what can I actually reach right now?* It probes live pods, records every connection attempt and its outcome, and lets you visualize the resulting connectivity graph. This makes it useful for verifying that isolation policies work as intended, for finding unexpected reachability before an attacker does, and for understanding the real communication footprint of an application.

## Get Started

### 1. Install Kubesonde

Grab the latest `kubesonde.yaml` installer manifest from the [releases page](https://github.com/kubesonde/kubesonde/releases) and apply it to your cluster:

```bash
kubectl apply -f kubesonde.yaml
```

This creates the `kubesonde-system` namespace, the Kubesonde CRDs, and the controller that runs the probes.

### 2. Create a Kubesonde resource

Once the controller is running, create a `Kubesonde` resource that describes what you want to probe. The following example targets every pod in the `default` namespace:

```yaml
apiVersion: security.kubesonde.io/v1
kind: Kubesonde
metadata:
  name: kubesonde-sample
spec:
  namespace: default
  probe: all
```

Save it as `probe.yaml` and apply it:

```bash
kubectl apply -f probe.yaml
```

The controller starts probing the target namespace. You can watch progress — including whether probing has quiesced — with:

```bash
kubectl get kubesondes
```

```
NAME                 NAMESPACE   PROBES   COMPLETE   AGE
kubesonde-sample     default     228      true       5m
```

### 3. Fetch the results

Port-forward the controller and download the recorded probes:

```bash
kubectl --namespace kubesonde-system port-forward deployment.apps/kubesonde-controller-manager 2709 &
curl localhost:2709/probes > probes.json
```

To avoid guessing how long probing takes, wait for the `Complete` condition before fetching:

```bash
kubectl wait --for=condition=Complete kubesonde/kubesonde-sample --timeout=300s
```

> :warning: Fetching results immediately after applying the resource may return empty or incomplete data. Wait for the `Complete` condition (or a few minutes, depending on the number of pods).

### 4. Visualize

The installer deploys the Kubesonde frontend alongside the controller, so you can explore results without leaving your cluster. Port-forward the frontend service (and the controller, which the UI reads from at `localhost:2709`), then open the UI in your browser:

```bash
kubectl --namespace kubesonde-system port-forward service/kubesonde-frontend 8088:8088 &
kubectl --namespace kubesonde-system port-forward deployment.apps/kubesonde-controller-manager 2709:2709 &
```

Then browse to [http://localhost:8088](http://localhost:8088).

Alternatively, if you exported a `probes.json` file, you can upload it to the [hosted Kubesonde website](https://kubesonde.jackops.dev) to explore the connectivity graph.

### Clean up

```bash
kubectl delete -f probe.yaml     # remove the scanner
kubectl delete -f kubesonde.yaml # remove the controller and CRDs
```

## Custom Resources

Kubesonde is driven by a single custom resource, `Kubesonde` (API group `security.kubesonde.io/v1`). Its spec supports the following fields:

| Field | Description |
| --- | --- |
| `namespace` | Target namespace whose pods will be probed. |
| `probe` | Default probing behavior: `all` probes every pod-to-pod pair, `none` probes nothing unless explicitly included. |
| `debuggerImage` | *(optional)* Override the image used for the debugger container. |
| `monitorImage` | *(optional)* Override the image used for the monitor container. |
| `exclude` | *(optional)* List of connections to skip during probing. |
| `include` | *(optional)* List of connections to explicitly probe, optionally with an expected outcome. |

Each `include` / `exclude` entry accepts `fromPodSelector`, `toPodSelector`, `port` (default `80`) and `protocol` (default `TCP`). An `include` entry may also set `expected` (`Allow` or `Deny`) to assert the expected outcome of that probe.

The resource also exposes status information: `probeCount`, a `complete` flag, and a `Complete` condition (so you can use `kubectl wait --for=condition=Complete`). Completeness is based on **quiescence**: probing is considered complete when the recorded probe count has been stable for a fixed window (30 seconds by default). This is a hint that it is a good time to fetch results — the controller keeps probing continuously.

You can also query completeness directly from the controller:

```bash
curl localhost:2709/probes/status
```

```json
{
  "complete": true,
  "count": 228,
  "window": 30000000000,
  "secondsSinceLastChange": 42.5
}
```

### Tuning probe concurrency (optional)

Kubesonde dispatches probes through a bounded worker pool that runs **10** probes in parallel by default. Override it with the `KUBESONDE_PROBE_WORKERS` environment variable on the controller (any positive integer; invalid or unset values fall back to `10`):

```yaml
env:
  - name: KUBESONDE_PROBE_WORKERS
    value: "20"
```

## Building Locally

Clone the repository — it contains the controller, the CRDs and the frontend:

```bash
git clone https://github.com/kubesonde/kubesonde.git
cd kubesonde
```

The project is organized as follows:

- `crd`: backend controller and the Kubesonde CRD
- `frontend`: the UI for analyzing probe outputs
- `examples`: sample output from Kubesonde
- `docs`: documentation and design notes

### Backend / controller

The controller lives in `crd/` and is built with a standard Kubebuilder-style `Makefile`:

```bash
cd crd
make build   # build the manager binary
make test    # run the tests
make run     # run the controller against your current kubecontext
```

Export the worker override when running locally if you want to tune concurrency:

```bash
KUBESONDE_PROBE_WORKERS=20 make run
```

To regenerate the consolidated installer manifest (the `kubesonde.yaml` shipped in releases):

```bash
make build-installer
```

### Frontend

The UI lives in `frontend/` and uses Vite:

```bash
cd frontend
nvm use      # switch to the Node version pinned in .nvmrc
npm install
npm start    # start the dev server
npm run build
npm test
```

## Contributing

Contributions to the project are welcome. Create a PR and let's discuss the changes.

## Credits

Logo from [Elisabetta Russo](stelladigitale.it) info@stelladigitale.it

## Publications

Kubesonde has been described and used in the following peer-reviewed papers:

- Jacopo Bufalino, Mario Di Francesco, Tuomas Aura. **["Analyzing Microservice Connectivity with Kubesonde"](https://dl.acm.org/doi/10.1145/3611643.3613899)**. ESEC/FSE 2023.
- Jacopo Bufalino, Jose Luis Martin-Navarro, Mario Di Francesco, Tuomas Aura. **["Inside Job: Defending Kubernetes Clusters Against Network Misconfigurations"](https://dl.acm.org/doi/10.1145/3749220)**. Proceedings of the ACM on Networking, 2025.

## License

Kubesonde is licensed under the Apache License, Version 2.0.
