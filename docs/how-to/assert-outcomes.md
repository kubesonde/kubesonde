# How to assert expected outcomes

Kubesonde can do more than observe — it can check observations against what you *expect*. Add an `expected` value to an `include` entry and Kubesonde records whether the observed outcome matched.

## Assert that a connection is allowed

```yaml
spec:
  namespace: shop
  probe: all
  include:
    - fromPodSelector: { matchLabels: { app: api } }
      toPodSelector:   { matchLabels: { app: db } }
      port: "80"
      protocol: TCP
      expected: Allow
```

## Assert that a connection is denied

```yaml
    - fromPodSelector: { matchLabels: { app: web } }
      toPodSelector:   { matchLabels: { app: db } }
      port: "80"
      protocol: TCP
      expected: Deny
```

`expected` accepts `Allow` or `Deny`. It is only valid on `include` entries.

## How to read the result

Fetch the probe results and look at the asserted connections:

```bash
kubectl -n kubesonde-system \
  port-forward deployment.apps/kubesonde-controller-manager 2709 &
curl -s localhost:2709/probes | less
```

Each asserted probe records the observed outcome alongside the expectation, so a mismatch is visible in the output and highlighted in the [UI](/how-to/visualize).

## Use it as a regression test

Because expectations live in the `Kubesonde` resource, keep the manifest in version control and re-apply it after any `NetworkPolicy` change:

```bash
kubectl apply -f shop-audit.yaml
kubectl wait --for=condition=Complete kubesonde/shop-audit --timeout=300s
```

A walkthrough that combines this with a real `NetworkPolicy` is in the [audit tutorial](/tutorials/audit-networkpolicy).
