# How to tune probe concurrency

Kubesonde dispatches probes through a bounded worker pool. By default it runs **10** probes in parallel.

## Override the worker count

Set the `KUBESONDE_PROBE_WORKERS` environment variable on the controller to any positive integer:

```yaml
env:
  - name: KUBESONDE_PROBE_WORKERS
    value: "20"
```

Invalid or unset values fall back to `10`.

## Apply it to a running controller

Patch the controller deployment:

```bash
kubectl -n kubesonde-system set env \
  deployment/kubesonde-controller-manager KUBESONDE_PROBE_WORKERS=20
```

Or edit the installer manifest before applying it.

## When running locally

Export the variable when running the controller from source:

```bash
KUBESONDE_PROBE_WORKERS=20 make run
```

## Choosing a value

- **Higher** concurrency probes a large namespace faster but puts more load on the cluster network and the API server.
- **Lower** concurrency is gentler on busy clusters.

Start at the default of `10` and increase only if scans of large namespaces are too slow.

See the [environment variables reference](/reference/environment) for all configuration knobs.
