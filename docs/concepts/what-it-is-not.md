# What Kubesonde is not

Kubesonde is a **diagnostic and auditing tool**. Being clear about what it does *not* do prevents dangerous misunderstandings about your cluster's security posture.

## Kubesonde does not enforce anything

Kubesonde is **not** a network policy engine, admission controller, or firewall. It does **not** block or modify traffic. Deploying Kubesonde changes nothing about what your pods are allowed to do — it only *observes* and *reports*.

## Kubesonde does not replace your CNI or policies

It does **not** replace tools like **Cilium**, **Calico**, or Kubernetes **`NetworkPolicies`**. Those tools *enforce* connectivity rules. Kubesonde *measures* the result of that enforcement. The two are complementary:

| Tool | Role |
| --- | --- |
| Cilium / Calico / NetworkPolicy | **Enforce** the connectivity you intend |
| Kubesonde | **Verify** that the enforcement matches your intent |

## Kubesonde is not a one-shot report you can ignore

Because connectivity depends on live cluster state, a scan reflects a moment in time. Kubesonde keeps probing (see [completeness](/concepts/completeness)), and you should re-run audits after changes to policies, workloads, or the CNI.

## What it *is*

A tool to find the gap between the connectivity you *think* you have and the connectivity you *actually* have — so you can:

- Verify isolation policies work as intended.
- Find unexpected reachability before an attacker does.
- Understand the real communication footprint of an application.

See [Intended vs. actual connectivity](/concepts/intended-vs-actual) for the motivation.
