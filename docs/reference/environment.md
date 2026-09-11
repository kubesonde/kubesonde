# Environment variables

Configuration knobs read by the Kubesonde controller.

## `KUBESONDE_PROBE_WORKERS`

Number of probes executed **concurrently** by the dispatcher's bounded worker pool.

| | |
| --- | --- |
| **Type** | Positive integer |
| **Default** | `10` |
| **Invalid/unset** | Falls back to `10` |

The default is kept modest to stay well within the kubelet's exec/attach concurrency limits. Raise it to probe large namespaces faster, at the cost of more load on the cluster.

Set it on the controller deployment:

```yaml
env:
  - name: KUBESONDE_PROBE_WORKERS
    value: "20"
```

Or patch a running controller:

```bash
kubectl -n kubesonde-system set env \
  deployment/kubesonde-controller-manager KUBESONDE_PROBE_WORKERS=20
```

Or when running from source:

```bash
KUBESONDE_PROBE_WORKERS=20 make run
```

See [How to tune probe concurrency](/how-to/tune-concurrency) for guidance on choosing a value.
