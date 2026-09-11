# Kubesonde resource reference

Kubesonde is driven by a single custom resource, `Kubesonde`, in the API group **`security.kubesonde.io/v1`**.

```yaml
apiVersion: security.kubesonde.io/v1
kind: Kubesonde
metadata:
  name: kubesonde-sample
spec:
  namespace: default
  probe: all
```

## `spec`

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `namespace` | string | — | Target namespace whose pods will be probed. |
| `probe` | string | — | Default probing behavior: `all` probes every pod-to-pod pair; `none` probes nothing unless explicitly included. |
| `debuggerImage` | string | `instrumentisto/nmap:latest` | Override the image used for the `debugger` (probe) container. |
| `monitorImage` | string | `ghcr.io/kubesonde/gonetstat:latest` | Override the image used for the `monitor` container. |
| `exclude` | list | *(optional)* | Connections to skip during probing. See [entry fields](#include-exclude-entry-fields). |
| `include` | list | *(optional)* | Connections to explicitly probe, optionally with an expected outcome. |

## `include` / `exclude` entry fields

Each entry in `include` and `exclude` accepts:

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `fromPodSelector` | string | — | Selector for the source pod(s). |
| `toPodSelector` | string | — | Selector for the destination pod(s). |
| `port` | string | `"80"` | Destination port to probe. |
| `protocol` | string | `TCP` | Transport protocol. |
| `expected` | string | *(include only)* | Asserts the expected outcome: `Allow` or `Deny`. |

::: warning `port` is a string
In the CRD schema `port` and `protocol` are **strings**. Quote numeric ports in YAML — `port: "8080"` — to avoid validation errors.
:::

### Example with include/exclude

```yaml
apiVersion: security.kubesonde.io/v1
kind: Kubesonde
metadata:
  name: shop-audit
spec:
  namespace: shop
  probe: all
  exclude:
    - fromPodSelector: { matchLabels: { app: web } }
      toPodSelector:   { matchLabels: { app: metrics } }
      port: "9090"
      protocol: TCP
  include:
    - fromPodSelector: { matchLabels: { app: web } }
      toPodSelector:   { matchLabels: { app: db } }
      port: "80"
      protocol: TCP
      expected: Deny
```

## `status`

The resource reports progress in its status:

| Field | Type | Description |
| --- | --- | --- |
| `complete` | bool | `true` when the recorded probe count has been stable for the completeness window (probing has quiesced). |
| `probeCount` | int | Number of probes recorded so far for this Kubesonde. |
| `lastProbeTime` | timestamp | When probing last ran. |
| `conditions` | list | Standard conditions. The `Complete` condition is set to `True` once probing quiesces, so you can use `kubectl wait --for=condition=Complete`. |

::: tip Completeness is a hint, not a stop
`complete: true` means results have stabilized enough to read — it does **not** mean Kubesonde stopped. The controller keeps probing continuously. See [Probe completeness](/concepts/completeness).
:::

## Accepted values

| Constant | Values |
| --- | --- |
| `probe` | `all`, `none` |
| `expected` | `Allow`, `Deny` |

## Watching a resource

```bash
kubectl get kubesondes
kubectl wait --for=condition=Complete kubesonde/kubesonde-sample --timeout=300s
kubectl explain kubesonde.spec
```
