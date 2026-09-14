# Security & privileges

Kubesonde probes real connectivity by injecting two [ephemeral containers](https://kubernetes.io/docs/concepts/workloads/pods/ephemeral-containers/) — a **debugger** (runs `nmap`/`wget`) and a **monitor** (reads the pod's listening sockets) — into the pods it targets. To do their job these containers need elevated privileges, so **Kubesonde is best run in a test or staging cluster first**, with the resulting policies then propagated to production.

## What Kubesonde requires from a pod

The debugger and monitor containers are injected with:

- `securityContext.privileged: true`
- effective **root** inside the target pod's namespaces (they read `/proc/net`, run raw-socket scans, and reach the pod's network stack)

This is intentional: to observe what a pod can *actually* reach on the network, the probe has to run from inside that pod with enough privilege to open raw sockets and inspect kernel networking state.

## Why this matters

Injecting a privileged, root-capable container into a workload is a meaningful privilege escalation surface. Anyone who can create a `Kubesonde` resource can effectively obtain root inside every pod in the target namespace. Treat the ability to create `Kubesonde` resources as equivalent to granting root on those workloads, and scope RBAC accordingly.

There is also a practical consequence: **hardened pods reject the probe.** A pod with a restrictive `securityContext` — for example `runAsNonRoot: true` — will refuse to start Kubesonde's ephemeral containers, and Kubernetes reports:

```
container has runAsNonRoot and image will run as root
reason: CreateContainerConfigError
```

When this happens the pod cannot be probed from the inside. Kubesonde still records outside-in probes against it (a `404`, for instance, means the port is open and answering), but the pod will be missing from the per-pod listening-socket data (`podNetworking`), which can make an open port look closed or unverified in the UI.

## Recommended workflow

Because of the privilege requirements, adopt Kubesonde in a **test cluster** and promote the *outputs* — not the tool — to production:

1. **Run Kubesonde in a test/staging cluster** that mirrors your production topology (same namespaces, workloads, and NetworkPolicies). Elevated privileges are acceptable here because the cluster is non-production and isolated.
2. **Scan and analyze** the intended-vs-actual connectivity, then author or refine your `NetworkPolicy` (or other controls) from the results. See [Auditing a NetworkPolicy](/tutorials/audit-networkpolicy).
3. **Propagate the resulting policies to production** through your normal GitOps/CD pipeline. Only the reviewed `NetworkPolicy` manifests move forward — Kubesonde and its privileged probes stay in the test cluster.
4. **Re-run in test on change.** When workloads or policies change, repeat the scan in the test cluster and promote the updated policies, rather than probing production directly.

::: tip
If you must run Kubesonde in a shared or sensitive cluster, restrict who can create `Kubesonde` resources via RBAC, and limit `spec.namespace` to namespaces you explicitly intend to probe. Expect hardened pods (`runAsNonRoot`, restricted PodSecurity standards) to be un-probeable from the inside.
:::
