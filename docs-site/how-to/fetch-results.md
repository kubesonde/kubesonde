# How to fetch and export results

Kubesonde records probes in the controller. You read them over the controller's HTTP API after port-forwarding.

## Wait until probing is complete

Fetching too early returns empty or partial data. Wait for the `Complete` condition:

```bash
kubectl wait --for=condition=Complete kubesonde/kubesonde-sample --timeout=300s
```

Or check the count directly:

```bash
kubectl get kubesondes
```

## Port-forward the controller

```bash
kubectl --namespace kubesonde-system \
  port-forward deployment.apps/kubesonde-controller-manager 2709 &
```

## Download the probes

```bash
curl -s localhost:2709/probes > probes.json
```

This JSON is what you upload to the [hosted UI](https://kubesonde.jackops.dev) or load into the [bundled frontend](/how-to/visualize).

## Check completeness programmatically

The controller also exposes a status endpoint:

```bash
curl -s localhost:2709/probes/status
```

```json
{
  "complete": true,
  "count": 228,
  "window": 30000000000,
  "secondsSinceLastChange": 42.5
}
```

- `complete` — whether the probe count has been stable for the quiescence window.
- `count` — number of probes recorded so far.
- `window` — the quiescence window in nanoseconds (30s = `30000000000`).
- `secondsSinceLastChange` — how long the count has been unchanged.

See [Probe completeness](/concepts/completeness) for what "complete" means and why the controller keeps probing afterward.

## Full API reference

See the [Controller HTTP API reference](/reference/http-api) for every endpoint.
