# Introduction

**Kubesonde is a tool that probes and visualizes the _actual_ network connectivity of applications running in a Kubernetes cluster**, so you can compare it against the network policies you meant to enforce. It works by instrumenting live pods at runtime and reporting every connection attempt it observes, rather than relying only on what your `NetworkPolicy` manifests declare.

Kubesonde answers a concrete question that manifests alone cannot:

> _From this pod, what can I actually reach right now?_

It probes live pods, records every connection attempt and its outcome, and lets you visualize the resulting connectivity graph.

## Why Kubesonde

Kubernetes network security is declarative: you write `NetworkPolicies` and trust that the cluster enforces them. In practice, the gap between *intended* and *actual* connectivity is where misconfigurations hide.

- **NetworkPolicies are hard to get right.** A missing selector, an overly broad `podSelector`, or a forgotten default-deny policy can silently leave pods reachable when they should be isolated. Reading the manifests does not tell you what traffic is really allowed.
- **The enforced state can differ from the declared state.** CNI plugins, admission controllers, and overlapping policies interact in non-obvious ways. The only way to know what a pod can actually reach is to try.
- **Auditing connectivity by hand does not scale.** In a namespace with dozens of microservices, manually reasoning about which pod can talk to which — and on which ports — is error-prone and quickly becomes impossible.

## What you can use it for

- **Verify isolation policies** work as intended after you write or change a `NetworkPolicy`.
- **Find unexpected reachability** before an attacker does.
- **Understand the real communication footprint** of an application you inherited or are onboarding.
- **Regression-test connectivity** by asserting expected `Allow`/`Deny` outcomes.

## What Kubesonde is not

Kubesonde is **not** a network policy engine, admission controller, or firewall. It does **not** block or modify traffic, and it does **not** replace tools like Cilium, Calico, or Kubernetes `NetworkPolicies`. It is a diagnostic and auditing tool for finding gaps between the connectivity you think you have and the connectivity you actually have.

See [What Kubesonde is not](/concepts/what-it-is-not) for the full boundary.

## Publications

Kubesonde has been described and used in the following peer-reviewed papers:

- Jacopo Bufalino, Mario Di Francesco, Tuomas Aura. **[Analyzing Microservice Connectivity with Kubesonde](https://dl.acm.org/doi/10.1145/3611643.3613899)**. ESEC/FSE 2023.
- Jacopo Bufalino, Jose Luis Martin-Navarro, Mario Di Francesco, Tuomas Aura. **[Inside Job: Defending Kubernetes Clusters Against Network Misconfigurations](https://dl.acm.org/doi/10.1145/3749220)**. Proceedings of the ACM on Networking, 2025.
