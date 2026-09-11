# How to target specific pods

By default `probe: all` probes every pod-to-pod pair in the target namespace. When you only care about certain connections — or want to skip noisy ones — use `include` and `exclude`.

## Probe everything (the default)

```yaml
spec:
  namespace: default
  probe: all
```

## Probe nothing except what you list

Set `probe: none` and add `include` entries. Only the listed connections are probed:

```yaml
spec:
  namespace: default
  probe: none
  include:
    - fromPodSelector: { matchLabels: { app: web } }
      toPodSelector:   { matchLabels: { app: api } }
      port: "80"
      protocol: TCP
```

## Skip specific connections

Keep `probe: all` but carve out connections you don't want probed with `exclude`:

```yaml
spec:
  namespace: default
  probe: all
  exclude:
    - fromPodSelector: { matchLabels: { app: web } }
      toPodSelector:   { matchLabels: { app: metrics } }
      port: "9090"
      protocol: TCP
```

## Entry fields

Each `include` / `exclude` entry accepts:

| Field | Default | Description |
| --- | --- | --- |
| `fromPodSelector` | — | Label selector for the source pods. |
| `toPodSelector` | — | Label selector for the destination pods. |
| `port` | `80` | Destination port to probe. |
| `protocol` | `TCP` | Transport protocol. |
| `expected` | — | *(include only)* Assert the outcome: `Allow` or `Deny`. See [assert outcomes](/how-to/assert-outcomes). |

See the full [Kubesonde resource reference](/reference/kubesonde-resource) for the authoritative schema.

::: tip
Selectors follow standard Kubernetes label-selector semantics, so `matchLabels` and `matchExpressions` both work. Target a single pod by selecting a label that is unique to it.
:::
