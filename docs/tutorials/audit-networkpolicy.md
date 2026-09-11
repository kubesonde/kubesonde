# Tutorial: Auditing a NetworkPolicy

::: info What you'll learn
You'll write a `NetworkPolicy` to isolate a database, then use Kubesonde's **expected outcomes** to assert — automatically — that the isolation actually holds. This turns a one-off scan into a repeatable policy test.
:::

**Time:** ~20 minutes · **You need:** the `shop` namespace from [your first scan](/tutorials/first-scan), a CNI that enforces `NetworkPolicy` (e.g. Calico or Cilium), and Kubesonde installed.

## Step 1 — State the intent

Our intended policy for the `shop` namespace is:

- `web` may reach `api` on port 80.
- Nothing may reach `db` except `api`.

Right now (from the previous tutorial) *everything* can reach *everything*. Let's fix that and verify it.

## Step 2 — Apply a default-deny + allow policy

```yaml
# db-isolation.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: db-allow-api-only
  namespace: shop
spec:
  podSelector:
    matchLabels:
      app: db
  policyTypes: ["Ingress"]
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: api
      ports:
        - port: 80
          protocol: TCP
```

```bash
kubectl apply -f db-isolation.yaml
```

## Step 3 — Encode your expectations in Kubesonde

Instead of just observing connectivity, tell Kubesonde what you *expect* and let it check. Use `include` entries with an `expected` outcome:

```yaml
# shop-audit.yaml
apiVersion: security.kubesonde.io/v1
kind: Kubesonde
metadata:
  name: shop-audit
spec:
  namespace: shop
  probe: all
  include:
    # api is allowed to reach db
    - fromPodSelector: { matchLabels: { app: api } }
      toPodSelector:   { matchLabels: { app: db } }
      port: "80"
      protocol: TCP
      expected: Allow
    # web must NOT reach db
    - fromPodSelector: { matchLabels: { app: web } }
      toPodSelector:   { matchLabels: { app: db } }
      port: "80"
      protocol: TCP
      expected: Deny
```

```bash
kubectl apply -f shop-audit.yaml
kubectl wait --for=condition=Complete kubesonde/shop-audit --timeout=300s
```

## Step 4 — Read the verdict

Fetch the results and inspect the two asserted probes:

```bash
kubectl -n kubesonde-system \
  port-forward deployment.apps/kubesonde-controller-manager 2709 &
curl -s localhost:2709/probes > shop-audit.json
```

For each `include` entry with an `expected` value, the probe result records both the observed outcome and whether it matched your expectation. A mismatch — for example `web → db` succeeding when you expected `Deny` — is a policy bug you want to catch here rather than in production.

::: tip
The connectivity graph in the [UI](/how-to/visualize) makes mismatches easy to spot visually: expected-but-blocked and unexpected-but-allowed edges stand out from the rest of the graph.
:::

## Step 5 — Iterate

If an assertion fails, adjust your `NetworkPolicy` and re-run. Because the expectations live in the `Kubesonde` resource, you can keep it in version control and re-apply it after any policy change as a lightweight connectivity regression test.

## Clean up

```bash
kubectl delete -f shop-audit.yaml
kubectl delete -f db-isolation.yaml
kubectl delete namespace shop
```

